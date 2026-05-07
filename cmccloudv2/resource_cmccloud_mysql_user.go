package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
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
		"allowHost":       d.Get("host").(string), //strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ","),
		"userPermissions": getMysqlUserPermission(d),
	}
	_, err := getClient(meta).MysqlInstance.CreateUser(d.Get("instance_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating mysql user: %v", err)
	}
	_, err = waitUntilDatabaseUserFound(d, meta)
	if err != nil {
		return fmt.Errorf("error creating mysql user: %v", err)
	}
	d.SetId(buildMysqlUserID(d.Get("instance_id").(string), d.Get("username").(string), d.Get("host").(string)))
	return resourceMysqlUserRead(d, meta)
}

func resourceMysqlUserRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, host, err := parseMysqlUserID(d.Id())
	if err != nil {
		return err
	}
	user, err := getClient(meta).DBv2.GetUser(instanceID, username, host)
	if err != nil {
		return fmt.Errorf("error retrieving mysql user %s/%s: %v", instanceID, username, err)
	}
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("username", user.Name)
	_ = d.Set("host", user.Host)
	return nil
}

func resourceMysqlUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("password") || d.HasChange("host") || d.HasChange("user_permissions") {
		// Lấy password và currentAllowHost là giá trị trước khi thay đổi
		var oldPassword string
		if d.HasChange("password") {
			old, _ := d.GetChange("password")
			oldPassword = old.(string)
		} else {
			oldPassword = d.Get("password").(string)
		}

		var currentAllowHost string
		if d.HasChange("host") {
			oldHosts, _ := d.GetChange("host")
			currentAllowHost = oldHosts.(string)
		} else {
			currentAllowHost = d.Get("host").(string)
		}
		// if d.HasChange("hosts") {
		// 	oldHosts, _ := d.GetChange("hosts")
		// 	currentAllowHost = strings.Join(getStringArrayFromTypeSet(oldHosts.(*schema.Set)), ",")
		// } else {
		// 	currentAllowHost = strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ",")
		// }

		params := map[string]interface{}{
			"username":         d.Get("username").(string),
			"newPassword":      d.Get("password").(string),
			"password":         oldPassword,
			"currentAllowHost": currentAllowHost,
			"newAllowHost":     d.Get("host").(string), //strings.Join(getStringArrayFromTypeSet(d.Get("hosts").(*schema.Set)), ","),
			"userPermissions":  getMysqlUserPermission(d),
		}

		_, err := client.MysqlInstance.UpdateUser(d.Get("instance_id").(string), params)
		if err != nil {
			return fmt.Errorf("error when update mysql user %s: %v", id, err)
		}
		_, err = waitUntilDatabaseUserFound(d, meta)
		if err != nil {
			return fmt.Errorf("error when update mysql user %s: %v", id, err)
		}
	}
	return resourcePostgresInstanceRead(d, meta)
}
func resourceMysqlUserDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, username, host, err := parseMysqlUserID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).DBv2.DeleteUser(instanceID, username, host)
	if err != nil {
		return fmt.Errorf("error deleting mysql user %s/%s/%s: %v", instanceID, username, host, err)
	}
	_, err = waitUntilDatabaseUserDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error deleting mysql user %s/%s/%s: %v", instanceID, username, host, err)
	}
	return nil
}

func resourceMysqlUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMysqlUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildMysqlUserID(instanceID string, username string, host string) string {
	return instanceID + "/user/" + username + "/" + host
}

func parseMysqlUserID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 4 || parts[0] == "" || parts[2] == "" || parts[3] == "" {
		return "", "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/user/<name>/<host>", id)
	}
	return parts[0], parts[2], parts[3], nil
}

func waitUntilDatabaseUserDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"false"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBv2.ListUsers(d.Get("instance_id").(string), map[string]string{})
	}, func(obj interface{}) string {
		users := obj.([]gocmcapiv2.DBv2User)
		for _, t := range users {
			if t.Name == d.Get("username").(string) && t.Host == d.Get("host").(string) {
				return "false"
			}
		}
		return "true"
	})
}

func waitUntilDatabaseUserFound(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"true"}, []string{"false"}, WaitConf{
		Timeout:    40 * time.Second,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DBv2.GetUser(d.Get("instance_id").(string), d.Get("username").(string), d.Get("host").(string))
	}, func(obj interface{}) string {
		user := obj.(gocmcapiv2.DBv2User)
		if user.Name != "" {
			return "true"
		}
		return "false"
	})
}
