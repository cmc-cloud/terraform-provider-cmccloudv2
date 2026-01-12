package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRedisInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisInstanceCreate,
		Read:   resourceRedisInstanceRead,
		Update: resourceRedisInstanceUpdate,
		Delete: resourceRedisInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRedisInstanceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        redisinstanceSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// sourceType := diff.Get("source_type").(string)
			// mode := diff.Get("mode").(string)

			// _, sourceIDSet := diff.GetOk("backup_id")
			// _, replicasSet := diff.GetOk("replicas")
			// if sourceType == "new" {
			// 	if sourceIDSet {
			// 		return fmt.Errorf("when source_type is 'new', 'backup_id' must not be set")
			// 	}
			// } else if sourceType == "backup" {
			// 	if !sourceIDSet {
			// 		return fmt.Errorf("when source_type is 'backup', 'backup_id' must be set")
			// 	}
			// }

			return nil
		},
	}
}

func resourceRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	datastores, err := client.RedisInstance.ListDatastore(map[string]string{})
	if err != nil {
		return fmt.Errorf("can't get list of datastore %v", err)
	}

	database_engine := d.Get("database_engine").(string)
	database_version := d.Get("database_version").(string)
	database_mode := d.Get("database_mode").(string)

	redisMode := ""
	datastoreId := ""
	datastoreCode := ""
	datastoreVersionId := ""
	datastoreModeId := ""
	for _, datastore := range datastores {
		// gocmcapiv2.Logo("datastore", datastore)
		if strings.EqualFold(database_engine, datastore.Name) {
			// gocmcapiv2.Logs("found datastore " + database_engine + " & " + datastore.Name)
			datastoreCode = datastore.Code
			datastoreId = datastore.ID
			for _, version := range datastore.VersionInfos {
				if strings.EqualFold(database_version, version.VersionName) {
					datastoreVersionId = version.ID
					for _, mode := range version.ModeInfo {
						if caseInsensitiveContains(mode.Name, database_mode) {
							datastoreModeId = mode.ID
							redisMode = mode.Code
						}
					}
				}
			}
			if datastoreVersionId == "" {
				return fmt.Errorf("not found database_version %s", database_version)
			}

			if datastoreModeId == "" {
				return fmt.Errorf("not found database_mode %s", database_mode)
			}
		}
	}

	if datastoreCode == "" {
		return fmt.Errorf("not found database_engine %s", database_engine)
	}

	_, replicasSet := d.GetOk("replicas")
	if redisMode == "cluster" {
		if !replicasSet {
			return fmt.Errorf("when `mode` is 'cluster', 'replicas' must be set")
		}
	} else {
		if replicasSet {
			return fmt.Errorf("when `mode` is not 'cluster', 'replicas' must not be set")
		}
	}

	err = checkSecurityGroupConflict(d, meta)
	if err != nil {
		return err
	}
	// if d.Get("redis_configuration_id").(string) != "" {
	// 	configuration, err := client.RedisConfiguration.Get(d.Get("redis_configuration_id").(string))
	// 	if err != nil {
	// 		return fmt.Errorf("error getting RedisDatabase configuration %s: %s", d.Get("redis_configuration_id").(string), err)
	// 	}

	// 	if configuration.DatastoreMode != database_mode {
	// 		return fmt.Errorf("Datastore mode of configuration is `%s` != database_mode `%s`", configuration.DatastoreMode, database_mode)
	// 	}
	// }

	params := map[string]interface{}{
		// "project":              client.Configs.ProjectId,
		// "region":               client.Configs.RegionId,
		"billingMode": d.Get("billing_mode").(string),

		"name": d.Get("name").(string),
		// "securityGroupIds": strings.Join(getStringArrayFromTypeSet(d.Get("security_group_ids").(*schema.Set)), ","),
		"flavorId":   d.Get("flavor_id").(string),
		"password":   d.Get("password").(string),
		"backupId":   d.Get("backup_id").(string),
		"volumeSize": d.Get("volume_size").(int),
		"volumeType": d.Get("volume_type").(string),
		// "volumeType":           d.Get("volume_type").(string),
		"groupConfigurationId": d.Get("redis_configuration_id").(string),
		"networkId":            subnet.NetworkID,
		"subnetId":             subnet.ID,
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		"datastore_type": datastoreId,
	}

	requestMetadata := map[string]interface{}{
		"password": d.Get("password").(string),
	}
	zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))

	switch redisMode {
	case "standalone":
		zone := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))[0]
		requestMetadata["zone"] = zone

	case "master_slave":
		requestMetadata["zones"] = zones
		requestMetadata["numOfSlaves"] = 2
		params["zones"] = zones

	case "cluster":
		zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))
		requestMetadata["zones"] = zones
		requestMetadata["numOfMasterServer"] = 3
		requestMetadata["replicas"] = d.Get("replicas").(int)

	default:
		// nếu redisMode không khớp case nào
		// có thể log cảnh báo hoặc bỏ qua
	}

	// params["zone"] = requestMetadata["zone"]

	jsonData, err := json.Marshal(requestMetadata)
	if err != nil {
		return fmt.Errorf("error creating RedisDatabase Instance: %s", err)
	}
	params["requestMetadata"] = string(jsonData)

	instance, err := client.RedisInstance.Create(params)
	if err != nil {
		return fmt.Errorf("error creating RedisDatabase Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)
	_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating RedisDatabase Instance: %s", err)
	}
	return resourceRedisInstanceRead(d, meta)
}

