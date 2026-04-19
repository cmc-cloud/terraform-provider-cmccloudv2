package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDBv2Backup(dbType string) *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2BackupCreate(d, meta, dbType)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2BackupRead(d, meta, dbType)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2BackupDelete(d, meta, dbType)
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return resourceDBv2BackupImport(d, meta, dbType)
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        dbv2BackupSchema(dbType),
	}
}

func resourceDBv2BackupCreate(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.DBv2Backup.Create(d.Get("instance_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating backup of [%s]: %v", d.Get("instance_id").(string), err)
	}
	d.SetId(backup.Data.BackupID)
	_, err = waitUntilDBv2BackupStatusChangedState(d, meta, []string{"COMPLETED"}, []string{"ERROR", "FAILED"}, dbType)
	if err != nil {
		return fmt.Errorf("error creating backup of [%s]: %v", d.Get("instance_id").(string), err)
	}
	return resourceDBv2BackupRead(d, meta, dbType)
}

func resourceDBv2BackupRead(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.DBv2Backup.Get(dbType, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving backup %s: %v", d.Id(), err)
	}
	_ = d.Set("name", backup.BackupName)
	_ = d.Set("size", backup.Size)
	_ = d.Set("created", backup.Created)

	return nil
}

func resourceDBv2BackupDelete(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DBv2Backup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete %s backup [%s]: %v", dbType, d.Id(), err)
	}
	_, err = waitUntilDBv2BackupDeleted(d, meta, dbType)
	if err != nil {
		return fmt.Errorf("error delete %s backup: %v", dbType, err)
	}
	return nil
}

func resourceDBv2BackupImport(d *schema.ResourceData, meta interface{}, dbType string) ([]*schema.ResourceData, error) {
	err := resourceDBv2BackupRead(d, meta, dbType)
	return []*schema.ResourceData{d}, err
}

func waitUntilDBv2BackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, dbType string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBv2Backup.Get(dbType, d.Get("instance_id").(string), d.Id())
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.DBv2Backup).Status)
	})
}

func waitUntilDBv2BackupDeleted(d *schema.ResourceData, meta interface{}, dbType string) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBv2Backup.Get(dbType, d.Get("instance_id").(string), d.Id())
	})
}
