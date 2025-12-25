package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        volumeSchema(),
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("size", func(old, new, meta interface{}) bool {
				return new.(int) < old.(int)
			}),
		),
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vol, err := client.Volume.Create(map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"size":         d.Get("size").(int),
		"type":         d.Get("type").(string),
		"zone_name":    d.Get("zone").(string),
		"backup_id":    d.Get("backup_id").(string),
		"billing_mode": d.Get("billing_mode").(string),
		"secret":       d.Get("secret_id").(string),
		"tags":         d.Get("tags").(*schema.Set).List(),
	})
	if err != nil {
		return fmt.Errorf("error creating volume: %s", err)
	}
	d.SetId(vol.ID)
	_, err = waitUntilVolumeStatusChangedState(d, meta, []string{"available"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("error creating volume: %s", err)
	}
	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	volume, err := client.Volume.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Volume %s: %v", d.Id(), err)
	}

	_ = d.Set("name", volume.Name)
	_ = d.Set("description", volume.Description)
	_ = d.Set("size", volume.Size)
	_ = d.Set("type", volume.VolumeType)
	_ = d.Set("zone", volume.AvailabilityZone)
	_ = d.Set("billing_mode", volume.BillingMode)
	_ = d.Set("status", volume.Status)
	setString(d, "secret_id", volume.EncryptionKeyID)
	_ = d.Set("tags", convertTagsToSet(volume.Tags))
	_ = d.Set("created_at", volume.CreatedAt)
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("tags") {
		// Resize Volume to new flavor
		_, err := client.Volume.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"tags":        d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("error when update Volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("size") {
		_, err := client.Volume.Resize(id, d.Get("size").(int))
		if err != nil {
			return fmt.Errorf("error when resize volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetVolumeBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("error when update billing mode of volume [%s]: %v", id, err)
		}
	}
	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Volume.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete volume: %v", err)
	}
	_, err = waitUntilVolumeDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete volume: %v", err)
	}
	return nil
}

func resourceVolumeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVolumeDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      5 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Volume.Get(id)
	})
}

func waitUntilVolumeStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Volume.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Volume).Status
	})
}
