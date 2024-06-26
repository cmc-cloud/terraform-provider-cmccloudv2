package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeAttachmentCreate,
		Read:   resourceVolumeAttachmentRead,
		Delete: resourceVolumeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeAttachmentImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        volumeAttachmentSchema(),
	}
}

func resourceVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server_id := d.Get("server_id").(string)
	_, err := client.Volume.Attach(d.Get("volume_id").(string), map[string]interface{}{
		"server_id":             server_id,
		"delete_on_termination": d.Get("delete_on_termination").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error when attach Volume %s to Server %s: %s", d.Get("volume_id").(string), server_id, err)
	}

	d.SetId(d.Get("volume_id").(string))

	_, err = waitUntilVolumeAttachedStateChanged(d, meta, server_id, []string{"", "Detached"}, []string{"Attached"})
	if err != nil {
		return fmt.Errorf("[ERROR] Error attach volume %s to server %s: %v", d.Id(), server_id, err)
	}
	return resourceVolumeAttachmentRead(d, meta)
}

func resourceVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	volumeID := d.Id()
	vol, err := client.Server.GetVolumeAttachmentDetail(d.Get("server_id").(string), volumeID)
	if err != nil {
		return fmt.Errorf("Error retrieving Volume Attachment %s: %v", d.Id(), err)
	}
	_ = d.Set("server_id", vol.ServerID)
	_ = d.Set("volume_id", volumeID)
	_ = d.Set("delete_on_terminated", vol.DeleteOnTermination)
	return nil
}

func resourceVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server_id := d.Get("server_id").(string)
	_, err := client.Volume.Detach(d.Id(), server_id)

	if err != nil {
		return fmt.Errorf("[ERROR] Error detaching volume %s from server %s: %v", d.Id(), server_id, err)
	}
	// wait until detached
	_, err = waitUntilVolumeAttachedStateChanged(d, meta, server_id, []string{"", "Attached"}, []string{"Detached"})
	if err != nil {
		return fmt.Errorf("[ERROR] Error detaching volume %s from server %s: %v", d.Id(), server_id, err)
	}
	return nil
}

func resourceVolumeAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeAttachmentRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVolumeAttachedStateChanged(d *schema.ResourceData, meta interface{}, server_id string, pendingStatus []string, targetStatus []string) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        pendingStatus,
		Target:         targetStatus,
		Refresh:        volumeAttachedStateRefreshfunc(d, meta, server_id),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          2 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 5,
	}
	return stateConf.WaitForState()
}

func volumeAttachedStateRefreshfunc(d *schema.ResourceData, meta interface{}, server_id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*CombinedConfig).goCMCClient()
		volume, err := client.Volume.Get(d.Id())
		if err != nil {
			fmt.Errorf("Error retrieving volume %s: %v", d.Id(), err)
			return nil, "", err
		}
		for _, attachment := range volume.Attachments {
			if attachment.ServerID == server_id {
				return volume, "Attached", nil
			}
		}
		return volume, "Detached", nil
	}
}
