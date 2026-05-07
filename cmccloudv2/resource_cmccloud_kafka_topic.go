package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceKafkaTopicCreate,
		Read:   resourceKafkaTopicRead,
		Update: resourceKafkaTopicUpdate,
		Delete: resourceKafkaTopicDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKafkaTopicImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kafkaTopicSchema(),
		// Rule: (do not allow decreasing partition_count)
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			oldRaw, newRaw := diff.GetChange("partition_count")

			oldVal := oldRaw.(int)
			newVal := newRaw.(int)

			if newVal < oldVal {
				return fmt.Errorf("partition_count of kafka topic [%s] cannot be decreased", diff.Get("name").(string))
			}

			return nil
		},
	}
}

func resourceKafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.KafkaInstance.CreateTopic(d.Get("instance_id").(string), d.Get("name").(string),
		d.Get("partition_count").(int), d.Get("replication_factor").(int))
	if err != nil {
		return fmt.Errorf("error creating Kafka Topic: %s", err)
	}
	_, err = waitUntilKafkaTopicFound(d, meta)
	if err != nil {
		return fmt.Errorf("error creating Kafka Topic: %s", err)
	}
	d.SetId(d.Get("instance_id").(string) + "/" + d.Get("name").(string))

	_, err = client.KafkaInstance.UpdateTopic(d.Get("instance_id").(string), d.Get("name").(string), d.Get("partition_count").(int), d.Get("rentation_day").(int))
	if err != nil {
		return fmt.Errorf("error when update kafka topic [%s]: %v", d.Get("name").(string), err)
	}

	return resourceKafkaTopicRead(d, meta)
}

func resourceKafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error retrieving Kafka Topic %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	return nil
}

func resourceKafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("rentation_day") || d.HasChange("partition_count") {
		_, err := client.KafkaInstance.UpdateTopic(d.Get("instance_id").(string), d.Get("name").(string), d.Get("partition_count").(int), d.Get("rentation_day").(int))
		if err != nil {
			return fmt.Errorf("error when update kafka topic [%s]: %v", d.Get("name").(string), err)
		}
	}

	return resourceKafkaTopicRead(d, meta)
}

func resourceKafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.KafkaInstance.DeleteTopic(d.Get("instance_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete kafka topic: %v", err)
	}
	_, err = waitUntilKafkaTopicDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete kafka topic: %v", err)
	}
	return nil
}

func resourceKafkaTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKafkaTopicRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKafkaTopicDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"false"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.ListTopics(d.Get("instance_id").(string), map[string]string{})
	}, func(obj interface{}) string {
		topics := obj.([]gocmcapiv2.KafkaTopic)
		for _, t := range topics {
			if t.Name == d.Get("name").(string) {
				return "false"
			}
		}
		return ""
	})
}

func waitUntilKafkaTopicFound(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"false"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Get("name").(string))
	}, func(obj interface{}) string {
		topic := obj.(gocmcapiv2.KafkaTopic)
		if topic.Name == d.Get("name").(string) {
			return "true"
		}
		return ""
	})
}
