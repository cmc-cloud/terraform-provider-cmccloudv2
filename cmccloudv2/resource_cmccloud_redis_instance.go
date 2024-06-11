package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
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
			sourceType := diff.Get("source_type").(string)
			_, sourceIDSet := diff.GetOk("backup_id")
			if sourceType == "new" {
				if sourceIDSet {
					return fmt.Errorf("When source_type is 'new', 'backup_id' must not be set")
				}
			} else if sourceType == "backup" {
				if !sourceIDSet {
					return fmt.Errorf("When source_type is 'backup', 'backup_id' must be set")
				}
			}
			return nil
		},
	}
}

func resourceRedisInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Subnet id is not valid %v", err)
	}

	datastores, _ := client.RedisInstance.ListDatastore(map[string]string{})
	database_engine := d.Get("database_engine").(string)
	database_version := d.Get("database_version").(string)
	database_mode := d.Get("database_mode").(string)

	redisMode := ""
	datastoreId := ""
	datastoreCode := ""
	datastoreVersionId := ""
	datastoreModeId := ""
	for _, datastore := range datastores {
		if strings.ToLower(database_engine) == strings.ToLower(datastore.Name) {
			datastoreCode = datastore.Code
			datastoreId = datastore.ID
			for _, version := range datastore.VersionInfos {
				if strings.ToLower(database_version) == strings.ToLower(version.VersionName) {
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
				return fmt.Errorf("Not found database_version %s", database_version)
			}

			if datastoreModeId == "" {
				return fmt.Errorf("Not found database_mode %s", database_mode)
			}
		}
	}

	if datastoreCode == "" {
		return fmt.Errorf("Not found database_engine %s", database_engine)
	}

	// if d.Get("redis_configuration_id").(string) != "" {
	// 	configuration, err := client.RedisConfiguration.Get(d.Get("redis_configuration_id").(string))
	// 	if err != nil {
	// 		return fmt.Errorf("Error getting RedisDatabase configuration %s: %s", d.Get("redis_configuration_id").(string), err)
	// 	}

	// 	if configuration.DatastoreMode != database_mode {
	// 		return fmt.Errorf("Datastore mode of configuration is `%s` != database_mode `%s`", configuration.DatastoreMode, database_mode)
	// 	}
	// }

	params := map[string]interface{}{
		"project":              client.Configs.ProjectId,
		"region":               client.Configs.RegionId,
		"name":                 d.Get("name").(string),
		"billing_mode":         d.Get("billing_mode").(string),
		"source_type":          d.Get("source_type").(string),
		"source_id":            d.Get("source_id").(string),
		"backupId":             d.Get("backup_id").(string),
		"slaveCount":           2,
		"volumeSize":           d.Get("volume_size").(int),
		"volumeType":           d.Get("volume_type").(string),
		"password":             d.Get("password").(string),
		"enableMonitor":        true,
		"enable_public_ip":     false,
		"is_public":            false,
		"replicate_count":      1,
		"vpcId":                subnet.NetworkID,
		"subnetId":             subnet.ID,
		"zoneMaster":           d.Get("zone_master").(string),
		"zoneSlaves":           d.Get("zone_slave").(string),
		"flavorId":             d.Get("flavor_id").(string),
		"groupConfigurationId": d.Get("redis_configuration_id").(string),
		"securityGroupIds":     strings.Join(getStringArrayFromTypeSet(d.Get("security_group_ids").(*schema.Set)), ","),
		"datastore": map[string]string{
			"datastoreCode":      datastoreCode,
			"datastoreVersionId": datastoreVersionId,
			"datastoreModeId":    datastoreModeId,
		},
		"datastore_type": datastoreId,
	}

	ip_master := d.Get("ip_master").(string)
	if ip_master != "" {
		params["master"] = map[string]interface{}{
			"ipAddressType": "manual",
			"ipAddress":     ip_master,
		}
	} else {
		params["master"] = map[string]interface{}{
			"ipAddressType": "auto",
			"ipAddress":     "",
		}
	}

	ip_slave1 := d.Get("ip_slave1").(string)
	ip_slave2 := d.Get("ip_slave2").(string)
	slaves := make([]map[string]interface{}, 2)

	if ip_slave1 != "" {
		slaves[0] = map[string]interface{}{
			"ipAddressType": "manual",
			"ipAddress":     ip_slave1,
		}
	} else {
		slaves[0] = map[string]interface{}{
			"ipAddressType": "auto",
			"ipAddress":     "",
		}
	}

	if ip_slave2 != "" {
		slaves[1] = map[string]interface{}{
			"ipAddressType": "manual",
			"ipAddress":     ip_slave2,
		}
	} else {
		slaves[1] = map[string]interface{}{
			"ipAddressType": "auto",
			"ipAddress":     "",
		}
	}

	requestMetadataRaw := map[string]interface{}{
		"zoneMaster": params["zoneMaster"],
		"password":   d.Get("password").(string),
	}
	if redisMode == "master_slave" {
		requestMetadataRaw["master"] = params["master"]
		requestMetadataRaw["slaves"] = slaves
		requestMetadataRaw["zoneSlaves"] = params["zoneSlaves"]
	} else if redisMode == "standalone" {
		requestMetadataRaw["ipAddressType"] = params["master"].(map[string]interface{})["ipAddressType"]
		requestMetadataRaw["ipAddress"] = params["master"].(map[string]interface{})["ipAddress"]
	} else if redisMode == "cluster" {
		requestMetadataRaw["numOfMasterServer"] = d.Get("master_count").(int)
		requestMetadataRaw["zoneSlaves"] = params["zoneSlaves"]
		requestMetadataRaw["replicas"] = params["replicas"]
	}

	jsonData, err := json.Marshal(requestMetadataRaw)
	if err != nil {
		return fmt.Errorf("Error creating RedisDatabase Instance: %s", err)
	}
	params["requestMetadataRaw"] = requestMetadataRaw
	params["requestMetadata"] = string(jsonData)

	delete(params, "master")
	delete(params, "slaves")
	// gocmcapiv2.Logs("redisMode = " + redisMode)
	// gocmcapiv2.Logo("params = ", params)
	// return fmt.Errorf("Error creating RedisDatabase Instance: %s", "test")

	instance, err := client.RedisInstance.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating RedisDatabase Instance: %s", err)
	}
	d.SetId(instance.Data.InstanceID)
	_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating RedisDatabase Instance: %s", err)
	}
	return resourceRedisInstanceRead(d, meta)
}

func resourceRedisInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.RedisInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving RedisDatabase Instance %s: %v", d.Id(), err)
	}

	billing_mode, err := client.BillingMode.GetBilingMode(d.Id(), "RDS")
	if billing_mode != "" {
		_ = d.Set("billing_mode", billing_mode)
	}
	_ = d.Set("id", instance.ID)
	_ = d.Set("name", instance.Name)
	// _ = d.Set("source_type", )
	// _ = d.Set("source_id", )
	// _ = d.Set("backup_id", )
	_ = d.Set("database_engine", instance.DatastoreName)
	_ = d.Set("database_version", instance.DatastoreVersion)
	_ = d.Set("database_mode", instance.DatastoreMode)
	_ = d.Set("zone_master", instance.DataDetail.MasterInfo.ZoneName)
	if len(instance.DataDetail.SlavesInfo) > 0 {
		setString(d, "zone_slave", instance.DataDetail.SlavesInfo[0].ZoneName)
	}

	var security_group_ids []string
	err = json.Unmarshal([]byte(instance.SecurityClientIds), &security_group_ids)
	if err != nil {
		fmt.Errorf("Error when get info of Redis Database Instance [%s]: %v", d.Id(), err)
	}

	_ = d.Set("security_group_ids", security_group_ids)
	// _ = d.Set("flavor_id",      instance.)
	// _ = d.Set("volume_type",           instance.)

	_ = d.Set("volume_size", instance.DataDetail.MasterInfo.VolumeSize)
	_ = d.Set("subnet_id", instance.SubnetID)
	// _ = d.Set("ip_master", )
	// _ = d.Set("ip_slave1", )
	// _ = d.Set("ip_slave2", )
	_ = d.Set("redis_configuration_id", instance.GroupConfigID)
	// _ = d.Set("password", )
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
			return fmt.Errorf("Error when update info of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetRedisInstanceBilingMode(d.Id(), d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when billing mode of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("password") {
		_, err := client.RedisInstance.SetPassword(id, d.Get("password").(string))
		if err != nil {
			return fmt.Errorf("Error when update password of Redis Database Instance [%s]: %v", id, err)
		}
		_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when update password of Redis Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("redis_configuration_id") {
		_, err := client.RedisInstance.SetConfigurationGroupId(id, d.Get("redis_configuration_id").(string))
		if err != nil {
			return fmt.Errorf("Error when set configuration group to %s of redis database instance %s: %v", d.Get("redis_configuration_id").(string), id, err)
		}
		_, err = waitUntilRedisInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when set configuration group to %s of redis database instance %s: %v", d.Get("redis_configuration_id").(string), id, err)
		}
	}
	if d.HasChange("security_group_ids") {
		removes, adds := getDiffSet(d.GetChange("security_group_ids"))

		for _, security_group_id := range removes.List() {
			_, err := client.RedisInstance.DetachSecurityGroupId(d.Id(), security_group_id.(string))
			if err != nil {
				return fmt.Errorf("Error detach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
			waitUntilRedisInstanceDetachFinished(d, meta, security_group_id.(string))
		}
		for _, security_group_id := range adds.List() {
			_, err := client.RedisInstance.AttachSecurityGroupId(d.Id(), security_group_id.(string))
			if err != nil {
				return fmt.Errorf("Error attach security group %s from %s: %v", security_group_id, d.Id(), err)
			}
			waitUntilRedisInstanceAttachFinished(d, meta, security_group_id.(string))
		}
	}

	return resourceRedisInstanceRead(d, meta)
}

func resourceRedisInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.RedisInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete redis database instance: %v", err)
	}
	_, err = waitUntilRedisInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete redis database instance: %v", err)
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
		return obj.(gocmcapiv2.RedisInstance).Status
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
