package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageCreate,
		Read:   resourceImageRead,
		Update: resourceImageUpdate,
		Delete: resourceImageDelete,
		Importer: &schema.ResourceImporter{
			State: resourceImageImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        imageSchema(),
	}
}

func resourceImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	image, err := client.Image.CreateFromVolume(d.Get("volume_id").(string), map[string]interface{}{
		"image_name":  d.Get("name").(string),
		"disk_format": d.Get("disk_format").(string),
		"force":       d.Get("force").(bool),
	})
	if err != nil {
		return fmt.Errorf("error creating image: %s", err)
	}
	d.SetId(image.ID)
	_, err = waitUntilImageStatusChangedState(d, meta, []string{"active"}, []string{"killed", "deleted"})
	if err != nil {
		return fmt.Errorf("error creating image: %s", err)
	}
	return resourceImageRead(d, meta)
}

func resourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	image, err := client.Image.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Image %s: %v", d.Id(), err)
	}

	_ = d.Set("name", image.Name)
	_ = d.Set("status", image.Status)
	_ = d.Set("created_at", image.CreatedAt)
	return nil
}

func resourceImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.Image.Update(id, map[string]interface{}{
			"name": d.Get("name").(string),
		})
		if err != nil {
			return fmt.Errorf("error when update Image [%s]: %v", id, err)
		}
	}
	return resourceImageRead(d, meta)
}

func resourceImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Image.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete image: %v", err)
	}
	_, err = waitUntilImageDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete image: %v", err)
	}
	return nil
}

func resourceImageImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceImageRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilImageDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Image.Get(id)
	})
}

func waitUntilImageStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      30 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Image.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Image).Status
	})
}
