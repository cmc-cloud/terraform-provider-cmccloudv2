package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func elbSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"network_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"public", "private"}, true),
		},
		"subnet_id": {
			Type:          schema.TypeString,
			Description:   "Required when network_type = private",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"bandwidth_mbps"},
		},
		"bandwidth_mbps": {
			Type:          schema.TypeInt,
			Description:   "Used when network_type = public",
			Optional:      true,
			ConflictsWith: []string{"subnet_id"},
		},
		"tags": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
