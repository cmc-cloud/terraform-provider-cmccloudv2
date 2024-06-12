package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createDatabaseInstanceSubnetsElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}
func databaseinstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"zone": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"source_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"source_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"datastore_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"datastore_version": {
			Type:     schema.TypeString,
			Required: true,
		},
		"volume_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"subnets": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Required: true,
			Elem: &schema.Resource{
				Schema: createDatabaseInstanceSubnetsElementSchema(),
			},
		},
		"enable_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
			Description: "Enable public ip on the database instance",
		},
		"is_public": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Whether the database service is exposed to the public",
		},
		"allowed_cidrs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "A list of IPv4, IPv6 or mix of both CIDRs that restrict access to the database service",
		},
		"allowed_host": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User allowed host",
		},
		"admin_user": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "admin",
		},
		"admin_password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"replicate_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Only use when source_type = instance",
		},
	}
}
