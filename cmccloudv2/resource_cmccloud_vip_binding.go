package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVIPBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceVIPBindingCreate,
		Read:   resourceVIPBindingRead,
		Delete: resourceVIPBindingDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVIPBindingImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        vipBindingSchema(),
	}
}

func resourceVIPBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.VIP.BindPort(d.Get("vip_id").(string), d.Get("port_id").(string))
	if err != nil {
		return fmt.Errorf("error when create vip binding: %s", err)
	}

	d.SetId(d.Get("vip_id").(string) + "_" + d.Get("port_id").(string))
	return resourceVIPBindingRead(d, meta)
}

func resourceVIPBindingRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	splitIdx := -1
	for i := 0; i < len(id); i++ {
		if id[i] == '_' {
			splitIdx = i
			break
		}
	}
	if splitIdx != -1 {
		vip_id := id[:splitIdx]
		port_id := id[splitIdx+1:]
		_ = d.Set("vip_id", vip_id)
		_ = d.Set("port_id", port_id)
	}

	return nil
}

func resourceVIPBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.VIP.UnbindPort(d.Get("vip_id").(string), d.Get("port_id").(string))

	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting vip binding %s: %v", d.Id(), err)
	}
	return nil
}

func resourceVIPBindingImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVIPBindingRead(d, meta)
	return []*schema.ResourceData{d}, err
}
