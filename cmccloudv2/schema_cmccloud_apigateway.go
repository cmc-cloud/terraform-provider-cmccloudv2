package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func apigatewaySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the ApiGateway instance",
		},
		"mode": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "standalone",
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"standalone",
			}, false),
			Description: "Mode of the ApiGateway instance, currently only `standalone` is supported",
		},
		"zones": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			MinItems:    1,
			Description: "Currently only 1 zone is supported",
			Required:    true,
			ForceNew:    true,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "Flavor id",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "Subnet id",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The volume type",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The size of the volume in GB",
		},
		"public_access": {
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
			Description: "If set to true, this api gateway can be access from public network",
		},
		"bandwidth": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      500,
			ValidateFunc: validation.IntAtLeast(500),
			ForceNew:     true,
			Description:  "Required if public_access is set to true",
		},
		"tags": tagSchema(),
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the ApiGateway instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the ApiGateway instance",
		},
	}
}
