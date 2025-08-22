package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIamUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeString,
			Description: "Filter by full username",
			Optional:    true,
			ForceNew:    true,
		},
		"email": {
			Type:        schema.TypeString,
			Description: "filter by email of user",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceIamUser() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceIamUserRead,
		Schema: datasourceIamUserSchema(),
	}
}

func dataSourceIamUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var alls []gocmcapiv2.IamUser

	params := map[string]string{
		"email":    d.Get("email").(string),
		"username": d.Get("username").(string),
	}
	alls, err := client.IamUser.List(params)
	if err != nil {
		return fmt.Errorf("error when get iam users %v", err)
	}
	if len(alls) > 0 {
		var filtered []gocmcapiv2.IamUser
		for _, user := range alls {
			if v := d.Get("email").(string); v != "" {
				if v != user.Email {
					continue
				}
			}
			if v := d.Get("username").(string); v != "" {
				if v != user.Username {
					continue
				}
			}
			filtered = append(filtered, user)
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

	user := alls[0]
	d.SetId(user.Username)
	_ = d.Set("username", user.Username)
	_ = d.Set("email", user.Email)
	return nil
}
