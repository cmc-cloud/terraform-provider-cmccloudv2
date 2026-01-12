package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func parameterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The key of the parameter",
			// ForceNew: true,
		},
		"value": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The value of the parameter",
			// ForceNew: true,
		},
		// "string_type": {
		// 	Type:     schema.TypeBool,
		// 	Optional: true,
		// 	// ForceNew: true,
		// },
	}
}
func databaseConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the database configuration",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the database configuration",
		},
		"datastore_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The type of the datastore",
		},
		"datastore_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The version of the datastore",
		},
		"parameters": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: parameterSchema(),
			},
			Description: "The parameters of the database configuration",
		},
	}
}