func checkSecurityGroupConflict(d *schema.ResourceData, meta interface{}) error {
	// security_group_ids := getStringArrayFromTypeSet(d.Get("security_group_ids").(*schema.Set))
	// if len(security_group_ids) > 1 {
	// 	firstValue := ""
	// 	for _, security_group_id := range security_group_ids {
	// 		group, err := meta.(*CombinedConfig).goCMCClient().SecurityGroup.Get(security_group_id)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if firstValue == "" {
	// 			firstValue = fmt.Sprintf("%t", group.Stateful)
	// 		}
	// 		if firstValue != fmt.Sprintf("%t", group.Stateful) {
	// 			return fmt.Errorf("invalid security_group_ids, all security groups must have the same stateful")
	// 		}
	// 	}
	// }
	return nil
}

func resourceRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.RedisInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving RedisDatabase Instance %s: %v", d.Id(), err)
	}

	billing_mode, _ := client.BillingMode.GetBilingMode(d.Id(), "RDS")
	if billing_mode != "" {
		_ = d.Set("billing_mode", billing_mode)
	}
	_ = d.Set("id", instance.ID)
	_ = d.Set("name", instance.Name)
	// _ = d.Set("backup_id", )
	_ = d.Set("database_engine", instance.DatastoreName)
	_ = d.Set("database_version", instance.DatastoreVersion)
	_ = d.Set("database_mode", instance.DatastoreMode)

	// var security_group_ids []string
	// err = json.Unmarshal([]byte(instance.SecurityClientIds), &security_group_ids)
	// if err != nil {
	// 	return fmt.Errorf("error when get info of Redis Database Instance [%s]: %v", d.Id(), err)
	// }
	// _ = d.Set("security_group_ids", security_group_ids)
	_ = d.Set("flavor_id", instance.FlavorID)
	// _ = d.Set("volume_type",           instance.)

	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("subnet_id", instance.SubnetID)
	if d.Get("redis_configuration_id").(string) == "" {
		// _, err := client.RedisConfiguration.Get(instance.GroupConfigID)
		// if err == nil {
		// 	// la default configuration => ko set
		// } else {
		// 	_ = d.Set("redis_configuration_id", instance.GroupConfigID)
		// }
	} else {
		_ = d.Set("redis_configuration_id", instance.GroupConfigID)
	}
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	return nil
}

func resourceRedisInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		_, err := client.RedisInstance.Update(id, map[string]interface{}{"name": d.Get("name").(string)})
		if err != nil {
			return fmt.Errorf("error when update info of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetRedisInstanceBilingMode(d.Id(), d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("error when billing mode of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("password") {
		_, err := client.RedisInstance.SetPassword(id, d.Get("password").(string))
		if err != nil {
			return fmt.Errorf("error when update password of Redis Database Instance [%s]: %v", id, err)
		}
		_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when update password of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("redis_configuration_id") {
		defaultTemplateId := d.Get("redis_configuration_id").(string)
		if d.Get("redis_configuration_id") == "" {
			defaultTemplates, err := client.RedisConfiguration.List(map[string]string{
				"page":          "1",
				"size":          "1000",
				"datastoreCode": "redis",
				"getDefault":    "true",
			})
			if err != nil {
				return fmt.Errorf("error when getting default redis configuration templates: %v", err)
			}
			database_engine := d.Get("database_engine").(string)
			database_version := d.Get("database_version").(string)
			database_mode := d.Get("database_mode").(string)

			for _, template := range defaultTemplates {
				if template.DatastoreName == database_engine && template.DatastoreVersion == database_version && template.DatastoreMode == database_mode {
					defaultTemplateId = template.ID
					if template.ID2 != "" {
						defaultTemplateId = template.ID2
					}
				}
			}
		}
		_, err := client.RedisInstance.SetConfigurationGroupId(id, defaultTemplateId)
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of redis database instance %s: %v", d.Get("redis_configuration_id").(string), id, err)
		}
		_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when set configuration group to %s of redis database instance %s: %v", d.Get("redis_configuration_id").(string), id, err)
		}
	}
	// if d.HasChange("security_group_ids") {
	// 	err := checkSecurityGroupConflict(d, meta)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	removes, adds := getDiffSet(d.GetChange("security_group_ids"))

	// 	for _, security_group_id := range removes.List() {
	// 		_, err := client.RedisInstance.DetachSecurityGroupId(d.Id(), security_group_id.(string))
	// 		if err != nil {
	// 			return fmt.Errorf("error detach security group %s from %s: %v", security_group_id, d.Id(), err)
	// 		}
	// 		_, err = waitUntilRedisInstanceDetachFinished(d, meta, security_group_id.(string))
	// 		if err != nil {
	// 			return fmt.Errorf("error detach security group %s from %s: %v", security_group_id, d.Id(), err)
	// 		}
	// 	}
	// 	for _, security_group_id := range adds.List() {
	// 		_, err := client.RedisInstance.AttachSecurityGroupId(d.Id(), security_group_id.(string))
	// 		if err != nil {
	// 			return fmt.Errorf("error attach security group %s from %s: %v", security_group_id, d.Id(), err)
	// 		}
	// 		_, err = waitUntilRedisInstanceAttachFinished(d, meta, security_group_id.(string))
	// 		if err != nil {
	// 			return fmt.Errorf("error attach security group %s from %s: %v", security_group_id, d.Id(), err)
	// 		}
	// 	}
	// }

	return resourceRedisInstanceRead(d, meta)
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.RedisInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete redis database instance: %v", err)
	}
	_, err = waitUntilRedisInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete redis database instance: %v", err)
	}
	return nil
}

func resourceRedisInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRedisInstanceRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilRedisInstanceJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING", "SHUTDOWN"}, []string{"ERROR"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RedisInstance.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.RedisInstance).Status)
	})
}

func waitUntilRedisInstanceDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RedisInstance.Get(id)
	})
}

func waitUntilRedisInstanceAttachFinished(d *schema.ResourceData, meta interface{}, security_group_id string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RedisInstance.Get(id)
	}, func(obj interface{}) string {
		instance := obj.(gocmcapiv2.RedisInstance)
		if strings.Contains(instance.SecurityClientIds, security_group_id) {
			return "true"
		}
		return ""
	})
}

func waitUntilRedisInstanceDetachFinished(d *schema.ResourceData, meta interface{}, security_group_id string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RedisInstance.Get(id)
	}, func(obj interface{}) string {
		instance := obj.(gocmcapiv2.RedisInstance)
		if !strings.Contains(instance.SecurityClientIds, security_group_id) {
			return "true"
		}
		return ""
	})
}
