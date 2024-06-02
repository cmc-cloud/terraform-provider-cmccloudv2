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
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateName,
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PING", "TCP", "TLS-HELLO", "UDP-CONNECT", "SCTP"}, false),
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
