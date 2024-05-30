package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ecsgroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"policy": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"soft-anti-affinity", "soft-affinity"}, true),
			Optional:     true,
			Default:      "soft-anti-affinity",
			ForceNew:     true,
		},
	}
}

func resourceEcsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceEcsGroupCreate,
		Read:   resourceEcsGroupRead,
		// Update: resourceEcsGroupUpdate,
		Delete: resourceEcsGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceEcsGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        ecsgroupSchema(),
	}
}

func resourceEcsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc, err := client.EcsGroup.Create(map[string]interface{}{
		"name":   d.Get("name").(string),
		"policy": d.Get("policy").(string),
	})
	if err != nil {
		return fmt.Errorf("Error creating EcsGroup: %s", err)
	}
	d.SetId(vpc.ID)
	return resourceEcsGroupRead(d, meta)
}

func resourceEcsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	vpc, err := client.EcsGroup.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving EcsGroup %s: %v", d.Id(), err)
	}

	_ = d.Set("id", vpc.ID)
	_ = d.Set("name", vpc.Name)
	_ = d.Set("policy", vpc.Policy)
	return nil
}

func resourceEcsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceEcsGroupRead(d, meta)
}

func resourceEcsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.EcsGroup.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud vpc: %v", err)
	}
	return nil
}

func resourceEcsGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceEcsGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}
