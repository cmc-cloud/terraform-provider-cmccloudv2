package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func parameterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
			// ForceNew: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
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
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"datastore_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"datastore_version": {
			Type:     schema.TypeString,
			Required: true,
		},
		"parameters": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: parameterSchema(),
			},
		},
	}
}
