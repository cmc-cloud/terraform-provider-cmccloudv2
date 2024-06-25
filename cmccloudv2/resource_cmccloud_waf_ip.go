package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWafIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafIPCreate,
		Read:   resourceWafIPRead,
		Update: resourceWafIPUpdate,
		Delete: resourceWafIPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafIPImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        wafipSchema(),
	}
}

func resourceWafIPCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"waf_id":      d.Get("waf_id").(string),
		"type":        d.Get("type").(string),
		"value":       d.Get("ip").(string),
		"format":      "IP",
		"description": d.Get("description").(string),
	}

	ip, err := getClient(meta).WafIP.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating waf ip: %s", err)
	}
	d.SetId(ip.ID)
	return resourceWafIPRead(d, meta)
}

func resourceWafIPRead(d *schema.ResourceData, meta interface{}) error {
	ip, err := getClient(meta).WafIP.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving WafIP %s: %v", d.Id(), err)
	}

	_ = d.Set("id", ip.ID)
	_ = d.Set("type", ip.Type)
	_ = d.Set("ip", ip.Value)
	_ = d.Set("description", ip.Description)

	return nil
}

func resourceWafIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := map[string]interface{}{
		"type":        d.Get("type").(string),
		"value":       d.Get("ip").(string),
		"format":      "IP",
		"description": d.Get("description").(string),
	}
	_, err := client.WafIP.Update(id, params)
	if err != nil {
		return fmt.Errorf("Error when update waf ip [%s]: %v", id, err)
	}

	return resourceWafIPRead(d, meta)
}
func resourceWafIPDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).WafIP.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete waf ip: %v", err)
	}
	_, err = waitUntilWafIPDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete waf ip: %v", err)
	}
	return nil
}

func resourceWafIPImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceWafIPRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilWafIPDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).WafIP.Get(id)
	})
}
