package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceAutoScalingGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"autoscaling_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the autoscaling group",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the autoscaling group, exact match",
			ForceNew:    true,
		},
		"configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Configuration id of the autoscaling group, exact match",
			ForceNew:    true,
		},
		"configuration_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Configuration name of the autoscaling group, exact match",
			ForceNew:    true,
		},
		"status": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"active", "init", "creating", "updating", "deleting", "resizing", "warning", "paused", "error", "critical"}, true),
			ForceNew:     true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceAutoScalingGroupRead,
		Schema: datasourceAutoScalingGroupSchema(),
	}
}

func dataSourceAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allAutoScalingGroups []gocmcapiv2.AutoScalingGroup
	if autoscalinggroup_id := d.Get("autoscaling_group_id").(string); autoscalinggroup_id != "" {
		autoscalinggroup, err := client.AutoScalingGroup.Get(autoscalinggroup_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve autoscaling group [%s]: %s", autoscalinggroup_id, err)
			}
		}
		allAutoScalingGroups = append(allAutoScalingGroups, autoscalinggroup)
	} else {
		params := map[string]string{
			"name": d.Get("name").(string),
		}
		autoscalinggroups, err := client.AutoScalingGroup.List(params)
		if err != nil {
			return fmt.Errorf("error when get autoscaling group %v", err)
		}
		allAutoScalingGroups = append(allAutoScalingGroups, autoscalinggroups...)
	}
	if len(allAutoScalingGroups) > 0 {
		var filteredAutoScalingGroups []gocmcapiv2.AutoScalingGroup
		for _, autoscalinggroup := range allAutoScalingGroups {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(autoscalinggroup.Name, v) {
					continue
				}
			}
			if v := d.Get("configuration_id").(string); v != "" {
				if !strings.EqualFold(autoscalinggroup.ProfileID, v) {
					continue
				}
			}
			if v := d.Get("configuration_name").(string); v != "" {
				if !strings.EqualFold(autoscalinggroup.ProfileName, v) {
					continue
				}
			}
			if v := d.Get("status").(string); v != "" {
				if !strings.EqualFold(autoscalinggroup.Status, v) {
					continue
				}
			}
			filteredAutoScalingGroups = append(filteredAutoScalingGroups, autoscalinggroup)
		}
		allAutoScalingGroups = filteredAutoScalingGroups
	}
	if len(allAutoScalingGroups) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allAutoScalingGroups) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allAutoScalingGroups)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeAutoScalingGroupAttributes(d, allAutoScalingGroups[0])
}

func dataSourceComputeAutoScalingGroupAttributes(d *schema.ResourceData, autoscalinggroup gocmcapiv2.AutoScalingGroup) error {
	log.Printf("[DEBUG] Retrieved autoscaling group %s: %#v", autoscalinggroup.ID, autoscalinggroup)
	d.SetId(autoscalinggroup.ID)
	_ = d.Set("name", autoscalinggroup.Name)
	_ = d.Set("created_at", autoscalinggroup.CreatedAt)
	return nil
}
