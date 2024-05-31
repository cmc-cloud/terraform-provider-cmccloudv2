package cmccloudv2

import (
	"fmt"
	"time"

	// "strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingScaleOutPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingScaleOutPolicyCreate,
		Read:   resourceAutoScalingScaleOutPolicyRead,
		Update: resourceAutoScalingScaleOutPolicyUpdate,
		Delete: resourceAutoScalingScaleOutPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingScaleOutPolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingScaleInOutPolicySchema(),
	}
}
func resourceAutoScalingScaleOutPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":         d.Get("name").(string),
		"scale_number": d.Get("scale_number").(int),
		"cooldown":     d.Get("cooldown").(int),
		"scale_type":   d.Get("scale_type").(string),
		"scale_action": "CLUSTER_SCALE_OUT",
	}
	res, err := client.AutoScalingPolicy.CreateScalePolicy(params)

	if err != nil {
		return fmt.Errorf("Error creating scale out policy: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingScaleOutPolicyRead(d, meta)
}

func resourceAutoScalingScaleOutPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	policy, err := client.AutoScalingPolicy.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving scale out policy %s: %v", d.Id(), err)
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("scale_number", policy.Spec.Properties.Adjustment.Number)
	_ = d.Set("cooldown", policy.Spec.Properties.Adjustment.Cooldown)
	_ = d.Set("scale_type", policy.Spec.Properties.Adjustment.Type)
	return nil
}

func resourceAutoScalingScaleOutPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("name") {
		params := map[string]interface{}{
			"name": d.Get("name").(string),
		}
		_, err := client.AutoScalingPolicy.UpdateScalePolicy(d.Id(), params)
		if err != nil {
			return fmt.Errorf("Error when update scale out policy [%s]: %v", d.Id(), err)
		}
	}
	return resourceAutoScalingScaleOutPolicyRead(d, meta)
}

func resourceAutoScalingScaleOutPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingPolicy.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete scale out policy: %v", err)
	}
	_, err = waitUntilAutoScalingPolicyDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete scale out policy: %v", err)
	}
	return nil
}

func resourceAutoScalingScaleOutPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingScaleOutPolicyRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilAutoScalingPolicyDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).AutoScalingPolicy.Get(id)
	})
}
