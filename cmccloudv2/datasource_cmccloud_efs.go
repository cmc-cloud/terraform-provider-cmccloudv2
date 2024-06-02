package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceEFSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"efs_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the efs",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter EFS by name, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"endpoint": {
			Type:        schema.TypeString,
			Description: "Filter EFS by endpoint, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "Filter EFS that has description contains this text",
			Computed:    true,
			ForceNew:    true,
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Description: "Filter EFS by vpc_id, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Description: "Filter EFS by subnet_id, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"capacity": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceEFS() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceEFSRead,
		Schema: datasourceEFSSchema(),
	}
}

func dataSourceEFSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allEFSs []gocmcapiv2.EFS
	if efs_id := d.Get("efs_id").(string); efs_id != "" {
		efs, err := client.EFS.Get(efs_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve efs [%s]: %s", efs_id, err)
			}
		}
		allEFSs = append(allEFSs, efs)
	} else {
		params := map[string]string{}
		efss, err := client.EFS.List(params)
		if err != nil {
			return fmt.Errorf("Error when get efses %v", err)
		}
		allEFSs = append(allEFSs, efss...)
	}
	if len(allEFSs) > 0 {
		var filteredEFSs []gocmcapiv2.EFS
		for _, efs := range allEFSs {
			if v := d.Get("name").(string); v != "" {
				if strings.ToLower(efs.Name) != strings.ToLower(v) {
					continue
				}
			}
			if v := d.Get("endpoint").(string); v != "" {
				if strings.ToLower(efs.Endpoint) != strings.ToLower(v) {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(efs.Description), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("vpc_id").(string); v != "" {
				if strings.ToLower(efs.VpcID) != strings.ToLower(v) {
					continue
				}
			}
			if v := d.Get("subnet_id").(string); v != "" {
				if strings.ToLower(efs.VpcID) != strings.ToLower(v) {
					continue
				}
			}
			filteredEFSs = append(filteredEFSs, efs)
		}
		allEFSs = filteredEFSs
	}
	if len(allEFSs) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allEFSs) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allEFSs)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeEFSAttributes(d, allEFSs[0])
}

func dataSourceComputeEFSAttributes(d *schema.ResourceData, efs gocmcapiv2.EFS) error {
	log.Printf("[DEBUG] Retrieved efs %s: %#v", efs.ID, efs)
	d.SetId(efs.ID)
	d.Set("name", efs.Name)
	d.Set("type", efs.Type)
	d.Set("status", efs.Status)
	d.Set("capacity", efs.Capacity)
	d.Set("protocol_type", efs.ProtocolType)
	d.Set("vpc_id", efs.VpcID)
	d.Set("subnet_id", efs.SubnetID)
	d.Set("created_at", efs.CreatedAt)
	d.Set("endpoint", efs.Endpoint)
	d.Set("shared_path", efs.SharedPath)
	d.Set("command_line", efs.CommandLine)
	return nil
}
