package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kubernetesv2Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateNameK8s,
			Description:  "Cluster name, only allow digits, characters and -",
			ForceNew:     true,
		},
		"zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"kubernetes_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"master_count": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"master_flavor_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		},
		// "max_node_count": {
		// 	Type:     schema.TypeInt,
		// 	Required: true,
		// 	ForceNew: true,
		// },
		"cidr_block_pod": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateIPCidrRange,
		},
		"cidr_block_service": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPCidrRange,
			ForceNew:     true,
		},
		"network_driver": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "calico",
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		},

		"enable_autohealing": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"enable_monitoring": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"enable_autoscale": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"autoscale_max_node": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"autoscale_max_ram_gb": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"autoscale_max_core": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
