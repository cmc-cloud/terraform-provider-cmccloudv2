package cmccloudv2

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDevopsProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevopsProjectCreate,
		Read:   resourceDevopsProjectRead,
		// Update: resourceDevopsProjectUpdate,
		Delete: resourceDevopsProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDevopsProjectImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        devopsProjectSchema(),
	}
}

// func resourceDevopsProjectUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return resourceDevopsProjectRead(d, meta)
// }

func resourceDevopsProjectCreate(d *schema.ResourceData, meta interface{}) error {
	devopsproject, err := getClient(meta).DevopsProject.Create(map[string]interface{}{
		"name": d.Get("name").(string),
	})
	if err != nil {
		return fmt.Errorf("error creating Devops Project : %s", err)
	}
	d.SetId(strconv.Itoa(devopsproject.ID))
	return resourceDevopsProjectRead(d, meta)
}

func resourceDevopsProjectRead(d *schema.ResourceData, meta interface{}) error {
	devopsproject, err := getClient(meta).DevopsProject.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Devops Project %s: %v", d.Id(), err)
	}

	_ = d.Set("name", devopsproject.Name)
	_ = d.Set("description", devopsproject.Description)
	_ = d.Set("is_default", devopsproject.IsDefault)
	_ = d.Set("created_at", devopsproject.CreatedAt)
	return nil
}

func resourceDevopsProjectDelete(d *schema.ResourceData, meta interface{}) error {
	// _, err := getClient(meta).DevopsProject.Delete(d.Id())
	// if err != nil {
	// 	return fmt.Errorf("error delete devops project: %v", err)
	// }
	return nil
}

func resourceDevopsProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDevopsProjectRead(d, meta)
	return []*schema.ResourceData{d}, err
}

// func waitUntilDevopsProjectDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
// 	return waitUntilResourceDeleted(d, meta, WaitConf{
// 		Delay:      10 * time.Second,
// 		MinTimeout: 20 * time.Second,
// 	}, func(id string) (any, error) {
// 		return getClient(meta).DevopsProject.Get(id)
// 	})
// }
