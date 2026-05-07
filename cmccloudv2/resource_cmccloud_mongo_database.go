package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMongoDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDatabaseCreate,
		Read:   resourceMongoDatabaseRead,
		Delete: resourceMongoDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDatabaseImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        mongoDatabaseSchema(),
	}
}

func resourceMongoDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	databaseName := d.Get("name").(string)
	params := map[string]interface{}{
		"databaseName": databaseName,
		"collections":  getStringArrayFromTypeSet(d.Get("collections").(*schema.Set)),
	}
	_, err := getClient(meta).MongoInstance.CreateDatabase(instanceID, params)
	if err != nil {
		return fmt.Errorf("error creating mongo database: %v", err)
	}
	_, err = waitUntilDatabaseFound(d, meta, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating mongo database %s/%s: %v", d.Get("instance_id").(string), d.Get("name").(string), err)
	}
	d.SetId(buildMongoDatabaseID(instanceID, databaseName))
	return resourceMongoDatabaseRead(d, meta)
}

func resourceMongoDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseMongoDatabaseID(d.Id())
	if err != nil {
		return err
	}
	db, err := getClient(meta).DBv2.GetDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error retrieving mongo database %s/%s: %v", instanceID, database, err)
	}
	_ = d.Set("id", buildMongoDatabaseID(instanceID, database))
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("name", db.Name)
	// _ = d.Set("collections", db.Owner)
	return nil
}

func resourceMongoDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parseMongoDatabaseID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).MongoInstance.DeleteDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error deleting mongo database %s/%s: %v", instanceID, database, err)
	}
	_, err = waitUntilDatabaseDeleted(d, meta, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error deleting mongo database %s/%s: %v", instanceID, database, err)
	}
	return nil
}

func resourceMongoDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceMongoDatabaseRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildMongoDatabaseID(instanceID string, database string) string {
	return instanceID + "/db/" + database
}

func parseMongoDatabaseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/db/<database>", id)
	}
	return parts[0], parts[2], nil
}
