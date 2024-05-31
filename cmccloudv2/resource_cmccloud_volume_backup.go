package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolumeBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeBackupCreate,
		Read:   resourceVolumeBackupRead,
		Update: resourceVolumeBackupUpdate,
		Delete: resourceVolumeBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVolumeBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        volumeBackupSchema(),
	}
}

func resourceVolumeBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.Volume.CreateBackup(d.Get("volume_id").(string), map[string]interface{}{
		"name":        d.Get("name").(string),
		"force":       d.Get("force").(bool),
		"incremental": d.Get("incremental").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error creating backup of volume [%s]: %v", d.Get("volume_id").(string), err)
	}
	d.SetId(backup.ID)
	_, err = waitUntilVolumeBackupStatusChangedState(d, meta, []string{"available"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("Error creating backup of volume [%s]: %v", d.Get("volume_id").(string), err)
	}
	return resourceVolumeBackupRead(d, meta)
}

func resourceVolumeBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.Backup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving backup %s: %v", d.Id(), err)
	}

	// real_size := float64(backup.RealSize) / (1024 * 1024 * 1024)
	// real_size_round, _ := fmt.Printf("%.2f", real_size)
	_ = d.Set("name", backup.Name)
	_ = d.Set("volume_id", backup.VolumeID)
	_ = d.Set("volume_size", backup.Size)
	_ = d.Set("real_size_gb", backup.RealSizeGB) // round 2 decimal
	_ = d.Set("incremental", backup.IsIncremental)
	_ = d.Set("status", backup.Status)
	_ = d.Set("created_at", backup.CreatedAt)
	if backup.Volume.Name != "" {
		_ = d.Set("volume_name", backup.Volume.Name)
	}

	return nil
}

func resourceVolumeBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.Backup.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Backup [%s]: %v", id, err)
		}
	}
	return resourceVolumeBackupRead(d, meta)
}

func resourceVolumeBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Backup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete volume backup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilVolumeBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete volume backup: %v", err)
	}
	return nil
}

func resourceVolumeBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVolumeBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVolumeBackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Backup.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Backup).Status
	})
}

func waitUntilVolumeBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Backup.Get(id)
	})
}
