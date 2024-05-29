package cmccloudv2

import (
	"fmt"

	// "strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingLBPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingLBPolicyCreate,
		Read:   resourceAutoScalingLBPolicyRead,
		Update: resourceAutoScalingLBPolicyUpdate,
		Delete: resourceAutoScalingLBPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingLBPolicyImport,
		},
		SchemaVersion: 1,
		Schema:        autoscalingLoadbalancerPolicySchema(),
	}
}
func resourceAutoScalingLBPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":                d.Get("name").(string),
		"lb_id":               d.Get("lb_id").(string),
		"lb_pool_id":          d.Get("lb_pool_id").(string),
		"lb_protocol_port":    d.Get("lb_protocol_port").(int),
		"as_configuration_id": d.Get("as_configuration_id").(string),
		"health_monitor_id":   d.Get("health_monitor_id").(string),
	}
	res, err := client.AutoScalingPolicy.CreateLBPolicy(params)

	if err != nil {
		return fmt.Errorf("Error creating lb policy: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingLBPolicyRead(d, meta)
}

func resourceAutoScalingLBPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	policy, err := client.AutoScalingPolicy.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving lb policy %s: %v", d.Id(), err)
	}

	_ = d.Set("name", policy.Name)
	_ = d.Set("lb_id", policy.Spec.Properties.Loadbalancer)
	_ = d.Set("lb_pool_id", policy.Spec.Properties.Pool.ID)
	_ = d.Set("lb_protocol_port", policy.Spec.Properties.Pool.ProtocolPort)
	_ = d.Set("health_monitor_id", policy.Spec.Properties.HealthMonitor.ID)
	return nil
}

func resourceAutoScalingLBPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	if d.HasChange("name") {
		params := map[string]interface{}{
			"name": d.Get("name").(string),
		}
		_, err := client.AutoScalingPolicy.UpdateLBPolicy(d.Id(), params) //Update(id, map[string]interface{}{"name": d.Get("name").(string)})
		if err != nil {
			return fmt.Errorf("Error when update lb policy [%s]: %v", d.Id(), err)
		}
	}
	return resourceAutoScalingLBPolicyRead(d, meta)
}

func resourceAutoScalingLBPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingPolicy.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete lb policy: %v", err)
	}
	return nil
}

func resourceAutoScalingLBPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingLBPolicyRead(d, meta)
	return []*schema.ResourceData{d}, err
}
