package cmccloudv2

import (
	"fmt"
	"strings"
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
	}
}

func resourceKafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{}
	instance, err := client.KafkaInstance.CreateTopic(d.Get("instance_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating Kafka Topic: %s", err)
	}
	d.SetId(instance.ID)
	_, err = waitUntilKafkaTopicJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating Kafka Topic: %s", err)
	}

	return resourceKafkaTopicRead(d, meta)
}

func resourceKafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Kafka Topic %s: %v", d.Id(), err)
	}
	_ = d.Set("id", instance.ID)
	_ = d.Set("name", instance.Name)
	// _ = d.Set("status", instance.Status)
	// _ = d.Set("created_at", instance.Created)
	return nil
}

func resourceKafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*CombinedConfig).goCMCClient()
	// id := d.Id()

	return resourceKafkaTopicRead(d, meta)
}

func resourceKafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.KafkaInstance.DeleteTopic(d.Get("instance_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete kafka database instance: %v", err)
	}
	_, err = waitUntilKafkaTopicDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete kafka database instance: %v", err)
	}
	return nil
}

func resourceKafkaTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKafkaTopicRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKafkaTopicJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING", "ACTIVE", "READY", "STABLE"}, []string{"ERROR", "SHUTDOWN", "FAILURE"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Id())
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.KafkaTopic).Status)
	})
}

func waitUntilKafkaTopicDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Id())
	})
}

// func waitUntilKafkaTopicAttachFinished(d *schema.ResourceData, meta interface{}, securityGroupId string) (interface{}, error) {
// 	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
// 		Timeout:    40 * time.Second,
// 		Delay:      5 * time.Second,
// 		MinTimeout: 5 * time.Second,
// 	}, func(id string) (any, error) {
// 		return getClient(meta).KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Id())
// 	}, func(obj interface{}) string {
// 		// instance := obj.(gocmcapiv2.KafkaTopic)
// 		// if strings.Contains(instance.SecurityClientIds, securityGroupId) {
// 		// 	return "true"
// 		// }
// 		return ""
// 	})
// }
// func waitUntilKafkaTopicDetachFinished(d *schema.ResourceData, meta interface{}, securityGroupId string) (interface{}, error) {
// 	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"error"}, WaitConf{
// 		Timeout:    40 * time.Second,
// 		Delay:      5 * time.Second,
// 		MinTimeout: 5 * time.Second,
// 	}, func(id string) (any, error) {
// 		return getClient(meta).KafkaInstance.GetTopic(d.Get("instance_id").(string), d.Id())
// 	}, func(obj interface{}) string {
// 		// instance := obj.(gocmcapiv2.KafkaTopic)
// 		// if !strings.Contains(instance.SecurityClientIds, securityGroupId) {
// 		// 	return "true"
// 		// }
// 		return ""
// 	})
// }
