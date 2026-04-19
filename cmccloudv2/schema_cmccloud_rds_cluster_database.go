package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func rdsClusterDatabaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "RdsCluster id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "RdsCluster database name",
		},
		"character_set": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "utf8mb4",
			ForceNew:    true,
			Description: "Character set, default is `utf8mb4`",
		},
		"collation": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "utf8mb4_unicode_ci",
			Description: "Collation, default is `utf8mb4_unicode_ci`",
		},
	}
}
