package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecurityGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        securityGroupSchema(),
		// CustomizeDiff: func(d *schema.ResourceDiff, meta interface{}) error {
		// 	rules := d.Get("rule").(*schema.Set).List()
		// 	for _, rule := range rules {
		// 		ruleMap := rule.(map[string]interface{})
		// 		// bỏ qua diff nếu chỉ là khác id
		// 		delete(ruleMap, "id")
		// 	}
		// 	return nil
		// },
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// Before creating the security group, make sure all rules are valid.
	// if err := checkRuleErrors(d, "rule"); err != nil {
	// 	return fmt.Errorf("invalid rule: %s", err)
	// }
	// If all rules are valid, proceed with creating the security group.
	group, err := client.SecurityGroup.Create(map[string]interface{}{
		"name":          d.Get("name").(string),
		"description":   d.Get("description").(string),
		"stateful":      d.Get("stateful").(bool),
		"default_rules": "false", // khong tao default rules
		// "tags":        d.Get("tags").(*schema.Set).List(),
	})
	if err != nil {
		return fmt.Errorf("error creating Security Group: %s", err)
	}
	d.SetId(group.ID)

	// get security group and delete all default rules
	sg, err := client.SecurityGroup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error receiving Security Group %s: %v", d.Id(), err)
	}
	for _, rule := range sg.Rules {
		_, err := client.SecurityGroup.DeleteRule(rule.ID)
		if err != nil {
			return fmt.Errorf("error deleting default rule %s from Security Group %s: %v", rule.ID, d.Id(), err)
		}
	}

	// Now that the security group has been created, iterate through each rule and create it
	// rawRules := d.Get("rule").(*schema.Set).List()
	// for i, rawRule := range rawRules {
	// 	rawRuleMap := rawRule.(map[string]interface{})
	// 	_, err := client.SecurityGroup.CreateRule(group.ID, map[string]interface{}{
	// 		"ether_type":      rawRuleMap["ether_type"].(string),
	// 		"direction":       rawRuleMap["direction"].(string),
	// 		"protocol":        rawRuleMap["protocol"].(string),
	// 		"port_range_min":  rawRuleMap["port_range_min"].(int),
	// 		"port_range_max":  rawRuleMap["port_range_max"].(int),
	// 		"cidr":            rawRuleMap["cidr"].(string),
	// 		"remote_group_id": rawRuleMap["remote_group_id"].(string),
	// 		"description":     rawRuleMap["description"].(string),
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("error creating Security Group Rule index %d rule: %s", (i + 1), err)
	// 	}
	// }

	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	sg, err := client.SecurityGroup.Get(id)
	if err != nil {
		return fmt.Errorf("error receiving Security Group %s: %v", d.Id(), err)
	}
	_ = d.Set("name", sg.Name)
	_ = d.Set("description", sg.Description)
	_ = d.Set("stateful", sg.Stateful)
	_ = d.Set("rule", convertSecurityGroupRules(sg.Rules))
	// fmt.Printf("error receiving Security Group %v", convertSecurityGroupRules(sg.Rules))

	gocmcapiv2.Logo("convertSecurityGroupRules(sg.Rules)", convertSecurityGroupRules(sg.Rules))
	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	// inboundJSON := "[]"
	// outboundJSON := "[]"

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("stateful") {
		_, err := client.SecurityGroup.Update(d.Id(), map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"stateful":    d.Get("stateful").(bool),
		})
		if err != nil {
			return fmt.Errorf("error updating Security Group: %s", err)
		}
	}

	if d.HasChange("rule") {
		rulesToRemove, rulesToAdd := getDiffSet(d.GetChange("rule"))

		log.Printf("[DEBUG] openstack_compute_secgroup_v2 %s rules to add: %v", d.Id(), rulesToAdd)
		log.Printf("[DEBUG] openstack_compute_secgroup_v2 %s rules to remove: %v", d.Id(), rulesToRemove)

		for _, rawRule := range rulesToRemove.List() {
			rawRuleMap := rawRule.(map[string]interface{})
			rule_id := rawRuleMap["id"].(string)
			_, err := client.SecurityGroup.DeleteRule(rule_id)
			if err != nil {
				if errors.Is(err, gocmcapiv2.ErrNotFound) {
					continue
				}

				return fmt.Errorf("error removing rule %s from security group %s: %s", rule_id, d.Id(), err)
			}
		}
		for _, rawRule := range rulesToAdd.List() {
			rawRuleMap := rawRule.(map[string]interface{})
			_, err := client.SecurityGroup.CreateRule(d.Id(), map[string]interface{}{
				"ether_type":      rawRuleMap["ether_type"].(string),
				"direction":       rawRuleMap["direction"].(string),
				"protocol":        rawRuleMap["protocol"].(string),
				"port_range_min":  rawRuleMap["port_range_min"].(int),
				"port_range_max":  rawRuleMap["port_range_max"].(int),
				"cidr":            rawRuleMap["cidr"].(string),
				"remote_group_id": rawRuleMap["remote_group_id"].(string),
				"description":     rawRuleMap["description"].(string),
			})
			if err != nil {
				return fmt.Errorf("error creating Security Group Rule %v: %v", rawRule, err)
			}
		}
	}
	return resourceSecurityGroupRead(d, meta)
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.SecurityGroup.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("error delete security group [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilSecurityGroupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete security group [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceSecurityGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceSecurityGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func convertSecurityGroupRules(rules []gocmcapiv2.SecurityGroupRule) []map[string]interface{} {
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		// ports := strconv.Itoa(rule.PortRangeMin)
		// if rule.PortRangeMax > 0 {
		// 	ports += "-" + strconv.Itoa(rule.PortRangeMax)
		// }
		remote_group_id := rule.RemoteGroupID
		if remote_group_id == nil {
			remote_group_id = ""
		}
		ether_type := rule.EtherType
		// if ether_type == nil {
		// 	ether_type = ""
		// }
		port_range_max := rule.PortRangeMax
		// if port_range_max == nil {
		// 	port_range_max = 0
		// }
		port_range_min := rule.PortRangeMin
		// if port_range_min == nil {
		// 	port_range_min = 0
		// }
		result[i] = map[string]interface{}{
			"id":              rule.ID,
			"protocol":        rule.Protocol,
			"direction":       rule.Direction,
			"cidr":            rule.CIDR,
			"remote_group_id": remote_group_id,
			"ether_type":      ether_type,
			"port_range_max":  port_range_max,
			"port_range_min":  port_range_min,
			// "ports":                 ports,
		}
	}
	return result
}

