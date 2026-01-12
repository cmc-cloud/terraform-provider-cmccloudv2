package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ecsgroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the ECS Group",
		},
		"policy": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"soft-anti-affinity", "soft-affinity"}, true),
			Optional:     true,
			Default:      "soft-anti-affinity",
			ForceNew:     true,
			Description:  "The policy of the ECS Group",
		},
	}
}
