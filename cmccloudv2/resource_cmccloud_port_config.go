package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePortConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortConfigCreate,
		Read:   resourcePortConfigRead,
		Update: resourcePortConfigUpdate,
		Delete: resourcePortConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePortConfigImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        portConfigSchema(),
	}
}

func getPortConfigParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{}

	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	// bool cần xử lý riêng
	//nolint:staticcheck // Need tri-state bool support
	if v, ok := d.GetOkExists("port_security_enabled"); ok {
		params["port_security_enabled"] = v.(bool)
	}

	//nolint:staticcheck // Need tri-state bool support
	if v, ok := d.GetOkExists("security_group_ids"); ok {
		params["security_groups"] = v.(*schema.Set).List()
	}

	//nolint:staticcheck // Need tri-state bool support
	if v, ok := d.GetOkExists("allowed_address_pairs"); ok {
		params["allowed_address_pairs"] = v.(*schema.Set).List()
	}

	return params
}
func resourcePortConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := getPortConfigParams(d)
	if len(params) == 0 {
		return nil
	}
	_, err := client.Port.Patch(d.Get("port_id").(string), params)
	if err != nil {
		return fmt.Errorf("error updating Port: %s", err)
	}
	d.SetId("port_config_" + d.Get("port_id").(string))

	return resourcePortConfigRead(d, meta)
}

func resourcePortConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := getPortConfigParams(d)
	if len(params) == 0 {
		return nil
	}
	_, err := client.Port.Patch(d.Get("port_id").(string), params)
	if err != nil {
		return fmt.Errorf("error when update port config [%s]: %v", id, err)
	}
	return resourcePortConfigRead(d, meta)
}
func resourcePortConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	port, err := client.Port.Get(d.Get("port_id").(string))
	if err != nil {
		return fmt.Errorf("error retrieving Port %s: %v", d.Get("port_id").(string), err)
	}
	_ = d.Set("name", port.Name)

	//nolint:staticcheck // Need tri-state bool support
	if _, ok := d.GetOkExists("security_group_ids"); ok {
		setTypeSet(d, "security_group_ids", port.SecurityGroups)
	}

	//nolint:staticcheck // Need tri-state bool support
	if _, ok := d.GetOkExists("port_security_enabled"); ok {
		_ = d.Set("port_security_enabled", port.PortSecurityEnabled)
	}

	//nolint:staticcheck // Need tri-state bool support
	if _, ok := d.GetOkExists("allowed_address_pairs"); ok {
		pairs := make([]map[string]interface{}, 0)
		for _, p := range port.AllowedAddressPairs {
			pair := map[string]interface{}{
				"ip_address": p.IPAddress,
			}
			if p.MacAddress != "" {
				pair["mac_address"] = p.MacAddress
			}
			pairs = append(pairs, pair)
		}
		_ = d.Set("allowed_address_pairs", pairs)
	}

	return nil
}

func resourcePortConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// khong lam gi ca
	// resource này chỉ để update config của port, ko xóa port
	return nil
}

func resourcePortConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePortConfigRead(d, meta)
	return []*schema.ResourceData{d}, err
}
