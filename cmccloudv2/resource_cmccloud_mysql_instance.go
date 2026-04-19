package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMysqlInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMysqlInstanceCreate,
		Read:   resourceMysqlInstanceRead,
		Update: resourceMysqlInstanceUpdate,
		Delete: resourceMysqlInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMysqlInstanceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mysqlInstanceSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			mode := strings.ToLower(diff.Get("mode").(string))
			secondaryQuantity := diff.Get("quantity_of_slave").(int)
			// replica_set mode requires quantity_of_slave
			if mode == "replica_set" {
				if secondaryQuantity <= 0 {
					return fmt.Errorf("`quantity_of_slave` is required and must be > 0 when mode is `%s`", mode)
				}
			}
			return nil
		},
	}
}

func resourceMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	datastores, err := client.MysqlInstance.ListDatastore(map[string]string{})
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
		"quantityOfSlave":      d.Get("quantity_of_slave").(int),
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		"datastore_type": datastoreTypeId,
		// "datastore_version": datastoreVersionId,
	}
	requestMetadata := map[string]interface{}{}
	switch mode {
	case "replica_set":
		requestMetadata["quantityOfSlave"] = d.Get("quantity_of_slave").(int)
		requestMetadata["zones"] = zones
	default:
		// Standalone node
		requestMetadata["zone"] = zones[0]
	}

	metadataJsonData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal requestMetadata: %s", err)
	}
	params["requestMetadata"] = string(metadataJsonData)

	instance, err := client.MysqlInstance.Create(params)
	if err != nil {
		return fmt.Errorf("error creating MysqlDatabase Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)

	_, err = client.Tag.UpdateTag(instance.Data.InstanceID, "MYSQL", d)
	if err != nil {
		fmt.Printf("error updating Mysql Database tags: %s\n", err)
	}

	_, err = waitUntilMysqlInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating Mysql Database Instance: %s", err)
	}
	return resourceMysqlInstanceRead(d, meta)
}

func resourceMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.MysqlInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Mysql Database Instance %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("version", instance.DatastoreVersion)
	if strings.Contains(strings.ToLower(instance.DatastoreMode), "replica") {
		_ = d.Set("mode", "replica_set")
	} else if strings.Contains(strings.ToLower(instance.DatastoreMode), "standalone") {
		_ = d.Set("mode", "standalone")
	}
	// _ = d.Set("flavor_id", )
	_ = d.Set("volume_size", instance.VolumeSize)
	// _ = d.Set("subnet_id", instance.SubnetID)
	_ = d.Set("configuration_id", instance.GroupConfigID)
	_ = d.Set("tags", convertTagsToSet(instance.Tags))
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.CreatedAt)
	return nil
}

func resourceMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("tags") {
		_, err := client.Tag.UpdateTag(id, "MYSQL", d)
		if err != nil {
			return fmt.Errorf("error when set mysql tags [%s]: %v", id, err)
		}
	}
	if d.HasChange("configuration_id") {
		_, err := client.MysqlInstance.SetConfigurationGroupId(id, d.Get("configuration_id").(string))
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of mysql database instance %s: %v", d.Get("configuration_id").(string), id, err)
		}
		_, err = waitUntilMysqlInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of mysql database instance %s: %v", d.Get("configuration_id").(string), id, err)
		}
	}

	if d.HasChange("volume_size") {
		_, err := client.MysqlInstance.ResizeVolume(id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of mysql database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
		_, err = waitUntilMysqlInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of mysql database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
	}

	if d.HasChange("flavor_id") {
		_, err := client.MysqlInstance.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of mysql database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
		_, err = waitUntilMysqlInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of mysql database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
	}
	return resourceMysqlInstanceRead(d, meta)
}

func resourceMysqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.MysqlInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete mysql database instance: %v", err)
	}
	_, err = waitUntilMysqlInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete mysql database instance: %v", err)
	}
	return nil
}

func resourceMysqlInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMysqlInstanceRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilMysqlInstanceJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING"}, []string{"ERROR", "SHUTDOWN"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MysqlInstance.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.MysqlInstance).Status)
	})
}

func waitUntilMysqlInstanceDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MysqlInstance.Get(id)
	})
}
