package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func postgresConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Postgres configuration",
		},
		"database_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The database version of the Postgres configuration. Example `15`, `16`, `17`",
		},
		"database_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"master_slave", "ha_cluster", "standalone"}, false),
			Description:  "The database mode of the Postgres configuration, `master_slave`, `ha_cluster`, `standalone`",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the Postgres configuration",
		},
		"parameters": {
			Type:        schema.TypeMap,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of the parameters for the Postgres configuration",
		},
	}
}
