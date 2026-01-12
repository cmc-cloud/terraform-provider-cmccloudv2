package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func elbPoolMemberSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the pool to attach the member to",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the member",
		},
		"address": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The IP address of the backend member server",
		},
		"protocol_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "The protocol port number for the backend member server.",
		},
		"weight": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 256),
			Description:  "The weight of a member determines the portion of requests or connections it services compared to the other members of the pool. For example, a member with a weight of 10 receives five times as many requests as a member with a weight of 2. A value of 0 means the member does not receive new connections but continues to service existing connections. A valid value is from 0 to 256. Default is 1.",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
			Description:  "The subnet ID the member service is accessible from.",
		},
		"monitor_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "An alternate IP address used for health monitoring a backend member",
		},
		"monitor_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "An alternate protocol port used for health monitoring a backend member",
		},
		"backup": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Is the member a backup? Backup members only receive traffic when all non-backup members are down.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The creation time of the member",
			Computed:    true,
		},
		"provisioning_status": {
			Type:        schema.TypeString,
			Description: "The state of the operation — in other words, whether member is still creating, updating, or deleting the resource.",
			Computed:    true,
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shows the real runtime health of the member.",
		},
	}
}
