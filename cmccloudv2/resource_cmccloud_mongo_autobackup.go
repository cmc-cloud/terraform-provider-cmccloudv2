package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoAutoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoAutoBackupCreate,
		Read:   resourceMongoAutoBackupRead,
		Update: resourceMongoAutoBackupUpdate,
		Delete: resourceMongoAutoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoAutoBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mongoAutoBackupSchema(),
	}
}

func resourceMongoAutoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	vol, err := client.MongoAutoBackup.Create(map[string]interface{}{
		"instanceId":       d.Get("instance_id").(string),
		"hour":             hour,
		"minute":           minute,
		"second":           0,
		"interval":         d.Get("interval").(int),
		"keepRecordBackup": d.Get("max_keep").(int),
		"timeZone":         "GMT+07:00",
	})

	if err != nil {
		return fmt.Errorf("error creating Mongo AutoBackup: %s", err)
	}
	d.SetId(vol.ID)

	return resourceMongoAutoBackupRead(d, meta)
}

func resourceMongoAutoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	autobackup, err := client.MongoAutoBackup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Mongo AutoBackup %s: %v", d.Id(), err)
	}

	scheduleTimeStr := fmt.Sprintf("%02d:%02d", autobackup.Hour, autobackup.Minute)
	_ = d.Set("schedule_time", scheduleTimeStr)
	_ = d.Set("interval", autobackup.IntervalNum)
	_ = d.Set("max_keep", autobackup.KeepRecordBackup)
	_ = d.Set("created", autobackup.Created)
	_ = d.Set("next_run", autobackup.NextBackupTime)
	return nil
}

func resourceMongoAutoBackupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	hour, minute, err := getHourAndMinuteFromConfig(d)
	if err != nil {
		return err
	}

	if d.HasChange("schedule_time") || d.HasChange("interval") || d.HasChange("max_keep") || d.HasChange("incremental") {
		_, err := client.MongoAutoBackup.Update(id, map[string]interface{}{
			"hour":     hour,
			"minute":   minute,
			"interval": d.Get("interval").(int),
			"max_keep": d.Get("max_keep").(int),
		})
		if err != nil {
			return fmt.Errorf("error when update Mongo AutoBackup [%s]: %v", id, err)
		}
	}
	return resourceMongoAutoBackupRead(d, meta)
}

func resourceMongoAutoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.MongoAutoBackup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete mongo autobackup: %v", err)
	}
	_, err = waitUntilMongoAutoBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete mongo autobackup: %v", err)
	}
	return nil
}

func resourceMongoAutoBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMongoAutoBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilMongoAutoBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MongoAutoBackup.Get(id)
	})
}
