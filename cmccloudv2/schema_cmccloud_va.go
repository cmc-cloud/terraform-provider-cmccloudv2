package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func vaSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The name of the VA",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"HOST_DISCOVERY", "BASIC_NETWORK_SCAN", "ADVANCED_SCAN", "WEB_APPLICATION_TESTS"}, false),
			Description:  "The type of the VA, HOST_DISCOVERY, BASIC_NETWORK_SCAN, ADVANCED_SCAN, WEB_APPLICATION_TESTS",
		},
		"schedule": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The schedule of the VA, example `2026-01-25 01:00:50`. If not set, it will run immediately",
		},
		"target": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The target of the VA, example `https://www.google.com`, `1.2.3.4`",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The description of the VA",
		},
		"report_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the report",
		},
		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the VA",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The created time of the VA",
		},
	}
}
