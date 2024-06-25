package cmccloudv2

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWafCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafCertCreate,
		Read:   resourceWafCertRead,
		// Update: resourceWafCertUpdate,
		Delete: resourceWafCertDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafCertImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        wafcertSchema(),
	}
}

func resourceWafCertCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":        d.Get("name").(string),
		"cert_data":   base64.StdEncoding.EncodeToString([]byte(d.Get("cert_data").(string))),
		"key_data":    base64.StdEncoding.EncodeToString([]byte(d.Get("key_data").(string))),
		"key_name":    d.Get("key_name").(string),
		"cert_name":   d.Get("cert_name").(string),
		"description": d.Get("description").(string),
	}
	cert, err := getClient(meta).WafCert.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating waf cert: %s", err)
	}
	d.SetId(cert.ID)
	return resourceWafCertRead(d, meta)
}

func resourceWafCertRead(d *schema.ResourceData, meta interface{}) error {
	cert, err := getClient(meta).WafCert.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving waf cert %s: %v", d.Id(), err)
	}

	timestamp := int64(cert.Created)
	t := time.Unix(timestamp, 0)
	// Định dạng thành chuỗi YYYY-MM-DD HH:mm:ss
	formattedTime := t.Format("2006-01-02 15:04:05")

	_ = d.Set("id", cert.ID)
	_ = d.Set("name", cert.Name)
	_ = d.Set("cert_name", cert.CertName)
	_ = d.Set("key_name", cert.KeyName)
	// _ = d.Set("cert_data", cert.CertData)
	// _ = d.Set("key_data", cert.KeyData)
	_ = d.Set("description", cert.Description)
	_ = d.Set("created_at", formattedTime)

	return nil
}

func resourceWafCertDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).WafCert.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete waf cert: %v", err)
	}
	_, err = waitUntilWafCertDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete waf cert: %v", err)
	}
	return nil
}

func resourceWafCertImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceWafCertRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilWafCertDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).WafCert.Get(id)
	})
}
