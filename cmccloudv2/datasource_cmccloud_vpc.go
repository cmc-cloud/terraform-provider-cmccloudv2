package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVPCSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the vpc",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by ip address of vpc",
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "filter by vpc description (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"cidr": {
			Type:        schema.TypeString,
			Description: "filter by vpc cidr",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			ForceNew:    true,
			Description: "Created time of vpc",
		},
	}
}

func datasourceVPC() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceVPCRead,
		Schema: datasourceVPCSchema(),
	}
}

func dataSourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allVPCs []gocmcapiv2.VPC
	if vpc_id := d.Get("vpc_id").(string); vpc_id != "" {
		vpc, err := client.VPC.Get(vpc_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve vpc [%s]: %s", vpc_id, err)
			}
		}
		allVPCs = append(allVPCs, vpc)
	} else {
		params := map[string]string{
			"name": d.Get("name").(string),
		}
		vpcs, err := client.VPC.List(params)
		if err != nil {
			return fmt.Errorf("error when get vpcs %v", err)
		}
		allVPCs = append(allVPCs, vpcs...)
	}
	if len(allVPCs) > 0 {
		var filteredVPCs []gocmcapiv2.VPC
		for _, vpc := range allVPCs {
			if v := d.Get("cidr").(string); v != "" {
				if v != vpc.Cidr {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(vpc.Description), strings.ToLower(v)) {
					continue
				}
			}
			filteredVPCs = append(filteredVPCs, vpc)
		}
		allVPCs = filteredVPCs
	}
	if len(allVPCs) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allVPCs) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allVPCs)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeVPCAttributes(d, allVPCs[0])
}

func dataSourceComputeVPCAttributes(d *schema.ResourceData, vpc gocmcapiv2.VPC) error {
	log.Printf("[DEBUG] Retrieved vpc %s: %#v", vpc.ID, vpc)
	d.SetId(vpc.ID)
	_ = d.Set("cidr", vpc.Cidr)
	_ = d.Set("name", vpc.Name)
	_ = d.Set("description", vpc.Description)
	_ = d.Set("created_at", vpc.CreatedAt)
	return nil
}
