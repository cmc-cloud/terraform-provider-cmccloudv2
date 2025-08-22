package cmccloudv2

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCDN() *schema.Resource {
	return &schema.Resource{
		Create: resourceCDNCreate,
		Read:   resourceCDNRead,
		Update: resourceCDNUpdate,
		Delete: resourceCDNDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCDNImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        cdnSchema(),
	}
}

func resourceCDNCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":     d.Get("name").(string),
		"domain":   d.Get("domain_or_ip").(string),
		"port":     strconv.Itoa(d.Get("port").(int)),
		"protocol": d.Get("protocol").(string),
		"vod":      "false",
	}
	cdn_id, err := getClient(meta).CDN.Create(params)

	if err != nil {
		return fmt.Errorf("error creating cdn site: %s", err)
	}

	if cdn_id == "" {
		return fmt.Errorf("error creating cdn site")
	}
	d.SetId(cdn_id)

	return resourceCDNRead(d, meta)
}

func resourceCDNRead(d *schema.ResourceData, meta interface{}) error {
	cdn, err := getClient(meta).CDN.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving cdn site %s: %v", d.Id(), err)
	}

	domain := strings.ReplaceAll(cdn.OriginServerURL, "https://www.", "")
	domain = strings.ReplaceAll(domain, "http://www.", "")
	domain = strings.ReplaceAll(domain, "https://", "")
	domain = strings.ReplaceAll(domain, "http://", "")
	domain = strings.Split(domain, ":")[0]

	_ = d.Set("id", cdn.ID)
	_ = d.Set("name", cdn.Name)
	_ = d.Set("domain_or_ip", domain)
	_ = d.Set("port", cdn.OriginSetting.Port)
	_ = d.Set("protocol", cdn.OriginSetting.Protocol)
	_ = d.Set("cdn_url", cdn.MultiCdnURL)
	_ = d.Set("multi_cdn_url", cdn.CdnURL)
	_ = d.Set("status", cdn.Status)
	_ = d.Set("updated_at", cdn.UpdatedAt)

	return nil
}

func resourceCDNUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	cdn, err := getClient(meta).CDN.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving cdn site %s: %v", d.Id(), err)
	}
	cdn.Name = d.Get("name").(string)

	params := map[string]interface{}{
		"name":              d.Get("name").(string),
		"origin_server_url": cdn.OriginServerURL,
		"port":              strconv.Itoa(d.Get("port").(int)),
		"protocol":          d.Get("protocol").(string),
		"vod":               "false",
	}
	_, err = client.CDN.Update(id, params)
	if err != nil {
		return fmt.Errorf("error when update cdn site [%s]: %v", id, err)
	}

	return resourceCDNRead(d, meta)
}

func resourceCDNDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).CDN.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete cdn site: %v", err)
	}
	_, err = waitUntilCDNDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete cdn site: %v", err)
	}
	return nil
}

func resourceCDNImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceCDNRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilCDNDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).CDN.Get(id)
	})
}
