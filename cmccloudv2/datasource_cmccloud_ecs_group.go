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

func datasourceEcsGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ecs_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the ecsgroup",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of ecsgroup (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"policy": {
			Type:         schema.TypeString,
			Description:  "filter by policy type",
			ValidateFunc: validation.StringInSlice([]string{"soft-anti-affinity", "soft-affinity"}, true),
			Optional:     true,
			ForceNew:     true,
		},
	}
}

func datasourceEcsGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceEcsGroupRead,
		Schema: datasourceEcsGroupSchema(),
	}
}

func dataSourceEcsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allEcsGroups []gocmcapiv2.EcsGroup
	if ecs_group_id := d.Get("ecs_group_id").(string); ecs_group_id != "" {
		ecsgroup, err := client.EcsGroup.Get(ecs_group_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve ecs group [%s]: %s", ecs_group_id, err)
			}
		}
		allEcsGroups = append(allEcsGroups, ecsgroup)
	} else {
		params := map[string]string{
			"policy": d.Get("policy").(string),
		}
		ecsgroups, err := client.EcsGroup.List(params)
		if err != nil {
			return fmt.Errorf("error when get ecs groups %v", err)
		}
		allEcsGroups = append(allEcsGroups, ecsgroups...)
	}
	if len(allEcsGroups) > 0 {
		var filteredEcsGroups []gocmcapiv2.EcsGroup
		for _, ecsgroup := range allEcsGroups {
			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(ecsgroup.Name), strings.ToLower(v)) {
					continue
				}
			}
			filteredEcsGroups = append(filteredEcsGroups, ecsgroup)
		}
		allEcsGroups = filteredEcsGroups
	}
	if len(allEcsGroups) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allEcsGroups) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allEcsGroups)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeEcsGroupAttributes(d, allEcsGroups[0])
}

func dataSourceComputeEcsGroupAttributes(d *schema.ResourceData, ecsgroup gocmcapiv2.EcsGroup) error {
	log.Printf("[DEBUG] Retrieved ecsgroup %s: %#v", ecsgroup.ID, ecsgroup)
	d.SetId(ecsgroup.ID)
	_ = d.Set("name", ecsgroup.Name)
	_ = d.Set("policy", ecsgroup.Policy)
	_ = d.Set("ecs_group_id", ecsgroup.ID)
	return nil
}
