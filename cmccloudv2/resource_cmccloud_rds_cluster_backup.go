package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRdsClusterBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsClusterBackupCreate,
		Read:   resourceRdsClusterBackupRead,
		Delete: resourceRdsClusterBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRdsClusterBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        rdsclusterBackupSchema(),
	}
}

func resourceRdsClusterBackupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.RdsCluster.CreateBackup(d.Get("cluster_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating backup of rds cluster [%s]: %v", d.Get("cluster_id").(string), err)
	}
	d.SetId(backup.ID)
	_, err = waitUntilRdsClusterBackupStatusChangedState(d, meta, []string{"SUCCEEDED"}, []string{"FAILED", "ERROR"}) // Starting / initializing / Running / InProgress / Succeeded / Completed / Failed / Error / Waiting / Deleting
	if err != nil {
		return fmt.Errorf("error creating backup of rds cluster [%s]: %v", d.Get("cluster_id").(string), err)
	}
	return resourceRdsClusterBackupRead(d, meta)
}

func resourceRdsClusterBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	backup, err := client.RdsCluster.GetBackup(d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving backup %s: %v", d.Id(), err)
	}
	_ = d.Set("name", backup.Name)
	_ = d.Set("size", backup.Size)
	_ = d.Set("created", backup.Created)

	return nil
}

func resourceRdsClusterBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.RdsCluster.DeleteBackup(d.Get("cluster_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete rds cluster backup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilRdsClusterBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete rds cluster backup: %v", err)
	}
	return nil
}

func resourceRdsClusterBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRdsClusterBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilRdsClusterBackupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RdsCluster.GetBackup(d.Get("cluster_id").(string), d.Id())
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.RdsClusterBackup).Status)
	})
}

func waitUntilRdsClusterBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RdsCluster.GetBackup(d.Get("cluster_id").(string), d.Id())
	})
}
