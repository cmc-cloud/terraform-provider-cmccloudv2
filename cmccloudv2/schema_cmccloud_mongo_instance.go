package cmccloudv2

import (
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
			Description: "Mode of the Mongo instance, allow values are replica_set, standalone",
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
		"tags": tagSchema(),
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
