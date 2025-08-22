package cmccloudv2

import (
	"fmt"
	"time"

	// "strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingDeletePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingDeletePolicyCreate,
		Read:   resourceAutoScalingDeletePolicyRead,
		Update: resourceAutoScalingDeletePolicyUpdate,
		Delete: resourceAutoScalingDeletePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingDeletePolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingDeletePolicySchema(),
	}
}
func resourceAutoScalingDeletePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"criteria":                d.Get("criteria").(string),
		"grace_period":            d.Get("grace_period").(int),
		"destroy_after_deletion":  d.Get("destroy_after_deletion").(bool),
		"reduce_desired_capacity": d.Get("reduce_desired_capacity").(bool),
		"lifecycle_timeout":       d.Get("lifecycle_timeout").(int),
		"lifecycle_hook_url":      d.Get("lifecycle_hook_url").(string),
	}
	res, err := client.AutoScalingPolicy.CreateDeletePolicy(params)

	if err != nil {
		return fmt.Errorf("error creating delete policy: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingDeletePolicyRead(d, meta)
}

func resourceAutoScalingDeletePolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	policy, err := client.AutoScalingPolicy.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving delete policy %s: %v", d.Id(), err)
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("criteria", policy.Spec.Properties.Criteria)
	_ = d.Set("grace_period", policy.Spec.Properties.GracePeriod)
	_ = d.Set("destroy_after_deletion", policy.Spec.Properties.DestroyAfterDeletion)
	_ = d.Set("reduce_desired_capacity", policy.Spec.Properties.ReduceDesiredCapacity)
	_ = d.Set("lifecycle_timeout", policy.Spec.Properties.Hooks.Timeout)
	_ = d.Set("lifecycle_hook_url", policy.Spec.Properties.Hooks.Params.URL)

	return nil
}

func resourceAutoScalingDeletePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("name") {
		params := map[string]interface{}{
			"name": d.Get("name").(string),
		}
		_, err := client.AutoScalingPolicy.UpdateDeletePolicy(d.Id(), params) //Update(id, map[string]interface{}{"name": d.Get("name").(string)})
		if err != nil {
			return fmt.Errorf("error when update delete policy [%s]: %v", d.Id(), err)
		}
	}
	return resourceAutoScalingDeletePolicyRead(d, meta)
}

func resourceAutoScalingDeletePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingPolicy.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete the delete policy: %v", err)
	}
	_, err = waitUntilAutoScalingPolicyDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete the delete policy: %v", err)
	}
	return nil
}

func resourceAutoScalingDeletePolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingDeletePolicyRead(d, meta)
	return []*schema.ResourceData{d}, err
}
