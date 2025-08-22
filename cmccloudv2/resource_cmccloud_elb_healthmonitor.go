package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBHealthMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBHealthMonitorCreate,
		Read:   resourceELBHealthMonitorRead,
		Update: resourceELBHealthMonitorUpdate,
		Delete: resourceELBHealthMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBHealthMonitorImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        elbhealthmonitorSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			_type := diff.Get("type").(string)
			if _type != "HTTP" && _type != "HTTPS" {
				if isSet(diff, "http_method") || isSet(diff, "expected_codes") || isSet(diff, "url_path") {
					return fmt.Errorf("`http_method`, `expected_codes`, `url_path` can be set only when `type` is HTTP or HTTPS ")
				}
			}
			if _type == "HTTP" {
				if !isSet(diff, "http_method") || !isSet(diff, "expected_codes") {
					return fmt.Errorf("`http_method`, `expected_codes` must be set when `type` is HTTP or HTTPS ")
				}
			}

			return nil
		},
	}
}
func resourceELBHealthMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	htype := d.Get("type").(string)
	params := map[string]interface{}{
		"name":             d.Get("name").(string),
		"pool_id":          d.Get("pool_id").(string),
		"type":             d.Get("type").(string),
		"max_retries_down": d.Get("max_retries_down").(int),
		"delay":            d.Get("delay").(int),
		"max_retries":      d.Get("max_retries").(int),
		"timeout":          d.Get("timeout").(int),
		"domain_name":      d.Get("domain_name").(string),
	}
	if htype == "HTTP" || htype == "HTTPS" {
		params["http_method"] = d.Get("http_method").(string)
		params["expected_codes"] = d.Get("expected_codes").(string)
		params["url_path"] = d.Get("url_path").(string)
	}
	// truoc khi tao health monitor can doi ELB het pending
	// err := waitUntilELBEditable(d, meta)
	// if err != nil {
	// 	return err
	// }
	healthmonitor, err := getClient(meta).ELB.CreateHealthMonitor(params)
	if err != nil {
		return fmt.Errorf("error creating ELB HealthMonitor: %s", err)
	}
	d.SetId(healthmonitor.ID)
	_, err = waitUntilELBHealthMonitorStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ELB HealthMonitor: %s", err)
	}
	return resourceELBHealthMonitorRead(d, meta)
}

func resourceELBHealthMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	htype := d.Get("type").(string)
	params := map[string]interface{}{
		"name":             d.Get("name").(string),
		"type":             d.Get("type").(string),
		"max_retries_down": d.Get("max_retries_down").(int),
		"delay":            d.Get("delay").(int),
		"max_retries":      d.Get("max_retries").(int),
		"timeout":          d.Get("timeout").(int),
		"domain_name":      d.Get("domain_name").(string),
	}
	if htype == "HTTP" || htype == "HTTPS" {
		params["http_method"] = d.Get("http_method").(string)
		params["expected_codes"] = d.Get("expected_codes").(string)
		params["url_path"] = d.Get("url_path").(string)
	}

	// truoc khi tao health monitor can doi ELB het pending
	// err := waitUntilELBEditable(d, meta)
	// if err != nil {
	// 	return err
	// }
	_, err := getClient(meta).ELB.UpdateHealthMonitor(d.Id(), params)

	if err != nil {
		return fmt.Errorf("error updating ELB HealthMonitor: %s", err)
	}
	_, err = waitUntilELBHealthMonitorStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error updating ELB HealthMonitor: %s", err)
	}
	return resourceELBHealthMonitorRead(d, meta)
}

func resourceELBHealthMonitorRead(d *schema.ResourceData, meta interface{}) error {
	healthmonitor, err := getClient(meta).ELB.GetHealthMonitor(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ELBHealthMonitor %s: %v", d.Id(), err)
	}
	if len(healthmonitor.Pools) > 0 {
		_ = d.Set("pool_id", healthmonitor.Pools[0].ID)
	}
	_ = d.Set("name", healthmonitor.Name)
	_ = d.Set("type", healthmonitor.Type)
	_ = d.Set("http_method", healthmonitor.HTTPMethod)
	_ = d.Set("expected_codes", healthmonitor.ExpectedCodes)
	_ = d.Set("max_retries_down", healthmonitor.MaxRetriesDown)
	_ = d.Set("delay", healthmonitor.Delay)
	_ = d.Set("max_retries", healthmonitor.MaxRetries)
	_ = d.Set("timeout", healthmonitor.Timeout)
	_ = d.Set("domain_name", healthmonitor.DomainName)
	_ = d.Set("operating_status", healthmonitor.OperatingStatus)
	_ = d.Set("provisioning_status", healthmonitor.ProvisioningStatus)
	// _ = d.Set("url_path", healthmonitor.URLPath)
	setString(d, "url_path", healthmonitor.URLPath)

	return nil
}

func resourceELBHealthMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).ELB.DeleteHealthMonitor(d.Id())

	if err != nil {
		return fmt.Errorf("error delete ELB HealthMonitor: %v", err)
	}
	_, err = waitUntilELBHealthMonitorDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete ELB HealthMonitor: %v", err)
	}
	return nil
}

func resourceELBHealthMonitorImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceELBHealthMonitorRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilELBHealthMonitorDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetHealthMonitor(id)
	})
}

func waitUntilELBHealthMonitorStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetHealthMonitor(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELBHealthMonitor).ProvisioningStatus
	})
}
