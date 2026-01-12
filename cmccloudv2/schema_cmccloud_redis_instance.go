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
			Description:  "The name of the Redis instance",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the Redis instance",
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
		"backup_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The ID of the redis backup",
		},
		"database_engine": {
			Type:             schema.TypeString,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Required:         true,
			ForceNew:         true,
			Description:      "The database engine of the Redis instance, match exactly, case-insensitive",
		},
		"database_version": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			ForceNew:         true,
			Description:      "The database version of the Redis instance, match exactly, case-insensitive",
		},
		"database_mode": {
			Type:             schema.TypeString,
			Required:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			ForceNew:         true,
			Description:      "The database mode of the Redis instance, case-insensitive",
		},
		"replicas": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Description: "The number of replicas of the Redis instance, required for Redis Cluster mode",
		},
		// "security_group_ids": { // securityGroups
		// 	Type: schema.TypeSet,
		// 	Elem: &schema.Schema{
		// 		Type:         schema.TypeString,
		// 		ValidateFunc: validateUUID,
		// 	},
		// 	Optional:    true,
		// 	Description: "The security group IDs of the Redis instance",
		// },
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The database flavor ID, using cmccloudv2_flavor_dbaas data source to get the flavor ID from flavor name",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The volume type of the volume",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
			Description: "The size of the volume in GB",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the subnet of the Redis instance",
		},
		"redis_configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the Redis configuration of the Redis instance",
		},
		"password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			Description:  "Password must be 8-32 characters long, include upper and lower case letters and numbers, and contain no special characters.",
			ValidateFunc: validateRedisPassword,
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the Redis instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the Redis instance",
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
