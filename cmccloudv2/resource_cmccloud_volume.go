package cmccloudv2

import (
	"fmt"
	"time"

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
		},
		SchemaVersion: 1,
		Schema:        volumeSchema(),

		// CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
		// 	// if the new size of the volume is smaller than the old one return an error since
		// 	// only expanding the volume is allowed
		// 	oldSize, newSize := diff.GetChange("size")
		// 	if newSize.(int) < oldSize.(int) {
		// 		return fmt.Errorf("volumes `size` can only be expanded and not shrunk")
		// 	}

		// 	return nil
		// },
		CustomizeDiff: customdiff.All(
			// customdiff.ValidateChange("size", func (old, new, meta interface{}) error {
			//     // If we are increasing "size" then the new value must be
			//     // a multiple of the old value.
			//     if new.(int) <= old.(int) {
			//         return fmt.Errorf("volumes `size` can only be expanded and not shrunk")
			//     }
			//     return nil
			// }),
			customdiff.ForceNewIfChange("size", func(old, new, meta interface{}) bool {
				// "size" can only increase in-place, so we must create a new resource
				// if it is decreased.
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
		"billing_mode": d.Get("billing_mode").(string),
		"tags":         d.Get("tags").(*schema.Set).List(),
	})
	if err != nil {
		return fmt.Errorf("Error creating Volume: %s", err)
	}
	d.SetId(vol.ID)

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	volume, err := client.Volume.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Volume %s: %v", d.Id(), err)
	}

	_ = d.Set("name", volume.Name)
	_ = d.Set("description", volume.Description)
	_ = d.Set("size", volume.Size)
	_ = d.Set("type", volume.VolumeType)
	_ = d.Set("zone", volume.AvailabilityZone)
	_ = d.Set("billing_mode", volume.BillingMode)
	_ = d.Set("status", volume.Status)
	_ = d.Set("tags", volume.Tags)
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
			return fmt.Errorf("Error when update Volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("size") {
		_, err := client.Volume.Resize(id, d.Get("size").(int))
		if err != nil {
			return fmt.Errorf("Error when resize volume [%s]: %v", id, err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetVolumeBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when update billing mode of volume [%s]: %v", id, err)
		}
	}
	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Volume.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud volume: %v", err)
	}
	return nil
}

func resourceVolumeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeRead(d, meta)
	return []*schema.ResourceData{d}, err
}
