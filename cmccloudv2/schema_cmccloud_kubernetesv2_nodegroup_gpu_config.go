package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kubernetesv2NodeGroupGpuConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"nodegroup_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"gpu_model": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"driver": {
			Type:     schema.TypeString,
			Required: true,
		},
		"strategy": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "single",
			ValidateFunc: validation.StringInSlice([]string{"single", "mixed"}, false),
		},
		"mig_profile": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"gpu_profiles": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"replicas": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
	}
}
