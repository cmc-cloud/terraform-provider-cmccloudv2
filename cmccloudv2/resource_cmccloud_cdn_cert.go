package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCDNCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceCDNCertCreate,
		Read:   resourceCDNCertRead,
		Update: resourceCDNCertUpdate,
		Delete: resourceCDNCertDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: resourceCDNCertImport,
		// },
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        cdnCertSchema(),
	}
}

func resourceCDNCertCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"cert_data": d.Get("cert_data").(string),
		"key_data":  d.Get("key_data").(string),
		"key_name":  d.Get("key_name").(string),
		"cert_name": d.Get("cert_name").(string),
	}
	cdn, err := getClient(meta).CDNCert.Create(params)

	if err != nil {
		return fmt.Errorf("Error creating cdn cert: %s", err)
	}
	d.SetId(cdn.ID)

	return resourceCDNCertRead(d, meta)
}

func resourceCDNCertRead(d *schema.ResourceData, meta interface{}) error {
	cdn, err := getClient(meta).CDNCert.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving cdn cert %s: %v", d.Id(), err)
	}

	_ = d.Set("id", cdn.ID)
	_ = d.Set("certificate_type", cdn.CertificateType)
	_ = d.Set("common_name", cdn.CommonName)
	_ = d.Set("expiration_date", cdn.ExpirationDate)
	_ = d.Set("status", cdn.Status)

	return nil
}

func resourceCDNCertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	params := map[string]interface{}{
		"cert_data": d.Get("cert_data").(string),
		"key_data":  d.Get("key_data").(string),
		"key_name":  d.Get("key_name").(string),
		"cert_name": d.Get("cert_name").(string),
	}
	_, err := client.CDNCert.Update(id, params)
	if err != nil {
		return fmt.Errorf("Error when update dns cert [%s]: %v", id, err)
	}

	return resourceCDNCertRead(d, meta)
}
func resourceCDNCertDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).CDNCert.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cdn cert: %v", err)
	}
	_, err = waitUntilCDNCertDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete cdn cert: %v", err)
	}
	return nil
}

// func resourceCDNCertImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
// 	err := resourceCDNCertRead(d, meta)
// 	return []*schema.ResourceData{d}, err
// }

func waitUntilCDNCertDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).CDNCert.Get(id)
	})
}
