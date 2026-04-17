package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func mongoUserSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mongo instance id",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Mongo username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Mongo user password",
		},
		"permissions": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "Mongo user permissions, e.g. CREATEDB, CREATEROLE, LOGIN, REPLICATION",
			Elem: &schema.Schema{
				Type: schema.TypeString,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := strings.ToUpper(v.(string))

					allowed := map[string]bool{
						"CREATEDB":    true,
						"CREATEROLE":  true,
						"LOGIN":       true,
						"REPLICATION": true,
					}

					if !allowed[value] {
						es = append(es, fmt.Errorf("%q must be one of CREATEDB, CREATEROLE, LOGIN, REPLICATION", k))
					}
					return
				},
			},
		},
	}
}
