package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceEFSCreate,
		Read:   resourceEFSRead,
		Update: resourceEFSUpdate,
		Delete: resourceEFSDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEFSImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        efsSchema(),
	}
}

func resourceEFSCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"billing_mode":  d.Get("billing_mode").(string),
		"name":          d.Get("name").(string),
		"capacity":      d.Get("capacity").(int),
		"subnet_id":     d.Get("subnet_id").(string),
		"description":   d.Get("description").(string),
		"type":          d.Get("type").(string),
		"protocol_type": d.Get("protocol_type").(string),
		"tags":          d.Get("tags").(*schema.Set).List(),
	}
	efs, err := getClient(meta).EFS.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating EFS: %s", err)
	}
	d.SetId(efs.ID)
	return resourceEFSRead(d, meta)
}

func resourceEFSRead(d *schema.ResourceData, meta interface{}) error {
	efs, err := getClient(meta).EFS.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving EFS %s: %v", d.Id(), err)
	}

	_ = d.Set("id", efs.ID)
	_ = d.Set("name", efs.Name)
	_ = d.Set("capacity", efs.Capacity)
	_ = d.Set("subnet_id", efs.SubnetID)
	_ = d.Set("type", efs.Type)
	_ = d.Set("protocol_type", efs.ProtocolType)
	_ = d.Set("endpoint", efs.Endpoint)
	_ = d.Set("shared_path", efs.SharedPath)
	_ = d.Set("command_line", efs.CommandLine)
	_ = d.Set("status", efs.Status)
	_ = d.Set("created_at", efs.CreatedAt)

	// các field optional phải set riêng
	setString(d, "description", efs.Description)
	setString(d, "billing_mode", efs.BillingMode)
	// setTypeSet(d, "tags", efs.Tags)

	return nil
}

func resourceEFSUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	if d.HasChange("name") || d.HasChange("description") || d.HasChange("tags") || d.HasChange("capacity") {
		_, err := getClient(meta).EFS.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"capacity":    d.Get("capacity").(int),
			"tags":        d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("Error when update EFS [%s]: %v", id, err)
		}
	}
	if d.HasChange("billing_mode") {
		_, err := getClient(meta).BillingMode.SetEFSBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when change biling mode of EFS [%s]: %v", id, err)
		}
	}

	return resourceEFSRead(d, meta)
}

func resourceEFSDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).EFS.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete efs: %v", err)
	}
	_, err = waitUntilEFSDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete efs: %v", err)
	}
	return nil
}

func resourceEFSImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEFSRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilEFSDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).EFS.Get(id)
	})
}
