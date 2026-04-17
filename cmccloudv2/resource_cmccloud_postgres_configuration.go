package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresConfigurationCreate,
		Read:   resourcePostgresConfigurationRead,
		Update: resourcePostgresConfigurationUpdate,
		Delete: resourcePostgresConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresConfigurationImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        postgresConfigurationSchema(),
	}
}

func resourcePostgresConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datastores, _ := client.PostgresInstance.ListDatastore(map[string]string{})
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
	configuration, err := client.PostgresConfiguration.Create(map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"datastoreModeId": datastoreModeId,
		"cacheEngine":     datastoreVersionId,
		"overridesConfig": d.Get("parameters").(map[string]interface{}),
	})
	if err != nil {
		return fmt.Errorf("error creating Postgres Configuration: %s", err)
	}
	d.SetId(configuration.ID)
	return resourcePostgresConfigurationRead(d, meta)
}

func resourcePostgresConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.PostgresConfiguration.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Postgres Configuration %s: %v", d.Id(), err)
	}

	_ = d.Set("id", configuration.ID)
	_ = d.Set("name", configuration.Name)
	_ = d.Set("description", configuration.Description)
	_ = d.Set("database_version", configuration.DatastoreVersion)
	_ = d.Set("database_mode", configuration.DatastoreMode)
	_ = d.Set("parameters", convertPostgresConfigurationParameters(configuration.Parameters))
	return nil
}

func convertPostgresConfigurationParameters(obj []gocmcapiv2.PostgresConfigurationParameter) map[string]interface{} {
	result := map[string]interface{}{}
	index := 0
	for _, param := range obj {
		result[param.Name] = param.Value
		index++
	}
	return result
}
func resourcePostgresConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		_, err := client.PostgresConfiguration.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		})
		if err != nil {
			return fmt.Errorf("error when update info of Postgres Configuration [%s]: %v", id, err)
		}
	}
	if d.HasChange("parameters") {
		_, err := client.PostgresConfiguration.UpdateParameters(id, d.Get("parameters").(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("error when update parameters of Postgres Configuration [%s]: %v", id, err)
		}
	}
	return resourcePostgresConfigurationRead(d, meta)
}

func resourcePostgresConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.PostgresConfiguration.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete postgres configuration: %v", err)
	}
	_, err = waitUntilPostgresConfigurationDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete postgres configuration: %v", err)
	}
	return nil
}

func resourcePostgresConfigurationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresConfigurationRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilPostgresConfigurationDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).PostgresConfiguration.Get(id)
	})
}
