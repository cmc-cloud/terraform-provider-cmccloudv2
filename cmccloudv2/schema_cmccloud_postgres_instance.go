package cmccloudv2

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func postgresInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Postgres instance",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the Postgres instance",
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
				"master_slave", "ha_cluster", "standalone",
			}, false),
			Description: "mode of the Postgres instance",
		},
		"version": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Description:      "The Postgres version of the Postgres instance",
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
		"configuration_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Configuration id, if empty, postgres instance will use default configuration template. Configuration template must have the same mode and database version",
		},
		"admin_password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ForceNew:     true,
			Description:  "Admin password",
			ValidateFunc: validatePostgresPassword,
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "5432",
			Description: "Postgres port",
		},
		"slave_count": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1,
			Description: "Slave count, required if mode is master_slave or ha_cluster",
		},
		"proxy_flavor_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     2,
			Description: "Flavor of Proxy, required if mode is ha_cluster",
		},
		"proxy_quantity": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      2,
			ValidateFunc: validation.IntBetween(2, 5),
			Description:  "Proxy quantity, required if mode is ha_cluster",
		},
		"retention_period": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(3, 7),
			ForceNew:     true,
			Description:  "Retention period in days, default is 3",
		},
		"tags": tagSchema(),
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the Postgres instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the Postgres instance",
		},
	}
}

// Custom validation function
func validatePostgresPassword(val interface{}, key string) (warns []string, errs []error) {
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
