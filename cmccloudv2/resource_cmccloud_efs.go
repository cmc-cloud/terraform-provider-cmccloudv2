package cmccloudv2

import (
	"errors"
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceEFSCreate,
		Read:   resourceEFSRead,
		Update: resourceEFSUpdate,
		Delete: resourceEFSDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEFSImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        efsSchema(),
	}
}

func resourceEFSCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	efs_id := d.Get("efs_id").(string)
	efs, err := client.EFS.Create(efs_id, map[string]interface{}{
		"name":             d.Get("name").(string),
		"ip_version":       d.Get("ip_version").(int),
		"enable_dhcp":      d.Get("enable_dhcp").(bool),
		"gateway_ip":       d.Get("gateway_ip").(string),
		"allocation_pools": d.Get("allocation_pools").([]interface{}),
		"host_routes":      d.Get("host_routes").([]interface{}),
		"dns_nameservers":  d.Get("dns_nameservers").([]interface{}),
		"tags":             d.Get("tags").(*schema.Set).List(),
		"cidr":             d.Get("cidr").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating EFS: %s", err)
	}
	d.SetId(efs.ID)
	return resourceEFSRead(d, meta)
}

func resourceEFSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	efs, err := client.EFS.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving EFS %s: %v", d.Id(), err)
	}

	_ = d.Set("id", efs.ID)
	_ = d.Set("name", efs.Name)
	_ = d.Set("efs_id", efs.VpcID)
	_ = d.Set("ip_version", efs.IpVersion)
	_ = d.Set("enable_dhcp", efs.EnableDhcp)
	_ = d.Set("gateway_ip", efs.GatewayIP)
	_ = d.Set("allocation_pools", convertAllocationPools(efs.AllocationPools))
	_ = d.Set("host_routes", convertHostRoutes(efs.HostRoutes))
	_ = d.Set("dns_nameservers", efs.DNSNameservers)
	_ = d.Set("tags", efs.Tags)
	_ = d.Set("cidr", efs.Cidr)
	return nil
}

func resourceEFSUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("cidr") || d.HasChange("efs_id") || d.HasChange("ip_version") {
		return errors.New("These fields 'cidr, efs_id, ip_version' cannot be changed after creation")
	}
	if d.HasChange("name") || d.HasChange("enable_dhcp") || d.HasChange("gateway_ip") || d.HasChange("allocation_pools") || d.HasChange("host_routes") || d.HasChange("dns_nameservers") || d.HasChange("tags") {
		_, err := client.EFS.Update(id, map[string]interface{}{
			"name":             d.Get("name").(string),
			"enable_dhcp":      d.Get("enable_dhcp").(bool),
			"gateway_ip":       d.Get("gateway_ip").(string),
			"allocation_pools": flatternAllocationPools(d.Get("allocation_pools").([]interface{})),
			"host_routes":      flatternHostRoutes(d.Get("host_routes").([]interface{})),
			"dns_nameservers":  d.Get("dns_nameservers").([]interface{}),
			"tags":             d.Get("tags").(*schema.Set).List(),
			"cidr":             d.Get("cidr").(string),
		})
		if err != nil {
			return fmt.Errorf("Error when rename EFS [%s]: %v", id, err)
		}
	}
	return resourceEFSRead(d, meta)
}

func resourceEFSDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EFS.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete efs: %v", err)
	}
	_, err = waitUntilEFSDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete efs: %v", err)
	}
	return nil
}

func resourceEFSImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEFSRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func flatternAllocationPools(pools []interface{}) []gocmcapiv2.AllocationPool {
	items := []gocmcapiv2.AllocationPool{}
	for _, pool := range pools {
		r := pool.(map[string]interface{})
		item := gocmcapiv2.AllocationPool{
			Start: r["start"].(string),
			End:   r["end"].(string),
		}
		items = append(items, item)
	}
	return items
}
func flatternHostRoutes(routes []interface{}) []gocmcapiv2.HostRoute {
	items := []gocmcapiv2.HostRoute{}
	for _, route := range routes {
		r := route.(map[string]interface{})
		item := gocmcapiv2.HostRoute{
			Destination: r["destination"].(string),
			NextHop:     r["nexthop"].(string),
		}
		items = append(items, item)
	}
	return items
}
func convertAllocationPools(pools []gocmcapiv2.AllocationPool) []map[string]interface{} {
	result := make([]map[string]interface{}, len(pools))
	for i, pool := range pools {
		result[i] = map[string]interface{}{
			"start": pool.Start,
			"end":   pool.End,
		}
	}
	return result
}
func convertHostRoutes(routes []gocmcapiv2.HostRoute) []map[string]interface{} {
	result := make([]map[string]interface{}, len(routes))
	for i, route := range routes {
		result[i] = map[string]interface{}{
			"destination": route.Destination,
			"nexthop":     route.NextHop,
		}
	}
	return result
}

func waitUntilEFSDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EFS.Get(id)
	})
}
