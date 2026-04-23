package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func keyvaultSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the KeyVault instance",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the KeyVault instance",
		},
		"mode": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"ha_cluster",
			}, false),
			Description: "Mode of the KeyVault instance",
		},
		"version": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Description:      "The KeyVault version of the KeyVault instance",
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
			ForceNew:     true,
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
			ForceNew:    true,
			Description: "The size of the volume in GB",
		},
		"slave_count": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(2),
			Description:  "Slave count, required if mode is master_slave or ha_cluster",
		},
		"proxy_flavor_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Flavor for proxy node",
		},
		"proxy_quantity": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(2, 5),
			Description:  "Number of proxy nodes",
		},
		"tags": tagSchema(),
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the KeyVault instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the KeyVault instance",
		},
	}
}
