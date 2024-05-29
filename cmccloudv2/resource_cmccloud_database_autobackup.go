package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabaseAutoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseAutoBackupCreate,
		Read:   resourceDatabaseAutoBackupRead,
		Update: resourceDatabaseAutoBackupUpdate,
		Delete: resourceDatabaseAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDatabaseAutoBackupImport,
		},
		SchemaVersion: 1,
		Schema:        databaseAutoBackupSchema(),
	}
}

func resourceDatabaseAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	schedule_time := d.Get("schedule_time").(string)
	parts := strings.Split(schedule_time, ":")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid schedule time [%s], correct format is HH:mm (24-h format), eg: 19:05", schedule_time)
	}
	hour := parts[0]
	minute := parts[1]

	vol, err := client.DatabaseAutoBackup.Create(map[string]interface{}{
		"name":        d.Get("name").(string),
		"instance_id": d.Get("instance_id").(string),
		"hour":        hour,
		"minute":      minute,
		"interval":    d.Get("interval").(int),
		"max_keep":    d.Get("max_keep").(int),
		"incremental": d.Get("incremental").(bool),
	})
	if err != nil {
		return fmt.Errorf("Error creating Database AutoBackup: %s", err)
	}
	d.SetId(vol.ID)

	return resourceDatabaseAutoBackupRead(d, meta)
}

func resourceDatabaseAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	autobackup, err := client.DatabaseAutoBackup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Database AutoBackup %s: %v", d.Id(), err)
	}

	_ = d.Set("name", autobackup.Name)
	_ = d.Set("schedule_time", autobackup.Time)
	_ = d.Set("interval", autobackup.Interval)
	_ = d.Set("max_keep", autobackup.MaxKeep)
	_ = d.Set("incremental", !autobackup.IsFullBackup)
	_ = d.Set("created_at", autobackup.Created)
	_ = d.Set("last_run", autobackup.LastRun)
	_ = d.Set("volume_size", autobackup.VolumeSize)
	return nil
}

func resourceDatabaseAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	schedule_time := d.Get("schedule_time").(string)
	parts := strings.Split(schedule_time, ":")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid schedule time [%s], correct format is HH:mm (24-h format), eg: 19:05", schedule_time)
	}
	hour := parts[0]
	minute := parts[1]

	if d.HasChange("name") || d.HasChange("schedule_time") || d.HasChange("interval") || d.HasChange("max_keep") || d.HasChange("incremental") {
		_, err := client.DatabaseAutoBackup.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"hour":        hour,
			"minute":      minute,
			"interval":    d.Get("interval").(int),
			"max_keep":    d.Get("max_keep").(int),
			"incremental": d.Get("incremental").(bool),
		})
		if err != nil {
			return fmt.Errorf("Error when update Database AutoBackup [%s]: %v", id, err)
		}
	}
	return resourceDatabaseAutoBackupRead(d, meta)
}

func resourceDatabaseAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DatabaseAutoBackup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete database autobackup: %v", err)
	}
	return nil
}

func resourceDatabaseAutoBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDatabaseAutoBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}
