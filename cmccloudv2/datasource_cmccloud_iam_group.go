package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIamGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of group",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceIamGroup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIamGroupRead,
		Schema: datasourceIamGroupSchema(),
	}
}

func dataSourceIamGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var alls []gocmcapiv2.IamGroup

	params := map[string]string{}
	alls, err := client.IamGroup.List(params)
	if err != nil {
		return fmt.Errorf("error when get iam groups %v", err)
	}
	if len(alls) > 0 {
		var filtered []gocmcapiv2.IamGroup
		for _, group := range alls {
			if v := d.Get("name").(string); v != "" {
				if v != group.Name {
					continue
				}
			}
			filtered = append(filtered, group)
		}
		alls = filtered
	}

	if len(alls) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(alls) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", alls)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	group := alls[0]
	d.SetId(group.Name)
	_ = d.Set("name", group.Name)
	_ = d.Set("description", group.Description)
	return nil
}
