package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func postgresDatabaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Postgres instance id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Postgres database name",
		},
		"owner": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Postgres username of owner",
		},
	}
}
