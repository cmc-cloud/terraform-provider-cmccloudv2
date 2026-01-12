package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func elbhealthmonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pool_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateUUID,
			Description:  "The ID of the pool to attach the health monitor to",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
			Description:  "The name of the health monitor",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PING", "TCP", "TLS-HELLO", "UDP-CONNECT", "SCTP"}, false),
			Description:  "The type of health monitor. A valid value is HTTP, HTTPS, PING, TCP, TLS-HELLO, UDP-CONNECT, or SCTP.",
		},
		"http_method": {
			Type:     schema.TypeString,
			Optional: true,
			// ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "OPTIONS", "PATCH", "CONNECT"}, false),
			Description:  "The HTTP method that the health monitor uses for requests",
		},
		"expected_codes": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The list of HTTP status codes expected in response from the member to declare it healthy. Specify one of the following values: A single value, such as 200. A list, such as 200, 202. A range, such as 200-204. The default is 200.",
		},
		"max_retries_down": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 10),
			Description:  "The number of allowed check failures before changing the operating status of the member to ERROR. A valid value is from 1 to 10. The default is 3.",
		},
		"delay": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The time, in seconds, between sending probes to members.",
		},
		"max_retries": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 10),
			Description:  "The number of successful checks before changing the operating status of the member to ONLINE. A valid value is from 1 to 10.",
		},
		"timeout": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The maximum time, in seconds, that a monitor waits to connect before it times out. This value must be less than the delay value.",
		},
		"url_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The HTTP URL path of the request sent by the monitor to test the health of a backend member. Must be a string that begins with a forward slash (/). The default URL path is /",
		},
		"domain_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The domain name, which be injected into the HTTP Host Header to the backend server for HTTP health check.",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The creation time of the health monitor",
		},
		"provisioning_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The state of the operation — in other words, whether health monitor is still creating, updating, or deleting the resource.",
		},
		"operating_status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Shows the real runtime health of the health monitor.",
		},
	}
}
