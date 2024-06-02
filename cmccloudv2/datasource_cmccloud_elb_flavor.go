package cmccloudv2

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceFlavorELBSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"flavor_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the flavor",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of flavor (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Filter flavor that description contains this text (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceFlavorForELB() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceFlavorELBRead,
		Schema: datasourceFlavorELBSchema(),
	}
}
func dataSourceFlavorELBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allFlavors []gocmcapiv2.ELBFlavor
	if flavor_id := d.Get("flavor_id").(string); flavor_id != "" {
		flavor, err := client.ELB.GetFlavor(flavor_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve flavor [%s]: %s", flavor_id, err)
			}
		}
		allFlavors = append(allFlavors, flavor)
	} else {
		flavors, err := client.ELB.ListFlavors()
		if err != nil {
			return fmt.Errorf("Error when get flavors %v", err)
		}
		allFlavors = append(allFlavors, flavors...)
	}
	if len(allFlavors) > 0 {
		var filteredFlavors []gocmcapiv2.ELBFlavor
		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if strings.ToLower(flavor.Name) != strings.ToLower(v) {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(flavor.Description), strings.ToLower(v)) {
					continue
				}
			}
			filteredFlavors = append(filteredFlavors, flavor)
		}
		allFlavors = filteredFlavors
	}
	if len(allFlavors) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allFlavors) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allFlavors)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeFlavorELBAttributes(d, allFlavors[0])
}

func dataSourceComputeFlavorELBAttributes(d *schema.ResourceData, flavor gocmcapiv2.ELBFlavor) error {
	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("description", flavor.Description)
	return nil
}
