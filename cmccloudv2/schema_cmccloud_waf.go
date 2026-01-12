package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func wafSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateDomainName,
			Description:  "Domain name",
		},
		"mode": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"DETECT", "BLOCK"}, false),
			Description:  "Mode",
		},
		// "type": {
		// 	Type:         schema.TypeString,
		// 	Required:     true,
		// 	ForceNew:     true,
		// 	ValidateFunc: validation.StringInSlice([]string{"DOMAIN", "IP"}, false),
		// },
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS"}, false),
			Description:  "Protocol",
		},
		"real_server": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Real server domain name or ip address",
		},
		"port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validatePortNumber,
		},
		"certificate_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"send_file": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"client_max_body_size": {
			Type:     schema.TypeInt,
			Default:  1,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"load_balance_enable": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"load_balance_method": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ip_hash", "least_conn", "default"}, false),
		},
		"load_balance_keepalive": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"report_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}
