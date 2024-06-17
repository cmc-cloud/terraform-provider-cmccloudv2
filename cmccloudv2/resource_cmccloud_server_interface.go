package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServerInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerInterfaceCreate,
		Read:   resourceServerInterfaceRead,
		Delete: resourceServerInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServerInterfaceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        serverInterfaceSchema(),
	}
}

func resourceServerInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"subnet_id": d.Get("subnet_id").(string),
	}
	if d.Get("ip_address").(string) != "" {
		subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))

		if err != nil {
			return fmt.Errorf("Error when getting subnet info: %v", err)
		}
		_, err = isIpBelongToCidr(d.Get("ip_address").(string), subnet.Cidr)
		if err != nil {
			return err
		}
		params["ip_address"] = d.Get("ip_address").(string)
	}
	inter, err := client.NetworkInterface.Create(d.Get("server_id").(string), params)

	if err != nil {
		return fmt.Errorf("Error when create Interface of Server %s: %s", d.Get("server_id").(string), err)
	}

	d.SetId(inter.ID)
	return resourceServerInterfaceRead(d, meta)
}

func resourceServerInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	inter, err := client.NetworkInterface.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Interface %s: %v", d.Id(), err)
	}
	// _ = d.Set("server_id", inter.ServerID)
	_ = d.Set("subnet_id", inter.FixedIps[0].SubnetID)
	// chua set thi moi set
	if d.Get("ip_address").(string) == "" || d.Get("ip_address").(string) == inter.FixedIps[0].IPAddress {
		setString(d, "ip_address", inter.FixedIps[0].IPAddress)
	}
	return nil
}

func resourceServerInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server_id := d.Get("server_id").(string)
	_, err := client.NetworkInterface.Delete(d.Id(), server_id)

	if err != nil {
		return fmt.Errorf("[ERROR] Error detaching volume %s from server %s: %v", d.Id(), server_id, err)
	}
	return nil
}

func resourceServerInterfaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceServerInterfaceRead(d, meta)
	return []*schema.ResourceData{d}, err
}
