package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseInstanceCreate,
		Read:   resourceDatabaseInstanceRead,
		Update: resourceDatabaseInstanceUpdate,
		Delete: resourceDatabaseInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDatabaseInstanceImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        databaseinstanceSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			sourceType := diff.Get("source_type").(string)
			_, sourceIDSet := diff.GetOk("source_id")
			if sourceType == "new" {
				if sourceIDSet {
					return fmt.Errorf("When source_type is 'new', 'source_id' must not be set")
				}
			} else {
				if !sourceIDSet {
					return fmt.Errorf("When source_type is not 'new', 'source_id' must be set")
				}
			}
			return nil
		},
	}
}

func resourceDatabaseInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.DatabaseInstance.Create(map[string]interface{}{
		"name":              d.Get("name").(string),
		"billing_mode":      d.Get("billing_mode").(string),
		"zone":              d.Get("zone").(string),
		"volume_type":       d.Get("volume_type").(string),
		"volume_size":       d.Get("volume_size").(int),
		"datastore_type":    d.Get("datastore_type").(string),
		"datastore_version": d.Get("datastore_version").(string),
		"admin_user":        d.Get("admin_user").(string),
		"admin_password":    d.Get("admin_password").(string),
		"flavor_id":         d.Get("flavor_id").(string),
		"enable_public_ip":  d.Get("enable_public_ip").(bool),
		"is_public":         d.Get("is_public").(bool),
		"allowed_cidrs":     d.Get("allowed_cidrs").([]interface{}),
		"allowed_host":      d.Get("allowed_host").(string),
		"source_type":       d.Get("source_type").(string),
		"source_id":         d.Get("source_id").(string),
		"replicate_count":   d.Get("replicate_count").(int),
		"subnets":           d.Get("subnets").([]interface{}),
	})
	if err != nil {
		return fmt.Errorf("Error creating Database Instance: %s", err)
	}
	d.SetId(configuration.ID)
	_, err = waitUntilDatabaseInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Database Instance: %s", err)
	}
	return resourceDatabaseInstanceRead(d, meta)
}

func resourceDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.DatabaseInstance.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Database Instance %s: %v", d.Id(), err)
	}

	_ = d.Set("id", instance.ID)
	_ = d.Set("name", instance.Name)
	// _ = d.Set("zone", instance.)
	_ = d.Set("billing_mode", instance.BillingMode)
	// _ = d.Set("volume_type",       instance.Volume)
	_ = d.Set("volume_size", instance.Volume.Size)
	_ = d.Set("datastore_type", instance.Datastore.Type)
	_ = d.Set("datastore_version", instance.Datastore.Version)
	// _ = d.Set("admin_user",        instance.A)
	// _ = d.Set("admin_password",    instance.)
	_ = d.Set("flavor_id", instance.Flavor.ID)
	// _ = d.Set("enable_public_ip",  instance.)
	_ = d.Set("is_public", instance.Access.IsPublic)

	// _ = d.Set("allowed_cidrs",     instance.A)
	// _ = d.Set("allowed_host",      instance.)
	// _ = d.Set("subnets",           instance.)
	return nil
}

func resourceDatabaseInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") {
		_, err := client.DatabaseInstance.Update(id, d.Get("parameters").(map[string]interface{}))
		if err != nil {
			return fmt.Errorf("Error when update info of Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("is_public") || d.HasChange("allowed_cidrs") {
		_, err := client.DatabaseInstance.UpdateInstanceAccessbility(id, map[string]interface{}{
			"is_public":     d.Get("is_public").(bool),
			"allowed_cidrs": d.Get("allowed_cidrs").([]interface{}),
		})
		if err != nil {
			return fmt.Errorf("Error when update accessibility of Database Instance [%s]: %v", id, err)
		}
		_, err = waitUntilDatabaseInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when update accessibility of Database Instance [%s]: %v", id, err)
		}
	}
	if d.HasChange("flavor_id") {
		_, err := client.DatabaseInstance.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("Error when resize Database Instance [%s] to flavor [%s]: %v", id, d.Get("flavor_id").(string), err)
		}
		_, err = waitUntilDatabaseInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when resize Database Instance [%s] to flavor [%s]: %v", id, d.Get("flavor_id").(string), err)
		}
	}
	if d.HasChange("volume_size") {
		_, err := client.DatabaseInstance.ResizeVolume(id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("Error when resize volume Database Instance [%s] to new size: %v", id, err)
		}
		_, err = waitUntilDatabaseInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when resize volume Database Instance [%s] to new size: %v", id, err)
		}
	}

	if d.HasChange("datastore_version") {
		_, err := client.DatabaseInstance.UpgradeDatastoreVersion(id, d.Get("datastore_version").(string))
		if err != nil {
			return fmt.Errorf("Error when upgrade datastore version of Database Instance [%s]: %v", id, err)
		}
		_, err = waitUntilDatabaseInstanceJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when upgrade datastore version of Database Instance [%s]: %v", id, err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetDatabaseInstanceBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when update billing mode of Database Instance [%s]: %v", id, err)
		}
	}
	return resourceDatabaseInstanceRead(d, meta)
}

func resourceDatabaseInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.DatabaseInstance.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete database instance: %v", err)
	}
	_, err = waitUntilDatabaseInstanceDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete database instance: %v", err)
	}
	return nil
}

func resourceDatabaseInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDatabaseInstanceRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDatabaseInstanceJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	// Pending:        []string{"BUILD", "NEW", "REBOOT", "RESIZE", "UPGRADE", "PROMOTE", "EJECT", "DETACH", "SHUTDOWN", "BACKUP"},
	// Target:         []string{"ACTIVE", "ERROR"},
	return waitUntilResourceStatusChanged(d, meta, []string{"ACTIVE", "SHUTDOWN"}, []string{"ERROR"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DatabaseInstance.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.DatabaseInstance).Status
	})
}

func waitUntilDatabaseInstanceDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DatabaseInstance.Get(id)
	})
}
