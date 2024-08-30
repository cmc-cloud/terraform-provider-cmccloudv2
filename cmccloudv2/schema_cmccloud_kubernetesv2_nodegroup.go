package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kubernetesv2NodeGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": { // clusterId
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
		},
		"name": { // nameNodeGroup
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			ForceNew:     true,
			Description:  "Name of node group, this can't be changed after created",
		},
		"zone": { // zone
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"flavor_id": { // flavorId
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		},
		"key_name": { // sshKeyName
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"volume_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"security_group_ids": { // securityGroups
			Type: schema.TypeList,
			// Set:  schema.HashString,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Optional: true,
			ForceNew: true,
		},
		"enable_autoscale": { // isAutoscale
			Type:     schema.TypeBool,
			Required: true,
		},
		"min_node": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"max_node": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"max_pods": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"cpu_threshold_percent": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(10),
		},
		"memory_threshold_percent": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(10),
		},
		"disk_threshold_percent": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(10),
		},
		"image_gpu_tag": { // sshKeyName
			Type:     schema.TypeString,
			Optional: true,
			Default:  "default",
			ForceNew: true,
		},
		"enable_autohealing": { // isAutoscale
			Type:     schema.TypeBool,
			Required: true,
		},
		"max_unhealthy_percent": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"node_startup_timeout_minutes": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// "labels": {
		// 	Type:     schema.TypeMap,
		// 	Optional: true,
		// },
		// "taint": {
		// 	Type:     schema.TypeMap,
		// 	Optional: true,
		// },
		// "kubernetes_labels": {
		// 	Type:     schema.TypeMap,
		// 	Optional: true,
		// },
	}
}
