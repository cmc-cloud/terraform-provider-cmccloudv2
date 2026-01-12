package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func keymanagementcontainerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			ForceNew:     true,
			Description:  "The name of the key management container",
		},
		// "type": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validation.StringInSlice([]string{"generic", "rsa", "certificate"}, true),
		// },
		"container_ref": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The reference of the key management container",
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
