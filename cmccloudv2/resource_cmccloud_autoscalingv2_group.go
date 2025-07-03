package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingV2Group() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingV2GroupCreate,
		Read:   resourceAutoScalingV2GroupRead,
		Update: resourceAutoScalingV2GroupUpdate,
		Delete: resourceAutoScalingV2GroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingV2GroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingGroupV2Schema(),
	}
}

func resourceAutoScalingV2GroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	datas := map[string]interface{}{
		"name":             d.Get("name").(string),
		"zone":             d.Get("zone").(string),
		"min_size":         d.Get("min_size").(int),
		"max_size":         d.Get("max_size").(int),
		"desired_capacity": d.Get("desired_capacity").(int),
		"configuration_id": d.Get("configuration_id").(string),
		"lb_pool_id":       d.Get("lb_pool_id").(string),
		"lb_protocol_port": d.Get("lb_protocol_port").(int),
		"cooldown":         d.Get("cooldown").(int),

		"scale_up_adjustment_type": d.Get("scale_up_adjustment_type").(string),
		"scale_up_cooldown":        d.Get("scale_up_cooldown").(int),
		"scale_up_adjustment":      d.Get("scale_up_adjustment").(int),

		"scale_down_adjustment_type": d.Get("scale_down_adjustment_type").(string),
		"scale_down_cooldown":        d.Get("scale_down_cooldown").(int),
		"scale_down_adjustment":      d.Get("scale_down_adjustment").(int),
	}
	res, err := client.AutoScalingV2Group.Create(datas)

	if err != nil {
		return fmt.Errorf("Error creating autoscalingv2 group: %v", err.Error())
	}
	d.SetId(res.ID)

	_, err = waitUntilAutoScalingV2GroupStatusChangedState(d, meta, []string{"CREATE_COMPLETE", "UPDATE_COMPLETE", "COMPLETE"}, []string{"CREATE_FAILED", "UPDATE_FAILED", "FAILED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating autoscalingv2 group: %v", err.Error())
	}
	return resourceAutoScalingV2GroupRead(d, meta)
}

func resourceAutoScalingV2GroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	group, err := client.AutoScalingV2Group.Get(d.Id())
	parameters := group.Parameters
	if err != nil {
		return fmt.Errorf("Error retrieving autoscaling group %s: %v", d.Id(), err)
	}
	_ = d.Set("name", group.Name)
	_ = d.Set("zone", group.Parameters.AvailabilityZone)
	_ = d.Set("min_size", parameters.MinSize)
	_ = d.Set("max_size", parameters.MaxSize)
	_ = d.Set("desired_capacity", parameters.DesiredCapacity)
	// _ = d.Set("configuration_id", parameters.ProfileID)

	setString(d, "lb_pool_id", parameters.LbPool)
	setInt(d, "lb_protocol_port", parameters.LbMemberPort)
	// _ = d.Set("lb_pool_id", parameters.LbPoolID)
	// _ = d.Set("lb_protocol_port", parameters.LbProtocolPort)
	_ = d.Set("cooldown", parameters.Cooldown)
	_ = d.Set("scale_up_adjustment_type", parameters.ScaleUpAdjustmentType)
	_ = d.Set("scale_up_cooldown", parameters.ScaleUpCooldown)
	_ = d.Set("scale_up_adjustment", parameters.ScaleUpAdjustment)
	_ = d.Set("scale_down_adjustment_type", parameters.ScaleDownAdjustmentType)
	_ = d.Set("scale_down_cooldown", parameters.ScaleDownCooldown)
	_ = d.Set("scale_down_adjustment", AbsInt(parameters.ScaleDownAdjustment))
	_ = d.Set("created", group.CreationTime)
	_ = d.Set("status", group.Status)
	_ = d.Set("status_reason", group.StatusReason)

	return nil
}
func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
func resourceAutoScalingV2GroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	datas := map[string]interface{}{
		"zone":             d.Get("zone").(string),
		"min_size":         d.Get("min_size").(int),
		"max_size":         d.Get("max_size").(int),
		"desired_capacity": d.Get("desired_capacity").(int),
		"configuration_id": d.Get("configuration_id").(string),
		"lb_pool_id":       d.Get("lb_pool_id").(string),
		"lb_protocol_port": d.Get("lb_protocol_port").(int),
		"cooldown":         d.Get("cooldown").(int),

		"scale_up_adjustment_type": d.Get("scale_up_adjustment_type").(string),
		"scale_up_cooldown":        d.Get("scale_up_cooldown").(int),
		"scale_up_adjustment":      d.Get("scale_up_adjustment").(int),

		"scale_down_adjustment_type": d.Get("scale_down_adjustment_type").(string),
		"scale_down_cooldown":        d.Get("scale_down_cooldown").(int),
		"scale_down_adjustment":      d.Get("scale_down_adjustment").(int),
	}
	_, err := client.AutoScalingV2Group.Update(id, datas)
	if err != nil {
		return fmt.Errorf("Error when update asv2 group: %v", err)
	}
	_, err = waitUntilAutoScalingV2GroupStatusChangedState(d, meta, []string{"CREATE_COMPLETE", "UPDATE_COMPLETE", "COMPLETE"}, []string{"CREATE_FAILED", "UPDATE_FAILED", "FAILED"}, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("Error when update asv2 group: %v", err)
	}
	return resourceAutoScalingV2GroupRead(d, meta)
}

func resourceAutoScalingV2GroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// destroy the autoscale autoscalinggroup
	_, err := client.AutoScalingV2Group.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete autoscale asv2 group: %v", err)
	}
	_, err = waitUntilAutoScalingV2GroupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete autoscale asv2 group: %v", err)
	}
	return nil
}

func resourceAutoScalingV2GroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingV2GroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilAutoScalingV2GroupStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).AutoScalingV2Group.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.AutoScalingV2Group).Status
	})
}

func waitUntilAutoScalingV2GroupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).AutoScalingV2Group.Get(id)
	})
}
