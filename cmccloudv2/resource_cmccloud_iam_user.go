package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamUserCreate,
		Read:   resourceIamUserRead,
		Update: resourceIamUserUpdate,
		Delete: resourceIamUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIamUserImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        iamUserSchema(),
	}
}

func resourceIamUserCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"username":   d.Get("short_name").(string),
		"first_name": d.Get("first_name").(string),
		"last_name":  d.Get("last_name").(string),
		"password":   d.Get("password").(string),
		"email":      d.Get("email").(string),
	}
	iamuser, err := getClient(meta).IamUser.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating iam user: %s", err)
	}

	d.SetId(iamuser.Username)
	if d.Get("enabled").(bool) == false {
		getClient(meta).IamUser.Disable(iamuser.Username)
	}
	return resourceIamUserRead(d, meta)
}

func resourceIamUserRead(d *schema.ResourceData, meta interface{}) error {
	iamuser, err := getClient(meta).IamUser.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving iam user %s: %v", d.Id(), err)
	}
	_ = d.Set("id", iamuser.Username)
	_ = d.Set("short_name", iamuser.ShortName)
	_ = d.Set("username", iamuser.Username)
	_ = d.Set("first_name", iamuser.FirstName)
	_ = d.Set("last_name", iamuser.LastName)
	_ = d.Set("email", iamuser.Email)
	_ = d.Set("enabled", iamuser.Enabled)
	// _ = d.Set("group_ids", iamuser.Gr)
	return nil
}

func resourceIamUserUpdate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamUser.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving iam user %s: %v", d.Id(), err)
	}

	if d.HasChange("first_name") || d.HasChange("last_name") {
		params := map[string]interface{}{
			"first_name": d.Get("first_name").(string),
			"last_name":  d.Get("last_name").(string),
			"email":      d.Get("email").(string),
		}
		_, err := getClient(meta).IamUser.Update(d.Id(), params)
		if err != nil {
			return fmt.Errorf("Error when update iam user info [%s]: %v", d.Id(), err)
		}
	}

	if d.HasChange("email") {
		_, err := getClient(meta).IamUser.UpdateEmail(d.Id(), d.Get("email").(string))
		if err != nil {
			return fmt.Errorf("Error when update iam user email [%s]: %v", d.Id(), err)
		}
	}

	if d.HasChange("password") {
		params := map[string]interface{}{
			"password": d.Get("password").(string),
		}
		_, err := getClient(meta).IamUser.SetPassword(d.Id(), params)
		if err != nil {
			return fmt.Errorf("Error when update iam user password [%s]: %v", d.Id(), err)
		}
	}

	if d.HasChange("enabled") {
		if d.Get("enabled").(bool) {
			getClient(meta).IamUser.Enable(d.Get("username").(string))
		} else {
			getClient(meta).IamUser.Disable(d.Get("username").(string))
		}
	}
	return resourceIamUserRead(d, meta)
}

func resourceIamUserDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamUser.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete iam user: %v", err)
	}
	_, err = waitUntilIamUserDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete iam user: %v", err)
	}
	return nil
}

func resourceIamUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilIamUserDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).IamUser.Get(id)
	})
}
