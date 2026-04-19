package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func rdsClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the RDS cluster.",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the Postgres instance",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The ID of the flavor for the RDS cluster.",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The volume size (in GiB) for the RDS cluster.",
		},
		"db_engine": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The database engine (e.g., postgres, mysql) for the RDS cluster.",
		},
		"db_version": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The version of the database engine.",
		},
		"mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"cluster"}, false),
			Description:  "The operation mode for the cluster, currently support `cluster` only.",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The subnet ID associated with the RDS cluster.",
		},
		"cluster_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of database instances in the cluster.",
		},
		"proxy_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of proxy nodes in the cluster.",
		},
		"enable_backup": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Enable automated daily backup.",
		},
		"enable_pitr": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable Point-in-Time Recovery (PITR).",
		},
		"backup_schedule": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Cron schedule string for automated backups.",
		},
		"backup_retention": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of backup retention days.",
		},
		"enable_storage_autoscaling": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Whether to enable storage autoscaling.",
		},
		"storage_autoscaling_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Storage autoscaling threshold (percent).",
		},
		"storage_autoscaling_increment": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Storage autoscaling increment (GiB) when triggered.",
		},
		"lb_vip_ipaddress": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Load Balancer VIP IP Address",
		},
		"tags": tagSchema(),
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Cluster status",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time",
		},
	}
}
