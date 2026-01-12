package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func efsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_mode": {
			Type:         schema.TypeString,
			Default:      "monthly",
			Optional:     true,
			ValidateFunc: validateBillingMode,
		},
		"capacity": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Capacity in GB",
			ValidateFunc: validation.All(
				validation.IntDivisibleBy(100),
				validation.IntAtLeast(200),
			),
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the EFS",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The ID of the subnet to attach the EFS to",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"high_performance", "hdd_standard", "hdd_premium", "ssd_standard", "ssd_premium"}, true),
			Description:  "The type of the EFS, high_performance is deprecated",
		},
		"protocol_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"nfs"}, true),
			Description:  "The protocol type of the EFS",
		},
		"tags": tagSchema(),
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the EFS. Changing this updates the EFS's description.",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The endpoint of the EFS, it is used to mount the EFS to the instance",
		},
		"shared_path": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The shared path of the EFS, it is used to mount the EFS to the instance",
		},
		"command_line": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The command line of the EFS, it is used to mount the EFS to the instance",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the EFS",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the EFS",
		},
	}
}
