package cmccloudv2

import (
	"errors"
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeSnapshotCreate,
		Read:   resourceVolumeSnapshotRead,
		Update: resourceVolumeSnapshotUpdate,
		Delete: resourceVolumeSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeSnapshotImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        volumeSnapshotSchema(),
	}
}

func resourceVolumeSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	snapshot, err := client.Volume.CreateSnapshot(d.Get("volume_id").(string), map[string]interface{}{
		"name":  d.Get("name").(string),
		"force": d.Get("force").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error creating Snapshot: %s", err)
	}
	d.SetId(snapshot.ID)
	waitUntilVolumeSnapshotCreated(d, meta)
	return resourceVolumeSnapshotRead(d, meta)
}

func resourceVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	snapshot, err := client.Snapshot.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Snapshot %s: %v", d.Id(), err)
	}

	_ = d.Set("name", snapshot.Name)
	_ = d.Set("volume_id", snapshot.VolumeID)
	_ = d.Set("volume_size", snapshot.Size)
	_ = d.Set("real_size_gb", snapshot.RealSizeGB)
	_ = d.Set("status", snapshot.Status)
	_ = d.Set("created_at", snapshot.CreatedAt)
	if snapshot.Volume.Name != "" {
		_ = d.Set("volume_name", snapshot.Volume.Name)
	}

	return nil
}

func resourceVolumeSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.Snapshot.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Snapshot [%s]: %v", id, err)
		}
	}
	return resourceVolumeSnapshotRead(d, meta)
}

func resourceVolumeSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Snapshot.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud volume: %v", err)
	}
	return nil
}

func resourceVolumeSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeSnapshotRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVolumeSnapshotCreated(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"creating"},
		Target:         []string{"avaiable", "error"},
		Refresh:        createVolumeSnapshotStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func createVolumeSnapshotStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).goCMCClient()
	return func() (interface{}, string, error) {
		backup, err := client.Snapshot.Get(d.Id())

		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			return backup, "", nil
		}

		return backup, backup.Status, nil

	}
}
func waitUntilVolumeSnapshotDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        deleteVolumeSnapshotStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func deleteVolumeSnapshotStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).goCMCClient()
	return func() (interface{}, string, error) {
		backup, err := client.Snapshot.Get(d.Id())

		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			return backup, "true", nil
		}

		return backup, "", nil

	}
}
