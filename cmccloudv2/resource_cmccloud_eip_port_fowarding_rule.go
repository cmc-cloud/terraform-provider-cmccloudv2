package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEIPPortForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPPortForwardingRuleCreate,
		Read:   resourceEIPPortForwardingRuleRead,
		Update: resourceEIPPortForwardingRuleUpdate,
		Delete: resourceEIPPortForwardingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEIPPortForwardingRuleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        createEipPortForwardingRuleElementSchema(),
	}
}

func resourceEIPPortForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	rule, err := client.EIP.CreatePortForwardingRule(d.Get("eip_id").(string), map[string]interface{}{
		"protocol":            d.Get("protocol").(string),
		"internal_ip_address": d.Get("internal_ip_address").(string),
		"internal_port_id":    d.Get("internal_port_id").(string),
		"internal_port":       d.Get("internal_port").(int),
		"external_port":       d.Get("external_port").(int),
		"internal_port_range": d.Get("internal_port_range").(string),
		"external_port_range": d.Get("external_port_range").(string),
		"description":         d.Get("description").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating EIP PortForwarding rule: %s", err)
	}
	d.SetId(rule.ID)

	return resourceEIPPortForwardingRuleRead(d, meta)
}

func resourceEIPPortForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	rule, err := client.EIP.GetPortForwardingRule(d.Get("eip_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving EIP Port Forwarding Rule %s: %v", d.Id(), err)
	}
	_ = d.Set("protocol", rule.Protocol)
	_ = d.Set("internal_ip_address", rule.InternalIPAddress)
	_ = d.Set("internal_port_id", rule.InternalPortID)
	_ = d.Set("internal_port", rule.InternalPort)
	_ = d.Set("external_port", rule.ExternalPort)
	_ = d.Set("internal_port_range", rule.InternalPortRange)
	_ = d.Set("external_port_range", rule.ExternalPortRange)
	_ = d.Set("description", rule.Description)

	return nil
}

func resourceEIPPortForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	params := map[string]interface{}{
		"protocol":            d.Get("protocol").(string),
		"internal_ip_address": d.Get("internal_ip_address").(string),
		"internal_port_id":    d.Get("internal_port_id").(string),
		"internal_port":       d.Get("internal_port").(int),
		"external_port":       d.Get("external_port").(int),
		"internal_port_range": d.Get("internal_port_range").(string),
		"external_port_range": d.Get("external_port_range").(string),
		"description":         d.Get("description").(string),
	}
	_, err := client.EIP.UpdatePortForwardingRule(d.Get("eip_id").(string), d.Id(), params)

	if err != nil {
		return fmt.Errorf("Error when update EIP Port Forwarding Rule [%s]: %v", id, err)
	}

	return resourceEIPPortForwardingRuleRead(d, meta)
}

func resourceEIPPortForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EIP.DeletePortForwardingRule(d.Get("eip_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error delete EIP Port Forwarding Rule: %v", err)
	}
	_, err = waitUntilEIPPortForwardingRuleDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete EIP Port Forwarding Rule [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceEIPPortForwardingRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEIPPortForwardingRuleRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilEIPPortForwardingRuleDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EIP.GetPortForwardingRule(d.Get("eip_id").(string), id)
	})
}
