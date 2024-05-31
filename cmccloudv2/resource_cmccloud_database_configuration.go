package cmccloudv2

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabaseConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseConfigurationCreate,
		Read:   resourceDatabaseConfigurationRead,
		Update: resourceDatabaseConfigurationUpdate,
		Delete: resourceDatabaseConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDatabaseConfigurationImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        databaseConfigurationSchema(),
	}
}

func resourceDatabaseConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.DatabaseConfiguration.Create(map[string]interface{}{
		"name":              d.Get("name").(string),
		"description":       d.Get("description").(string),
		"datastore_type":    d.Get("datastore_type").(string),
		"datastore_version": d.Get("datastore_version").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating Database Configuration: %s", err)
	}
	d.SetId(configuration.ID)
	parameters := convertParametersJsonString(d.Get("parameters").(*schema.Set))
	_, err = client.DatabaseConfiguration.UpdateParameters(d.Id(), map[string]interface{}{"parameters": parameters})
	if err != nil {
		return fmt.Errorf("Error when update parameters of Database Configuration [%s]: %v", configuration.ID, err)
	}

	return resourceDatabaseConfigurationRead(d, meta)
}

func convertParameters(obj gocmcapiv2.ArrayOrMap) []map[string]interface{} {
	if obj.IsObject {
		result := make([]map[string]interface{}, len(obj.Object))
		index := 0
		for key, value := range obj.Object {
			result[index] = map[string]interface{}{
				"key":   key,
				"value": value,
			}
			index++
		}
		return result
	}
	// da mang rong =>
	return []map[string]interface{}{}
}

// func setToMap(s *schema.Set) map[string]interface{} {
// 	result := make(map[string]interface{})
// 	for _, item := range s.List() {
// 		itemMap := item.(map[string]interface{})
// 		key := itemMap["key"].(string)
// 		value := itemMap["value"].(string)
// 		result[key] = value
// 	}
// 	return result
// }

func resourceDatabaseConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.DatabaseConfiguration.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Database Configuration %s: %v", d.Id(), err)
	}

	_ = d.Set("id", configuration.ID)
	_ = d.Set("name", configuration.Name)
	_ = d.Set("description", configuration.Description)
	_ = d.Set("datastore_type", configuration.DatastoreName)
	_ = d.Set("datastore_version", configuration.DatastoreVersionName)
	_ = d.Set("parameters", convertParameters(configuration.Parameters))
	return nil
}

func resourceDatabaseConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		_, err := client.DatabaseConfiguration.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		})
		if err != nil {
			return fmt.Errorf("Error when update info of Database Configuration [%s]: %v", id, err)
		}
	}
	if d.HasChange("parameters") {
		parameters := convertParametersJsonString(d.Get("parameters").(*schema.Set))
		_, err := client.DatabaseConfiguration.UpdateParameters(id, map[string]interface{}{"parameters": parameters})
		if err != nil {
			return fmt.Errorf("Error when update parameters of Database Configuration [%s]: %v", id, err)
		}
	}
	return resourceDatabaseConfigurationRead(d, meta)
}

func convertParametersJsonString(params *schema.Set) string {
	results := make(map[string]interface{})
	for _, param := range params.List() {
		_param := param.(map[string]interface{})
		key := _param["key"].(string)
		val := _param["value"]

		if isStringType, ok := _param["string_type"].(bool); !ok || !isStringType {
			// check if value can be converted into int
			if valueInt, err := strconv.Atoi(val.(string)); err == nil {
				val = valueInt
				// check if value can be converted into bool
			} else if valueBool, err := strconv.ParseBool(val.(string)); err == nil {
				val = valueBool
			}
		}
		results[key] = val
	}
	// return results
	jsonData, err := json.Marshal(results)
	if err != nil {
		fmt.Errorf("Error converting map to JSON: %s", err)
	}
	return string(jsonData)
}
func resourceDatabaseConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DatabaseConfiguration.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete database configuration: %v", err)
	}
	_, err = waitUntilDatabaseConfigurationDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete database configuration: %v", err)
	}
	return nil
}

func resourceDatabaseConfigurationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDatabaseConfigurationRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDatabaseConfigurationDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DatabaseConfiguration.Get(id)
	})
}
