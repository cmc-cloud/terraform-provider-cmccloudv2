package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceELBSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"elb_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the elb",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of elb, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"ip_address": {
			Type:        schema.TypeString,
			Description: "Filter by ip address of elb, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceELB() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceELBRead,
		Schema: datasourceELBSchema(),
	}
}

func dataSourceELBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allELBs []gocmcapiv2.ELB
	if elb_id := d.Get("elb_id").(string); elb_id != "" {
		elb, err := client.ELB.Get(elb_id)
		if err != nil {
			return fmt.Errorf("unable to retrieve elb [%s]: %s", elb_id, err)
		}
		allELBs = append(allELBs, elb)
	} else {
		params := map[string]string{}
		elbs, err := client.ELB.List(params)
		if err != nil {
			return fmt.Errorf("error when get elbs %v", err)
		}
		allELBs = append(allELBs, elbs...)
	}
	if len(allELBs) > 0 {
		var filteredELBs []gocmcapiv2.ELB
		for _, elb := range allELBs {
			if v := d.Get("ip_address").(string); v != "" {
				if v != elb.VipAddress {
					continue
				}
			}
			if v := d.Get("name").(string); v != "" {
				if v != elb.Name {
					continue
				}
			}
			filteredELBs = append(filteredELBs, elb)
		}
		allELBs = filteredELBs
	}
	if len(allELBs) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allELBs) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allELBs)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeELBAttributes(d, allELBs[0])
}

func dataSourceComputeELBAttributes(d *schema.ResourceData, elb gocmcapiv2.ELB) error {
	d.SetId(elb.ID)
	_ = d.Set("ip_address", elb.VipAddress)
	_ = d.Set("name", elb.Name)
	_ = d.Set("created_at", elb.CreatedAt)
	return nil
}
