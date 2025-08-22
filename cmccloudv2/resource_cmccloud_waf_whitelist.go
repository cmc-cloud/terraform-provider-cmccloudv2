package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWafWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafWhitelistCreate,
		Read:   resourceWafWhitelistRead,
		Update: resourceWafWhitelistUpdate,
		Delete: resourceWafWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafWhitelistImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        wafwhitelistSchema(),
	}
}

// func getWhitelistMz(d *schema.ResourceData) string {
// 	var mz []string
// 	if d.Get("match_request_body").(bool) {
// 		mz = append(mz, "BODY")
// 	}
// 	if d.Get("match_get_arguments").(bool) {
// 		mz = append(mz, "ARGS")
// 	}
// 	if d.Get("match_http_headers").(bool) {
// 		mz = append(mz, "HEADERS")
// 	}
// 	if d.Get("match_filename").(bool) {
// 		mz = append(mz, "FILE_EXT")
// 	}
// 	if d.Get("match_url").(bool) {
// 		mz = append(mz, "URL")
// 	}
// 	if d.Get("match_name_check").(bool) {
// 		mz = append(mz, "NAME")
// 	}

// 	if d.Get("match_header_var").(string) != "" {
// 		mz = append(mz, d.Get("match_header_var").(string))
// 	}
// 	return strings.Join(mz, "|")
// }

func resourceWafWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"waf_id":           d.Get("waf_id").(string),
		"wl":               d.Get("wl").(string),
		"description":      d.Get("description").(string),
		"mz":               getRuleMz(d),
		"whitelist_type":   "BasicRule",
		"whitelist_set_id": "",
		"rmks":             "",
		"active":           true,
		"negative":         false,
		"timestamp":        time.Now(),
	}

	whitelist, err := getClient(meta).WafWhitelist.Create(params)

	if err != nil {
		return fmt.Errorf("error creating waf whitelist: %s", err)
	}
	d.SetId(whitelist.ID)
	return resourceWafWhitelistRead(d, meta)
}

func resourceWafWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	whitelist, err := getClient(meta).WafWhitelist.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving waf whitelist %s: %v", d.Id(), err)
	}

	_ = d.Set("id", whitelist.ID)
	_ = d.Set("wl", whitelist.Wl)
	_ = d.Set("description", whitelist.Description)
	_ = d.Set("waf_id", whitelist.WafID)

	setBool(d, "match_request_body", strings.Contains(whitelist.Mz, "BODY"))
	setBool(d, "match_get_arguments", strings.Contains(whitelist.Mz, "ARGS"))
	setBool(d, "match_http_headers", strings.Contains(whitelist.Mz, "HEADERS"))
	setBool(d, "match_filename", strings.Contains(whitelist.Mz, "FILE_EXT"))
	setBool(d, "match_url", strings.Contains(whitelist.Mz, "URL"))
	setBool(d, "match_name_check", strings.Contains(whitelist.Mz, "NAME"))

	if v, ok := d.GetOk("match_header_var"); ok && v.(string) != "" {
		// if v, ok := d.GetOkExists("match_header_var"); ok && v.(string) != "" {
		if strings.Contains(whitelist.Mz, "Cookie") {
			_ = d.Set("match_header_var", "Cookie")
		}
		if strings.Contains(whitelist.Mz, "Content-Type") {
			_ = d.Set("match_header_var", "Content-Type")
		}
		if strings.Contains(whitelist.Mz, "User-Agent") {
			_ = d.Set("match_header_var", "User-Agent")
		}
		if strings.Contains(whitelist.Mz, "Accept-Encoding") {
			_ = d.Set("match_header_var", "Accept-Encoding")
		}
		if strings.Contains(whitelist.Mz, "Connection") {
			_ = d.Set("match_header_var", "Connection")
		}
	}

	return nil
}

func resourceWafWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := map[string]interface{}{
		"waf_id":           d.Get("waf_id").(string),
		"wl":               d.Get("wl").(string),
		"description":      d.Get("description").(string),
		"mz":               getRuleMz(d),
		"whitelist_type":   "BasicRule",
		"whitelist_set_id": "",
		"rmks":             "",
		"active":           true,
		"negative":         false,
		"timestamp":        time.Now(),
	}
	_, err := client.WafWhitelist.Update(id, params)
	if err != nil {
		return fmt.Errorf("error when update waf whitelist [%s]: %v", id, err)
	}

	return resourceVPCRead(d, meta)
}
func resourceWafWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).WafWhitelist.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete waf whitelist: %v", err)
	}
	_, err = waitUntilWafWhitelistDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete waf whitelist: %v", err)
	}
	return nil
}

func resourceWafWhitelistImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceWafWhitelistRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilWafWhitelistDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).WafWhitelist.Get(id)
	})
}
