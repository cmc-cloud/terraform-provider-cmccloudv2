package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresInstanceCreate,
		Read:   resourcePostgresInstanceRead,
		Update: resourcePostgresInstanceUpdate,
		Delete: resourcePostgresInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresInstanceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        postgresInstanceSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			mode := strings.ToLower(diff.Get("mode").(string))
			proxyQuantity := diff.Get("proxy_quantity").(int)
			proxyFlavorID := diff.Get("proxy_flavor_id").(string)

			// master_slave or ha_cluster requires proxy_quantity
			if mode == "master_slave" || mode == "ha_cluster" {
				if proxyQuantity <= 0 {
					return fmt.Errorf("`proxy_quantity` is required and must be > 0 when mode is `%s`", mode)
				}
			}

			// ha_cluster requires both proxy_quantity and proxy_flavor_id
			if mode == "ha_cluster" {
				if strings.TrimSpace(proxyFlavorID) == "" {
					return fmt.Errorf("`proxy_flavor_id` is required when mode is `ha_cluster`")
				}
				if proxyQuantity <= 0 {
					return fmt.Errorf("`proxy_quantity` is required and must be > 0 when mode is `ha_cluster`")
				}
			}

			return nil
		},
	}
}

func resourcePostgresInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	datastores, err := client.PostgresInstance.ListDatastore(map[string]string{})
	if err != nil {
		return fmt.Errorf("can't get list of datastore %v", err)
	}

	version := d.Get("version").(string)
	mode := d.Get("mode").(string)

	datastoreVersionId, datastoreModeId, datastoreCode, datastoreTypeId, _, err := findPostgresDatastoreInfo(datastores, version, mode)
	if err != nil {
		return err
	}
	zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))
	params := map[string]interface{}{
		"project":              client.Configs.ProjectId,
		"region":               client.Configs.RegionId,
		"name":                 d.Get("name").(string),
		"billingMode":          d.Get("billing_mode").(string),
		"zone":                 zones[0],
		"zones":                zones,
		"flavorId":             d.Get("flavor_id").(string),
		"backupId":             d.Get("backup_id").(string),
		"volumeSize":           d.Get("volume_size").(int),
		"volumeType":           d.Get("volume_type").(string),
		"groupConfigurationId": d.Get("configuration_id").(string),
		"networkId":            subnet.NetworkID,
		"subnetId":             subnet.ID,
		"quantityOfSlave":      d.Get("slave_count").(int),
		"port":                 d.Get("port").(string),
		"adminPassword":        d.Get("admin_password").(string),
		"proxyQuantity":        d.Get("proxy_quantity").(int),
		"proxyFlavorId":        d.Get("proxy_flavor_id").(string),
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		"datastore_type": datastoreTypeId,
		// "datastore_version": datastoreVersionId,
	}
	backupId := d.Get("backup_id").(string)
	if backupId != "" {
		params["source_type"] = "restore"
		params["source_id"] = backupId
	}

	requestMetadata := map[string]interface{}{
		"port":            d.Get("port").(string),
		"adminPassword":   d.Get("admin_password").(string),
		"enablePitr":      true,
		"retentionPeriod": d.Get("retention_period").(int),
		"createType":      "",
		"restoreInfo": map[string]interface{}{
			"type":        "",
			"backupId":    backupId,
			"instanceId":  "",
			"timeRestore": "",
		},
	}
	if backupId != "" {
		requestMetadata["createType"] = "restore"
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

	instance, err := client.PostgresInstance.Create(params)
	if err != nil {
		return fmt.Errorf("error creating PostgresDatabase Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)
	_, err = waitUntilPostgresInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating PostgresDatabase Instance: %s", err)
	}
	return resourcePostgresInstanceRead(d, meta)
}

func resourcePostgresInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.PostgresInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving PostgresDatabase Instance %s: %v", d.Id(), err)
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

	_ = d.Set("flavor_id", instance.FlavorInfo.ID)
	// _ = d.Set("slave_count", instance.QuantityOfSlave)
	// _ = d.Set("proxy_quantity", instance.ProxyQuantity)
	// _ = d.Set("proxy_flavor_id", instance.ProxyFlavorID)
	if len(instance.Connections) > 0 {
		_ = d.Set("port", instance.Connections[0].Port)
	}
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("subnet_id", instance.SubnetID)
	_ = d.Set("configuration_id", instance.GroupConfigID)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	return nil
}

func resourcePostgresInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("configuration_id") {
		_, err := client.PostgresInstance.SetConfigurationGroupId(id, d.Get("configuration_id").(string))
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of postgres database instance %s: %v", d.Get("configuration_id").(string), id, err)
		}
		_, err = waitUntilPostgresInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of postgres database instance %s: %v", d.Get("configuration_id").(string), id, err)
		}
	}

	if d.HasChange("volume_size") {
		_, err := client.PostgresInstance.ResizeVolume(id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of postgres database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
		_, err = waitUntilPostgresInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of postgres database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
	}

	if d.HasChange("flavor_id") {
		_, err := client.PostgresInstance.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of postgres database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
		_, err = waitUntilPostgresInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of postgres database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
	}
	return resourcePostgresInstanceRead(d, meta)
}

func resourcePostgresInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.PostgresInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete postgres database instance: %v", err)
	}
	_, err = waitUntilPostgresInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete postgres database instance: %v", err)
	}
	return nil
}

func resourcePostgresInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresInstanceRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilPostgresInstanceJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING"}, []string{"ERROR", "SHUTDOWN"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresInstance.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.PostgresInstance).Status)
	})
}

func waitUntilPostgresInstanceDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresInstance.Get(id)
	})
}
