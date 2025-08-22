package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBPoolCreate,
		Read:   resourceELBPoolRead,
		Update: resourceELBPoolUpdate,
		Delete: resourceELBPoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBPoolImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        elbpoolSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			session_persistence := diff.Get("session_persistence").(string)
			if session_persistence == "APP_COOKIE" {
				cookie_name, ok := diff.GetOk("cookie_name")

				if !ok || cookie_name == "" {
					return fmt.Errorf("'cookie_name' must be set when session_persistence is `APP_COOKIE`")
				}
			}

			return nil
		},
	}
}

func resourceELBPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ELB.UpdatePool(d.Id(), map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"algorithm":           d.Get("algorithm").(string),
		"session_persistence": d.Get("session_persistence").(string),
		"cookie_name":         d.Get("cookie_name").(string),
		"tls_enabled":         d.Get("tls_enabled").(bool),
		"tls_ciphers":         d.Get("tls_ciphers").(string),
		"tls_versions":        d.Get("tls_versions").([]interface{}),
	})
	if err != nil {
		return fmt.Errorf("error updating ELB Pool: %s", err)
	}
	_, err = waitUntilELBPoolStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error updating ELB Pool: %s", err)
	}
	return resourceELBPoolRead(d, meta)
}

func resourceELBPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	elbpool, err := client.ELB.CreatePool(d.Get("elb_id").(string), map[string]interface{}{
		"name":                d.Get("name").(string),
		"description":         d.Get("description").(string),
		"protocol":            d.Get("protocol").(string),
		"algorithm":           d.Get("algorithm").(string),
		"session_persistence": d.Get("session_persistence").(string),
		"cookie_name":         d.Get("cookie_name").(string),
		"tls_enabled":         d.Get("tls_enabled").(bool),
		"tls_ciphers":         d.Get("tls_ciphers").(string),
		"tls_versions":        d.Get("tls_versions").([]interface{}),
	})
	if err != nil {
		return fmt.Errorf("error creating ELB Pool: %s", err)
	}
	d.SetId(elbpool.ID)
	_, err = waitUntilELBPoolStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ELB Pool: %s", err)
	}
	return resourceELBPoolRead(d, meta)
}

func resourceELBPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	elbpool, err := client.ELB.GetPool(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ELBPool %s: %v", d.Id(), err)
	}

	if len(elbpool.Loadbalancers) > 0 {
		_ = d.Set("elb_id", elbpool.Loadbalancers[0].ID)
	}
	if len(elbpool.Listeners) > 0 {
		_ = d.Set("listener_id", elbpool.Listeners[0].ID)
	}
	_ = d.Set("name", elbpool.Name)
	_ = d.Set("description", elbpool.Description)
	_ = d.Set("protocol", elbpool.Protocol)
	_ = d.Set("algorithm", elbpool.LbAlgorithm)
	_ = d.Set("session_persistence", elbpool.SessionPersistence)
	_ = d.Set("cookie_name", elbpool.SessionPersistence.CookieName)
	_ = d.Set("tls_enabled", elbpool.TLSEnabled)
	_ = d.Set("tls_ciphers", elbpool.TLSCiphers)
	_ = d.Set("tls_versions", elbpool.TLSVersions)
	_ = d.Set("operating_status", elbpool.OperatingStatus)
	_ = d.Set("provisioning_status", elbpool.ProvisioningStatus)

	return nil
}

func resourceELBPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ELB.DeletePool(d.Id())

	if err != nil {
		return fmt.Errorf("error delete ELB Pool: %v", err)
	}
	_, err = waitUntilELBPoolDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete ELB Pool: %v", err)
	}
	return nil
}

func resourceELBPoolImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceELBPoolRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilELBPoolDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetPool(id)
	})
}

func waitUntilELBPoolStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetPool(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELBPool).ProvisioningStatus
	})
}
