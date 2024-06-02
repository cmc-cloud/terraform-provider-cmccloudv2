package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func elbpoolSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"elb_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"listener_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "TCP", "UDP", "PROXY", "PROXYV2,SCTP"}, false),
		},
		"algorithm": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "SOURCE_IP_PORT"}, false),
		},
		"session_persistence": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"NONE", "SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE"}, false),
		},
		"cookie_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tls_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tls_ciphers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tls_versions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"SSLv3", "TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"}, false),
			},
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
