package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceAutoScalingConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autoscaling_configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the autoscaling configuration",
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceAutoScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceAutoScalingConfigurationRead,
		Schema: datasourceAutoScalingConfigurationSchema(),
	}
}

func dataSourceAutoScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allAutoScalingConfigurations []gocmcapiv2.AutoScalingConfiguration
	if autoscalingconfiguration_id := d.Get("autoscaling_configuration_id").(string); autoscalingconfiguration_id != "" {
		autoscalingconfiguration, err := client.AutoScalingConfiguration.Get(autoscalingconfiguration_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve autoscaling configuration [%s]: %s", autoscalingconfiguration_id, err)
			}
		}
		allAutoScalingConfigurations = append(allAutoScalingConfigurations, autoscalingconfiguration)
	} else {
		params := map[string]string{
			"name": d.Get("name").(string),
		}
		autoscalingconfigurations, err := client.AutoScalingConfiguration.List(params)
		if err != nil {
			return fmt.Errorf("error when get autoscaling configuration %v", err)
		}
		allAutoScalingConfigurations = append(allAutoScalingConfigurations, autoscalingconfigurations...)
	}
	if len(allAutoScalingConfigurations) > 0 {
		var filteredAutoScalingConfigurations []gocmcapiv2.AutoScalingConfiguration
		for _, autoscalingconfiguration := range allAutoScalingConfigurations {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(autoscalingconfiguration.Name, v) {
					continue
				}
			}
			filteredAutoScalingConfigurations = append(filteredAutoScalingConfigurations, autoscalingconfiguration)
		}
		allAutoScalingConfigurations = filteredAutoScalingConfigurations
	}
	if len(allAutoScalingConfigurations) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allAutoScalingConfigurations) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allAutoScalingConfigurations)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeAutoScalingConfigurationAttributes(d, allAutoScalingConfigurations[0])
}

func dataSourceComputeAutoScalingConfigurationAttributes(d *schema.ResourceData, autoscalingconfiguration gocmcapiv2.AutoScalingConfiguration) error {
	log.Printf("[DEBUG] Retrieved autoscaling configuration %s: %#v", autoscalingconfiguration.ID, autoscalingconfiguration)
	d.SetId(autoscalingconfiguration.ID)
	_ = d.Set("name", autoscalingconfiguration.Name)
	_ = d.Set("created_at", autoscalingconfiguration.CreatedAt)
	return nil
}
