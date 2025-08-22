package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
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
		return fmt.Errorf("error creating snapshot: %s", err)
	}
	d.SetId(snapshot.ID)
	_, err = waitUntilVolumeSnapshotStatusChangedState(d, meta, []string{"available"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("error creating snapshot: %s", err)
	}

	return resourceVolumeSnapshotRead(d, meta)
}

func resourceVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	snapshot, err := client.Snapshot.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Snapshot %s: %v", d.Id(), err)
	}

	_ = d.Set("name", snapshot.Name)
	_ = d.Set("volume_id", snapshot.VolumeID)
	_ = d.Set("volume_size", snapshot.Size)
	_ = d.Set("real_size_gb", snapshot.RealSizeGB)
	_ = d.Set("status", snapshot.Status)
	_ = d.Set("created_at", snapshot.CreatedAt)
	_ = d.Set("volume_name", snapshot.Volume.Name)

	return nil
}

func resourceVolumeSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.Snapshot.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("error when rename volume snapshot [%s]: %v", id, err)
		}
	}
	return resourceVolumeSnapshotRead(d, meta)
}

func resourceVolumeSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Snapshot.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete volume snapshot: %v", err)
	}
	_, err = waitUntilVolumeSnapshotDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete volume snapshot: %v", err)
	}
	return nil
}

func resourceVolumeSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeSnapshotRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVolumeSnapshotStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Snapshot.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Snapshot).Status
	})
}

func waitUntilVolumeSnapshotDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 15 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Snapshot.Get(id)
	})
}
