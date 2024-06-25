package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVA() *schema.Resource {
	return &schema.Resource{
		Create: resourceVACreate,
		Read:   resourceVARead,
		// Update: resourceVAUpdate,
		Delete: resourceVADelete,
		Importer: &schema.ResourceImporter{
			State: resourceVAImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        vaSchema(),
	}
}

func resourceVACreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"scan_name":   d.Get("name").(string),
		"scan_type":   d.Get("type").(string),
		"target":      d.Get("target").(string),
		"description": d.Get("description").(string),
	}

	if d.Get("schedule").(string) != "" {
		params["schedule"] = true
		// Phân tích chuỗi thời gian thành đối tượng time.Time

		layout := "2006-01-02 15:04:05"
		t, err := time.Parse(layout, d.Get("schedule").(string))
		if err != nil {
			return fmt.Errorf("Error parsing schedule time:", err)
		}

		// Chuyển đổi đối tượng time.Time thành Unix timestamp
		timestamp := t.Unix()
		params["timestamp"] = timestamp
	} else {
		params["schedule"] = false
	}

	va, err := getClient(meta).VA.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating VA: %s", err)
	}
	d.SetId(va.ID)
	return resourceVARead(d, meta)
}

func resourceVARead(d *schema.ResourceData, meta interface{}) error {
	va, err := getClient(meta).VA.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving VA %s: %v", d.Id(), err)
	}

	_ = d.Set("id", va.ID)
	_ = d.Set("name", va.Name)
	_ = d.Set("type", va.Type)
	setString(d, "schedule", va.Schedule)
	setString(d, "description", va.Description)
	_ = d.Set("status", va.Status)
	_ = d.Set("report_id", va.ReportID)
	_ = d.Set("created_at", va.CreatedAt)

	return nil
}

func resourceVADelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).VA.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete va: %v", err)
	}
	_, err = waitUntilVADeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete va: %v", err)
	}
	return nil
}

func resourceVAImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceVARead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilVADeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).VA.Get(id)
	})
}
