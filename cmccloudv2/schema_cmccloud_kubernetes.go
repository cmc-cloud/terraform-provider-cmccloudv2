package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createLabelsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"kube_tag": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"network_driver": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "calico",
		},
		"calico_ipv4pool": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"kube_dashboard_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"metrics_server_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"npd_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},
		"auto_healing_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},
		"auto_scaling_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},
		"max_node_count": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			// Default:     50,
			Description: "Required when auto_scaling_enabled = 'true'",
		},
		"min_node_count": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			// Default:     1,
			Description: "Required when auto_scaling_enabled = 'true'",
		},
	}
}

func kubernetesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
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
		"default_master": {
			Type:     schema.TypeList, // Use TypeList or TypeSet depending on your use case
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"node_count": {
						Type:     schema.TypeInt,
						Required: true,
						ForceNew: true,
					},
					"flavor_id": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
						ForceNew:     true,
					},
					"billing_mode": {
						Type:         schema.TypeString,
						ValidateFunc: validateBillingMode,
						Default:      "monthly",
						Optional:     true,
					},
				},
			},
		},
		"default_worker": {
			Type:     schema.TypeList, // Use TypeList or TypeSet depending on your use case
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"flavor_id": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validateUUID,
					},
					"node_count": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"billing_mode": {
						Type:         schema.TypeString,
						ValidateFunc: validateBillingMode,
						Default:      "hourly",
						Optional:     true,
					},
				},
			},
		},
		"keypair": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"labels": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: createLabelsSchema(),
			},
		},

		// "custom_labels": {
		// 	Type:     schema.TypeMap,
		// 	Elem: &schema.Schema{
		// 	  	Type: schema.TypeString,
		// 	},
		// },

		"create_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Timeout in minutes",
			Default:     120,
		},
	}
}
