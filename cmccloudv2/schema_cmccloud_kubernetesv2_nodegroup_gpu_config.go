package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func kubernetesv2NodeGroupGpuConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the Kubernetes cluster to attach the GPU config to",
		},
		"nodegroup_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the node group to attach the GPU config to",
		},
		"gpu_model": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The GPU model of the GPU config",
		},
		"driver": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The driver of the GPU config",
		},
		"strategy": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "single",
			ValidateFunc: validation.StringInSlice([]string{"single", "mixed"}, false),
			Description:  "The strategy of the GPU config, a valid value is single or mixed",
		},
		"mig_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The MIG profile of the GPU config",
		},
		"gpu_profiles": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "The name of the GPU profile",
					},
					"replicas": {
						Type:        schema.TypeInt,
						Required:    true,
						Description: "The replicas of the GPU profile",
					},
				},
			},
			Description: "The GPU profiles of the GPU config",
		},
	}
}
