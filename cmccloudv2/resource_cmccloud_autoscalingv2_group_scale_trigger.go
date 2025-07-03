package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingV2ScaleTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingV2ScaleTriggerCreate,
		Read:   resourceAutoScalingV2ScaleTriggerRead,
		Update: resourceAutoScalingV2ScaleTriggerUpdate,
		Delete: resourceAutoScalingV2ScaleTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingV2ScaleTriggerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingGroupV2ScaleTriggerSchema(),
	}
}

func resourceAutoScalingV2ScaleTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	datas := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"action":      d.Get("action").(string),
		"function":    d.Get("function").(string),
		"metric":      d.Get("metric").(string),
		"comparator":  d.Get("comparator").(string),
		"threadhold":  d.Get("threadhold").(int),
		"interval":    d.Get("interval").(int),
		"count":       d.Get("times").(int),
		"enabled":     d.Get("enabled").(bool),
	}
	res, err := client.AutoScalingV2ScaleTrigger.Create(d.Get("group_id").(string), datas)

	if err != nil {
		return fmt.Errorf("Error creating autoscalingv2 group scale trigger: %v", err.Error())
	}
	gocmcapiv2.Logo("************========= AutoScalingV2ScaleTrigger create = : ", res.ID)
	d.SetId(res.ID)

	return resourceAutoScalingV2ScaleTriggerRead(d, meta)
}

func resourceAutoScalingV2ScaleTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	trigger, err := client.AutoScalingV2ScaleTrigger.Get(d.Get("group_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving autoscaling group scale trigger%s: %v", d.Id(), err)
	}
	_ = d.Set("name", trigger.Name)
	_ = d.Set("description", trigger.Description)
	_ = d.Set("action", trigger.Action)
	_ = d.Set("function", trigger.Function)
	_ = d.Set("metric", trigger.Metric)
	_ = d.Set("comparator", trigger.Comparator)
	_ = d.Set("threadhold", trigger.Threadhold)
	_ = d.Set("interval", trigger.Interval)
	_ = d.Set("times", trigger.Count)
	_ = d.Set("enabled", trigger.Enabled)

	return nil
}

func resourceAutoScalingV2ScaleTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	datas := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"action":      d.Get("action").(string),
		"function":    d.Get("function").(string),
		"metric":      d.Get("metric").(string),
		"comparator":  d.Get("comparator").(string),
		"threadhold":  d.Get("threadhold").(int),
		"interval":    d.Get("interval").(int),
		"count":       d.Get("times").(int),
		"enabled":     d.Get("enabled").(bool),
	}
	_, err := client.AutoScalingV2ScaleTrigger.Update(d.Get("group_id").(string), id, datas)
	if err != nil {
		return fmt.Errorf("Error when update asv2 group scale trigger: %v", err)
	}
	return resourceAutoScalingV2ScaleTriggerRead(d, meta)
}

func resourceAutoScalingV2ScaleTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// destroy the autoscale autoscalinggroup
	_, err := client.AutoScalingV2ScaleTrigger.Delete(d.Get("group_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error delete autoscale asv2 group scale trigger: %v", err)
	}
	_, err = waitUntilAutoscalingV2GroupScaleTriggerDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete autoscale asv2 group scale trigger: %v", err)
	}
	return nil
}

func resourceAutoScalingV2ScaleTriggerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingV2ScaleTriggerRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilAutoscalingV2GroupScaleTriggerDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).AutoScalingV2ScaleTrigger.Get(d.Get("group_id").(string), id)
	})
}
