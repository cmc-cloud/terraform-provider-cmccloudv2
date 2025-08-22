package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIamRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name",
			Optional:    true,
			ForceNew:    true,
		},
		"title": {
			Type:        schema.TypeString,
			Description: "filter by title",
			Optional:    true,
			ForceNew:    true,
		},
		"region_id": {
			Type:        schema.TypeString,
			Description: "filter by region id",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceIamRole() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIamRoleRead,
		Schema: datasourceIamRoleSchema(),
	}
}

func dataSourceIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var alls []gocmcapiv2.IamRole

	params := map[string]string{
		"name":      d.Get("name").(string),
		"title":     d.Get("title").(string),
		"region_id": d.Get("region_id").(string),
	}
	alls, err := client.IamRole.List(params)
	if err != nil {
		return fmt.Errorf("error when get iam roles %v", err)
	}

	if len(alls) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(alls) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", alls)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	role := alls[0]
	d.SetId(role.ID)
	_ = d.Set("name", role.Name)
	_ = d.Set("title", role.Title)
	return nil
}
