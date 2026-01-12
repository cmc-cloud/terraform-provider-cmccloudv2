package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createServerVolumesElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"delete_on_termination": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}
func serverSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "Name of server",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
			Description:  "Name of billing mode, monthly or hourly",
		},
		"zone": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "Name of zone, example: AZ1, AZ2",
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "The flavor ID of the desired flavor for the server. Changing this resizes the existing server",
		},
		"source_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Type of source: image, snapshot or volume. Changing this replaces the existing server.",
		},
		"source_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "ID of source image, snapshot or snapshot. Changing this replaces the existing server.",
		},
		"volume_name": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Name of root volume",
			ValidateFunc: validation.NoZeroValues,
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Size of root volume in GB",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Name of volume type of root volume",
		},
		"delete_on_termination": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "Set to false if you want to keep the root volume after the server is deleted",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
			Description:  "ID of subnet",
		},
		"ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "IP address of server, if not provided, a new IP address will be assigned",
		},
		"volumes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createServerVolumesElementSchema(),
			},
			Description: "The volumes for server",
		},
		"security_group_names": {
			Type: schema.TypeSet,
			Set:  schema.HashString,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Description: "An array of one or more security group names to associate with the server. Changing this results in adding/removing security groups from the existing server",
		},
		"tags": tagSchema(),
		"ecs_group_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "ID of ECS group, server will be placed according to the rules defined in this ECS group",
		},
		"key_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The name of a key pair to put on the server. The key pair must already be created and associated with the current account. Changing this creates a new server",
		},
		"user_data": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "User data of server, it will be executed when the server is created",
		},
		"password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "Password of server, it will be used to login to the server. Changing this changes the root password on the existing server.",
		},
		"interface_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The interface ID (port id) of the server",
		},
		"created": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the server",
		},
		"vm_state": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "active",
			ValidateFunc: validation.StringInSlice([]string{
				"active", "stopped",
			}, true),
			DiffSuppressFunc: suppressPowerStateDiffs,
			Description:      "State of server, active or stopped",
		},
	}
}

// suppressPowerStateDiffs will allow a state of "error" or "migrating" even though we don't
// allow them as a user input.
func suppressPowerStateDiffs(_, old, _ string, _ *schema.ResourceData) bool {
	if old == "error" || old == "migrating" {
		return true
	}

	return false
}
