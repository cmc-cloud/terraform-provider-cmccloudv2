package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoBackupCreate,
		Read:   resourceMongoBackupRead,
		Delete: resourceMongoBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mongoBackupSchema(),
	}
}

func resourceMongoBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.MongoInstance.CreateBackup(d.Get("instance_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating backup of mongo [%s]: %v", d.Get("mongo_id").(string), err)
	}
	// d.SetId(backup.ID)
	_, err = waitUntilMongoBackupStatusChangedState(d, meta, []string{"available"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("error creating backup of mongo [%s]: %v", d.Get("mongo_id").(string), err)
	}
	return resourceMongoBackupRead(d, meta)
}

func resourceMongoBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.MongoInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving backup %s: %v", d.Id(), err)
	}
	_ = d.Set("name", backup.BackupName)
	_ = d.Set("size", backup.Size)
	_ = d.Set("created", backup.Created)

	return nil
}

func resourceMongoBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.MongoInstance.DeleteBackup(d.Id())

	if err != nil {
		return fmt.Errorf("error delete mongo backup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilMongoBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete mongo backup: %v", err)
	}
	return nil
}

func resourceMongoBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMongoBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilMongoBackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MongoInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Backup).Status
	})
}

func waitUntilMongoBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MongoInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	})
}
