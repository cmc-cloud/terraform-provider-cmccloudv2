package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceEIPCreate,
		Read:   resourceEIPRead,
		Update: resourceEIPUpdate,
		Delete: resourceEIPDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEIPImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        eipSchema(),
	}
}

func resourceEIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	eip, err := client.EIP.Create(map[string]interface{}{
		"project":            client.Configs.ProjectId,
		"description":        d.Get("description").(string),
		"dns_domain":         d.Get("dns_domain").(string),
		"dns_name":           d.Get("dns_name").(string),
		"tags":               d.Get("tags").(*schema.Set).List(),
		"billing_mode":       d.Get("billing_mode").(string),
		"domestic_bandwidth": d.Get("domestic_bandwidth").(int),
		"inter_bandwidth":    d.Get("inter_bandwidth").(int),
	})
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}
	d.SetId(eip.ID)

	_, err = waitUntilEIPStatusChangedState(d, meta, []string{"ACTIVE", "DOWN"}, []string{"ERROR"})
	if err != nil {
		return fmt.Errorf("Error creating EIP: %s", err)
	}
	return resourceEIPRead(d, meta)
}

func resourceEIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	eip, err := client.EIP.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving EIP %s: %v", d.Id(), err)
	}
	_ = d.Set("description", eip.Description)
	_ = d.Set("dns_domain", eip.DNSDomain)
	_ = d.Set("dns_name", eip.DNSName)
	_ = d.Set("tags", eip.Tags)
	_ = d.Set("billing_mode", eip.BillingMode)
	_ = d.Set("domestic_bandwidth", eip.DomesticBandwidthMbps)
	_ = d.Set("inter_bandwidth", eip.InterBandwidthMbps)
	// _ = d.Set("port_forwardings", convertPortForwardings(eip.PortForwardings))
	_ = d.Set("created_at", eip.CreatedAt)
	_ = d.Set("status", eip.Status)
	_ = d.Set("eip_address", eip.FloatingIPAddress)

	return nil
}

func resourceEIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("description") || d.HasChange("domestic_bandwidth") || d.HasChange("inter_bandwidth") || d.HasChange("tags") {
		_, err := client.EIP.Update(id, map[string]interface{}{
			"description":        d.Get("description").(string),
			"domestic_bandwidth": d.Get("domestic_bandwidth").(int),
			"inter_bandwidth":    d.Get("inter_bandwidth").(int),
			"tags":               d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("Error when update EIP [%s]: %v", id, err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetEIPBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when update billing mode of EIP [%s]: %v", id, err)
		}
	}

	if d.HasChange("dns_domain") {
		return fmt.Errorf("You can't not change dns_domain after eip created")
	}

	if d.HasChange("dns_name") {
		return fmt.Errorf("You can't not change dns_name after eip created")
	}

	return resourceEIPRead(d, meta)
}

func resourceEIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EIP.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete EIP: %v", err)
	}
	_, err = waitUntilEIPDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete EIP: %v", err)
	}
	return nil
}

func resourceEIPImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEIPRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilEIPDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EIP.Get(id)
	})
}

func waitUntilEIPStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 10 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EIP.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.EIP).Status
	})
}
