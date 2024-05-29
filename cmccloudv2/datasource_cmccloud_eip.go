package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceEIPSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"eip_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the eip",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Description: "Filter by ip address of eip",
			Optional:    true,
			ForceNew:    true,
		},
		"fixed_ip_address": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"device_id": {
			Type:        schema.TypeString,
			Description: "Filter by device_id that eip assosiate with",
			Optional:    true,
			ForceNew:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "filter by eip status (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},

		"description": {
			Type:        schema.TypeString,
			Description: "filter by eip that contains this text (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},

		"dns_name": {
			Type:        schema.TypeString,
			Description: "filter by dns_name that contains this text (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},

		"dns_domain": {
			Type:        schema.TypeString,
			Description: "filter by dns_domain that contains this text (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceEIP() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceEIPRead,
		Schema: datasourceEIPSchema(),
	}
}

func dataSourceEIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allEIPs []gocmcapiv2.EIP
	if eip_id := d.Get("eip_id").(string); eip_id != "" {
		eip, err := client.EIP.Get(eip_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve eip [%s]: %s", eip_id, err)
			}
		}
		allEIPs = append(allEIPs, eip)
	} else {
		params := map[string]string{
			"fixed_ip_address": d.Get("fixed_ip_address").(string),
		}
		eips, err := client.EIP.List(params)
		if err != nil {
			return fmt.Errorf("Error when get eips %v", err)
		}
		allEIPs = append(allEIPs, eips...)
	}
	if len(allEIPs) > 0 {
		var filteredEIPs []gocmcapiv2.EIP
		for _, eip := range allEIPs {
			if v := d.Get("ip_address").(string); v != "" {
				if v != eip.FloatingIPAddress {
					continue
				}
			}
			if v := d.Get("device_id").(string); v != "" {
				if eip.PortDetails.DeviceID != v {
					continue
				}
			}
			if v := d.Get("status").(string); v != "" {
				if strings.ToLower(eip.Status) != strings.ToLower(v) {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(eip.Description), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("dns_name").(string); v != "" {
				if !strings.Contains(strings.ToLower(eip.DNSName), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("dns_domain").(string); v != "" {
				if !strings.Contains(strings.ToLower(eip.DNSDomain), strings.ToLower(v)) {
					continue
				}
			}

			filteredEIPs = append(filteredEIPs, eip)
		}
		allEIPs = filteredEIPs
	}
	if len(allEIPs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allEIPs) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allEIPs)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeEIPAttributes(d, allEIPs[0])
}

func dataSourceComputeEIPAttributes(d *schema.ResourceData, eip gocmcapiv2.EIP) error {
	log.Printf("[DEBUG] Retrieved eip %s: %#v", eip.ID, eip)
	d.SetId(eip.ID)
	d.Set("ip_address", eip.FloatingIPAddress)
	d.Set("fix_ip_address", eip.FixedIPAddress)
	d.Set("device_id", eip.PortDetails.DeviceID)
	d.Set("status", eip.Status)
	d.Set("description", eip.Description)
	d.Set("dns_name", eip.DNSName)
	d.Set("dns_domain", eip.DNSDomain)
	d.Set("created_at", eip.CreatedAt)
	return nil
}
