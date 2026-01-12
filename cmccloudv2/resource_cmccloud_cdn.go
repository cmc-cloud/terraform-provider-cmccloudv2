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
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			originType := diff.Get("origin_type").(string)
			vod := diff.Get("vod").(bool)
			switch originType {
			case "s3":
				if !isSet(diff, "s3_access_key") || !isSet(diff, "s3_secret_key") || !isSet(diff, "s3_bucket_name") || !isSet(diff, "s3_region") || !isSet(diff, "s3_endpoint") {
					return fmt.Errorf("when `origin_type` is 's3', `s3_access_key, s3_secret_key, s3_bucket_name, s3_region, s3_endpoint must be set")
				}
				if !isSet(diff, "domain_or_ip") || !isSet(diff, "protocol") || !isSet(diff, "port") || !isSet(diff, "origin_path") {
					return fmt.Errorf("when `origin_type` is 's3', `domain_or_ip, protocol, port, origin_path must not be set")
				}
			case "host":
				if isSet(diff, "s3_access_key") || isSet(diff, "s3_secret_key") || isSet(diff, "s3_bucket_name") || isSet(diff, "s3_region") || isSet(diff, "s3_endpoint") {
					return fmt.Errorf("when `origin_type` is 'host', `s3_access_key, s3_secret_key, s3_bucket_name, s3_region, s3_endpoint must not be set")
				}
				if !isSet(diff, "domain_or_ip") || !isSet(diff, "protocol") || !isSet(diff, "port") || !isSet(diff, "origin_path") {
					return fmt.Errorf("when `origin_type` is 'host', `domain_or_ip, protocol, port, origin_path must be set")
				}
			}
			if vod {
				if isSet(diff, "origin_path") {
					return fmt.Errorf("when `vod` is 'true', `origin_path must not be set")
				}
			}

			return nil
		},
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
	cdnID, err := getClient(meta).CDN.Create(params)

	if err != nil {
		return fmt.Errorf("error creating cdn site: %s", err)
	}

	if cdnID == "" {
		return fmt.Errorf("error creating cdn site")
	}
	d.SetId(cdnID)

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
