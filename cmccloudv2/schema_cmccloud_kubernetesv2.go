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
		"node_mask_cidr": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 32),
			ForceNew:     true,
		},
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
		"network_driver": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "calico",
			ValidateFunc: validation.StringInSlice([]string{"calico", "cilium"}, false),
			ForceNew:     true,
		},
		"network_driver_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "overlay",
			ValidateFunc: validation.StringInSlice([]string{"overlay", "native-routing"}, false),
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
