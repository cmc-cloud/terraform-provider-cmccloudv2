package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceKeyManagementSecretSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Id of the container",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter KeyManagementSecret by name, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Filter KeyManagementSecret by type, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		// "algorithm": {
		// 	Type:        schema.TypeString,
		// 	Description: "Filter KeyManagementSecret by algorithm, match exactly",
		// 	Optional:    true,
		// 	ForceNew:    true,
		// },
		// "bit_length": {
		// 	Type:        schema.TypeString,
		// 	Description: "Filter KeyManagementSecret by bit_length, match exactly",
		// 	Optional:    true,
		// 	ForceNew:    true,
		// },
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created time of KeyManagementSecret",
		},
	}
}

func datasourceKeyManagementSecret() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceKeyManagementSecretRead,
		Schema: datasourceKeyManagementSecretSchema(),
	}
}

func dataSourceKeyManagementSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allKeyManagementSecrets []gocmcapiv2.KeyManagementSecret
	allKeyManagementSecrets, err := client.KeyManagement.GetSecrets(d.Get("container_id").(string))
	if err != nil {
		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			d.SetId("")
			return fmt.Errorf("unable to retrieve container [%s]: %s", d.Get("container_id").(string), err)
		}
	}
	if len(allKeyManagementSecrets) > 0 {
		var filteredKeyManagementSecrets []gocmcapiv2.KeyManagementSecret
		for _, secret := range allKeyManagementSecrets {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(secret.Name, v) {
					continue
				}
			}
			if v := d.Get("type").(string); v != "" {
				if !strings.EqualFold(secret.SecretType, v) {
					continue
				}
			}
			filteredKeyManagementSecrets = append(filteredKeyManagementSecrets, secret)
		}
		allKeyManagementSecrets = filteredKeyManagementSecrets
	}
	if len(allKeyManagementSecrets) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allKeyManagementSecrets) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allKeyManagementSecrets)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeKeyManagementSecretAttributes(d, allKeyManagementSecrets[0])
}

func dataSourceComputeKeyManagementSecretAttributes(d *schema.ResourceData, secret gocmcapiv2.KeyManagementSecret) error {
	log.Printf("[DEBUG] Retrieved container %s: %#v", secret.ID, secret)
	d.SetId(secret.ID)
	_ = d.Set("name", secret.Name)
	_ = d.Set("type", secret.SecretType)
	_ = d.Set("created_at", secret.Created)
	return nil
}
