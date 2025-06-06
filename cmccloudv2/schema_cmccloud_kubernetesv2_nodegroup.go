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
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
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
			Default:      110,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 256),
		},
		"init_current_node": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "Node count for initing node group",
		},
		"node_metadatas": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"value": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"effect": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"type": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Description: "Array of node metadata objects. Each object must have key, value, type, and optionally effect.",
		},
		// "node_metadatas": {
		// 	Type:     schema.TypeList,
		// 	Optional: true,
		// 	Description: "List of metadata objects for nodes. Each object can have arbitrary key-value pairs, e.g. [{\"key\": \"group\", \"value\": \"cmccloud\", \"type\": \"label\"}, {\"key\": \"thotd\", \"value\": \"thotd123\", \"effect\": \"NoSchedule\", \"type\": \"taint\"}]",
		// 	Elem: &schema.Schema{
		// 		Type: schema.TypeMap,
		// 		Elem: &schema.Schema{Type: schema.TypeString},
		// 	},
		// },
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
