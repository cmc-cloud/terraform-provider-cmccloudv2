package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApiGatewayBackup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiGatewayBackupCreate,
		Read:   resourceApiGatewayBackupRead,
		Delete: resourceApiGatewayBackupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceApiGatewayBackupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        apigatewayBackupSchema(),
	}
}

func resourceApiGatewayBackupCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.ApiGateway.CreateBackup(d.Get("instance_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating ApiGateway Backup: %s", err)
	}
	d.SetId(instance.Data.ID)

	_, err = waitUntilApiGatewayBackupJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ApiGatewayBackup: %s", err)
	}
	return resourceApiGatewayBackupRead(d, meta)
}

func resourceApiGatewayBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.ApiGateway.GetBackup(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ApiGatewayBackup %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("instance_id", instance.Instance.ID)
	_ = d.Set("size", instance.Size)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.CreatedAt)
	return nil
}

func resourceApiGatewayBackupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ApiGateway.DeleteBackup(d.Id())

	if err != nil {
		return fmt.Errorf("error delete Api Gateway backup: %v", err)
	}
	_, err = waitUntilApiGatewayBackupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete Api Gateway backup: %v", err)
	}
	return nil
}

func resourceApiGatewayBackupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceApiGatewayBackupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilApiGatewayBackupJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"COMPLETED"}, []string{"ERROR", "FAILED"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ApiGateway.GetBackup(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.ApiGatewayBackup).Status)
	})
}

func waitUntilApiGatewayBackupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ApiGateway.GetBackup(id)
	})
}
