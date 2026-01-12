package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceSubnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the subnet",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Description: "Filter by vpc id",
			Optional:    true,
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by ip address of subnet",
			Optional:    true,
			ForceNew:    true,
		},
		"cidr": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Filter by cidr of subnet",
		},
		"gateway_ip": {
			Type:        schema.TypeString,
			Description: "Filter by gateway_ip",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			ForceNew:    true,
			Description: "Created time of subnet",
		},
	}
}

func datasourceSubnet() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceSubnetRead,
		Schema: datasourceSubnetSchema(),
	}
}

func dataSourceSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allSubnets []gocmcapiv2.Subnet
	if subnet_id := d.Get("subnet_id").(string); subnet_id != "" {
		subnet, err := client.Subnet.Get(subnet_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve subnet [%s]: %s", subnet_id, err)
			}
		}
		allSubnets = append(allSubnets, subnet)
	} else {
		params := map[string]string{
			"vpc_id": d.Get("vpc_id").(string),
		}
		subnets, err := client.Subnet.List(params)
		if err != nil {
			return fmt.Errorf("error when get subnets %v", err)
		}
		allSubnets = append(allSubnets, subnets...)
	}
	if len(allSubnets) > 0 {
		var filteredSubnets []gocmcapiv2.Subnet
		for _, subnet := range allSubnets {
			if v := d.Get("cidr").(string); v != "" {
				if v != subnet.Cidr {
					continue
				}
			}
			if v := d.Get("gateway_ip").(string); v != "" {
				if subnet.GatewayIP != v {
					continue
				}
			}
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(subnet.Name, v) {
					continue
				}
			}
			filteredSubnets = append(filteredSubnets, subnet)
		}
		allSubnets = filteredSubnets
	}
	if len(allSubnets) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allSubnets) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allSubnets)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeSubnetAttributes(d, allSubnets[0])
}

func dataSourceComputeSubnetAttributes(d *schema.ResourceData, subnet gocmcapiv2.Subnet) error {
	log.Printf("[DEBUG] Retrieved subnet %s: %#v", subnet.ID, subnet)
	d.SetId(subnet.ID)
	_ = d.Set("name", subnet.Name)
	_ = d.Set("cidr", subnet.Cidr)
	_ = d.Set("gateway_ip", subnet.GatewayIP)
	_ = d.Set("vpc_id", subnet.VpcID)
	_ = d.Set("created_at", subnet.CreatedAt)
	return nil
}
