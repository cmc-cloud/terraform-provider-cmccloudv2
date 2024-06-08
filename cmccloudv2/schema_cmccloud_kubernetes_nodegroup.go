package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func kubernetesNodeGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			ForceNew:     true,
			Description:  "Name of node group, this can't be changed after created",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
		},
		"node_count": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"min_node_count": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"max_node_count": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "hourly",
			Optional:     true,
		},
		"docker_volume_size": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"docker_volume_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"zone": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
