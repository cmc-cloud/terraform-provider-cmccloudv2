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
		},
		// "type": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validation.StringInSlice([]string{"generic", "rsa", "certificate"}, true),
		// },
		"container_ref": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
