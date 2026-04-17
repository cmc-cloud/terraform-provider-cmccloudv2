package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRdsClusterUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsClusterUserCreate,
		Read:   resourceRdsClusterUserRead,
		Update: resourceRdsClusterUserUpdate,
		Delete: resourceRdsClusterUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRdsClusterUserImport,
		},
		SchemaVersion: 1,
		Schema:        rdsClusterUserSchema(),
	}
}

func resourceRdsClusterUserCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":      d.Get("name").(string),
		"password":  d.Get("password").(string),
		"host":      d.Get("host").(string),
		"databases": getStringArrayFromTypeSet(d.Get("databases").(*schema.Set)),
	}

	_, err := getClient(meta).RdsCluster.CreateUser(d.Get("cluster_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating rds cluster user: %v", err)
	}
	d.SetId(buildRdsClusterUserID(d.Get("cluster_id").(string), d.Get("name").(string), d.Get("host").(string)))
	return resourceRdsClusterUserRead(d, meta)
}

func resourceRdsClusterUserRead(d *schema.ResourceData, meta interface{}) error {
	clusterId, name, _, err := parseRdsClusterUserID(d.Id())
	if err != nil {
		return err
	}
	user, err := getClient(meta).RdsCluster.GetUser(clusterId, name)
	if err != nil {
		return fmt.Errorf("error retrieving rds cluster user %s/%s: %v", clusterId, name, err)
	}
	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("name", user.Name)
	_ = d.Set("host", user.Host)
	_ = d.Set("databases", user.Databases)
	return nil
}

func resourceRdsClusterUserUpdate(d *schema.ResourceData, meta interface{}) error {
	clusterId, name, oldHost, err := parseRdsClusterUserID(d.Id())
	if err != nil {
		return err
	}
	if d.HasChange("password") || d.HasChange("databases") {
		params := map[string]interface{}{
			"password":  d.Get("password").(string),
			"new_host":  d.Get("host").(string),
			"databases": getStringArrayFromTypeSet(d.Get("databases").(*schema.Set)),
		}
		if len(params) > 0 {
			_, err := getClient(meta).RdsCluster.UpdateUser(clusterId, d.Get("name").(string), oldHost, params)
			if err != nil {
				return fmt.Errorf("error updating rds cluster user %s/%s: %v", clusterId, name, err)
			}
		}
	}
	return resourceRdsClusterUserRead(d, meta)
}

func resourceRdsClusterUserDelete(d *schema.ResourceData, meta interface{}) error {
	clusterId, name, host, err := parseRdsClusterUserID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).RdsCluster.DeleteUser(clusterId, name, host)
	if err != nil {
		return fmt.Errorf("error deleting rds cluster user %s/%s: %v", clusterId, name, err)
	}
	return nil
}

func resourceRdsClusterUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRdsClusterUserRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildRdsClusterUserID(clusterId string, name string, host string) string {
	return clusterId + "/user/" + name + "/" + host
}

func parseRdsClusterUserID(id string) (string, string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 4 || parts[0] == "" || parts[2] == "" || parts[3] == "" {
		return "", "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/user/<name>/<host>", id)
	}
	return parts[0], parts[2], parts[3], nil
}
