package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDBv2AutoBackup(dbType string) *schema.Resource {
	return &schema.Resource{
		Create: resourceDBv2AutoBackupCreate,
		Read:   resourceDBv2AutoBackupRead,
		Update: resourceDBv2AutoBackupUpdate,
		Delete: resourceDBv2AutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDBv2AutoBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        dbv2AutoBackupSchema(),
	}
}

func resourceDBv2AutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	vol, err := client.DBv2AutoBackup.Create(map[string]interface{}{
		"instanceId":       d.Get("instance_id").(string),
		"hour":             hour,
		"minute":           minute,
		"second":           0,
		"interval":         d.Get("interval").(int),
		"keepRecordBackup": d.Get("max_keep").(int),
		"timeZone":         "GMT+07:00",
	})

	if err != nil {
		return fmt.Errorf("error creating AutoBackup: %s", err)
	}
	d.SetId(vol.Data.ID)

	return resourceDBv2AutoBackupRead(d, meta)
}

func resourceDBv2AutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	autobackup, err := client.DBv2AutoBackup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving AutoBackup %s: %v", d.Id(), err)
	}

	scheduleTimeStr := fmt.Sprintf("%02d:%02d", autobackup.Hour, autobackup.Minute)
	_ = d.Set("schedule_time", scheduleTimeStr)
	_ = d.Set("interval", autobackup.IntervalNum)
	_ = d.Set("max_keep", autobackup.KeepRecordBackup)
	_ = d.Set("created", autobackup.Created)
	_ = d.Set("next_run", autobackup.NextBackupTime)
	return nil
}

func resourceDBv2AutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	if d.HasChange("schedule_time") || d.HasChange("interval") || d.HasChange("max_keep") {
		_, err := client.DBv2AutoBackup.Update(id, map[string]interface{}{
			"hour":     hour,
			"minute":   minute,
			"interval": d.Get("interval").(int),
			"max_keep": d.Get("max_keep").(int),
		})
		if err != nil {
			return fmt.Errorf("error when update AutoBackup [%s]: %v", id, err)
		}
	}
	return resourceDBv2AutoBackupRead(d, meta)
}

func resourceDBv2AutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DBv2AutoBackup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete autobackup: %v", err)
	}
	_, err = waitUntilDBv2AutoBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete autobackup: %v", err)
	}
	return nil
}

func resourceDBv2AutoBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDBv2AutoBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDBv2AutoBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBv2AutoBackup.Get(id)
	})
}
