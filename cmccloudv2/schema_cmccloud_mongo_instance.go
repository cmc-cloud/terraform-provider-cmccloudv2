package cmccloudv2

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func mongoInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Mongo instance",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the Mongo instance",
		},
		"backup_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The backup id when you want to restore from backup",
		},
		"mode": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"replica_set", "standalone",
			}, false),
			Description: "mode of the Mongo instance, allow values are replica_set, standalone",
		},
		"version": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Description:      "The Mongo version of the Mongo instance",
		},
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
		"flavor_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "Flavor id",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "Subnet id",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The volume type",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The size of the volume in GB",
		},
		"quantity_of_secondary": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(2),
			Description:  "The quantity of secondary, required and must be >= 2 when mode is `replica_set`",
		},
		"configuration_id": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:    true,
			Description: "The id of configuration group",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the Mongo instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the Mongo instance",
		},
	}
}

// Custom validation function
func validateMongoPassword(val interface{}, key string) (warns []string, errs []error) {
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
