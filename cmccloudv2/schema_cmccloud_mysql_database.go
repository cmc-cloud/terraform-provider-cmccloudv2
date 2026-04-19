package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func mysqlDatabaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mysql database id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mysql database name",
		},
		"character_set": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "utf8mb4",
			ForceNew:    true,
			Description: "Mysql character set, default is `utf8mb4`",
		},
		"collation": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "utf8mb4_unicode_ci",
			Description: "Mysql collation, default is `utf8mb4_unicode_ci`",
		},
	}
}
