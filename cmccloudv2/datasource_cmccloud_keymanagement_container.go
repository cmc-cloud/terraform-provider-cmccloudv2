package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceKeyManagementContainerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the container",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter KeyManagementContainer by name, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Filter KeyManagementContainer by type, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceKeyManagementContainer() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceKeyManagementContainerRead,
		Schema: datasourceKeyManagementContainerSchema(),
	}
}

func dataSourceKeyManagementContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allKeyManagementContainers []gocmcapiv2.KeyManagementContainer
	if container_id := d.Get("container_id").(string); container_id != "" {
		container, err := client.KeyManagement.Get(container_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve container [%s]: %s", container_id, err)
			}
		}
		allKeyManagementContainers = append(allKeyManagementContainers, container)
	} else {
		params := map[string]string{
			"page": "1",
			"size": "1000",
		}
		containers, err := client.KeyManagement.List(params)
		if err != nil {
			return fmt.Errorf("error when get containeres %v", err)
		}
		allKeyManagementContainers = append(allKeyManagementContainers, containers...)
	}
	if len(allKeyManagementContainers) > 0 {
		var filteredKeyManagementContainers []gocmcapiv2.KeyManagementContainer
		for _, container := range allKeyManagementContainers {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(container.Name, v) {
					continue
				}
			}
			if v := d.Get("type").(string); v != "" {
				if !strings.EqualFold(container.Type, v) {
					continue
				}
			}
			filteredKeyManagementContainers = append(filteredKeyManagementContainers, container)
		}
		allKeyManagementContainers = filteredKeyManagementContainers
	}
	if len(allKeyManagementContainers) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allKeyManagementContainers) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allKeyManagementContainers)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeKeyManagementContainerAttributes(d, allKeyManagementContainers[0])
}

func dataSourceComputeKeyManagementContainerAttributes(d *schema.ResourceData, container gocmcapiv2.KeyManagementContainer) error {
	log.Printf("[DEBUG] Retrieved container %s: %#v", container.ID, container)
	d.SetId(container.ID)
	_ = d.Set("name", container.Name)
	_ = d.Set("type", container.Type)
	_ = d.Set("created_at", container.Created)
	return nil
}
