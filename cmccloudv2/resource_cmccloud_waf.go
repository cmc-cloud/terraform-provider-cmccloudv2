package cmccloudv2

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWaf() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCreate,
		Read:   resourceWafRead,
		Update: resourceWafUpdate,
		Delete: resourceWafDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        wafSchema(),
	}
}

func getParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"domain":               d.Get("domain").(string),
		"realserver":           d.Get("real_server").(string),
		"mode":                 d.Get("mode").(string),
		"protocol":             d.Get("protocol").(string),
		"port":                 strconv.Itoa(d.Get("port").(int)),
		"certificate_id":       d.Get("certificate_id").(string),
		"sendfile":             d.Get("send_file").(bool),
		"client_max_body_size": d.Get("client_max_body_size").(int),
		"description":          d.Get("description").(string),
	}

	if isValidIP(d.Get("real_server").(string)) {
		params["type"] = "IP"
	} else {
		params["type"] = "DOMAIN"
	}

	if d.Get("certificate_id").(string) != "" {
		params["ssl"] = true
	} else {
		params["ssl"] = false
	}
	return params
}
func resourceWafCreate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	account, err := client.Account.Get()
	if err != nil {
		return fmt.Errorf("error getting account info: %s", err)
	}
	params := getParams(d)
	params["user_id"] = account.ID
	waf, err := getClient(meta).Waf.Create(params)

	if err != nil {
		return fmt.Errorf("error creating waf: %s", err)
	}
	d.SetId(waf.ID)

	// update lb config
	params = map[string]interface{}{
		"wafId":        waf.ID,
		"lb_enable":    d.Get("load_balance_enable").(bool),
		"lb_keepalive": d.Get("load_balance_keepalive").(int),
		"lb_method":    d.Get("load_balance_method").(string),
	}
	_, err = getClient(meta).Waf.UpdateLoadBalance(d.Id(), params)
	if err != nil {
		return fmt.Errorf("error creating waf: %s", err)
	}

	return resourceWafRead(d, meta)
}

func resourceWafRead(d *schema.ResourceData, meta interface{}) error {
	waf, err := getClient(meta).Waf.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Waf %s: %v", d.Id(), err)
	}

	port, _ := strconv.Atoi(waf.Port)

	_ = d.Set("id", waf.ID)
	_ = d.Set("domain", waf.Domain)
	_ = d.Set("mode", waf.Mode)
	_ = d.Set("real_server", waf.Realserver)
	_ = d.Set("status", waf.Status)
	_ = d.Set("protocol", waf.Protocol)
	_ = d.Set("port", port)
	// _ = d.Set("certificate_id", port)

	_ = d.Set("created_at", waf.Created)
	_ = d.Set("send_file", waf.Sendfile)
	setString(d, "description", waf.Description)
	setInt(d, "client_max_body_size", waf.ClientMaxBodySize)

	return nil
}

func resourceWafUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChanges("domain", "mode", "protocol", "real_server", "port", "certificate_id", "send_file", "client_max_body_size", "description") {
		params := getParams(d)
		_, err := client.Waf.Update(id, params)
		if err != nil {
			return fmt.Errorf("error when update waf [%s]: %v", id, err)
		}
	}
	if d.HasChanges("load_balance_enable", "load_balance_keepalive", "load_balance_method") {
		params := map[string]interface{}{
			"wafId":        d.Id(),
			"lb_enable":    d.Get("load_balance_enable").(bool),
			"lb_keepalive": d.Get("load_balance_keepalive").(int),
			"lb_method":    d.Get("load_balance_method").(string),
		}
		_, err := getClient(meta).Waf.UpdateLoadBalance(d.Id(), params)
		if err != nil {
			return fmt.Errorf("error updating waf load balance: %s", err)
		}
	}
	return resourceWafRead(d, meta)
}
func resourceWafDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).Waf.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete waf: %v", err)
	}
	_, err = waitUntilWafDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete waf: %v", err)
	}
	return nil
}

func resourceWafImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceWafRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilWafDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Waf.Get(id)
	})
}
