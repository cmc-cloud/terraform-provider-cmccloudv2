package cmccloudv2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func createHostRouteElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"destination": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The destination IP address or CIDR block used for the host route. Changing this updates the destination of the existing host route.",
		},
		"nexthop": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The next hop IP address used by the router for forwarding traffic to the destination. Changing this updates the next hop of the existing host route.",
		},
	}
}
func createSubnetAllocationPoolsElementSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The start IP address of the IP address range for the allocation pool. Changing this updates the start of the existing allocation pool.",
		},
		"end": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The end IP address of the IP address range for the allocation pool. Changing this updates the end of the existing allocation pool.",
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
			Description:  "ID of VPC",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the subnet. Changing this updates the name of the existing subnet.",
		},
		"ip_version": {
			Type:     schema.TypeInt,
			Optional: true,
			// ForceNew: true,
			Default:     4,
			Description: "IP version, either 4 (default) or 6",
		},
		"enable_dhcp": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "The administrative state of the network. Changing this value enables or disables the DHCP capabilities of the existing subnet. Defaults to true",
		},
		"gateway_ip": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateIPAddress,
			Description:  "The IP address of the gateway. Changing this updates the gateway IP of the existing subnet.",
		},
		"allocation_pools": {
			Type:     schema.TypeList, // TypeSet => (where ordering doesn’t matter), TypeList (where ordering matters).
			Optional: true,
			Elem: &schema.Resource{
				Schema: createSubnetAllocationPoolsElementSchema(),
			},
			Description: "A block declaring the start and end range of the IP addresses available for use with DHCP in this subnet. Multiple allocation_pool blocks can be declared, providing the subnet with more than one range of IP addresses to use with DHCP. However, each IP range must be from the same CIDR that the subnet is part of",
		},
		"host_routes": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: createHostRouteElementSchema(),
			},
			Optional:    true,
			Description: "A list of host routes to be injected into the metadata of the server. Multiple host_routes blocks can be declared",
		},
		"dns_nameservers": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Optional:    true,
			Description: "A list of DNS name server names used by hosts in this subnet. Changing this updates the DNS name servers for the existing subnet.",
		},
		"tags": tagSchema(),
		"cidr": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			StateFunc: func(val interface{}) string {
				return strings.TrimSpace(val.(string))
			},
			ValidateFunc: validateIPCidrRange,
			Description:  "CIDR representing IP range for this subnet, based on IP version",
		},
	}
}
