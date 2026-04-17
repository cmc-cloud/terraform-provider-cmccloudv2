package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func kafkaTopicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateName,
			Description:  "The name of the Kafka topic",
		},
		"partition_count": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of partitions in a topic. More partitions = more parallel processing (higher throughput).",
		},
		"replication_factor": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Number of copies of each partition across brokers. Higher = better fault tolerance and data safety.",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the Kafka topic",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the Kafka topic",
		},
	}
}
