package cmccloudv2

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceStorageGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageGatewayCreate,
		Read:   resourceStorageGatewayRead,
		Update: resourceStorageGatewayUpdate,
		Delete: resourceStorageGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStorageGatewayImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        storageGatewaySchema(),
	}
}

func resourceStorageGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	storage, err := client.StorageGateway.Create(map[string]interface{}{
		"name":          d.Get("name").(string),
		"description":   d.Get("description").(string),
		"protocol_type": strings.ToLower(d.Get("protocol_type").(string)),
		"subnet_id":     d.Get("subnet_id").(string),
		"bucket":        d.Get("bucket").(string),
		"tags":          d.Get("tags").(*schema.Set).List(),
	})
	if err != nil {
		return fmt.Errorf("error creating StorageGateway: %s", err)
	}
	d.SetId(storage.ID)
	_, err = waitUntilStorageGatewayStatusChangedState(d, meta, []string{"ACTIVE"}, []string{"ERROR", "REMOVED"})
	if err != nil {
		return fmt.Errorf("create StorageGateway failed: %v", err)
	}

	return resourceStorageGatewayRead(d, meta)
}

func resourceStorageGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	storage, err := client.StorageGateway.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving StorageGateway %s: %v", d.Id(), err)
	}

	_ = d.Set("id", storage.ID)
	_ = d.Set("name", storage.Name)
	_ = d.Set("description", storage.Description)
	_ = d.Set("protocol_type", storage.ProtocolType)
	_ = d.Set("subnet_id", storage.SubnetID)
	_ = d.Set("bucket", storage.Bucket)
	_ = d.Set("command_line", storage.CommandLine)
	_ = d.Set("shared_path", storage.SharedPath)
	_ = d.Set("created", storage.CreatedAt)
	_ = d.Set("tags", convertTagsToSet(storage.Tags))
	return nil
}

func resourceStorageGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("tags") || d.HasChange("bucket") {
		_, err := client.StorageGateway.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"bucket":      d.Get("bucket").(string),
			"tags":        d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("error when update StorageGateway [%s]: %v", id, err)
		}
	}

	if d.HasChange("cidr") {
		return errors.New("the 'cidr' field cannot be changed after creation")
	}
	return resourceStorageGatewayRead(d, meta)
}

func resourceStorageGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.StorageGateway.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete storage: %v", err)
	}
	_, err = waitUntilStorageGatewayDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete storage: %v", err)
	}
	return nil
}

func resourceStorageGatewayImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceStorageGatewayRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilStorageGatewayStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).StorageGateway.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.StorageGateway).Status)
	})
}

func waitUntilStorageGatewayDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).StorageGateway.Get(id)
	})
}
