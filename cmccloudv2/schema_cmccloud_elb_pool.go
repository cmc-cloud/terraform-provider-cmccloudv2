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
			Description:  "The ID of the ELB to attach the pool to",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the pool",
		},
		"listener_id": {
			Type:        schema.TypeString,
			Description: "The ID of the listener that the pool attached to",
			Computed:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A description of the pool. Changing this updates the pool's description.",
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PROXY", "PROXYV2", "SCTP", "TCP", "UDP"}, false),
			Description:  "The protocol for which this pool and its members listen. A valid value is HTTP, HTTPS, PROXY, PROXYV2, SCTP, TCP, or UDP.",
		},
		"algorithm": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "SOURCE_IP_PORT"}, false),
			Description:  "The algorithm used to distribute traffic to the pool members. A valid value is ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP, or SOURCE_IP_PORT.",
		},
		"session_persistence": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "NONE",
			ValidateFunc: validation.StringInSlice([]string{"NONE", "SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE"}, false),
			Description:  "Session persistence type for the pool. One of NONE, APP_COOKIE, HTTP_COOKIE, or SOURCE_IP.",
		},
		"cookie_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the cookie to use for session persistence. Only applicable to the APP_COOKIE session persistence type where it is required.",
		},
		"tls_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "When true connections to backend member servers will use TLS encryption. Default is false.",
		},
		"tls_ciphers": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "List of ciphers in OpenSSL format (colon-separated)",
		},
		"tls_versions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"SSLv3", "TLSv1", "TLSv1.1", "TLSv1.2", "TLSv1.3"}, false),
			},
			Description: "A list of TLS protocol versions. Available versions: SSLv3, TLSv1, TLSv1.1, TLSv1.2, TLSv1.3",
		},
		"created_at": {
			Type:        schema.TypeString,
			Description: "The creation time of the pool",
			Computed:    true,
		},
		"provisioning_status": {
			Type:        schema.TypeString,
			Description: "The state of the operation — in other words, whether pool is still creating, updating, or deleting the resource.",
			Computed:    true,
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shows the real runtime health of the pool.",
		},
	}
}
