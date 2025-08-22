package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func keymanagementtokenSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"container_ids": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateUUID,
			},
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"expiration": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"token": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceKeyManagementToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeyManagementTokenCreate,
		Read:   resourceKeyManagementTokenRead,
		Update: resourceKeyManagementTokenUpdate,
		Delete: resourceKeyManagementTokenDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeyManagementTokenImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Create: schema.DefaultTimeout(2 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        keymanagementtokenSchema(),
	}
}

func resourceKeyManagementTokenUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("expiration") {
		_, err := getClient(meta).KeyManagement.RenewToken(d.Id(), d.Get("expiration").(string))
		if err != nil {
			return fmt.Errorf("error renewing KeyManagement Token %s: %v", d.Id(), err)
		}
	}
	return resourceKeyManagementTokenRead(d, meta)
}

func resourceKeyManagementTokenCreate(d *schema.ResourceData, meta interface{}) error {
	container_ids := d.Get("container_ids").([]interface{})
	containerUUids := make([]interface{}, len(container_ids))
	for index, container_id := range container_ids {
		containerUUids[index] = map[string]interface{}{
			"containerUuid": container_id.(string),
		}
	}
	params := map[string]interface{}{
		"containerUuids": containerUUids,
	}
	if d.Get("expiration").(string) != "" {
		params["expiredDate"] = d.Get("expiration").(string)
	}
	if d.Get("description").(string) != "" {
		params["description"] = d.Get("description").(string)
	}
	token, err := getClient(meta).KeyManagement.CreateToken(params)
	if err != nil {
		return fmt.Errorf("error creating KeyManagement Token: %s", err)
	}
	d.SetId(token.Data.ID)
	return resourceKeyManagementTokenRead(d, meta)
}

func resourceKeyManagementTokenRead(d *schema.ResourceData, meta interface{}) error {
	token, err := getClient(meta).KeyManagement.GetToken(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving KeyManagement Token %s: %v", d.Id(), err)
	}

	_ = d.Set("token", token.Token)
	setString(d, "description", token.Description)
	setString(d, "expiration", token.ExpireDateTime)
	_ = d.Set("created_at", token.CreatedTime)
	_ = d.Set("token", token.Token)
	return nil
}

func resourceKeyManagementTokenDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).KeyManagement.DeleteToken(d.Id())
	if err != nil {
		return fmt.Errorf("error delete KeyManagement Token: %v", err)
	}
	return nil
}

func resourceKeyManagementTokenImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKeyManagementTokenRead(d, meta)
	return []*schema.ResourceData{d}, err
}

// func waitUntilKeyManagementTokenDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
// 	return waitUntilResourceDeleted(d, meta, WaitConf{
// 		Delay:      10 * time.Second,
// 		MinTimeout: 20 * time.Second,
// 	}, func(id string) (any, error) {
// 		return getClient(meta).KeyManagement.GetToken(id)
// 	})
// }
