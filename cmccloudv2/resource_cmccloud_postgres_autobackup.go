package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresAutoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresAutoBackupCreate,
		Read:   resourcePostgresAutoBackupRead,
		Update: resourcePostgresAutoBackupUpdate,
		Delete: resourcePostgresAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresAutoBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        postgresAutoBackupSchema(),
	}
}
func resourcePostgresAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	vol, err := client.PostgresAutoBackup.Create(map[string]interface{}{
		"instanceId":       d.Get("instance_id").(string),
		"hour":             hour,
		"minute":           minute,
		"second":           0,
		"interval":         d.Get("interval").(int),
		"keepRecordBackup": d.Get("max_keep").(int),
		"timeZone":         "GMT+07:00",
	})

	if err != nil {
		return fmt.Errorf("error creating Postgres AutoBackup: %s", err)
	}
	d.SetId(vol.ID)

	return resourcePostgresAutoBackupRead(d, meta)
}

func resourcePostgresAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	autobackup, err := client.PostgresAutoBackup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Postgres AutoBackup %s: %v", d.Id(), err)
	}

	scheduleTimeStr := fmt.Sprintf("%02d:%02d", autobackup.Hour, autobackup.Minute)
	_ = d.Set("schedule_time", scheduleTimeStr)
	_ = d.Set("interval", autobackup.IntervalNum)
	_ = d.Set("max_keep", autobackup.KeepRecordBackup)
	_ = d.Set("created", autobackup.Created)
	_ = d.Set("next_run", autobackup.NextBackupTime)
	return nil
}

func resourcePostgresAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	if d.HasChange("schedule_time") || d.HasChange("interval") || d.HasChange("max_keep") || d.HasChange("incremental") {
		_, err := client.PostgresAutoBackup.Update(id, map[string]interface{}{
			"hour":     hour,
			"minute":   minute,
			"interval": d.Get("interval").(int),
			"max_keep": d.Get("max_keep").(int),
		})
		if err != nil {
			return fmt.Errorf("error when update Postgres AutoBackup [%s]: %v", id, err)
		}
	}
	return resourcePostgresAutoBackupRead(d, meta)
}

func resourcePostgresAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.PostgresAutoBackup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete postgres autobackup: %v", err)
	}
	_, err = waitUntilPostgresAutoBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete postgres autobackup: %v", err)
	}
	return nil
}

func resourcePostgresAutoBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresAutoBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilPostgresAutoBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresAutoBackup.Get(id)
	})
}
