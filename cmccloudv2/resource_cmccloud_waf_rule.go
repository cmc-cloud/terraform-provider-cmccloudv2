package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWafRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleCreate,
		Read:   resourceWafRuleRead,
		Update: resourceWafRuleUpdate,
		Delete: resourceWafRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRuleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        wafruleSchema(),
	}
}

func getRuleMz(d *schema.ResourceData) string {
	var mz []string
	if d.Get("match_request_body").(bool) {
		mz = append(mz, "BODY")
	}
	if d.Get("match_get_arguments").(bool) {
		mz = append(mz, "ARGS")
	}
	if d.Get("match_http_headers").(bool) {
		mz = append(mz, "HEADERS")
	}
	if d.Get("match_filename").(bool) {
		mz = append(mz, "FILE_EXT")
	}
	if d.Get("match_url").(bool) {
		mz = append(mz, "URL")
	}
	if d.Get("match_name_check").(bool) {
		mz = append(mz, "NAME")
	}
	if d.Get("match_header_cookie").(bool) {
		mz = append(mz, "$HEADERS_VAR:Cookie")
	}
	if d.Get("match_header_content_type").(bool) {
		mz = append(mz, "$HEADERS_VAR:Content-Type")
	}
	if d.Get("match_header_user_agent").(bool) {
		mz = append(mz, "$HEADERS_VAR:User-Agent")
	}
	if d.Get("match_header_accept_encoding").(bool) {
		mz = append(mz, "$HEADERS_VAR:Accept-Encoding")
	}
	if d.Get("match_header_connection").(bool) {
		mz = append(mz, "$HEADERS_VAR:Connection")
	}
	return strings.Join(mz, "|")
}
func resourceWafRuleCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"waf_id":      d.Get("waf_id").(string),
		"msg":         d.Get("message").(string),
		"detection":   d.Get("detection").(string),
		"score":       d.Get("action").(string),
		"description": d.Get("description").(string),
		"mz":          getRuleMz(d),
		"rule_type":   "BasicRule",
		"ruleset_id":  "",
		"rmks":        "",
		"active":      true,
		"negative":    false,
		"timestamp":   time.Now().Unix(),
	}

	rule, err := getClient(meta).WafRule.Create(params)

	if err != nil {
		return fmt.Errorf("error creating waf rule: %s", err)
	}
	d.SetId(rule.ID)
	return resourceWafRuleRead(d, meta)
}

func resourceWafRuleRead(d *schema.ResourceData, meta interface{}) error {
	rule, err := getClient(meta).WafRule.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving waf rule %s: %v", d.Id(), err)
	}

	_ = d.Set("id", rule.ID)
	_ = d.Set("message", rule.Msg)
	_ = d.Set("detection", rule.Detection)
	_ = d.Set("description", rule.Description)
	_ = d.Set("action", rule.Score)
	_ = d.Set("waf_id", rule.WafID)

	setBool(d, "match_request_body", strings.Contains(rule.Mz, "BODY"))
	setBool(d, "match_get_arguments", strings.Contains(rule.Mz, "ARGS"))
	setBool(d, "match_http_headers", strings.Contains(rule.Mz, "HEADERS"))
	setBool(d, "match_filename", strings.Contains(rule.Mz, "FILE_EXT"))
	setBool(d, "match_url", strings.Contains(rule.Mz, "URL"))
	setBool(d, "match_name_check", strings.Contains(rule.Mz, "NAME"))
	setBool(d, "match_header_cookie", strings.Contains(rule.Mz, "Cookie"))
	setBool(d, "match_header_content_type", strings.Contains(rule.Mz, "Content-Type"))
	setBool(d, "match_header_user_agent", strings.Contains(rule.Mz, "User-Agent"))
	setBool(d, "match_header_accept_encoding", strings.Contains(rule.Mz, "Accept-Encoding"))
	setBool(d, "match_header_connection", strings.Contains(rule.Mz, "Connection"))

	return nil
}

func resourceWafRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	rule, err := getClient(meta).WafRule.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving waf rule %s: %v", d.Id(), err)
	}

	params := map[string]interface{}{
		"id":          id,
		"mz":          getRuleMz(d),
		"waf_id":      d.Get("waf_id").(string),
		"msg":         d.Get("message").(string),
		"detection":   d.Get("detection").(string),
		"score":       d.Get("action").(string),
		"description": d.Get("description").(string),
		"rule_type":   "BasicRule",
		"ruleset_id":  "",
		"rmks":        "",
		"active":      true,
		"negative":    false,
		"timestamp":   rule.Timestamp,
	}
	_, err = client.WafRule.Update(id, params)
	if err != nil {
		return fmt.Errorf("error when update waf rule [%s]: %v", id, err)
	}

	return resourceWafRuleRead(d, meta)
}
func resourceWafRuleDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).WafRule.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete waf rule: %v", err)
	}
	_, err = waitUntilWafRuleDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete waf rule: %v", err)
	}
	return nil
}

func resourceWafRuleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceWafRuleRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilWafRuleDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).WafRule.Get(id)
	})
}
