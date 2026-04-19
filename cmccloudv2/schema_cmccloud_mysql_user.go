package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func mysqlUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mysql instance id",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mysql username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Mysql user password",
		},
		"hosts": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "`%` to allow all IPs, specific IP address (e.g., 192.168.1.1) or multiple IP addresses",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"user_permissions": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "List of user permissions",
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"database": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Database name to which the permissions apply. Set it to `*` if you want to allow all databases",
					},
					"table": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Table name to which the permissions apply. Set it to `*` if you want to allow all tables",
					},
					"permissions": {
						Type:        schema.TypeSet,
						Required:    true,
						Description: "List of permissions (alter, create, delete, drop, insert, select, update, index, create view, trigger, event, references). Set it to `*` if you want to allow permissions",
						Elem: &schema.Schema{
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"*",
								"alter",
								"create",
								"delete",
								"drop",
								"insert",
								"select",
								"update",
								"index",
								"create view",
								"trigger",
								"event",
								"references",
							}, false), // false = case-sensitive
						},
					},
				},
			},
		},
	}
}
