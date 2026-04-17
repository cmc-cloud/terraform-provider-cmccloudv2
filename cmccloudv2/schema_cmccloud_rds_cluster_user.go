package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func rdsClusterUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "RdsCluster id",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "RdsCluster username",
		},
		"host": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Host from which the user can connect (default is '%')",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Password for user. Must satisfy all requirements: minimum length 12, at least one uppercase character, at least one lowercase character, at least one number, at least one special character (! % & ( ) * + - . < = > ? @ [ ] ^ _ { } #)",
		},
		"databases": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Databases that this user has access to",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
