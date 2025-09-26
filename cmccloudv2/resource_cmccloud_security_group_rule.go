package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupRuleCreate,
		Read:   resourceSecurityGroupRuleRead,
		Delete: resourceSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityGroupRuleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        securityGroupRuleWithGroupIdSchema(),
	}
}

func resourceSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// Before creating the security group, make sure all rules are valid.
	if err := checkRuleError(d); err != nil {
		return fmt.Errorf("invalid rule: %s", err)
	}
	// If all rules are valid, proceed with creating the security group.
	rule, err := client.SecurityGroup.CreateRule(d.Get("security_group_id").(string), map[string]interface{}{
		"ether_type":      d.Get("ether_type").(string),
		"direction":       d.Get("direction").(string),
		"protocol":        d.Get("protocol").(string),
		"port_range_min":  d.Get("port_range_min").(int),
		"port_range_max":  d.Get("port_range_max").(int),
		"cidr":            d.Get("cidr").(string),
		"remote_group_id": d.Get("remote_group_id").(string),
		"description":     d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("error creating Security Group Rule: %s", err)
	}
	d.SetId(rule.ID)
	return resourceSecurityGroupRuleRead(d, meta)
}

func resourceSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	rule, err := client.SecurityGroup.GetRule(id)
	if err != nil {
		return fmt.Errorf("error receiving Security Group Rule %s: %v", d.Id(), err)
	}
	_ = d.Set("ether_type", rule.EtherType)
	_ = d.Set("direction", rule.Direction)
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("port_range_max", rule.PortRangeMax)
	_ = d.Set("port_range_min", rule.PortRangeMin)
	_ = d.Set("cidr", rule.CIDR)
	_ = d.Set("remote_group_id", rule.RemoteGroupID)
	_ = d.Set("description", rule.Description)
	return nil
}

func resourceSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.SecurityGroup.DeleteRule(d.Id())
	if err != nil {
		return fmt.Errorf("error delete security group rule [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilSecurityGroupRuleDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete security group rule [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceSecurityGroupRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceSecurityGroupRuleRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func checkRuleError(d *schema.ResourceData) error {
	// only one of cidr, from_group_id, or self can be set
	cidr := d.Get("cidr").(string)
	groupID := d.Get("remote_group_id").(string)
	port_range_min := d.Get("port_range_min").(int)
	port_range_max := d.Get("port_range_max").(int)
	errorMessage := fmt.Errorf("only one of cidr or remote_group_id can be set")

	// if cidr is set, from_group_id and self cannot be set
	if cidr != "" {
		if groupID != "" {
			return errorMessage
		}
	}

	// if from_group_id is set, cidr and self cannot be set
	if groupID != "" {
		if cidr != "" {
			return errorMessage
		}
	}

	if port_range_min != 0 && port_range_max != 0 && port_range_min > port_range_max {
		if cidr != "" {
			return fmt.Errorf("port_range_max must be >= port_range_min")
		}
	}

	return nil
}

func waitUntilSecurityGroupRuleDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      5 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).SecurityGroup.GetRule(id)
	})
}
