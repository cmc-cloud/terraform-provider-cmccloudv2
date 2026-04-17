package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourcePostgresConfigurationSchema() map[string]*schema.Schema {
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

func datasourcePostgresConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePostgresConfigurationRead,
		Schema: datasourcePostgresConfigurationSchema(),
	}
}

func dataSourcePostgresConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allPostgresConfigurations []gocmcapiv2.PostgresConfiguration
	if configurationId := d.Get("configuration_id").(string); configurationId != "" {
		var configuration gocmcapiv2.PostgresConfiguration
		var err error
		if d.Get("is_default").(bool) {
			configuration, err = client.PostgresConfiguration.GetDefaultConfiguration(configurationId)
		} else {
			configuration, err = client.PostgresConfiguration.Get(configurationId)
		}
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve postgres configuration [%s]: %s", configurationId, err)
			}
		}
		allPostgresConfigurations = append(allPostgresConfigurations, configuration)
	} else {
		params := map[string]string{
			// "name":          d.Get("name").(string),
			// "database_mode": d.Get("database_mode").(string),
			"page":          "1",
			"size":          "1000",
			"datastoreCode": "postgresql",
		}
		if d.Get("is_default").(bool) {
			params["getDefault"] = "true"
		}
		configurations, err := client.PostgresConfiguration.List(params)
		if err != nil {
			return fmt.Errorf("error when get postgres configuration %v", err)
		}
		allPostgresConfigurations = append(allPostgresConfigurations, configurations...)
	}
	if len(allPostgresConfigurations) > 0 {
		var filteredPostgresConfigurations []gocmcapiv2.PostgresConfiguration
		for _, configuration := range allPostgresConfigurations {
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
			filteredPostgresConfigurations = append(filteredPostgresConfigurations, configuration)
		}
		allPostgresConfigurations = filteredPostgresConfigurations
	}
	if len(allPostgresConfigurations) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allPostgresConfigurations) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allPostgresConfigurations)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputePostgresConfigurationAttributes(d, allPostgresConfigurations[0])
}

func dataSourceComputePostgresConfigurationAttributes(d *schema.ResourceData, configuration gocmcapiv2.PostgresConfiguration) error {
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
