package cmccloudv2

import (
	"fmt"
	"time"

	// "strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingScaleInPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingScaleInPolicyCreate,
		Read:   resourceAutoScalingScaleInPolicyRead,
		Update: resourceAutoScalingScaleInPolicyUpdate,
		Delete: resourceAutoScalingScaleInPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingScaleInPolicyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingScaleInOutPolicySchema(),
	}
}
func resourceAutoScalingScaleInPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":         d.Get("name").(string),
		"scale_number": d.Get("scale_number").(int),
		"cooldown":     d.Get("cooldown").(int),
		"scale_type":   d.Get("scale_type").(string),
		"scale_action": "CLUSTER_SCALE_IN",
	}
	res, err := client.AutoScalingPolicy.CreateScalePolicy(params)

	if err != nil {
		return fmt.Errorf("Error creating scale in policy: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingScaleInPolicyRead(d, meta)
}

func resourceAutoScalingScaleInPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	policy, err := client.AutoScalingPolicy.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving scale in policy %s: %v", d.Id(), err)
	}
	_ = d.Set("name", policy.Name)
	_ = d.Set("scale_number", policy.Spec.Properties.Adjustment.Number)
	_ = d.Set("cooldown", policy.Spec.Properties.Adjustment.Cooldown)
	_ = d.Set("scale_type", policy.Spec.Properties.Adjustment.Type)
	return nil
}

func resourceAutoScalingScaleInPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("name") {
		params := map[string]interface{}{
			"name": d.Get("name").(string),
		}
		_, err := client.AutoScalingPolicy.UpdateScalePolicy(d.Id(), params) //Update(id, map[string]interface{}{"name": d.Get("name").(string)})
		if err != nil {
			return fmt.Errorf("Error when update scale in policy [%s]: %v", d.Id(), err)
		}
	}
	return resourceAutoScalingScaleInPolicyRead(d, meta)
}

func resourceAutoScalingScaleInPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingPolicy.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete scale in policy: %v", err)
	}
	_, err = waitUntilAutoScalingPolicyDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete scale in policy: %v", err)
	}
	return nil
}

func resourceAutoScalingScaleInPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingScaleInPolicyRead(d, meta)
	return []*schema.ResourceData{d}, err
}
