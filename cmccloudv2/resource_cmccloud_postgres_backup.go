package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresBackupCreate,
		Read:   resourcePostgresBackupRead,
		Delete: resourcePostgresBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        postgresBackupSchema(),
	}
}

func resourcePostgresBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.PostgresInstance.CreateBackup(d.Get("instance_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating backup of postgres [%s]: %v", d.Get("postgres_id").(string), err)
	}
	// d.SetId(backup.ID)
	_, err = waitUntilPostgresBackupStatusChangedState(d, meta, []string{"available"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("error creating backup of postgres [%s]: %v", d.Get("postgres_id").(string), err)
	}
	return resourcePostgresBackupRead(d, meta)
}

func resourcePostgresBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.PostgresInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving backup %s: %v", d.Id(), err)
	}
	_ = d.Set("name", backup.BackupName)
	_ = d.Set("size", backup.Size)
	_ = d.Set("type", backup.Type)
	_ = d.Set("created", backup.Created)

	return nil
}

func resourcePostgresBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.PostgresInstance.DeleteBackup(d.Id())

	if err != nil {
		return fmt.Errorf("error delete postgres backup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilPostgresBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete postgres backup: %v", err)
	}
	return nil
}

func resourcePostgresBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilPostgresBackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Backup).Status
	})
}

func waitUntilPostgresBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresInstance.GetBackup(d.Get("instance_id").(string), d.Id())
	})
}
