package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKeyManagementContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyManagementContainerCreate,
		Read:   resourceKeyManagementContainerRead,
		// Update: resourceKeyManagementContainerUpdate,
		Delete: resourceKeyManagementContainerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyManagementContainerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        keymanagementcontainerSchema(),
	}
}

// func resourceKeyManagementContainerUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return resourceKeyManagementContainerRead(d, meta)
// }

func resourceKeyManagementContainerCreate(d *schema.ResourceData, meta interface{}) error {
	container, err := getClient(meta).KeyManagement.Create(map[string]interface{}{
		"name": d.Get("name").(string),
		"type": "generic", // d.Get("type").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating KeyManagementContainer: %s", err)
	}
	d.SetId(container.Data.ID)
	return resourceKeyManagementContainerRead(d, meta)
}

func resourceKeyManagementContainerRead(d *schema.ResourceData, meta interface{}) error {
	container, err := getClient(meta).KeyManagement.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving KeyManagement Container %s: %v", d.Id(), err)
	}

	_ = d.Set("name", container.Name)
	_ = d.Set("container_ref", container.ContainerRef)
	_ = d.Set("created_at", container.Created)
	return nil
}

func resourceKeyManagementContainerDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).KeyManagement.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("Error delete KeyManagement Container: %v", err)
	}
	return nil
}

func resourceKeyManagementContainerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKeyManagementContainerRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKeyManagementContainerDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KeyManagement.Get(id)
	})
}
