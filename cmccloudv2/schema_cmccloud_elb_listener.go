package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func elblistenerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"elb_id": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the ELB to attach the listener to",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the listener",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the listener. Changing this updates the listener's description.",
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "SCTP", "TCP", "TERMINATED_HTTPS", "UDP"}, false),
			Description:  "The protocol for which this listener listens. A valid value is HTTP, HTTPS, SCTP, TCP, TERMINATED_HTTPS, or UDP.",
		},
		"protocol_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "The port on which the listener listens. A valid value is between 1 and 65535.",
		},
		"default_pool_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the default pool to use for the listener.",
		},
		"sni_container_refs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A list of SSL certificate container references. This is required when protocol is HTTPS or TERMINATED_HTTPS.",
			Elem: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The ID of the SSL certificate container.",
			},
		},
		"default_tls_container_ref": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the default SSL certificate container. This is required when protocol is HTTPS or TERMINATED_HTTPS.",
		},

		"timeout_client_data": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Frontend client inactivity timeout in milliseconds.",
			// Default:  50000,
		},
		"timeout_tcp_inspect": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Timeout for TCP packet inspection in milliseconds.",
			// Default:  0,
		},
		"timeout_member_connect": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Backend member connection timeout in milliseconds. Default: 5000.",
			Default:     5000,
		},
		"timeout_member_data": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     50000,
			Description: "Backend member inactivity timeout in milliseconds. Default: 50000.",
		},
		"connection_limit": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     -1,
			Description: "The maximum number of connections permitted for this listener. Default value is -1 which represents infinite connections or a default value defined by the provider driver.",
		},
		"allowed_cidrs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:        schema.TypeString,
				Description: "The CIDR block to allow traffic from.",
			},
			Description: "A list of IPv4, IPv6 or mix of both CIDRs. The default is all allowed. When a list of CIDRs is provided, the default switches to deny all.",
		},
		"x_forwarded_for": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to include the X-Forwarded-For header in the request.",
			// Default:  false,
		},
		"x_forwarded_port": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to include the X-Forwarded-Port header in the request.",
			// Default:  false,
		},
		"x_forwarded_proto": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to include the X-Forwarded-Proto header in the request.",
			// Default:  false,
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the listener",
		},
		"provisioning_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of the operation — in other words, whether listener is still creating, updating, or deleting the resource.",
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shows the real runtime health of the listener.",
		},
	}
}
