package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKafkaInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKafkaInstanceCreate,
		Read:   resourceKafkaInstanceRead,
		Update: resourceKafkaInstanceUpdate,
		Delete: resourceKafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKafkaInstanceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kafkaInstanceSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			mode := strings.ToLower(diff.Get("mode").(string))
			brokerQuantity := diff.Get("broker_quantity").(int)

			// master_slave or ha_cluster requires proxy_quantity
			if mode == "cluster" {
				if brokerQuantity <= 0 {
					return fmt.Errorf("`broker_quantity` is required and must be > 0 when mode is `%s`", mode)
				}
			}

			// ha_cluster requires both proxy_quantity and proxy_flavor_id
			if mode == "single_node" {
				if brokerQuantity > 0 {
					return fmt.Errorf("`broker_quantity` must not be set when mode is `single_node`")
				}
			}

			return nil
		},
	}
}

func resourceKafkaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	datastores, err := client.KafkaInstance.ListDatastore()
	if err != nil {
		return fmt.Errorf("can't get list of datastore %v", err)
	}

	version := d.Get("version").(string)
	mode := d.Get("mode").(string)

	datastoreVersionId, datastoreModeId, datastoreCode, _, _, err := findPostgresDatastoreInfo(datastores, version, mode)
	if err != nil {
		return err
	}

	// Lấy list authenInformation từ thuộc tính users
	var authenInformation []map[string]interface{}
	if users, ok := d.GetOk("users"); ok {
		usersList := users.([]interface{})
		for _, user := range usersList {
			if userMap, ok := user.(map[string]interface{}); ok {
				auth := map[string]interface{}{
					"username": userMap["username"],
					"password": userMap["password"],
				}
				authenInformation = append(authenInformation, auth)
			}
		}
	}
	if authenInformation == nil {
		authenInformation = []map[string]interface{}{}
	}
	requestMetadata := map[string]interface{}{
		"enableBasicAuth":   d.Get("enable_basic_authen").(bool),
		"authenInformation": authenInformation,
	}
	zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))
	switch mode {
	case "cluster":
		requestMetadata["quantityOfBroker"] = d.Get("broker_quantity").(int)
		requestMetadata["zones"] = zones
	default:
		// Standalone node
		requestMetadata["quantityOfBroker"] = 1
		requestMetadata["zone"] = zones[0]
	}

	metadataJsonData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("failed to marshal requestMetadata: %s", err)
	}

	params := map[string]interface{}{
		// "project":          client.Configs.ProjectId,
		// "region":           client.Configs.RegionId,
		"name":        d.Get("name").(string),
		"billingMode": d.Get("billing_mode").(string),
		// "zone":             zones[0],
		// "zones":            zones,
		"flavorId":   d.Get("flavor_id").(string),
		"volumeSize": d.Get("volume_size").(int),
		// "volumeType":       d.Get("volume_type").(string),
		"networkId": subnet.NetworkID,
		"subnetId":  subnet.ID,
		// "quantityOfBroker": d.Get("broker_quantity").(int),
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		// "datastore_type": datastoreTypeId,
	}
	params["requestMetadata"] = string(metadataJsonData)

	instance, err := client.KafkaInstance.Create(params)
	if err != nil {
		return fmt.Errorf("error creating Kafka Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)
	_, err = waitUntilKafkaInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating Kafka Instance: %s", err)
	}

	securityGroupIds := d.Get("security_group_ids").(*schema.Set)
	for _, security_group_id := range securityGroupIds.List() {
		_, err = waitUntilKafkaInstanceAttachFinished(d, meta, security_group_id.(string))
		if err != nil {
			return fmt.Errorf("error attach security group %s to %s: %v", security_group_id, d.Id(), err)
		}
	}

	return resourceKafkaInstanceRead(d, meta)
}

func resourceKafkaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.KafkaInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Kafka Instance %s: %v", d.Id(), err)
	}
	_ = d.Set("id", instance.ID)
	_ = d.Set("name", instance.Name)
	_ = d.Set("version", instance.DatastoreVersion)
	if strings.Contains(strings.ToLower(instance.DatastoreMode), "cluster") {
		_ = d.Set("mode", "cluster")
	} else {
		_ = d.Set("mode", "single_node")
	}

	var securityGroupIds []string
	err = json.Unmarshal([]byte(instance.SecurityClientIds), &securityGroupIds)
	if err != nil {
		return fmt.Errorf("error when get info of Kafka Instance [%s]: %v", d.Id(), err)
	}
	_ = d.Set("security_group_ids", securityGroupIds)
	_ = d.Set("flavor_id", instance.FlavorInfo.ID)
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("subnet_id", instance.SubnetID)
	_ = d.Set("broker_quantity", instance.QuantityOfNodes)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	return nil
}

func resourceKafkaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("security_group_ids") {
		err := checkSecurityGroupConflict(d, meta)
		if err != nil {
			return err
		}

		removes, adds := getDiffSet(d.GetChange("security_group_ids"))

		for _, security_group_id := range removes.List() {
			_, err := client.KafkaInstance.DetachSecurityGroupId(d.Id(), security_group_id.(string))
			if err != nil {
				return fmt.Errorf("error detach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
			_, err = waitUntilKafkaInstanceDetachFinished(d, meta, security_group_id.(string))
			if err != nil {
				return fmt.Errorf("error detach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
		}
		for _, security_group_id := range adds.List() {
			_, err := client.KafkaInstance.AttachSecurityGroupId(d.Id(), security_group_id.(string))
			if err != nil {
				return fmt.Errorf("error attach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
			_, err = waitUntilKafkaInstanceAttachFinished(d, meta, security_group_id.(string))
			if err != nil {
				return fmt.Errorf("error attach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
		}
	}

	if d.HasChange("volume_size") {
		_, err := client.KafkaInstance.ResizeVolume(id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of kafka database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
		_, err = waitUntilKafkaInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize volume to %s of kafka database instance %s: %v", d.Get("volume_size").(string), id, err)
		}
	}

	if d.HasChange("flavor_id") {
		_, err := client.KafkaInstance.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of kafka database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
		_, err = waitUntilKafkaInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when resize flavor to %s of kafka database instance %s: %v", d.Get("flavor_id").(string), id, err)
		}
	}

	return resourceKafkaInstanceRead(d, meta)
}

func checkSecurityGroupConflict(d *schema.ResourceData, meta interface{}) error {
	securityGroupIds := getStringArrayFromTypeSet(d.Get("security_group_ids").(*schema.Set))
	if len(securityGroupIds) > 1 {
		firstValue := ""
		for _, security_group_id := range securityGroupIds {
			group, err := meta.(*CombinedConfig).goCMCClient().SecurityGroup.Get(security_group_id)
			if err != nil {
				return err
			}
			if firstValue == "" {
				firstValue = fmt.Sprintf("%t", group.Stateful)
			}
			if firstValue != fmt.Sprintf("%t", group.Stateful) {
				return fmt.Errorf("invalid security_group_ids, all security groups must have the same stateful")
			}
		}
	}
	return nil
}
func resourceKafkaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.KafkaInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete kafka database instance: %v", err)
	}
	_, err = waitUntilKafkaInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete kafka database instance: %v", err)
	}
	return nil
}

func resourceKafkaInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKafkaInstanceRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKafkaInstanceJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING", "ACTIVE", "READY", "STABLE"}, []string{"ERROR", "SHUTDOWN", "FAILURE"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.KafkaInstance).Status)
	})
}

func waitUntilKafkaInstanceDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.Get(id)
	})
}

func waitUntilKafkaInstanceAttachFinished(d *schema.ResourceData, meta interface{}, securityGroupId string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.Get(id)
	}, func(obj interface{}) string {
		instance := obj.(gocmcapiv2.KafkaInstance)
		if strings.Contains(instance.SecurityClientIds, securityGroupId) {
			return "true"
		}
		return ""
	})
}
func waitUntilKafkaInstanceDetachFinished(d *schema.ResourceData, meta interface{}, securityGroupId string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.Get(id)
	}, func(obj interface{}) string {
		instance := obj.(gocmcapiv2.KafkaInstance)
		if !strings.Contains(instance.SecurityClientIds, securityGroupId) {
			return "true"
		}
		return ""
	})
}
