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
			Description:  "The ID of the Kubernetes cluster to attach the node group to",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
			Description:  "The billing mode of the node group. A valid value is monthly or hourly.",
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
			Description:  "The zone of the node group",
		},
		"flavor_id": { // flavorId
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
			Description:  "The flavor ID of the node group",
		},
		"key_name": { // sshKeyName
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The SSH key name of the node group",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The volume type of the node group",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The volume size of the node group",
		},
		"security_group_ids": { // securityGroups
			Type: schema.TypeList,
			// Set:  schema.HashString,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Optional:    true,
			ForceNew:    true,
			Description: "The security group IDs of the node group",
		},
		"enable_autoscale": { // isAutoscale
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The enable autoscale of the node group",
		},
		"min_node": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "The min node of the node group, if enable_autoscale is true, this is required",
		},
		"max_node": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "The max node of the node group, if enable_autoscale is true, this is required",
		},
		"max_pods": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      110,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 256),
			Description:  "The max pods of the node group, if enable_autoscale is true, this is required",
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
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The key of the node metadata",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The value of the node metadata",
					},
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The type of the node metadata",
					},
					"effect": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The effect of the node metadata",
					},
				},
			},
			Description: "Array of node metadata objects. Each object must have key, value, type, and optionally effect.",
		},
		// "cidr_block_pod": {
		// 	Type:         schema.TypeString,
		// 	Required:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validateIPCidrRange,
		// },
		// "ntp_enabled": {
		// 	Type:     schema.TypeBool,
		// 	Optional: true,
		// 	Default:  false,
		// },
		"ntp_servers": {
			Type:        schema.TypeList, // dùng List để giữ thứ tự
			Optional:    true,
			ForceNew:    true,
			Description: "List of NTP servers (order matters).",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"host": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "NTP server hostname or IP.",
					},
					"port": {
						Type:        schema.TypeInt,
						Optional:    true,
						Default:     123,
						Description: "NTP server port.",
					},
					"protocol": {
						Type:         schema.TypeString,
						Optional:     true,
						Default:      "udp",
						Description:  "Protocol used (udp or tcp).",
						ValidateFunc: validation.StringInSlice([]string{"udp", "tcp"}, false),
					},
				},
			},
		},
		"image_gpu_tag": { // sshKeyName
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "default",
			ForceNew:    true,
			Description: "The image GPU tag of the node group",
		},
		"enable_autohealing": { // isAutoscale
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The enable autohealing of the node group",
		},
		"max_unhealthy_percent": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
			Description:  "The max unhealthy percent of the node group, if enable_autohealing is true, this is required",
		},
		"node_startup_timeout_minutes": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 100),
			Description:  "The node startup timeout minutes of the node group, if enable_autohealing is true, this is required",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the subnet to attach the node group to",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the node group",
		},
	}
}
