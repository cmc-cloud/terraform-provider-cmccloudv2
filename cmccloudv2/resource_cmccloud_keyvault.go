package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKeyVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyVaultCreate,
		Read:   resourceKeyVaultRead,
		Update: resourceKeyVaultUpdate,
		Delete: resourceKeyVaultDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyVaultImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        keyvaultSchema(),
	}
}

func resourceKeyVaultCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	datastores, err := client.KeyVault.ListDatastore(map[string]string{})
	if err != nil {
		return fmt.Errorf("can't get list of datastore %v", err)
	}

	version := d.Get("version").(string)
	mode := d.Get("mode").(string)

	datastoreVersionId, datastoreModeId, datastoreCode, datastoreTypeId, _, err := findDatastoreInfo(datastores, version, mode)
	if err != nil {
		return err
	}
	zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))
	params := map[string]interface{}{
		"name":            d.Get("name").(string),
		"billingMode":     d.Get("billing_mode").(string),
		"zone":            zones[0],
		"zones":           zones,
		"flavorId":        d.Get("flavor_id").(string),
		"volumeSize":      d.Get("volume_size").(int),
		"volumeType":      d.Get("volume_type").(string),
		"networkId":       subnet.NetworkID,
		"subnetId":        subnet.ID,
		"quantityOfSlave": d.Get("slave_count").(int),
		"proxyQuantity":   d.Get("proxy_quantity").(int),
		"proxyFlavorId":   d.Get("proxy_flavor_id").(string),
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		"datastore_type": datastoreTypeId,
		// "datastore_version": datastoreVersionId,
	}

	requestMetadata := map[string]interface{}{
		"createType": "",
		"restoreInfo": map[string]interface{}{
			"type":        "",
			"backupId":    "",
			"instanceId":  "",
			"timeRestore": "",
		},
	}

	switch mode {
	case "master_slave":
		// Multi-AZ deployments
		requestMetadata["zones"] = zones
		requestMetadata["quantityOfSlave"] = d.Get("slave_count").(int)
	case "ha_cluster":
		requestMetadata["zones"] = zones
		requestMetadata["quantityOfSlave"] = d.Get("slave_count").(int)
		requestMetadata["proxyQuantity"] = d.Get("proxy_quantity").(int)
		requestMetadata["proxyFlavorId"] = d.Get("proxy_flavor_id").(string)
	default:
		// Standalone node
		requestMetadata["zone"] = zones[0]
	}

	metadataJsonData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal requestMetadata: %s", err)
	}
	params["requestMetadata"] = string(metadataJsonData)

	instance, err := client.KeyVault.Create(params)
	if err != nil {
		return fmt.Errorf("error creating KeyVault Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)

	_, err = client.Tag.UpdateTag(instance.Data.InstanceID, "KeyVault", d)
	if err != nil {
		fmt.Printf("error updating KeyVault tags: %s\n", err)
	}

	_, err = waitUntilKeyVaultJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating KeyVault Instance: %s", err)
	}
	return resourceKeyVaultRead(d, meta)
}

func resourceKeyVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.KeyVault.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving KeyVault Instance %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("version", instance.DatastoreVersion)
	if strings.Contains(strings.ToLower(instance.DatastoreMode), "master") {
		_ = d.Set("mode", "master_slave")
	} else if strings.Contains(strings.ToLower(instance.DatastoreMode), "cluster") {
		_ = d.Set("mode", "ha_cluster")
	} else {
		_ = d.Set("mode", "standalone")
	}

	// _ = d.Set("flavor_id", instance.FlavorInfo.ID)
	// _ = d.Set("slave_count", instance.QuantityOfSlave)
	// _ = d.Set("proxy_quantity", instance.ProxyQuantity)
	// _ = d.Set("proxy_flavor_id", instance.ProxyFlavorID)
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("subnet_id", instance.SubnetID)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	return nil
}

func resourceKeyVaultUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("tags") {
		_, err := client.Tag.UpdateTag(id, "KeyVault", d)
		if err != nil {
			return fmt.Errorf("error when set keyvault tags [%s]: %v", id, err)
		}
	}
	// if d.HasChange("volume_size") {
	// 	_, err := client.KeyVault.ResizeVolume(id, d.Get("volume_size").(int))
	// 	if err != nil {
	// 		return fmt.Errorf("error when resize volume to %s of keyvault instance %s: %v", d.Get("volume_size").(string), id, err)
	// 	}
	// 	_, err = waitUntilKeyVaultJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf("error when resize volume to %s of keyvault instance %s: %v", d.Get("volume_size").(string), id, err)
	// 	}
	// }

	// if d.HasChange("flavor_id") {
	// 	_, err := client.KeyVault.Resize(id, d.Get("flavor_id").(string))
	// 	if err != nil {
	// 		return fmt.Errorf("error when resize flavor to %s of keyvault instance %s: %v", d.Get("flavor_id").(string), id, err)
	// 	}
	// 	_, err = waitUntilKeyVaultJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
	// 	if err != nil {
	// 		return fmt.Errorf("error when resize flavor to %s of keyvault instance %s: %v", d.Get("flavor_id").(string), id, err)
	// 	}
	// }
	return resourceKeyVaultRead(d, meta)
}

func resourceKeyVaultDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.KeyVault.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete keyvault instance: %v", err)
	}
	_, err = waitUntilKeyVaultDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete keyvault instance: %v", err)
	}
	return nil
}

func resourceKeyVaultImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKeyVaultRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKeyVaultJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING"}, []string{"ERROR", "SHUTDOWN"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KeyVault.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.KeyVault).Status)
	})
}

func waitUntilKeyVaultDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KeyVault.Get(id)
	})
}
