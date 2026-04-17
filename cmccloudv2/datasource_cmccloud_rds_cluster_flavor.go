package cmccloudv2

import (
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceRdsClusterFlavorSchema() map[string]*schema.Schema {
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
		"cpu": {
			Type:        schema.TypeInt,
			Description: "Filter by number of vcpus of flavor",
			Optional:    true,
			ForceNew:    true,
		},
		"ram": {
			Type:        schema.TypeInt,
			Description: "Filter by ram size (GB) of flavor",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceFlavorForRdsCluster() *schema.Resource {
	return &schema.Resource{
		Read:   datasourceFlavorForRdsClusterRead,
		Schema: datasourceRdsClusterFlavorSchema(),
	}
}

func datasourceFlavorForRdsClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allFlavors []gocmcapiv2.RdsClusterFlavor
	var err error

	allFlavors, err = client.RdsCluster.ListFlavors()
	if err != nil {
		return fmt.Errorf("error when get flavors %v", err)
	}
	if len(allFlavors) > 0 {
		var filteredFlavors []gocmcapiv2.RdsClusterFlavor
		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(flavor.Name, v) {
					continue
				}
			}
			if v, ok := d.GetOk("cpu"); ok {
				if v.(int) != flavor.Vcpus {
					continue
				}
			}
			if v, ok := d.GetOk("ram"); ok {
				if 1024*v.(int) != flavor.RAM {
					continue
				}
			}
			filteredFlavors = append(filteredFlavors, flavor)
		}
		allFlavors = filteredFlavors
	}
	if len(allFlavors) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allFlavors) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allFlavors)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeRdsClusterFlavorAttributes(d, allFlavors[0])
}

func dataSourceComputeRdsClusterFlavorAttributes(d *schema.ResourceData, flavor gocmcapiv2.RdsClusterFlavor) error {
	log.Printf("[DEBUG] Retrieved flavor %s: %#v", flavor.ID, flavor)
	d.SetId(flavor.ID)
	_ = d.Set("name", flavor.Name)
	_ = d.Set("cpu", flavor.Vcpus)
	_ = d.Set("ram", flavor.RAM)
	return nil
}
