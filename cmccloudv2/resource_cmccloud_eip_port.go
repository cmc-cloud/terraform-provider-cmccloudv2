package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEIPPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPPortCreate,
		Read:   resourceEIPPortRead,
		Delete: resourceEIPPortDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEIPPortImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        eipportSchema(),
	}
}

func resourceEIPPortCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EIP.AttachPort(d.Get("eip_id").(string), d.Get("port_id").(string), d.Get("fix_ip_address").(string))
	if err != nil {
		return fmt.Errorf("error when attach eip %s to port %s: %s", d.Get("eip_id").(string), d.Get("port_id").(string), err)
	}

	d.SetId(d.Get("eip_id").(string))
	// wait cho den khi attached
	_, err = waitUntilEIPPortAttachedStateChanged(d, meta, []string{"Attached"}, []string{})
	if err != nil {
		return fmt.Errorf("error when attach eip %s to port %s: %s", d.Get("eip_id").(string), d.Get("port_id").(string), err)
	}
	return resourceEIPPortRead(d, meta)
}

func resourceEIPPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	eip, err := client.EIP.Get(d.Get("eip_id").(string))
	if err != nil {
		return fmt.Errorf("error retrieving eip detail %s: %v", d.Id(), err)
	}
	_ = d.Set("eip_id", eip.ID)
	_ = d.Set("port_id", eip.PortID)
	_ = d.Set("fixed_ip_address", eip.FixedIPAddress)
	return nil
}

func resourceEIPPortDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EIP.DetachPort(d.Id())

	if err != nil {
		return fmt.Errorf("[ERROR] Error detaching port from eip %s: %v", d.Id(), err)
	}
	// wait until detached
	_, err = waitUntilEIPPortAttachedStateChanged(d, meta, []string{"Detached"}, []string{})
	if err != nil {
		return fmt.Errorf("[ERROR] Error detaching port from eip %s: %v", d.Id(), err)
	}
	return nil
}

func resourceEIPPortImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEIPPortRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilEIPPortAttachedStateChanged(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      5 * time.Second,
		MinTimeout: 10 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EIP.Get(id)
	}, func(obj interface{}) string {
		eip := obj.(gocmcapiv2.EIP)
		if eip.PortID != "" {
			return "Attached"
		}
		return "Detached"
	})
}
