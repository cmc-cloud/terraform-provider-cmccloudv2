package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceRedisConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the redis configuration",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of redis configuration (case-insenitive)",
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

func datasourceRedisConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceRedisConfigurationRead,
		Schema: datasourceRedisConfigurationSchema(),
	}
}

func dataSourceRedisConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allRedisConfigurations []gocmcapiv2.RedisConfiguration
	if configuration_id := d.Get("configuration_id").(string); configuration_id != "" {
		configuration, err := client.RedisConfiguration.Get(configuration_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve redis configuration [%s]: %s", configuration_id, err)
			}
		}
		allRedisConfigurations = append(allRedisConfigurations, configuration)
	} else {
		params := map[string]string{
			// "name":          d.Get("name").(string),
			// "database_mode": d.Get("database_mode").(string),
			"page":          "1",
			"size":          "1000",
			"datastoreCode": "redis",
		}
		configurations, err := client.RedisConfiguration.List(params)
		if err != nil {
			return fmt.Errorf("error when get redis configuration %v", err)
		}
		allRedisConfigurations = append(allRedisConfigurations, configurations...)
	}
	if len(allRedisConfigurations) > 0 {
		var filteredRedisConfigurations []gocmcapiv2.RedisConfiguration
		for _, configuration := range allRedisConfigurations {
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
			filteredRedisConfigurations = append(filteredRedisConfigurations, configuration)
		}
		allRedisConfigurations = filteredRedisConfigurations
	}
	if len(allRedisConfigurations) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allRedisConfigurations) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allRedisConfigurations)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeRedisConfigurationAttributes(d, allRedisConfigurations[0])
}

func dataSourceComputeRedisConfigurationAttributes(d *schema.ResourceData, configuration gocmcapiv2.RedisConfiguration) error {
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
