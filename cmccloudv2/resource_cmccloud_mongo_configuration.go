package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoConfigurationCreate,
		Read:   resourceMongoConfigurationRead,
		Update: resourceMongoConfigurationUpdate,
		Delete: resourceMongoConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoConfigurationImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mongoConfigurationSchema(),
	}
}

func resourceMongoConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datastores, _ := client.MongoInstance.ListDatastore(map[string]string{})
	databaseVersion := d.Get("database_version").(string)
	databaseMode := d.Get("database_mode").(string)

	datastoreVersionId, datastoreModeId, _, _, _, err := findPostgresDatastoreInfo(datastores, databaseVersion, databaseMode)
	if err != nil {
		return err
	}

	// datastoreModeId := ""
	// datastoreVersionId := ""
	// for _, datastore := range datastores {
	// 	versions := make([]string, len(datastore.VersionInfos))
	// 	for index, version := range datastore.VersionInfos {
	// 		versions[index] = version.VersionName
	// 		if version.VersionName == databaseVersion {
	// 			datastoreVersionId = version.ID
	// 			modes := make([]string, len(version.ModeInfo))
	// 			for i, mode := range version.ModeInfo {
	// 				modes[i] = mode.Code
	// 				if strings.Contains(mode.Code, databaseMode) {
	// 					datastoreModeId = mode.ID
	// 				}
	// 			}
	// 			if datastoreModeId == "" {
	// 				return fmt.Errorf("not found database_mode `%s`, must be one of %v", databaseMode, modes)
	// 			}
	// 		}
	// 	}
	// 	if datastoreVersionId == "" {
	// 		return fmt.Errorf("not found database_version `%s`, must be one of %v", databaseVersion, versions)
	// 	}
	// }

	// if datastoreModeId == "" || datastoreVersionId == "" {
	// 	return fmt.Errorf("not found database_mode / database version")
	// }
	configuration, err := client.MongoConfiguration.Create(map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"datastoreModeId": datastoreModeId,
		"cacheEngine":     datastoreVersionId,
		"overridesConfig": d.Get("parameters").(map[string]interface{}),
	})
	if err != nil {
		return fmt.Errorf("error creating Mongo Configuration: %s", err)
	}
	d.SetId(configuration.ID)
	return resourceMongoConfigurationRead(d, meta)
}

func resourceMongoConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.MongoConfiguration.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Mongo Configuration %s: %v", d.Id(), err)
	}

	_ = d.Set("id", configuration.ID)
	_ = d.Set("name", configuration.Name)
	_ = d.Set("description", configuration.Description)
	_ = d.Set("database_version", configuration.DatastoreVersion)
	_ = d.Set("database_mode", configuration.DatastoreMode)
	_ = d.Set("parameters", convertMongoConfigurationParameters(configuration.Parameters))
	return nil
}

func convertMongoConfigurationParameters(obj []gocmcapiv2.MongoConfigurationParameter) map[string]interface{} {
	result := map[string]interface{}{}
	index := 0
	for _, param := range obj {
		result[param.Name] = param.Value
		index++
	}
	return result
}
func resourceMongoConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		_, err := client.MongoConfiguration.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		})
		if err != nil {
			return fmt.Errorf("error when update info of Mongo Configuration [%s]: %v", id, err)
		}
	}
	if d.HasChange("parameters") {
		_, err := client.MongoConfiguration.UpdateParameters(id, d.Get("parameters").(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("error when update parameters of Mongo Configuration [%s]: %v", id, err)
		}
	}
	return resourceMongoConfigurationRead(d, meta)
}

func resourceMongoConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.MongoConfiguration.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete mongo configuration: %v", err)
	}
	_, err = waitUntilMongoConfigurationDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete mongo configuration: %v", err)
	}
	return nil
}

func resourceMongoConfigurationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMongoConfigurationRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilMongoConfigurationDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).MongoConfiguration.Get(id)
	})
}
