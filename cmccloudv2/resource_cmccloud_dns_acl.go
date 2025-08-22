package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDnsAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsAclCreate,
		Read:   resourceDnsAclRead,
		Update: resourceDnsAclUpdate,
		Delete: resourceDnsAclDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsAclImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        dnsAclSchema(),
	}
}

func resourceDnsAclCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"domain":      d.Get("domain").(string),
		"source_ip":   d.Get("source_ip").(string),
		"record_type": d.Get("record_type").(string),
		"action":      d.Get("action").(string),
	}
	acl, err := getClient(meta).DnsAcl.Create(d.Get("zone_id").(string), params)

	if err != nil {
		return fmt.Errorf("error creating acl: %s", err)
	}
	d.SetId(acl.ID)

	return resourceDnsAclRead(d, meta)
}

func resourceDnsAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := map[string]interface{}{
		"domain":      d.Get("domain").(string),
		"source_ip":   d.Get("source_ip").(string),
		"record_type": d.Get("record_type").(string),
		"action":      d.Get("action").(string),
	}
	_, err := client.DnsAcl.Update(d.Get("zone_id").(string), id, params)
	if err != nil {
		return fmt.Errorf("error when update dns acl [%s]: %v", id, err)
	}

	return resourceDnsAclRead(d, meta)
}
func resourceDnsAclRead(d *schema.ResourceData, meta interface{}) error {
	acl, err := getClient(meta).DnsAcl.Get(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving DnsAcl %s: %v", d.Id(), err)
	}

	_ = d.Set("id", acl.ID)
	_ = d.Set("domain", acl.Domain)
	_ = d.Set("source_ip", acl.SourceIP)
	_ = d.Set("record_type", acl.RecordType)
	_ = d.Set("action", acl.Action)
	_ = d.Set("created_at", acl.CreatedAt)
	_ = d.Set("updated_at", acl.UpdatedAt)

	return nil
}

func resourceDnsAclDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).DnsAcl.Delete(d.Get("zone_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete acl: %v", err)
	}
	_, err = waitUntilDnsAclDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete acl: %v", err)
	}
	return nil
}

func resourceDnsAclImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDnsAclRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDnsAclDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DnsAcl.Get(d.Get("zone_id").(string), id)
	})
}
