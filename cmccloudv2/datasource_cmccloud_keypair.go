package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceKeypairSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of keypair (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Filter by type of keypair (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			ForceNew:    true,
			Description: "Created time of keypair",
		},
		"fingerprint": {
			Type:        schema.TypeString,
			Computed:    true,
			ForceNew:    true,
			Description: "The fingerprint of the keypair",
		},
	}
}

func datasourceKeypair() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceKeypairRead,
		Schema: datasourceKeypairSchema(),
	}
}

func dataSourceKeypairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allKeypairs []gocmcapiv2.Keypair
	if name := d.Get("name").(string); name != "" {
		keypair, err := client.Keypair.Get(name)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve keypair [%s]: %s", name, err)
			}
		}
		allKeypairs = append(allKeypairs, keypair)
	} else {
		params := map[string]string{}
		keypairs, err := client.Keypair.List(params)
		if err != nil {
			return fmt.Errorf("error when get keypairs %v", err)
		}
		allKeypairs = append(allKeypairs, keypairs...)
	}
	if len(allKeypairs) > 0 {
		var filteredKeypairs []gocmcapiv2.Keypair
		for _, keypair := range allKeypairs {
			if v := d.Get("type").(string); v != "" {
				if !strings.EqualFold(keypair.Name, v) {
					continue
				}
			}
			filteredKeypairs = append(filteredKeypairs, keypair)
		}
		allKeypairs = filteredKeypairs
	}
	if len(allKeypairs) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allKeypairs) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allKeypairs)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeKeypairAttributes(d, allKeypairs[0])
}

func dataSourceComputeKeypairAttributes(d *schema.ResourceData, keypair gocmcapiv2.Keypair) error {
	log.Printf("[DEBUG] Retrieved keypair %s: %#v", keypair.Type, keypair)
	d.SetId(keypair.Name)
	_ = d.Set("name", keypair.Name)
	_ = d.Set("type", keypair.Type)
	_ = d.Set("fingerprint", keypair.Fingerprint)
	return nil
}
