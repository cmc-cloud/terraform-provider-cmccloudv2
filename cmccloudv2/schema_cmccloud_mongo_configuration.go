package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func mongoConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Mongo configuration",
		},
		"database_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The database version of the Mongo configuration. Example `6.0`",
		},
		"database_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"replica_set", "standalone"}, false),
			Description:  "The database mode of the Mongo configuration, `replica_set`, `standalone`",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the Mongo configuration",
		},
		"parameters": {
			Type:        schema.TypeMap,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of the parameters for the Mongo configuration",
		},
	}
}
