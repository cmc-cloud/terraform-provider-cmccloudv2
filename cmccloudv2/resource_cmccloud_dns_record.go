package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsRecordCreate,
		Read:   resourceDnsRecordRead,
		Update: resourceDnsRecordUpdate,
		Delete: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsRecordImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        dnsRecordSchema(),
	}
}

func resourceDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"domain":           d.Get("domain").(string),
		"type":             d.Get("type").(string),
		"ttl":              d.Get("ttl").(int),
		"loadbalance_type": d.Get("load_balance_type").(string),
		"detail":           flatternRecordIps(d, d.Get("ips").([]interface{})),
	}
	record, err := getClient(meta).DnsRecord.Create(d.Get("zone_id").(string), params)

	if err != nil {
		return fmt.Errorf("error creating record: %s", err)
	}
	d.SetId(record.ID)

	return resourceDnsRecordRead(d, meta)
}

func resourceDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := map[string]interface{}{
		"domain":           d.Get("domain").(string),
		"type":             d.Get("type").(string),
		"ttl":              d.Get("ttl").(int),
		"loadbalance_type": d.Get("load_balance_type").(string),
		"detail":           flatternRecordIps(d, d.Get("ips").([]interface{})),
	}
	_, err := client.WafRule.Update(id, params)
	if err != nil {
		return fmt.Errorf("error when update dns record [%s]: %v", id, err)
	}

	return resourceDnsRecordRead(d, meta)
}
func resourceDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	record, err := getClient(meta).DnsRecord.Get(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving DnsRecord %s: %v", d.Id(), err)
	}

	_ = d.Set("id", record.ID)
	_ = d.Set("domain", record.Domain)
	_ = d.Set("type", record.Type)
	_ = d.Set("ttl", record.TTL)
	_ = d.Set("load_balance_type", record.LoadbalanceType)
	_ = d.Set("created_at", record.CreatedAt)
	_ = d.Set("updated_at", record.UpdatedAt)
	_ = d.Set("ips", convertRecordIps(d, record.Detail))

	return nil
}

func flatternRecordIps(d *schema.ResourceData, ips []interface{}) []map[string]interface{} {
	load_balance_type := d.Get("load_balance_type").(string)
	result := make([]map[string]interface{}, len(ips))
	for i, ip := range ips {
		r := ip.(map[string]interface{})
		// intval, _ := strconv.Atoi(r["weight"].(string))
		result[i] = map[string]interface{}{
			"content": r["ip"].(string),
			"ttl":     d.Get("ttl").(int),
		}
		if load_balance_type != "none" {
			result[i]["weight"] = r["weight"].(int)
		} else {
			result[i]["weight"] = nil
		}
	}
	return result
}
func convertRecordIps(d *schema.ResourceData, ips []gocmcapiv2.DnsRecordIP) []map[string]interface{} {
	load_balance_type := d.Get("load_balance_type").(string)
	result := make([]map[string]interface{}, len(ips))
	for i, ip := range ips {
		result[i] = map[string]interface{}{
			"ip": ip.Content,
		}
		if load_balance_type != "none" {
			result[i]["weight"] = ip.Weight
		}
	}
	return result
}
func resourceDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).DnsRecord.Delete(d.Get("zone_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete record: %v", err)
	}
	_, err = waitUntilDnsRecordDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete record: %v", err)
	}
	return nil
}

func resourceDnsRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceDnsRecordRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilDnsRecordDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).DnsRecord.Get(d.Get("zone_id").(string), id)
	})
}
