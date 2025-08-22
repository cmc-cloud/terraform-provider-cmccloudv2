package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamUserServerPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamUserServerPermissionCreate,
		Read:   resourceIamUserServerPermissionRead,
		Update: resourceIamUserServerPermissionUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceIamUserServerPermissionImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Delete: resourceIamUserServerPermissionDelete,
		Schema: iamUserServerPermissionSchema(),
	}
}

func resourceIamUserServerPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamUser.SetServerPermission(d.Get("user_name").(string), map[string]interface{}{
		"server_id":    d.Get("server_id").(string),
		"blocked":      d.Get("blocked").(bool),
		"allow_view":   d.Get("allow_view").(bool),
		"allow_edit":   d.Get("allow_edit").(bool),
		"allow_create": d.Get("allow_create").(bool),
		"allow_delete": d.Get("allow_delete").(bool),
	})

	if err != nil {
		return fmt.Errorf("error setting iam user server permission: %s", err)
	}

	d.SetId(fmt.Sprintf("%s+%s", d.Get("user_name").(string), d.Get("server_id").(string)))
	return resourceIamUserServerPermissionRead(d, meta)
}

func resourceIamUserServerPermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamUser.SetServerPermission(d.Get("user_name").(string), map[string]interface{}{
		"server_id":    d.Get("server_id").(string),
		"blocked":      d.Get("blocked").(bool),
		"allow_view":   d.Get("allow_view").(bool),
		"allow_edit":   d.Get("allow_edit").(bool),
		"allow_create": d.Get("allow_create").(bool),
		"allow_delete": d.Get("allow_delete").(bool),
	})

	if err != nil {
		return fmt.Errorf("error setting iam user server permission: %s", err)
	}
	return resourceIamUserServerPermissionRead(d, meta)
}
func resourceIamUserServerPermissionRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	parts := make([]string, 2)
	splitIdx := -1
	for i := 0; i < len(id); i++ {
		if id[i] == '+' {
			splitIdx = i
			break
		}
	}
	if splitIdx != -1 {
		parts[0] = id[:splitIdx]
		parts[1] = id[splitIdx+1:]
		_ = d.Set("user_name", parts[0])
		_ = d.Set("server_id", parts[1])
	}
	permissions, err := getClient(meta).IamUser.GetServerPermission(d.Get("user_name").(string))
	if err != nil {
		return fmt.Errorf("error read iam user permission: %v", err)
	}
	var found bool
	for _, perm := range permissions {
		if perm.ServerID == parts[1] {
			_ = d.Set("blocked", perm.Blocked)
			_ = d.Set("allow_view", perm.AllowView)
			_ = d.Set("allow_edit", perm.AllowEdit)
			_ = d.Set("allow_create", perm.AllowCreate)
			_ = d.Set("allow_delete", perm.AllowDelete)
			found = true
			break
		}
	}
	if !found {
		d.SetId("")
	}
	return nil
}

func resourceIamUserServerPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamUser.SetServerPermission(d.Get("user_name").(string), map[string]interface{}{
		"server_id":    d.Get("server_id").(string),
		"blocked":      false,
		"allow_view":   false,
		"allow_edit":   false,
		"allow_create": false,
		"allow_delete": false,
	})

	if err != nil {
		return fmt.Errorf("error delete iam user permission: %v", err)
	}
	return nil
}

func resourceIamUserServerPermissionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamUserServerPermissionRead(d, meta)
	return []*schema.ResourceData{d}, err
}
