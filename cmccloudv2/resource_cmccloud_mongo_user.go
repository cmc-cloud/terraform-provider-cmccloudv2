package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoUserCreate,
		Read:   resourceMongoUserRead,
		Update: resourceMongoUserUpdate,
		Delete: resourceMongoUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoUserImport,
		},
		SchemaVersion: 1,
		Schema:        mongoUserSchema(),
	}
}

func resourceMongoUserCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	if password == "" {
		return fmt.Errorf("`password` is required when creating mongo user")
	}

	params := map[string]interface{}{
		"username":    username,
		"password":    password,
		"permissions": []string{},
	}
	if permissions, ok := d.GetOk("permissions"); ok {
		params["permissions"] = getStringArrayFromTypeSet(permissions.(*schema.Set))
	}

	_, err := getClient(meta).MongoInstance.CreateUser(instanceID, params)
	if err != nil {
		return fmt.Errorf("error creating mongo user: %v", err)
	}
	d.SetId(buildMongoUserID(instanceID, username))
	return resourceMongoUserRead(d, meta)
}

func resourceMongoUserRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parseMongoUserID(d.Id())
	if err != nil {
		return err
	}
	user, err := getClient(meta).MongoInstance.GetUser(instanceID, username)
	if err != nil {
		return fmt.Errorf("error retrieving mongo user %s/%s: %v", instanceID, username, err)
	}
	_ = d.Set("id", buildMongoUserID(instanceID, username))
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("username", user.Username)
	_ = d.Set("permissions", user.Permissions)
	return nil
}

func resourceMongoUserUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parseMongoUserID(d.Id())
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
			_, err := getClient(meta).MongoInstance.UpdateUser(instanceID, params)
			if err != nil {
				return fmt.Errorf("error updating mongo user %s/%s: %v", instanceID, username, err)
			}
		}
	}
	return resourceMongoUserRead(d, meta)
}

func resourceMongoUserDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parseMongoUserID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).MongoInstance.DeleteUser(instanceID, username)
	if err != nil {
		return fmt.Errorf("error deleting mongo user %s/%s: %v", instanceID, username, err)
	}
	return nil
}

func resourceMongoUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMongoUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildMongoUserID(instanceID string, username string) string {
	return instanceID + "/user/" + username
}

func parseMongoUserID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/user/<username>", id)
	}
	return parts[0], parts[2], nil
}
