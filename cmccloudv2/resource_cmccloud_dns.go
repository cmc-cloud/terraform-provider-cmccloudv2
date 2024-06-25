package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsCreate,
		Read:   resourceDnsRead,
		// Update: resourceDnsUpdate,
		Delete: resourceDnsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        dnsZoneSchema(),
	}
}

func resourceDnsCreate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	account, err := client.Account.Get()
	if err != nil {
		return fmt.Errorf("Error getting account info: %s", err)
	}
	params := map[string]interface{}{
		"zone": d.Get("domain").(string),
		"type": d.Get("type").(string),
	}
	params["user_id"] = account.ID
	zone, err := getClient(meta).Dns.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating zone: %s", err)
	}
	d.SetId(zone.ID)

	return resourceDnsRead(d, meta)
}

func resourceDnsRead(d *schema.ResourceData, meta interface{}) error {
	zone, err := getClient(meta).Dns.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Dns %s: %v", d.Id(), err)
	}

	_ = d.Set("id", zone.ID)
	_ = d.Set("domain", zone.Zone)
	_ = d.Set("type", zone.Type)
	_ = d.Set("created_at", zone.CreatedAt)
	_ = d.Set("updated_at", zone.UpdatedAt)

	return nil
}

func resourceDnsDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).Dns.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete zone: %v", err)
	}
	_, err = waitUntilDnsDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete zone: %v", err)
	}
	return nil
}

func resourceDnsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDnsRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDnsDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Dns.Get(id)
	})
}
