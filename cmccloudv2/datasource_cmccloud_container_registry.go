package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceContainerRegistryRepositorySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"devops_project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"container_registry_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the container registry",
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceContainerRegistryRepository() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceContainerRegistryRepositoryRead,
		Schema: datasourceContainerRegistryRepositorySchema(),
	}
}

func dataSourceContainerRegistryRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	project_id := d.Get("devops_project_id").(string)
	var allContainerRegistryRepositorys []gocmcapiv2.ContainerRegistryRepository
	if registry_id := d.Get("container_registry_id").(string); registry_id != "" {
		registry, err := client.ContainerRegistry.Get(project_id, registry_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve container registry [%s]: %s", registry_id, err)
			}
		}
		allContainerRegistryRepositorys = append(allContainerRegistryRepositorys, registry)
	} else {
		params := map[string]string{
			"q":    d.Get("name").(string),
			"page": "1",
			"size": "1000",
		}
		registrys, err := client.ContainerRegistry.List(project_id, params)
		if err != nil {
			return fmt.Errorf("Error when get container registry %v", err)
		}
		allContainerRegistryRepositorys = append(allContainerRegistryRepositorys, registrys...)
	}
	if len(allContainerRegistryRepositorys) > 0 {
		var filteredContainerRegistryRepositorys []gocmcapiv2.ContainerRegistryRepository
		for _, registry := range allContainerRegistryRepositorys {
			if v := d.Get("name").(string); v != "" {
				if strings.ToLower(registry.Name) != strings.ToLower(v) {
					continue
				}
			}
			filteredContainerRegistryRepositorys = append(filteredContainerRegistryRepositorys, registry)
		}
		allContainerRegistryRepositorys = filteredContainerRegistryRepositorys
	}
	if len(allContainerRegistryRepositorys) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allContainerRegistryRepositorys) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allContainerRegistryRepositorys)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeContainerRegistryRepositoryAttributes(d, allContainerRegistryRepositorys[0])
}

func dataSourceComputeContainerRegistryRepositoryAttributes(d *schema.ResourceData, registry gocmcapiv2.ContainerRegistryRepository) error {
	log.Printf("[DEBUG] Retrieved container registry %s: %#v", registry.ID, registry)
	d.SetId(strconv.Itoa(registry.ID))
	d.Set("name", registry.Name)
	d.Set("uri", registry.URI)
	d.Set("project_id", registry.ProjectId)
	d.Set("created_at", registry.CreatedAt)
	return nil
}
