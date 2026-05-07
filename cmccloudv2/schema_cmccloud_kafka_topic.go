package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kafkaTopicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The id of the Kafka instance",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Kafka topic",
		},
		"partition_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      6,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Number of partitions in a topic. More partitions = more parallel processing (higher throughput).",
		},
		"replication_factor": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			ForceNew:    true,
			Description: "Number of copies of each partition across brokers. Higher = better fault tolerance and data safety.",
		},
		"rentation_day": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 99999),
			Description:  "Number of partitions in a topic. More partitions = more parallel processing (higher throughput).",
		},
	}
}
