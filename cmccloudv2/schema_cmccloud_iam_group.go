package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func iamGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUserName,
			Description:  "The name of the IAM group.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The descriptionof the IAM group.",
		},
	}
}
