package cmccloudv2

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceContainerRegistryRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerRegistryRepositoryCreate,
		Read:   resourceContainerRegistryRepositoryRead,
		// Update: resourceContainerRegistryRepositoryUpdate,
		Delete: resourceContainerRegistryRepositoryDelete,
		Importer: &schema.ResourceImporter{
			State: resourceContainerRegistryRepositoryImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        containerRegistrySchema(),
	}
}

// func resourceContainerRegistryRepositoryUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return resourceContainerRegistryRepositoryRead(d, meta)
// }

func resourceContainerRegistryRepositoryCreate(d *schema.ResourceData, meta interface{}) error {
	devops_project_id := d.Get("devops_project_id").(string)
	registry, err := getClient(meta).ContainerRegistry.Create(devops_project_id, map[string]interface{}{
		"name": d.Get("name").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating Repository: %s", err)
	}
	d.SetId(strconv.Itoa(registry.ID))
	return resourceContainerRegistryRepositoryRead(d, meta)
}

func resourceContainerRegistryRepositoryRead(d *schema.ResourceData, meta interface{}) error {
	devops_project_id := d.Get("devops_project_id").(string)
	registry, err := getClient(meta).ContainerRegistry.Get(devops_project_id, d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Repository %s: %v", d.Id(), err)
	}

	_ = d.Set("name", registry.Name)
	_ = d.Set("uri", registry.URI)
	_ = d.Set("devops_project_id", strconv.Itoa(registry.ProjectId))
	_ = d.Set("created_at", registry.CreatedAt)
	return nil
}

func resourceContainerRegistryRepositoryDelete(d *schema.ResourceData, meta interface{}) error {
	devops_project_id := d.Get("devops_project_id").(string)
	_, err := getClient(meta).ContainerRegistry.Delete(devops_project_id, d.Id())
	if err != nil {
		return fmt.Errorf("Error delete Repository: %v", err)
	}
	return nil
}

func resourceContainerRegistryRepositoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceContainerRegistryRepositoryRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilContainerRegistryRepositoryDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	devops_project_id := d.Get("devops_project_id").(string)
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ContainerRegistry.Get(devops_project_id, id)
	})
}
