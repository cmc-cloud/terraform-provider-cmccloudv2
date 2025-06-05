package cmccloudv2

import (
	"fmt"

	// "strconv"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBL7policy() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBL7policyCreate,
		Read:   resourceELBL7policyRead,
		Update: resourceELBL7policyUpdate,
		Delete: resourceELBL7policyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBL7policyImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        lbL7policySchema(),
	}
}

func resourceELBL7policyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	action := d.Get("action").(string)
	datas := map[string]interface{}{
		"listener_id":      d.Get("listener_id").(string),
		"name":             d.Get("name").(string),
		"action":           d.Get("action").(string),
		"position":         d.Get("position").(int),
		"redirect_url":     d.Get("redirect_url").(string),
		"redirect_pool_id": d.Get("redirect_pool_id").(string),
		"redirect_prefix":  d.Get("redirect_prefix").(string),
	}
	if action == "REDIRECT_TO_URL" || action == "REDIRECT_PREFIX" {
		datas["redirect_http_code"] = d.Get("redirect_http_code").(int)
	}
	l7policy, err := client.ELB.CreateL7Policy(datas)
	if err != nil {
		return fmt.Errorf("Error creating L7Policy: %s", err)
	}

	_, err = waitUntilELBL7PolicyStatusChangedState(l7policy.ID, d, meta, []string{"ACTIVE"}, []string{"ERROR", "DELETED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating ELB L7 Policy: %s", err)
	}
	d.SetId(l7policy.ID)
	return resourceELBL7policyRead(d, meta)
}

func resourceELBL7policyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	l7policy, err := client.ELB.GetL7Policy(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading L7Policy: %s", err)
	}

	d.Set("name", l7policy.Name)
	d.Set("action", l7policy.Action)
	d.Set("listener_id", l7policy.ListenerID)
	d.Set("position", l7policy.Position)

	if l7policy.Action == "REDIRECT_TO_POOL" {
		d.Set("redirect_pool_id", l7policy.RedirectPoolID)
	}

	if l7policy.Action == "REDIRECT_PREFIX" {
		d.Set("redirect_prefix", l7policy.RedirectPrefix)
	}

	if l7policy.Action == "REDIRECT_TO_URL" || l7policy.Action == "REDIRECT_PREFIX" {
		d.Set("redirect_url", l7policy.RedirectURL)
	}

	if l7policy.Action == "REDIRECT_TO_URL" || l7policy.Action == "REDIRECT_PREFIX" {
		d.Set("redirect_http_code", l7policy.RedirectHTTPCode)
	}

	return nil
}

func resourceELBL7policyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	action := d.Get("action").(string)
	datas := map[string]interface{}{
		"name":             d.Get("name").(string),
		"action":           action,
		"position":         d.Get("position").(int),
		"redirect_url":     d.Get("redirect_url").(string),
		"redirect_pool_id": d.Get("redirect_pool_id").(string),
		"redirect_prefix":  d.Get("redirect_prefix").(string),
	}
	if action == "REDIRECT_TO_URL" || action == "REDIRECT_PREFIX" {
		datas["redirect_http_code"] = d.Get("redirect_http_code").(int)
	}
	_, err := client.ELB.UpdateL7Policy(d.Id(), datas)
	if err != nil {
		return fmt.Errorf("Error updating L7Policy: %s", err)
	}
	_, err = waitUntilELBL7PolicyStatusChangedState(d.Id(), d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error when change flavor of ELB [%s]: %v", d.Id(), err)
	}
	return resourceELBL7policyRead(d, meta)
}

func resourceELBL7policyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ELB.DeleteL7Policy(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete L7 Policy: %v", err)
	}
	_, err = waitUntilEELBL7PolicyDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete L7 Policy: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceELBL7policyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceELBL7policyRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func waitUntilEELBL7PolicyDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetL7Policy(id)
	})
}

func waitUntilELBL7PolicyStatusChangedState(resource_id string, d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceIdStatusChanged(resource_id, d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetL7Policy(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.L7Policy).ProvisioningStatus
	})
}
