package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceDBaaSConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the postgres configuration",
		},
		"is_default": {
			Type:        schema.TypeBool,
			Description: "If true, it is default template of system. If false, it is custom template",
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of postgres configuration (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"database_mode": {
			Type:        schema.TypeString,
			Description: "filter by database_mode",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceMongoConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceMongoConfigurationRead,
		Schema: datasourceDBaaSConfigurationSchema(),
	}
}
func datasourceRedisConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceRedisConfigurationRead,
		Schema: datasourceDBaaSConfigurationSchema(),
	}
}
func datasourceMysqlConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceMysqlConfigurationRead,
		Schema: datasourceDBaaSConfigurationSchema(),
	}
}
func datasourcePostgresConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePostgresConfigurationRead,
		Schema: datasourceDBaaSConfigurationSchema(),
	}
}

func dataSourceMongoConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceDBaaSConfigurationRead(d, meta, "mongodb")
}
func datasourceRedisConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceDBaaSConfigurationRead(d, meta, "redis")
}
func dataSourceMysqlConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceDBaaSConfigurationRead(d, meta, "mysql")
}
func dataSourcePostgresConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceDBaaSConfigurationRead(d, meta, "postgresql")
}

func dataSourceDBaaSConfigurationRead(d *schema.ResourceData, meta interface{}, engine string) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allDBaaSConfigurations []gocmcapiv2.DBv2Configuration
	if configurationId := d.Get("configuration_id").(string); configurationId != "" {
		var configuration gocmcapiv2.DBv2Configuration
		var err error
		if d.Get("is_default").(bool) {
			configuration, err = client.DBaaSConfiguration.GetDefaultConfiguration(configurationId)
		} else {
			configuration, err = client.DBaaSConfiguration.Get(configurationId)
		}
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve postgres configuration [%s]: %s", configurationId, err)
			}
		}
		allDBaaSConfigurations = append(allDBaaSConfigurations, configuration)
	} else {
		params := map[string]string{
			// "name":          d.Get("name").(string),
			// "database_mode": d.Get("database_mode").(string),
			"page":          "1",
			"size":          "1000",
			"datastoreCode": engine,
		}
		if d.Get("is_default").(bool) {
			params["getDefault"] = "true"
		}
		configurations, err := client.DBaaSConfiguration.List(params)
		if err != nil {
			return fmt.Errorf("error when get %s configuration %v", engine, err)
		}
		allDBaaSConfigurations = append(allDBaaSConfigurations, configurations...)
	}
	if len(allDBaaSConfigurations) > 0 {
		var filteredDBaaSConfigurations []gocmcapiv2.DBv2Configuration
		for _, configuration := range allDBaaSConfigurations {
			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(configuration.Name), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("database_mode").(string); v != "" {
				if !strings.Contains(strings.ToLower(configuration.DatastoreMode), strings.ToLower(v)) {
					continue
				}
			}
			filteredDBaaSConfigurations = append(filteredDBaaSConfigurations, configuration)
		}
		allDBaaSConfigurations = filteredDBaaSConfigurations
	}
	if len(allDBaaSConfigurations) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allDBaaSConfigurations) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allDBaaSConfigurations)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeDBaaSConfigurationAttributes(d, allDBaaSConfigurations[0])
}

func dataSourceComputeDBaaSConfigurationAttributes(d *schema.ResourceData, configuration gocmcapiv2.DBv2Configuration) error {
	log.Printf("[DEBUG] Retrieved configuration %s: %#v", configuration.ID, configuration)

	_ = d.Set("name", configuration.Name)
	_ = d.Set("database_mode", configuration.DatastoreMode)
	if configuration.ID2 != "" {
		_ = d.Set("configuration_id", configuration.ID2)
		d.SetId(configuration.ID2)
	} else {
		d.SetId(configuration.ID)
		_ = d.Set("configuration_id", configuration.ID)
	}
	return nil
}
