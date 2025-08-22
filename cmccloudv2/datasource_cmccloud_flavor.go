package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceFlavorSchema() map[string]*schema.Schema {
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
func datasourceFlavorDBSchema() map[string]*schema.Schema {
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
		"disk": {
			Type:        schema.TypeInt,
			Description: "Filter by disk size (GB) of flavor",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceFlavorForEC() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceFlavorECRead,
		Schema: datasourceFlavorSchema(),
	}
}
func datasourceFlavorForDB() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceFlavorDBRead,
		Schema: datasourceFlavorDBSchema(),
	}
}

func datasourceFlavorForK8s() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceFlavorK8sRead,
		Schema: datasourceFlavorDBSchema(),
	}
}
func dataSourceFlavorECRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceFlavorRead(d, meta, "EC")
}
func dataSourceFlavorDBRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceFlavorRead(d, meta, "DBaas")
}
func dataSourceFlavorK8sRead(d *schema.ResourceData, meta interface{}) error {
	return dataSourceFlavorRead(d, meta, "K8s")
}
func dataSourceFlavorRead(d *schema.ResourceData, meta interface{}, flavor_type string) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allFlavors []gocmcapiv2.Flavor
	if flavor_id := d.Get("flavor_id").(string); flavor_id != "" {
		flavor, err := client.Flavor.Get(flavor_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve flavor [%s]: %s", flavor_id, err)
			}
		}
		allFlavors = append(allFlavors, flavor)
	} else {
		params := map[string]string{}
		if d.Get("name").(string) != "" {
			params["name"] = d.Get("name").(string)
		}
		if v, ok := d.GetOk("cpu"); ok && v.(int) != 0 {
			params["cpu"] = strconv.Itoa(d.Get("cpu").(int))
		}
		if v, ok := d.GetOk("ram"); ok && v.(int) != 0 {
			params["ram"] = strconv.Itoa(d.Get("ram").(int))
		}
		if v, ok := d.GetOk("disk"); ok && v.(int) != 0 {
			params["disk"] = strconv.Itoa(d.Get("disk").(int))
		}
		flavors, err := client.Flavor.List(params)
		if err != nil {
			return fmt.Errorf("error when get flavors %v", err)
		}
		allFlavors = append(allFlavors, flavors...)
	}
	if len(allFlavors) > 0 {
		var filteredFlavors []gocmcapiv2.Flavor
		for _, flavor := range allFlavors {
			// check type
			if flavor_type == "DBaas" && !flavor.ExtraSpecs.IsDatabaseFlavor {
				continue
			}
			if flavor_type == "K8s" && !flavor.ExtraSpecs.IsK8sFlavor {
				continue
			}
			if flavor_type == "EC" && (flavor.ExtraSpecs.IsK8sFlavor || flavor.ExtraSpecs.IsDatabaseFlavor) {
				continue
			}
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
			if v, ok := d.GetOk("disk"); ok {
				if v.(int) != flavor.Disk {
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

	return dataSourceComputeFlavorAttributes(d, allFlavors[0], flavor_type)
}

func dataSourceComputeFlavorAttributes(d *schema.ResourceData, flavor gocmcapiv2.Flavor, flavor_type string) error {
	log.Printf("[DEBUG] Retrieved flavor %s: %#v", flavor.ID, flavor)
	d.SetId(flavor.ID)
	_ = d.Set("name", flavor.Name)
	_ = d.Set("cpu", flavor.Vcpus)
	_ = d.Set("ram", flavor.RAM/1024)
	if flavor.Disk > 0 {
		_ = d.Set("disk", flavor.Disk)
	}
	return nil
}
