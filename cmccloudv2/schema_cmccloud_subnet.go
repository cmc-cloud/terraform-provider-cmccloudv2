package cmccloudv2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func createHostRouteElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination": {
			Type:     schema.TypeString,
			Required: true,
		},
		"nexthop": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
func createSubnetAllocationPoolsElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start": {
			Type:     schema.TypeString,
			Required: true,
		},
		"end": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
func subnetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vpc_id": {
			Type:     schema.TypeString,
			Required: true,
			// ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"ip_version": {
			Type:     schema.TypeInt,
			Optional: true,
			// ForceNew: true,
			Default: 4,
		},
		"enable_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"gateway_ip": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateIPAddress,
		},
		"allocation_pools": {
			Type:     schema.TypeList, // TypeSet => (where ordering doesn’t matter), TypeList (where ordering matters).
			Optional: true,
			Elem: &schema.Resource{
				Schema: createSubnetAllocationPoolsElementSchema(),
			},
		},
		"host_routes": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: createHostRouteElementSchema(),
			},
			Optional: true,
		},
		"dns_nameservers": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"tags": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional: true,
		},
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
			// ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
			ValidateFunc: validateIPCidrRange,
		},
	}
}