// func checkRuleErrors(d *schema.ResourceData, field string) error {
// 	rawRules := d.Get(field).(*schema.Set).List()

// 	for index, rawRule := range rawRules {
// 		rawRuleMap := rawRule.(map[string]interface{})

// 		// only one of cidr, from_group_id, or self can be set
// 		cidr := rawRuleMap["cidr"].(string)
// 		groupID := rawRuleMap["remote_group_id"].(string)
// 		port_range_min := rawRuleMap["port_range_min"].(int)
// 		port_range_max := rawRuleMap["port_range_max"].(int)
// 		errorMessage := fmt.Errorf("rule.%d: only one of cidr or remote_group_id can be set", index)

// 		// if cidr is set, from_group_id and self cannot be set
// 		if cidr != "" {
// 			if groupID != "" {
// 				return errorMessage
// 			}
// 		}

// 		// if from_group_id is set, cidr and self cannot be set
// 		if groupID != "" {
// 			if cidr != "" {
// 				return errorMessage
// 			}
// 		}

// 		if port_range_min != 0 && port_range_max != 0 && port_range_min > port_range_max {
// 			if cidr != "" {
// 				return fmt.Errorf("rule.%d: port_range_max must be >= port_range_min", index)
// 			}
// 		}
// 	}

// 	return nil
// }

func waitUntilSecurityGroupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      5 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).SecurityGroup.Get(id)
	})
}
