package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMysqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceMysqlDatabaseCreate,
		Read:   resourceMysqlDatabaseRead,
		Delete: resourceMysqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMysqlDatabaseImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mysqlDatabaseSchema(),
	}
}

func resourceMysqlDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"databaseName":  d.Get("name").(string),
		"character_set": d.Get("character_set").(string),
		"collation":     d.Get("collation").(string),
	}
	_, err := getClient(meta).MysqlInstance.CreateDatabase(d.Get("instance_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating mysql database: %v", err)
	}

	_, err = waitUntilDatabaseFound(d, meta, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating mysql database %s/%s: %v", d.Get("instance_id").(string), d.Get("name").(string), err)
	}
	d.SetId(buildMysqlDatabaseID(d.Get("instance_id").(string), d.Get("name").(string)))
	return resourceMysqlDatabaseRead(d, meta)
}

func resourceMysqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseMysqlDatabaseID(d.Id())
	if err != nil {
		return err
	}
	db, err := getClient(meta).DBv2.GetDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error retrieving mysql database %s/%s: %v", instanceID, database, err)
	}
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("name", db.Name)
	return nil
}

func resourceMysqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseMysqlDatabaseID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).MysqlInstance.DeleteDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error deleting mysql database %s/%s: %v", instanceID, database, err)
	}
	_, err = waitUntilDatabaseDeleted(d, meta, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error deleting mysql database %s/%s: %v", instanceID, database, err)
	}
	return nil
}

func resourceMysqlDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMysqlDatabaseRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildMysqlDatabaseID(instanceID string, database string) string {
	return instanceID + "/db/" + database
}

func parseMysqlDatabaseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/db/<database>", id)
	}
	return parts[0], parts[2], nil
}
