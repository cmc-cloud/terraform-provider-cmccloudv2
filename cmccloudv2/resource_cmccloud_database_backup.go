package cmccloudv2

import (
	"errors"
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	backup, err := client.Volume.CreateBackup(d.Get("instance_id").(string), map[string]interface{}{
		"name":        d.Get("name").(string),
		"incremental": d.Get("incremental").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error creating backup of database instance [%s]: %v", d.Get("instance_id").(string), err)
	}
	d.SetId(backup.ID)
	waitUntilDatabaseBackupCreated(d, meta)
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
		_, err := client.Backup.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename Backup [%s]: %v", id, err)
		}
	}
	return resourceDatabaseBackupRead(d, meta)
}

func resourceDatabaseBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Backup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud Backup [%s]: %v", d.Id(), err)
	}
	waitUntilDatabaseBackupDeleted(d, meta)
	return nil
}

func resourceDatabaseBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDatabaseBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDatabaseBackupCreated(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"creating"},
		Target:         []string{"avaiable", "error"},
		Refresh:        createDatabaseBackupStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func createDatabaseBackupStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).goCMCClient()
	return func() (interface{}, string, error) {
		backup, err := client.Backup.Get(d.Id())

		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			return backup, "", nil
		}

		return backup, backup.Status, nil

	}
}
func waitUntilDatabaseBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        deleteDatabaseBackupStateRefreshFunc(d, meta),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          30 * time.Second,
		MinTimeout:     20 * time.Second,
		NotFoundChecks: 50,
	}
	return stateConf.WaitForState()
}

func deleteDatabaseBackupStateRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*CombinedConfig).goCMCClient()
	return func() (interface{}, string, error) {
		backup, err := client.Backup.Get(d.Id())

		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			return backup, "true", nil
		}

		return backup, "", nil

	}
}
