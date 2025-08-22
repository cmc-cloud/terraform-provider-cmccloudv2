package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamGroupCreate,
		Read:   resourceIamGroupRead,
		Update: resourceIamGroupUpdate,
		Delete: resourceIamGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIamGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        iamGroupSchema(),
	}
}

func resourceIamGroupCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
	}
	group, err := getClient(meta).IamGroup.Create(params)

	if err != nil {
		return fmt.Errorf("error creating iam group: %s", err)
	}

	// d.SetId(d.Get("name").(string))
	d.SetId(group.ID)
	return resourceIamGroupRead(d, meta)
}

func resourceIamGroupRead(d *schema.ResourceData, meta interface{}) error {
	iamgroup, err := getClient(meta).IamGroup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving iam group %s: %v", d.Id(), err)
	}
	_ = d.Set("id", iamgroup.ID)
	_ = d.Set("name", iamgroup.Name)
	_ = d.Set("description", iamgroup.Description)
	return nil
}

func resourceIamGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamGroup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving iam group %s: %v", d.Id(), err)
	}

	if d.HasChange("description") {
		params := map[string]interface{}{
			"description": d.Get("description").(string),
		}
		_, err := getClient(meta).IamGroup.Update(d.Get("name").(string), params)
		if err != nil {
			return fmt.Errorf("error when update iam group info [%s]: %v", d.Id(), err)
		}
	}
	return resourceIamGroupRead(d, meta)
}

func resourceIamGroupDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamGroup.Delete(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("error delete iam group: %v", err)
	}
	_, err = waitUntilIamGroupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete iam group: %v", err)
	}
	return nil
}

func resourceIamGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilIamGroupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).IamGroup.Get(id)
	})
}
