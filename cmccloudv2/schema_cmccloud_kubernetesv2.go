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
			Description:  "The billing mode of the Kubernetes cluster. A valid value is monthly or hourly.",
		},
		"zone": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The zone of the Kubernetes cluster",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The ID of the subnet to attach the Kubernetes cluster to",
		},
		"kubernetes_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The version of the Kubernetes cluster",
		},
		"master_count": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The number of master nodes in the Kubernetes cluster",
		},
		"master_flavor_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
			Description:  "The flavor name of the master node in the Kubernetes cluster",
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
			Description:  "The CIDR block of the pod in the Kubernetes cluster",
		},
		"cidr_block_service": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateIPCidrRange,
			ForceNew:     true,
			Description:  "The CIDR block of the service in the Kubernetes cluster",
		},
		"node_mask_cidr": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      24,
			ValidateFunc: validation.IntBetween(0, 32),
			ForceNew:     true,
			Description:  "The node mask CIDR of the Kubernetes cluster",
		},
		// "ntp_enabled": {
		// 	Type:     schema.TypeBool,
		// 	Optional: true,
		// 	Default:  false,
		// },
		"ntp_servers": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Description: "List of NTP servers in the Kubernetes cluster",
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
			Description:  "The network driver of the Kubernetes cluster",
		},
		"network_driver_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "overlay",
			ValidateFunc: validation.StringInSlice([]string{"overlay", "native-routing"}, false),
			ForceNew:     true,
			Description:  "The network driver mode of the Kubernetes cluster",
		},
		"enable_autohealing": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The enable autohealing addon for the Kubernetes cluster",
		},

		"enable_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The enable monitoring addon for the Kubernetes cluster",
		},
		"enable_autoscale": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The enable autoscale addon for the Kubernetes cluster",
		},
		"autoscale_max_node": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The autoscale max node of the Kubernetes cluster",
		},
		"autoscale_max_ram_gb": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The autoscale max ram gb of the Kubernetes cluster",
		},
		"autoscale_max_core": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The autoscale max core of the Kubernetes cluster",
		},
		"state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of the Kubernetes cluster",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the Kubernetes cluster",
		},
	}
}
