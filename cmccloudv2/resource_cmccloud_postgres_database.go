package cmccloudv2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePostgresDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourcePostgresDatabaseCreate,
		Read:   resourcePostgresDatabaseRead,
		Update: resourcePostgresDatabaseUpdate,
		Delete: resourcePostgresDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePostgresDatabaseImport,
		},
		SchemaVersion: 1,
		Schema:        postgresDatabaseSchema(),
	}
}

func resourcePostgresDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	instanceID := d.Get("instance_id").(string)
	databaseName := d.Get("name").(string)
	owner := d.Get("owner").(string)
	body := map[string]interface{}{
		"databaseName": databaseName,
		"owner":        owner,
	}
	params := map[string]interface{}{
		"command": "create_database",
		"body":    body,
	}

	_, err := getClient(meta).PostgresInstance.CreateDatabase(instanceID, params)
	if err != nil {
		return fmt.Errorf("error creating postgres user: %v", err)
	}
	d.SetId(buildPostgresDatabaseID(instanceID, databaseName))
	return resourcePostgresDatabaseRead(d, meta)
}

func resourcePostgresDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parsePostgresDatabaseID(d.Id())
	if err != nil {
		return err
	}
	db, err := getClient(meta).PostgresInstance.GetDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error retrieving postgres user %s/%s: %v", instanceID, database, err)
	}
	_ = d.Set("id", buildPostgresDatabaseID(instanceID, database))
	_ = d.Set("instance_id", instanceID)
	_ = d.Set("name", db.Name)
	_ = d.Set("owner", db.Owner)
	return nil
}

func resourcePostgresDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parsePostgresDatabaseID(d.Id())
	if err != nil {
		return err
	}
	if d.HasChange("owner") {
		params := map[string]interface{}{
			"databaseName": database,
			"owner":        d.Get("owner").(string),
		}
		if len(params) > 0 {
			_, err := getClient(meta).PostgresInstance.UpdateDatabase(instanceID, params)
			if err != nil {
				return fmt.Errorf("error updating postgres user %s/%s: %v", instanceID, database, err)
			}
		}
	}
	return resourcePostgresDatabaseRead(d, meta)
}

func resourcePostgresDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	instanceID, database, err := parsePostgresDatabaseID(d.Id())
	if err != nil {
		return err
	}
	_, err = getClient(meta).PostgresInstance.DeleteDatabase(instanceID, database)
	if err != nil {
		return fmt.Errorf("error deleting postgres user %s/%s: %v", instanceID, database, err)
	}
	return nil
}

func resourcePostgresDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourcePostgresDatabaseRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func buildPostgresDatabaseID(instanceID string, database string) string {
	return instanceID + "/db/" + database
}

func parsePostgresDatabaseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 3 || parts[0] == "" || parts[2] == "" {
		return "", "", fmt.Errorf("invalid id `%s`, expected format: <instance_id>/db/<database>", id)
	}
	return parts[0], parts[2], nil
}
