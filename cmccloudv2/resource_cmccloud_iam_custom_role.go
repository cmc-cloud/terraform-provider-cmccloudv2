package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamCustomRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamCustomRoleCreate,
		Read:   resourceIamCustomRoleRead,
		Update: resourceIamCustomRoleUpdate,
		Delete: resourceIamCustomRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIamCustomRoleImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        iamCustomRoleSchema(),
	}
}

func resourceIamCustomRoleCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"content":     d.Get("content").(string),
	}
	customrole, err := getClient(meta).IamCustomRole.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating iam custom role: %s", err)
	}

	d.SetId(customrole.ID)
	return resourceIamCustomRoleRead(d, meta)
}

func resourceIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	customrole, err := getClient(meta).IamCustomRole.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving custom role %s: %v", d.Id(), err)
	}
	_ = d.Set("id", customrole.ID)
	_ = d.Set("name", customrole.Name)
	_ = d.Set("description", customrole.Description)
	_ = d.Set("content", customrole.Content)
	_ = d.Set("created", customrole.Created)
	return nil
}

func resourceIamCustomRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamCustomRole.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving custom role %s: %v", d.Id(), err)
	}

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("content") {
		params := map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"content":     d.Get("content").(string),
		}
		_, err := getClient(meta).IamCustomRole.Update(d.Id(), params)
		if err != nil {
			return fmt.Errorf("Error when update custom role [%s]: %v", d.Id(), err)
		}
	}
	return resourceIamCustomRoleRead(d, meta)
}

func resourceIamCustomRoleDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamCustomRole.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete custom role: %v", err)
	}
	_, err = waitUntilIamCustomRoleDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete custom role: %v", err)
	}
	return nil
}

func resourceIamCustomRoleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamCustomRoleRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilIamCustomRoleDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).IamCustomRole.Get(id)
	})
}
