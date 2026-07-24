package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kafkaConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Redis configuration",
		},
		"database_engine": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The database engine of the Kafka configuration, Kafka",
		},
		"database_version": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The database version of the Kafka configuration. Example `3.7`",
		},
		"database_mode": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"single_node", "cluster"}, false),
			Description:  "The database mode of the Kafka configuration, `single_node`, `cluster`",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the Kafka configuration",
		},
		"parameters": {
			Type:        schema.TypeMap,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of the parameters for the Kafka configuration",
		},
	}
}
