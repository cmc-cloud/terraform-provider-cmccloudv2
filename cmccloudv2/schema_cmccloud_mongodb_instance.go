package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func mongodbinstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateBillingMode,
		},
		"source_type": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"new", "backup", "instance",
			}, false),
			ForceNew: true,
		},
		"source_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"backup_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"database_engine": {
			Type:             schema.TypeString,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Required:         true,
			ForceNew:         true,
			Description:      "match exactly, case-insensitive",
		},
		"database_version": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			ForceNew:         true,
			Description:      "match exactly",
		},
		"database_mode": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			ForceNew:         true,
			Description:      "search by text, case-insensitive",
		},
		"zone_primary": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		},
		"zone_secondaries": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.NoZeroValues,
			ForceNew:     true,
		},
		"slave_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"security_group_ids": { // securityGroups
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Optional: true,
		},
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"volume_type": {
			Type:         schema.TypeString, // name
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"mongodb_configuration_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			Description:  "Password must be 8-32 characters long, include upper and lower case letters and numbers, and contain no special characters.",
			ValidateFunc: validateRedisPassword,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
