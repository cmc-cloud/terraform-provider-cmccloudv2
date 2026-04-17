package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kafkaInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Kafka instance",
		},
		"billing_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateBillingMode,
			Description:  "The billing mode of the Kafka instance",
		},
		"mode": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"cluster", "single_node",
			}, false),
			Description: "mode of the Kafka instance",
		},
		"version": {
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: caseInsensitiveDiffSuppress,
			Description:      "The Kafka version of the Kafka instance",
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
		// "volume_type": {
		// 	Type:        schema.TypeString,
		// 	Required:    true,
		// 	ForceNew:    true,
		// 	Description: "The volume type",
		// },
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The size of the volume in GB",
		},
		"broker_quantity": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(3),
			Description:  "The quantity of broker, required if mode is cluster and must be >= 3",
		},
		"enable_basic_authen": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "True if you want to enable basic authen",
		},
		"users": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Username for user",
					},
					"password": {
						Type:      schema.TypeString,
						Required:  true,
						Sensitive: true,
						// ValidateFunc: validateKafkaPassword,
						Description: "Password for user",
					},
				},
			},
			Description: "A list of user objects containing username and password. Required if enable_basic_authen is true",
		},
		"security_group_ids": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
			Optional:    true,
			Description: "Set of security group IDs that you want to attach to the kafka",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the Kafka instance",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the Kafka instance",
		},
	}
}
