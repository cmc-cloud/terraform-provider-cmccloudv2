package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func opensearchSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the OpenSearch cluster",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the OpenSearch instance",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The ID of the OpenSearch node flavor",
		},
		"dashboard_flavor_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the dashboard node flavor",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The volume size (in GiB) for the OpenSearch cluster",
		},
		"version": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The version of OpenSearch",
		},
		"admin_password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The admin password for OpenSearch access",
		},
		"node_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of data nodes in the OpenSearch cluster",
		},

		"enable_isolate_master": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Configuring master nodes to only manage the cluster, without handling data or user requests",
		},

		"master_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of master nodes in the OpenSearch cluster. Required if enable_isolate_master is true",
		},
		// "enable_drain_nodes": {
		// 	Type:        schema.TypeBool,
		// 	Optional:    true,
		// 	Default:     true,
		// 	Description: "Enable draining nodes for graceful maintenance",
		// },
		"dashboard_replicas": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of dashboard replicas",
		},
		"enable_snapshot": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Whether to enable automatic snapshots",
		},
		"snapshot_creation_cron": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "0 2 * * *",
			Description: "Cron expression for snapshot creation schedule",
		},
		"snapshot_timezone": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Asia/Ho_Chi_Minh",
			Description: "Timezone for snapshot scheduling",
		},
		"rentation_max_age": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     14,
			Description: "Maximum retention age for snapshots (days)",
		},
		"rentation_min_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3,
			Description: "Minimum number of snapshots to retain",
		},
		"rentation_max_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10,
			Description: "Maximum number of snapshots to retain",
		},
		"lb_subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Subnet ID for the Load Balancer",
		},
		// "enable_lb_internal": {
		// 	Type:        schema.TypeBool,
		// 	Optional:    true,
		// 	Default:     true,
		// 	Description: "Enable internal Load Balancer",
		// },
		"enable_storage_autoscaling": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Enable storage autoscaling",
		},
		"storage_autoscaling_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     70,
			Description: "Threshold (percent) to trigger storage autoscaling",
		},
		"storage_autoscaling_increment": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10,
			Description: "Increment size (GiB) for each autoscale event",
		},
		"storage_autoscaling_max": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     5000,
			Description: "Maximum total storage (GiB) after autoscaling",
		},
		"api_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "API domain name. Changes to this field do NOT trigger resource update. If not set, it will be api-<name>.internal",
			ForceNew:    false, // explicitly mark no update on change
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return false //old != new // always suppress diff => do not trigger update
			},
		},
		"dashboard_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Dashboard domain name. Changes to this field do NOT trigger resource update. If not set, it will be dashboard-<name>.internal",
			ForceNew:    false,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return false //old != new
			},
		},
		"tags": tagSchema(),
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "OpenSearch status",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time",
		},
	}
}
