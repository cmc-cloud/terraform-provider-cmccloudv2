package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	// "strconv"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingGroupCreate,
		Read:   resourceAutoScalingGroupRead,
		Update: resourceAutoScalingGroupUpdate,
		Delete: resourceAutoScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingGroupSchema(),
	}
}

func resourceAutoScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	fmt.Errorf("resourceAutoScalingGroupCreate %s\n", d.Get("name").(string))
	datas := map[string]interface{}{
		"name":                d.Get("name").(string),
		"min_size":            0, //d.Get("min_size").(int),
		"max_size":            d.Get("max_size").(int),
		"desired_capacity":    0, //d.Get("desired_capacity").(int),
		"as_configuration_id": d.Get("as_configuration_id").(string),
	}
	res, err := client.AutoScalingGroup.Create(datas)

	if err != nil {
		return fmt.Errorf("Error creating autoscaling group: %v", err.Error())
	}
	gocmcapiv2.Logs("set id " + res.ID)
	d.SetId(res.ID)

	waitUntilAsGroupChangeState(d, meta, res.ID, []string{"INIT"}, []string{"ACTIVE", "WARNING"})

	// attach policies
	policies, ok := d.GetOk("policies")
	if ok {
		for _, policy_id := range policies.(*schema.Set).List() {
			action, _ := client.AutoScalingPolicy.AttachToASGroup(policy_id.(string), res.ID)
			waitUntilAsActionFinished(d, meta, action.ActionID)
		}
	}

	_, err = client.AutoScalingGroup.UpdateCapacity(res.ID, map[string]interface{}{
		"min_size":         d.Get("min_size").(int),
		"max_size":         d.Get("max_size").(int),
		"desired_capacity": d.Get("desired_capacity").(int),
		"strict":           true,
	})
	if err != nil {
		return fmt.Errorf("Error when update autoscaling group capacity [%s]: %v", res.ID, err)
	}

	return resourceAutoScalingGroupRead(d, meta)
}

func resourceAutoScalingGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	autoscalinggroup, err := client.AutoScalingGroup.Get(d.Id())
	if err != nil {
		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			log.Printf("[WARN] CMC Cloud AutoScalingGroup with id = (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving autoscalinggroup %s: %v", d.Id(), err)
	}
	// _ = d.Set("project_id", autoscalinggroup.ProjectID)
	_ = d.Set("name", autoscalinggroup.Name)
	_ = d.Set("min_size", autoscalinggroup.MinSize)
	_ = d.Set("max_size", autoscalinggroup.MaxSize)
	_ = d.Set("desired_capacity", autoscalinggroup.DesiredCapacity)
	_ = d.Set("as_configuration_id", autoscalinggroup.ProfileID)
	_ = d.Set("nodes", autoscalinggroup.Nodes)
	_ = d.Set("policies", autoscalinggroup.Policies)
	return nil
}

// func genPolicy(params map[string]interface{}) []map[string]interface{} {
// 	values := make([]map[string]interface{}, 1)
// 	values[0] = params
// 	return values
// }

func resourceAutoScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("policies") {
		removed, added := getDiffSet(d.GetChange("policies"))
		for _, remove_id := range removed.List() {
			// Logic xử lý phần tử bị xóa
			log.Printf("Detach policy id [%s]", remove_id)
			res, _ := client.AutoScalingPolicy.DetachFromASGroup(remove_id.(string), d.Id())
			waitUntilAsActionFinished(d, meta, res.ActionID)
		}
		for _, add_id := range added.List() {
			// Logic xử lý phần tử add them
			log.Printf("Attach policy id [%s]", add_id)
			res, _ := client.AutoScalingPolicy.AttachToASGroup(add_id.(string), d.Id())
			waitUntilAsActionFinished(d, meta, res.ActionID)
		}
	}

	if d.HasChange("name") || d.HasChange("as_configuration_id") {
		_, err := client.AutoScalingGroup.Update(id, map[string]interface{}{
			"name":                d.Get("name").(string),
			"as_configuration_id": d.Get("as_configuration_id").(string),
		})
		if err != nil {
			return fmt.Errorf("Error when update autoscalinggroup [%s]: %v", id, err)
		}
	}
	if d.HasChange("min_size") || d.HasChange("max_size") || d.HasChange("desired_capacity") {
		_, err := client.AutoScalingGroup.UpdateCapacity(id, map[string]interface{}{
			"min_size":         d.Get("min_size").(int),
			"max_size":         d.Get("max_size").(int),
			"desired_capacity": d.Get("desired_capacity").(int),
			"strict":           true,
		})
		if err != nil {
			return fmt.Errorf("Error when update autoscalinggroup capacity [%s]: %v", id, err)
		}
	}
	return resourceAutoScalingGroupRead(d, meta)
}

func resourceAutoScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// destroy the autoscale autoscalinggroup
	_, err := client.AutoScalingGroup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete autoscale autoscalinggroup: %v", err)
	}
	return nil
}

func resourceAutoScalingGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilAsGroupChangeState(d *schema.ResourceData, meta interface{}, id string, pendingStatus []string, targetStatus []string) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with id (%s) to be "+strings.Join(targetStatus, ","), id)
	stateConf := &resource.StateChangeConf{
		Pending:        pendingStatus,
		Target:         targetStatus,
		Refresh:        asGroupStateRefreshfunc(d, meta, id),
		Timeout:        d.Timeout(schema.TimeoutCreate),
		Delay:          20 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func asGroupStateRefreshfunc(d *schema.ResourceData, meta interface{}, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*CombinedConfig).goCMCClient()
		server, err := client.AutoScalingGroup.Get(d.Id())
		if err != nil {
			fmt.Errorf("Error retrieving AS Group %s: %v", id, err)
			return nil, "", err
		}
		return server, server.Status, nil
	}
}

func waitUntilAsActionFinished(d *schema.ResourceData, meta interface{}, action_id string) (interface{}, error) {
	log.Printf("[INFO] Waiting for action with id (%s) to be finished", action_id)
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"RUNNING", "WAITING", "SUSPENDED", "READY"},
		Target:         []string{"FAILED", "SUCCEEDED", "CANCELLED"},
		Refresh:        waitUntilAsActionRefreshfunc(d, meta, action_id),
		Timeout:        10 * time.Minute,
		Delay:          3 * time.Second,
		MinTimeout:     2 * time.Second,
		NotFoundChecks: 10,
	}
	return stateConf.WaitForState()
}

func waitUntilAsActionRefreshfunc(d *schema.ResourceData, meta interface{}, action_id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*CombinedConfig).goCMCClient()
		server, err := client.AutoScalingGroup.GetAction(action_id)
		if err != nil {
			fmt.Errorf("Error retrieving AS Group Action %s: %v", action_id, err)
			return nil, "", err
		}
		return server, server.Status, nil
	}
}
