package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func mysqlConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Mysql configuration",
		},
		"database_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The database version of the Mysql configuration. Example `8.0`",
		},
		"database_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"replica_set", "standalone"}, false),
			Description:  "The database mode of the Mysql configuration, `replica_set`, `standalone`",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the Mysql configuration",
		},
		"parameters": {
			Type:        schema.TypeMap,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of the parameters for the Mysql configuration",
		},
	}
}
