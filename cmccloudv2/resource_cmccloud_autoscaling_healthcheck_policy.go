package cmccloudv2

import (
	"fmt"
	"time"

	// "strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingHealthCheckPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingHealthCheckPolicyCreate,
		Read:   resourceAutoScalingHealthCheckPolicyRead,
		Update: resourceAutoScalingHealthCheckPolicyUpdate,
		Delete: resourceAutoScalingHealthCheckPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingHealthCheckPolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingHealthCheckPolicySchema(),
	}
}
func resourceAutoScalingHealthCheckPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":                  d.Get("name").(string),
		"health_check_interval": d.Get("interval").(int),
		"health_check_period":   d.Get("period").(int),
	}
	res, err := client.AutoScalingPolicy.CreateHealthCheckPolicy(params)

	if err != nil {
		return fmt.Errorf("Error creating health check policy: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingHealthCheckPolicyRead(d, meta)
}

func resourceAutoScalingHealthCheckPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	policy, err := client.AutoScalingPolicy.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving healcheck policy %s: %v", d.Id(), err)
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("interval", policy.Spec.Properties.Detection.Interval)
	_ = d.Set("period", policy.Spec.Properties.Detection.NodeUpdateTimeout)
	return nil
}

func resourceAutoScalingHealthCheckPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("name") {
		params := map[string]interface{}{
			"name": d.Get("name").(string),
		}
		_, err := client.AutoScalingPolicy.UpdateHealthCheckPolicy(d.Id(), params) //Update(id, map[string]interface{}{"name": d.Get("name").(string)})
		if err != nil {
			return fmt.Errorf("Error when update healcheck policy [%s]: %v", d.Id(), err)
		}
	}
	return resourceAutoScalingHealthCheckPolicyRead(d, meta)
}

func resourceAutoScalingHealthCheckPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingPolicy.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete healcheck policy: %v", err)
	}
	_, err = waitUntilAutoScalingPolicyDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete healcheck policy: %v", err)
	}
	return nil
}

func resourceAutoScalingHealthCheckPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingHealthCheckPolicyRead(d, meta)
	return []*schema.ResourceData{d}, err
}
