package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMysqlUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMysqlUserCreate,
		Read:   resourceMysqlUserRead,
		Update: resourceMysqlUserUpdate,
		Delete: resourceMysqlUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMysqlUserImport,
		},
		SchemaVersion: 1,
		Schema:        mysqlUserSchema(),
	}
}

func getMysqlUserPermission(d *schema.ResourceData) []map[string]interface{} {
	userPermissions := []map[string]interface{}{}
	if v, ok := d.GetOk("user_permissions"); ok {
		for _, up := range v.(*schema.Set).List() {
			upMap := up.(map[string]interface{})
			dbName := ""
			if v, ok := upMap["database"]; ok {
				dbName = v.(string)
			}
			table := ""
			if v, ok := upMap["table"]; ok {
				table = v.(string)
			}
			permissions := []string{}
			if perms, ok := upMap["permissions"]; ok && perms != nil {
				for _, perm := range perms.(*schema.Set).List() {
					permissions = append(permissions, perm.(string))
				}
				// if there a one permission = *, set permissions = full permissions
				for _, perm := range permissions {
					if perm == "*" {
						permissions = []string{
							"alter",
							"create",
							"delete",
							"drop",
							"insert",
							"select",
							"update",
							"index",
							"create view",
							"trigger",
							"event",
							"references",
						}
						break
					}
				}
			}

			userPermissions = append(userPermissions, map[string]interface{}{
				"databaseName": dbName,
				"table":        table,
				"permissions":  permissions,
			})
		}
	}
	return userPermissions
}
func resourceMysqlUserCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"username":        d.Get("username").(string),
		"password":        d.Get("password").(string),
		"allowHost":       strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ","),
		"userPermissions": getMysqlUserPermission(d),
	}
	_, err := getClient(meta).MysqlInstance.CreateUser(d.Get("instance_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating mysql user: %v", err)
	}
	d.SetId(buildMysqlUserID(d.Get("instance_id").(string), d.Get("username").(string)))
	return resourceMysqlUserRead(d, meta)
}

func resourceMysqlUserRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, err := parseMysqlUserID(d.Id())
	if err != nil {
		return err
	}
	user, err := getClient(meta).MysqlInstance.GetUser(instanceID, username)
	if err != nil {
		return fmt.Errorf("error retrieving mysql user %s/%s: %v", instanceID, username, err)
	}
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("username", user.Name)
	return nil
}

func resourceMysqlUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("password") || d.HasChange("hosts") || d.HasChange("user_permissions") {
		// Lấy password và currentAllowHost là giá trị trước khi thay đổi
		var oldPassword string
		if d.HasChange("password") {
			old, _ := d.GetChange("password")
			oldPassword = old.(string)
		} else {
			oldPassword = d.Get("password").(string)
		}

		var currentAllowHost string
		if d.HasChange("hosts") {
			oldHosts, _ := d.GetChange("hosts")
			currentAllowHost = strings.Join(getStringArrayFromTypeSet(oldHosts.(*schema.Set)), ",")
		} else {
			currentAllowHost = strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ",")
		}

		params := map[string]interface{}{
			"username":         d.Get("username").(string),
			"newPassword":      d.Get("password").(string),
			"password":         oldPassword,
			"currentAllowHost": currentAllowHost,
			"newAllowHost":     strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ","),
			"userPermissions":  getMysqlUserPermission(d),
		}

		_, err := client.MysqlInstance.UpdateUser(id, params)
		if err != nil {
			return fmt.Errorf("error when update mysql user %s: %v", id, err)
		}
		_, err = waitUntilPostgresInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error when update mysql user %s: %v", id, err)
		}
	}
	return resourcePostgresInstanceRead(d, meta)
}
func resourceMysqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, user, err := parseMysqlUserID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).MysqlInstance.DeleteUser(instanceID, user)
	if err != nil {
		return fmt.Errorf("error deleting mysql user %s/%s: %v", instanceID, user, err)
	}
	return nil
}

func resourceMysqlUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMysqlUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildMysqlUserID(instanceID string, username string) string {
	return instanceID + "/user/" + username
}

func parseMysqlUserID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/user/<username>", id)
	}
	return parts[0], parts[2], nil
}
