package cmccloudv2

import (
	"errors"
	"fmt"
	"log"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourcePortSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the port",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Description: "Filter by ip address of port",
			Optional:    true,
			ForceNew:    true,
		},
		"server_id": {
			Type:        schema.TypeString,
			Description: "Filter by server id",
			Optional:    true,
			ForceNew:    true,
		},
		"is_public": {
			Type:        schema.TypeBool,
			Description: "True if ip address is public, False if ip address is private",
			Optional:    true,
			ForceNew:    true,
		},
		"ips": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "List of ip address of this port",
			Computed:    true,
		},
	}
}

func datasourcePort() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourcePortRead,
		Schema: datasourcePortSchema(),
	}
}

func dataSourcePortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allPorts []gocmcapiv2.Port
	if port_id := d.Get("port_id").(string); port_id != "" {
		port, err := client.Port.Get(port_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve port [%s]: %s", port_id, err)
			}
		}
		allPorts = append(allPorts, port)
	} else {
		params := map[string]string{
			"server_id": d.Get("server_id").(string),
		}
		ports, err := client.Port.List(params)
		if err != nil {
			return fmt.Errorf("error when get ports %v", err)
		}
		allPorts = append(allPorts, ports...)
	}
	if len(allPorts) > 0 {
		var filteredPorts []gocmcapiv2.Port
		for _, port := range allPorts {
			found := false
			if v := d.Get("ip_address").(string); v != "" {
				for _, ip := range port.FixedIps {
					if v == ip.IPAddress {
						found = true
					}
				}
				if !found {
					continue
				}
			}
			found = false
			//nolint:staticcheck // Need tri-state bool support
			if v, ok := d.GetOkExists("is_public"); ok {
				for _, ip := range port.FixedIps {
					if IsPublicIP(ip.IPAddress) == v.(bool) {
						found = true
					}
				}
				if !found {
					continue
				}
			}
			filteredPorts = append(filteredPorts, port)
		}
		allPorts = filteredPorts
	}
	if len(allPorts) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allPorts) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: ", allPorts)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputePortAttributes(d, allPorts[0])
}

func dataSourceComputePortAttributes(d *schema.ResourceData, port gocmcapiv2.Port) error {
	log.Printf("[DEBUG] Retrieved port %s: %#v", port.ID, port)
	d.SetId(port.ID)
	_ = d.Set("port_id", port.ID)
	ips := []string{}
	for _, ip := range port.FixedIps {
		ips = append(ips, ip.IPAddress)
	}
	_ = d.Set("ips", ips)

	return nil
}
