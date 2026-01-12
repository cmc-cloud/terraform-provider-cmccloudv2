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
			Description:  "The ID of the subnet",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specific IP address for the server interface",
		},
	}
}
func databaseinstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the database instance",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The ID of the flavor",
		},
		"zone": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The zone of the database instance",
		},
		"source_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The source type of the database instance",
		},
		"source_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the source",
		},
		"datastore_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The type of the datastore",
		},
		"datastore_version": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The version of the datastore",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The type of the volume",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The size of the volume in GB",
		},
		"subnets": {
			Type:     schema.TypeList, // TypeList => (where ordering doesn’t matter), TypeList (where ordering matters).
			Required: true,
			Elem: &schema.Resource{
				Schema: createDatabaseInstanceSubnetsElementSchema(),
			},
			Description: "The subnets of the database instance",
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
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Password of the admin user",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
			Description:  "The billing mode of the database instance, can be monthly or hourly",
		},
		"replicate_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Only use when source_type = instance",
		},
	}
}
