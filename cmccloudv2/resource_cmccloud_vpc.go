package cmccloudv2

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCCreate,
		Read:   resourceVPCRead,
		Update: resourceVPCUpdate,
		Delete: resourceVPCDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVPCImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        vpcSchema(),
	}
}

func resourceVPCCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc, err := client.VPC.Create(map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"billing_mode": d.Get("billing_mode").(string),
		"cidr":         d.Get("cidr").(string),
		"tags":         d.Get("tags").(*schema.Set).List(),
	})
	if err != nil {
		return fmt.Errorf("Error creating VPC: %s", err)
	}
	d.SetId(vpc.ID)
	return resourceVPCRead(d, meta)
}

func resourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc, err := client.VPC.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving VPC %s: %v", d.Id(), err)
	}

	_ = d.Set("id", vpc.ID)
	_ = d.Set("name", vpc.Name)
	_ = d.Set("description", vpc.Description)
	_ = d.Set("billing_mode", vpc.BillingMode)
	_ = d.Set("cidr", vpc.Cidr)
	_ = d.Set("created", vpc.CreatedAt)
	_ = d.Set("tags", convertTagsToSet(vpc.Tags))
	// gocmcapiv2.Logo("tags", vpc.Tags)
	return nil
}

func resourceVPCUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("tags") {
		_, err := client.VPC.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"tags":        d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("Error when update VPC [%s]: %v", id, err)
		}
	}

	if d.HasChange("cidr") {
		return errors.New("The 'cidr' field cannot be changed after creation")
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetVPCBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when update billing mode of VPC [%s]: %v", id, err)
		}
	}
	return resourceVPCRead(d, meta)
}

func resourceVPCDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.VPC.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete vpc: %v", err)
	}
	_, err = waitUntilVPCDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete vpc: %v", err)
	}
	return nil
}

func resourceVPCImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVPCRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVPCDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).VPC.Get(id)
	})
}
