package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceVIPCreate,
		Read:   resourceVIPRead,
		Delete: resourceVIPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVIPImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        vipSchema(),
	}
}

func resourceVIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vip, err := client.VIP.Create(map[string]interface{}{
		// "name":       d.Get("name").(string),
		"subnet_id": d.Get("subnet_id").(string),
		"ip":        d.Get("ip_address").(string),
	})
	if err != nil {
		return fmt.Errorf("error when create vip: %s", err)
	}

	d.SetId(vip.ID)
	return resourceVIPRead(d, meta)
}

func resourceVIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vip, err := client.VIP.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving vip detail %s: %v", d.Id(), err)
	}
	_ = d.Set("subnet_id", vip.FixedIps[0].SubnetID)
	// setString(d, "ip_address", vip.FixedIps[0].IPAddress)
	_ = d.Set("ip_address", vip.FixedIps[0].IPAddress)
	_ = d.Set("name", vip.Name)
	_ = d.Set("status", vip.Status)
	_ = d.Set("created_at", vip.CreatedAt)
	return nil
}

func resourceVIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.VIP.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting vip %s: %v", d.Id(), err)
	}
	return nil
}

func resourceVIPImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVIPRead(d, meta)
	return []*schema.ResourceData{d}, err
}
