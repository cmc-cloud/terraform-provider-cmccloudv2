package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func createServerNicsElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// "security_group_ids": &schema.Schema{
		// 	Type: schema.TypeList,
		// 	Elem: &schema.Schema{
		// 		Type:         schema.TypeString,
		// 		ValidateFunc: validateUUID,
		// 	},
		// 	Optional: true,
		// },
	}
}

func createServerVolumesElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"delete_on_termination": {
			Type: schema.TypeBool,
			// Optional: true,
			// Default:  true,
			Computed: true,
		},
		// "size": {
		// 	Type:     schema.TypeInt,
		// 	Required: true,
		// },
		// "name": {
		// 	Type:     schema.TypeString,
		// 	Optional: true,
		// },
		// "type": {
		// 	Type:        schema.TypeString,
		// 	Description: "Volume type, eg: highio/commonio",
		// 	Required:    true,
		// },
		// "status": {
		// 	Type:        schema.TypeString,
		// 	Description: "Volume status",
		// 	Computed:    true,
		// },
		// "created_at": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },
		// "attachment_id": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },
	}
}
func serverSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			ValidateFunc: validateBillingMode,
			Default:      "monthly",
			Optional:     true,
		},
		"zone": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
		},
		"source_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"source_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"volume_type": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"delete_on_termination": {
			Type:     schema.TypeBool,
			Required: false,
			Default:  true,
			ForceNew: true,
		},
		// "nics": {
		// 	Type: schema.TypeList,
		// 	// Type: schema.TypeSet,
		// 	// Set: schema.HashResource(&schema.Resource{
		// 	// 	Schema: createServerNicsElementSchema(),
		// 	// }),
		// 	Required: true,
		// 	ForceNew: true,
		// 	Elem: &schema.Resource{
		// 		Schema: createServerNicsElementSchema(),
		// 	},
		// },
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			ForceNew:     true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"volumes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: createServerVolumesElementSchema(),
			},
		},
		"security_group_names": {
			Type: schema.TypeSet,
			Set:  schema.HashString,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"tags": tagSchema(),
		"ecs_group_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"key_name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"user_data": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.NoZeroValues,
		},
		"interface_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vm_state": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "active",
			ValidateFunc: validation.StringInSlice([]string{
				"active", "stopped",
			}, true),
			DiffSuppressFunc: suppressPowerStateDiffs,
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
