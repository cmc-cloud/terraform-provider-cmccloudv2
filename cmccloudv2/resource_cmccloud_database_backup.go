package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseBackupCreate,
		Read:   resourceDatabaseBackupRead,
		Update: resourceDatabaseBackupUpdate,
		Delete: resourceDatabaseBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDatabaseBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        databaseBackupSchema(),
	}
}

func resourceDatabaseBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.DatabaseInstance.CreateBackup(d.Get("instance_id").(string), map[string]interface{}{
		"name":        d.Get("name").(string),
		"incremental": d.Get("incremental").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error creating backup of database instance [%s]: %v", d.Get("instance_id").(string), err)
	}
	d.SetId(backup.ID)
	_, err = waitUntilDatabaseBackupStatusChangedState(d, meta, []string{"avaiable"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("Error creating backup of database instance [%s]: %v", d.Get("instance_id").(string), err)
	}
	return resourceDatabaseBackupRead(d, meta)
}

func resourceDatabaseBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.DatabaseBackup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving database instance backup %s: %v", d.Id(), err)
	}

	_ = d.Set("name", backup.Name)
	_ = d.Set("instance_id", backup.InstanceID)
	_ = d.Set("volume_size", backup.Size)
	_ = d.Set("real_size_gb", backup.RealSizeGB) // round 2 decimal
	_ = d.Set("status", backup.Status)
	_ = d.Set("created_at", backup.Created)
	if len(backup.ParentID) == 0 {
		_ = d.Set("incremental", false)
	} else {
		_ = d.Set("incremental", true)
	}

	return nil
}

func resourceDatabaseBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.DatabaseBackup.Update(id, map[string]interface{}{
			"name": d.Get("name").(string),
		})
		if err != nil {
			return fmt.Errorf("Error when rename database instance backup [%s]: %v", id, err)
		}
	}
	return resourceDatabaseBackupRead(d, meta)
}

func resourceDatabaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Backup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete database instance backup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilDatabaseBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete database instance backup [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceDatabaseBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDatabaseBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDatabaseBackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 10 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DatabaseBackup.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.DatabaseBackup).Status
	})
}

func waitUntilDatabaseBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DatabaseBackup.Get(id)
	})
}
