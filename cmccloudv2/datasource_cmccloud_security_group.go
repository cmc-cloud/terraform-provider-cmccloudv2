package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceSecurityGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"security_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the ecsgroup",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of ecsgroup, match exactly (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Filter by name of ecsgroup (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceSecurityGroupRead,
		Schema: datasourceSecurityGroupSchema(),
	}
}

func dataSourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allSecurityGroups []gocmcapiv2.SecurityGroup
	if security_group_id := d.Get("security_group_id").(string); security_group_id != "" {
		ecsgroup, err := client.SecurityGroup.Get(security_group_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve security group [%s]: %s", security_group_id, err)
			}
		}
		allSecurityGroups = append(allSecurityGroups, ecsgroup)
	} else {
		params := map[string]string{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
		}
		ecsgroups, err := client.SecurityGroup.List(params)
		if err != nil {
			return fmt.Errorf("error when get security groups %v", err)
		}
		allSecurityGroups = append(allSecurityGroups, ecsgroups...)
	}
	if len(allSecurityGroups) > 0 {
		var filteredSecurityGroups []gocmcapiv2.SecurityGroup
		for _, ecsgroup := range allSecurityGroups {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(ecsgroup.Name, v) {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(ecsgroup.Description), strings.ToLower(v)) {
					continue
				}
			}
			filteredSecurityGroups = append(filteredSecurityGroups, ecsgroup)
		}
		allSecurityGroups = filteredSecurityGroups
	}
	if len(allSecurityGroups) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allSecurityGroups) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allSecurityGroups)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeSecurityGroupAttributes(d, allSecurityGroups[0])
}

func dataSourceComputeSecurityGroupAttributes(d *schema.ResourceData, ecsgroup gocmcapiv2.SecurityGroup) error {
	log.Printf("[DEBUG] Retrieved ecsgroup %s: %#v", ecsgroup.ID, ecsgroup)
	d.SetId(ecsgroup.ID)
	_ = d.Set("name", ecsgroup.Name)
	_ = d.Set("description", ecsgroup.Description)
	_ = d.Set("security_group_id", ecsgroup.ID)
	return nil
}
