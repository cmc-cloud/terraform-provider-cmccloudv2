package cmccloudv2

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func redisinstanceSchema() map[string]*schema.Schema {
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
		// "mode": {
		// 	Type:     schema.TypeString,
		// 	Required: true,
		// 	ValidateFunc: validation.StringInSlice([]string{
		// 		"master/slave", "cluster", "standalone",
		// 	}, false),
		// 	ForceNew: true,
		// },
		"zones": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			MinItems:    1,
			Description: "If mode is standalone, only the first zone is accepted",
			Required:    true,
			ForceNew:    true,
		},
		// "source_type": {
		// 	Type:     schema.TypeString,
		// 	Required: true,
		// 	ValidateFunc: validation.StringInSlice([]string{
		// 		"new", "backup", "instance",
		// 	}, false),
		// 	ForceNew: true,
		// },
		// "source_id": {
		// 	Type:     schema.TypeString,
		// 	Optional: true,
		// 	ForceNew: true,
		// },
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
		// "zone_master": {
		// 	Type:         schema.TypeString,
		// 	Required:     true,
		// 	ValidateFunc: validation.NoZeroValues,
		// 	ForceNew:     true,
		// },
		// "zone_slave": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ValidateFunc: validation.NoZeroValues,
		// 	ForceNew:     true,
		// },
		"security_group_ids": { // securityGroups
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Optional: true,
		},
		"flavor_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		// "volume_type": {
		// 	Type:         schema.TypeString, // name
		// 	Required:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validation.NoZeroValues,
		// },
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
		// "ip_master": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validateIPAddress,
		// },
		// "ip_slave1": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validateIPAddress,
		// },
		// "ip_slave2": {
		// 	Type:         schema.TypeString,
		// 	Optional:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validateIPAddress,
		// },
		"redis_configuration_id": {
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

// Custom validation function
func validateRedisPassword(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)

	// Check for valid length
	if len(v) < 8 || len(v) > 32 {
		errs = append(errs, fmt.Errorf("%q must be between 8 and 32 characters long, got: %d", key, len(v)))
		return warns, errs
	}

	// Check for at least one lowercase letter
	if matched, _ := regexp.MatchString(`[a-z]`, v); !matched {
		errs = append(errs, fmt.Errorf("%q must contain at least one lowercase letter", key))
		return warns, errs
	}

	// Check for at least one uppercase letter
	if matched, _ := regexp.MatchString(`[A-Z]`, v); !matched {
		errs = append(errs, fmt.Errorf("%q must contain at least one uppercase letter", key))
		return warns, errs
	}

	// Check for at least one digit
	if matched, _ := regexp.MatchString(`\d`, v); !matched {
		errs = append(errs, fmt.Errorf("%q must contain at least one digit", key))
		return warns, errs
	}

	// Check for special characters (must not contain any)
	if matched, _ := regexp.MatchString(`[^A-Za-z0-9]`, v); matched {
		errs = append(errs, fmt.Errorf("%q must not contain special characters", key))
		return warns, errs
	}

	return warns, errs
}
