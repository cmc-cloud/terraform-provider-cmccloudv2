package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIamProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of project",
			Optional:    true,
			ForceNew:    true,
		},
		// "region_id": {
		// 	Type:        schema.TypeString,
		// 	Description: "filter by region id",
		// 	Optional:    true,
		// 	ForceNew:    true,
		// },
	}
}

func datasourceIamProject() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIamProjectRead,
		Schema: datasourceIamProjectSchema(),
	}
}

func dataSourceIamProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var alls []gocmcapiv2.IamProject

	params := map[string]string{}
	alls, err := client.IamProject.List(params)
	if err != nil {
		return fmt.Errorf("error when get iam projects %v", err)
	}
	if len(alls) > 0 {
		var filtered []gocmcapiv2.IamProject
		for _, project := range alls {
			if v := d.Get("name").(string); v != "" {
				if v != project.Name {
					continue
				}
			}
			// if v := d.Get("region_id").(string); v != "" {
			// 	if v != project.RegionId {
			// 		continue
			// 	}
			// }
			filtered = append(filtered, project)
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

	project := alls[0]
	d.SetId(project.ID)
	_ = d.Set("name", project.Name)
	// d.Set("region_id", project.RegionId)
	return nil
}
