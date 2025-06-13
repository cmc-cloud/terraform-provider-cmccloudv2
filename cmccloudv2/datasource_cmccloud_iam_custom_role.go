package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIamCustomRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by full custom role name",
			Required:    true,
			ForceNew:    true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceIamCustomRole() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIamCustomRoleRead,
		Schema: datasourceIamCustomRoleSchema(),
	}
}

func dataSourceIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var alls []gocmcapiv2.IamCustomRole

	params := map[string]string{
		"name": d.Get("name").(string),
	}
	alls, err := client.IamCustomRole.List(params)
	if err != nil {
		return fmt.Errorf("Error when get iam custom role %v", err)
	}
	if len(alls) > 0 {
		var filtered []gocmcapiv2.IamCustomRole
		for _, user := range alls {
			if v := d.Get("name").(string); v != "" {
				if v != user.Name {
					continue
				}
			}
			filtered = append(filtered, user)
		}
		alls = filtered
	}

	if len(alls) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(alls) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", alls)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	role := alls[0]
	d.SetId(role.ID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("created", role.Created)
	return nil
}
