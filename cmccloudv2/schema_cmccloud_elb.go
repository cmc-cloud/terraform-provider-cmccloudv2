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
			Description:  "Name of billing mode, monthly or hourly",
		},
		"zone": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The availability zone for the ELB. Changing this creates a new ELB",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "The flavor ID of the desired flavor for the ELB. Changing this resizes the existing ELB",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the ELB",
		},
		"network_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"public", "private"}, true),
			Description:  "The network type of the ELB, public or private. Changing this creates a new ELB",
		},
		"subnet_id": {
			Type:          schema.TypeString,
			Description:   "Required when network_type is private",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"bandwidth_mbps"},
		},
		"vip_address": {
			Type:          schema.TypeString,
			Description:   "Only available when network_type is private",
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
		"tags": tagSchema(),
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the ELB. Changing this updates the ELB's description.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the ELB",
		},
		"provisioning_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of the control-plane operation — in other words, whether load balancer is still creating, updating, or deleting the resource.",
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shows the real runtime health of the load balancer data plane.",
		},
	}
}
