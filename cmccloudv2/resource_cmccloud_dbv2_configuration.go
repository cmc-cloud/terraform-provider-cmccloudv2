package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDBv2Configuration(dbType string) *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2ConfigurationCreate(d, meta, dbType)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2ConfigurationRead(d, meta, dbType)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2ConfigurationDelete(d, meta, dbType)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			return resourceDBv2ConfigurationUpdate(d, meta, dbType)
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return resourceDBv2ConfigurationImport(d, meta, dbType)
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema: func() map[string]*schema.Schema {
			switch dbType {
			case "Redis":
				return redisconfigurationSchema()
			case "Postgres":
				return postgresConfigurationSchema()
			case "Mysql":
				return mysqlConfigurationSchema()
			case "Mongo":
				return mongoConfigurationSchema()
			default:
				return dbv2ConfigurationSchema(dbType)
			}
		}(),
	}
}

func resourceDBv2ConfigurationCreate(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datastores, _ := client.DBaaSConfiguration.ListDatastore(dbType, map[string]string{})
	databaseVersion := d.Get("database_version").(string)
	databaseMode := d.Get("database_mode").(string)

	datastoreVersionId, datastoreModeId, _, _, _, err := findDatastoreInfo(datastores, databaseVersion, databaseMode)
	if err != nil {
		return err
	}

	configuration, err := client.DBaaSConfiguration.Create(map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"datastoreModeId": datastoreModeId,
		"cacheEngine":     datastoreVersionId,
		"overridesConfig": d.Get("parameters").(map[string]interface{}),
	})
	if err != nil {
		return fmt.Errorf("error creating %s configuration: %s", dbType, err)
	}
	d.SetId(configuration.ID)
	return resourceDBv2ConfigurationRead(d, meta, dbType)
}

func resourceDBv2ConfigurationRead(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.DBaaSConfiguration.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving %s Configuration %s: %v", dbType, d.Id(), err)
	}

	_ = d.Set("id", configuration.ID)
	_ = d.Set("name", configuration.Name)
	_ = d.Set("description", configuration.Description)
	_ = d.Set("database_version", configuration.DatastoreVersion)
	_ = d.Set("database_mode", configuration.DatastoreMode)
	_ = d.Set("parameters", convertDBv2ConfigurationParameters(configuration.Parameters))
	return nil
}

func resourceDBv2ConfigurationUpdate(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		_, err := client.DBaaSConfiguration.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		})
		if err != nil {
			return fmt.Errorf("error when update info of %s configuration [%s]: %v", dbType, id, err)
		}
	}
	if d.HasChange("parameters") {
		_, err := client.DBaaSConfiguration.UpdateParameters(id, d.Get("parameters").(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("error when update parameters of %s Configuration [%s]: %v", dbType, id, err)
		}
	}
	return resourceDBv2ConfigurationRead(d, meta, dbType)
}

func resourceDBv2ConfigurationDelete(d *schema.ResourceData, meta interface{}, dbType string) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DBaaSConfiguration.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete %s configuration: %v", dbType, err)
	}
	_, err = waitUntilDBaaSConfigurationDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete %s configuration: %v", dbType, err)
	}
	return nil
}

func resourceDBv2ConfigurationImport(d *schema.ResourceData, meta interface{}, dbType string) ([]*schema.ResourceData, error) {
	err := resourceDBv2ConfigurationRead(d, meta, dbType)
	return []*schema.ResourceData{d}, err
}

func convertDBv2ConfigurationParameters(obj []gocmcapiv2.DBv2ConfigurationParameter) map[string]interface{} {
	result := map[string]interface{}{}
	index := 0
	for _, param := range obj {
		result[param.Name] = param.Value
		index++
	}
	return result
}
func waitUntilDBaaSConfigurationDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBaaSConfiguration.Get(id)
	})
}
