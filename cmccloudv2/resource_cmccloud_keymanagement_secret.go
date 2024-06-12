package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKeyManagementSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyManagementSecretCreate,
		Read:   resourceKeyManagementSecretRead,
		// Update: resourceKeyManagementSecretUpdate,
		Delete: resourceKeyManagementSecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyManagementSecretImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        keymanagementsecretSchema(),
	}
}

// func resourceKeyManagementSecretUpdate(d *schema.ResourceData, meta interface{}) error {
// 	return resourceKeyManagementSecretRead(d, meta)
// }

func resourceKeyManagementSecretCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":       d.Get("name").(string),
		"secretType": d.Get("type").(string),
		"content":    d.Get("content").(string),
	}
	if d.Get("expiration").(string) != "" {
		params["expireTime"] = d.Get("expiration").(string)
	}
	secret, err := getClient(meta).KeyManagement.CreateSecret(map[string]interface{}{
		"containerUuid": d.Get("container_id").(string),
		"secretDetails": []interface{}{params},
	})
	if err != nil {
		return fmt.Errorf("Error creating KeyManagement Secret: %s", err)
	}
	d.SetId(secret.Data.Secrets[0].ID)
	return resourceKeyManagementSecretRead(d, meta)
}

func resourceKeyManagementSecretRead(d *schema.ResourceData, meta interface{}) error {
	container, err := getClient(meta).KeyManagement.GetSecret(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving KeyManagement Secret %s: %v", d.Id(), err)
	}

	_ = d.Set("name", container.Name)
	// _ = d.Set("expiration", container.Expiration)
	_ = d.Set("type", container.SecretType)
	// _ = d.Set("content", container.)
	// _ = d.Set("secret_ref", container.SecretRef)
	_ = d.Set("created_at", container.Created)
	return nil
}

func resourceKeyManagementSecretDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).KeyManagement.DeleteSecret(d.Id())
	if err != nil {
		return fmt.Errorf("Error delete KeyManagement Secret: %v", err)
	}
	return nil
}

func resourceKeyManagementSecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKeyManagementSecretRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKeyManagementSecretDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).KeyManagement.GetSecret(id)
	})
}
