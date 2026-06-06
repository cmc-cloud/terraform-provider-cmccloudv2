package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func validateGroupName(v interface{}, k string) (warnings []string, errors []error) {
	re := `^[a-zA-Z][a-zA-Z0-9_-]*$`
	return validateRegexp(re)(v, k)
}
func iamGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateGroupName,
			Description:  "The name of the IAM group.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The descriptionof the IAM group.",
		},
	}
}
