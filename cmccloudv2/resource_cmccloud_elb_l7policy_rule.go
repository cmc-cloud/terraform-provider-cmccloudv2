package cmccloudv2

import (
	"fmt"

	// "strconv"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBL7policyRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBL7policyRuleCreate,
		Read:   resourceELBL7policyRuleRead,
		Update: resourceELBL7policyRuleUpdate,
		Delete: resourceELBL7policyRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBL7policyRuleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        lbL7policyRuleSchema(),
		CustomizeDiff: validateCompareTypeWithType(),
	}
}

func validateCompareTypeWithType() schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		t, ok := d.GetOk("type")
		c, ok2 := d.GetOk("compare_type")

		if !ok || !ok2 {
			return nil
		}
		typ := t.(string)
		compare := c.(string)

		if typ == "FILE_TYPE" {
			if compare != "EQUAL_TO" && compare != "REGEX" {
				return fmt.Errorf("compare_type must be REGEX or EQUAL_TO when type is FILE_TYPE")
			}
		}
		return nil
	}
}
func resourceELBL7policyRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datas := map[string]interface{}{
		"type":         d.Get("type").(string),
		"compare_type": d.Get("compare_type").(string),
		"value":        d.Get("value").(string),
		"key":          d.Get("key").(string),
		"invert":       d.Get("invert").(bool),
	}
	rule, err := client.ELB.CreateL7PolicyRule(d.Get("l7policy_id").(string), datas)
	if err != nil {
		return fmt.Errorf("error creating L7Policy Rule: %s", err)
	}

	_, err = waitUntilELBL7PolicyRuleStatusChangedState(rule.ID, d, meta, []string{"ACTIVE"}, []string{"ERROR", "DELETED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ELB L7 Policy Rule: %s", err)
	}
	d.SetId(rule.ID)
	return resourceELBL7policyRuleRead(d, meta)
}

func resourceELBL7policyRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	rule, err := client.ELB.GetL7PolicyRule(d.Get("l7policy_id").(string), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading L7Policy Rule: %s", err)
	}

	_ = d.Set("type", rule.Type)
	_ = d.Set("compare_type", rule.CompareType)
	_ = d.Set("value", rule.Value)
	_ = d.Set("key", rule.Key)
	_ = d.Set("invert", rule.Invert)
	_ = d.Set("provisioning_status", rule.ProvisioningStatus)
	_ = d.Set("operating_status", rule.OperatingStatus)
	_ = d.Set("created", rule.CreatedAt)

	return nil
}

func resourceELBL7policyRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datas := map[string]interface{}{
		"type":         d.Get("type").(string),
		"compare_type": d.Get("compare_type").(string),
		"value":        d.Get("value").(string),
		"key":          d.Get("key").(string),
		"invert":       d.Get("invert").(bool),
	}
	_, err := client.ELB.UpdateL7PolicyRule(d.Get("l7policy_id").(string), d.Id(), datas)
	if err != nil {
		return fmt.Errorf("error updating L7Policy Rule: %s", err)
	}
	_, err = waitUntilELBL7PolicyRuleStatusChangedState(d.Id(), d, meta, []string{"ACTIVE"}, []string{"ERROR", "DELETED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error updating ELB L7 Policy Rule: %s", err)
	}
	return resourceELBL7policyRuleUpdate(d, meta)
}

func resourceELBL7policyRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ELB.DeleteL7PolicyRule(d.Get("l7policy_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete L7 Policy Rule: %v", err)
	}
	_, err = waitUntilEELBL7PolicyRuleDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete L7 Policy Rule: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceELBL7policyRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceELBL7policyRuleRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func waitUntilEELBL7PolicyRuleDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetL7PolicyRule(d.Get("l7policy_id").(string), id)
	})
}
func waitUntilELBL7PolicyRuleStatusChangedState(resource_id string, d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceIdStatusChanged(resource_id, d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetL7PolicyRule(d.Get("l7policy_id").(string), id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.L7PolicyRule).ProvisioningStatus
	})
}
