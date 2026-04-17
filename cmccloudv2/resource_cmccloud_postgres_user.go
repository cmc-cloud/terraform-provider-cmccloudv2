package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresUserCreate,
		Read:   resourcePostgresUserRead,
		Update: resourcePostgresUserUpdate,
		Delete: resourcePostgresUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresUserImport,
		},
		SchemaVersion: 1,
		Schema:        postgresUserSchema(),
	}
}

func resourcePostgresUserCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	if password == "" {
		return fmt.Errorf("`password` is required when creating postgres user")
	}

	params := map[string]interface{}{
		"username":    username,
		"password":    password,
		"permissions": []string{},
	}
	if permissions, ok := d.GetOk("permissions"); ok {
		params["permissions"] = getStringArrayFromTypeSet(permissions.(*schema.Set))
	}

	_, err := getClient(meta).PostgresInstance.CreateUser(instanceID, params)
	if err != nil {
		return fmt.Errorf("error creating postgres user: %v", err)
	}
	d.SetId(buildPostgresUserID(instanceID, username))
	return resourcePostgresUserRead(d, meta)
}

func resourcePostgresUserRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parsePostgresUserID(d.Id())
	if err != nil {
		return err
	}
	user, err := getClient(meta).PostgresInstance.GetUser(instanceID, username)
	if err != nil {
		return fmt.Errorf("error retrieving postgres user %s/%s: %v", instanceID, username, err)
	}
	_ = d.Set("id", buildPostgresUserID(instanceID, username))
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("username", user.Username)
	_ = d.Set("permissions", user.Permissions)
	return nil
}

func resourcePostgresUserUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parsePostgresUserID(d.Id())
	if err != nil {
		return err
	}
	if d.HasChange("password") || d.HasChange("permissions") {
		params := map[string]interface{}{
			"username":    username,
			"password":    d.Get("password").(string),
			"permissions": getStringArrayFromTypeSet(d.Get("permissions").(*schema.Set)),
		}
		if len(params) > 0 {
			_, err := getClient(meta).PostgresInstance.UpdateUser(instanceID, params)
			if err != nil {
				return fmt.Errorf("error updating postgres user %s/%s: %v", instanceID, username, err)
			}
		}
	}
	return resourcePostgresUserRead(d, meta)
}

func resourcePostgresUserDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parsePostgresUserID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).PostgresInstance.DeleteUser(instanceID, username)
	if err != nil {
		return fmt.Errorf("error deleting postgres user %s/%s: %v", instanceID, username, err)
	}
	return nil
}

func resourcePostgresUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildPostgresUserID(instanceID string, username string) string {
	return instanceID + "/user/" + username
}

func parsePostgresUserID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/user/<username>", id)
	}
	return parts[0], parts[2], nil
}
