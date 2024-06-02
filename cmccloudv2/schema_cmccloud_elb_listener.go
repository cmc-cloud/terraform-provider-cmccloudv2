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
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "SCTP", "TCP", "TERMINATED_HTTPS", "UDP"}, false),
		},
		"protocol_port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},
		"default_pool_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateUUID,
		},
		"sni_container_refs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"default_tls_container_ref": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"timeout_client_data": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Default value: 50000",
			// Default:  50000,
		},
		"timeout_tcp_inspect": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Default value: 0",
			// Default:  0,
		},
		"timeout_member_connect": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Default value: 5000",
			Default:     5000,
		},
		"timeout_member_data": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     50000,
			Description: "Default value: 50000",
		},
		"connection_limit": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  -1,
		},
		"allowed_cidrs": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"x_forwarded_for": {
			Type:     schema.TypeBool,
			Optional: true,
			// Default:  false,
		},
		"x_forwarded_port": {
			Type:     schema.TypeBool,
			Optional: true,
			// Default:  false,
		},
		"x_forwarded_proto": {
			Type:     schema.TypeBool,
			Optional: true,
			// Default:  false,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"operating_status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
