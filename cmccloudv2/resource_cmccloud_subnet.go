package cmccloudv2

import (
	"errors"
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubnetCreate,
		Read:   resourceSubnetRead,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSubnetImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        subnetSchema(),
	}
}

func resourceSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc_id := d.Get("vpc_id").(string)
	subnet, err := client.VPC.CreateSubnet(vpc_id, map[string]interface{}{
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
		return fmt.Errorf("error creating Subnet: %s", err)
	}
	d.SetId(subnet.ID)
	return resourceSubnetRead(d, meta)
}

func resourceSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Subnet %s: %v", d.Id(), err)
	}

	_ = d.Set("id", subnet.ID)
	_ = d.Set("name", subnet.Name)
	_ = d.Set("vpc_id", subnet.VpcID)
	_ = d.Set("ip_version", subnet.IpVersion)
	_ = d.Set("enable_dhcp", subnet.EnableDhcp)
	_ = d.Set("gateway_ip", subnet.GatewayIP)
	_ = d.Set("allocation_pools", convertAllocationPools(subnet.AllocationPools))
	_ = d.Set("host_routes", convertHostRoutes(subnet.HostRoutes))
	_ = d.Set("dns_nameservers", subnet.DNSNameservers)
	_ = d.Set("tags", convertTagsToSet(subnet.Tags))
	_ = d.Set("cidr", subnet.Cidr)
	return nil
}

func resourceSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("cidr") || d.HasChange("vpc_id") || d.HasChange("ip_version") {
		return errors.New("these fields 'cidr, vpc_id, ip_version' cannot be changed after creation")
	}
	if d.HasChange("name") || d.HasChange("enable_dhcp") || d.HasChange("gateway_ip") || d.HasChange("allocation_pools") || d.HasChange("host_routes") || d.HasChange("dns_nameservers") || d.HasChange("tags") {
		_, err := client.Subnet.Update(id, map[string]interface{}{
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
			return fmt.Errorf("error when rename Subnet [%s]: %v", id, err)
		}
	}
	return resourceSubnetRead(d, meta)
}

func resourceSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Subnet.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete subnet: %v", err)
	}
	_, err = waitUntilSubnetDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete subnet: %v", err)
	}
	return nil
}

func resourceSubnetImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceSubnetRead(d, meta)
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

func waitUntilSubnetDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Subnet.Get(id)
	})
}
