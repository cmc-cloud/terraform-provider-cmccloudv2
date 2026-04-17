package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRdsClusterDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsClusterDatabaseCreate,
		Read:   resourceRdsClusterDatabaseRead,
		Delete: resourceRdsClusterDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRdsClusterDatabaseImport,
		},
		SchemaVersion: 1,
		Schema:        rdsClusterDatabaseSchema(),
	}
}

func resourceRdsClusterDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":          d.Get("name").(string),
		"character_set": d.Get("character_set").(string),
		"collation":     d.Get("collation").(string),
	}
	_, err := getClient(meta).RdsCluster.CreateDatabase(d.Get("cluster_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating rds cluster database: %v", err)
	}
	d.SetId(buildRdsClusterDatabaseID(d.Get("cluster_id").(string), d.Get("name").(string)))
	return resourceRdsClusterDatabaseRead(d, meta)
}

func resourceRdsClusterDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseRdsClusterDatabaseID(d.Id())
	if err != nil {
		return err
	}
	db, err := getClient(meta).RdsCluster.GetDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error retrieving rds cluster database %s/%s: %v", instanceID, database, err)
	}
	_ = d.Set("cluster_id", instanceID)
	_ = d.Set("name", db.Name)
	return nil
}

func resourceRdsClusterDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseRdsClusterDatabaseID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).RdsCluster.DeleteDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error deleting rds cluster database %s/%s: %v", instanceID, database, err)
	}
	return nil
}

func resourceRdsClusterDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRdsClusterDatabaseRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildRdsClusterDatabaseID(instanceID string, database string) string {
	return instanceID + "/db/" + database
}

func parseRdsClusterDatabaseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/db/<database>", id)
	}
	return parts[0], parts[2], nil
}
