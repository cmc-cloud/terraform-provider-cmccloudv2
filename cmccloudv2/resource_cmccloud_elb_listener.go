package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBListenerCreate,
		Read:   resourceELBListenerRead,
		Update: resourceELBListenerUpdate,
		Delete: resourceELBListenerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBListenerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        elblistenerSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			protocol := diff.Get("protocol").(string)
			if protocol == "UDP" || protocol == "SCTP" {
				_, timeout_client_dataSet := diff.GetOkExists("timeout_client_data")
				_, timeout_tcp_inspectSet := diff.GetOkExists("timeout_tcp_inspect")

				if timeout_client_dataSet || timeout_tcp_inspectSet {
					return fmt.Errorf("When protocol is 'UDP' or 'SCTP', 'timeout_client_data' & 'timeout_tcp_inspect' must not be set")
				}
			}
			if protocol != "HTTP" && protocol != "TERMINATED_HTTPS" {
				_, x_forwarded_forSet := diff.GetOkExists("x_forwarded_for")
				_, x_forwarded_portSet := diff.GetOkExists("x_forwarded_port")
				_, x_forwarded_protoSet := diff.GetOkExists("x_forwarded_proto")

				if x_forwarded_forSet || x_forwarded_portSet || x_forwarded_protoSet {
					return fmt.Errorf("'x_forwarded_for' & 'x_forwarded_port' & 'x_forwarded_proto' only avaiable when protocol is HTTP or TERMINATED_HTTPS")
				}
			}
			if protocol == "TERMINATED_HTTPS" {
				default_tls_container_ref, ok := diff.GetOk("default_tls_container_ref")

				if !ok || default_tls_container_ref == "" {
					return fmt.Errorf("'default_tls_container_ref' must be set when protocol is TERMINATED_HTTPS")
				}
			}

			return nil
		},
	}
}

func resourceELBListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":                      d.Get("name").(string),
		"description":               d.Get("description").(string),
		"default_pool_id":           d.Get("default_pool_id").(string),
		"sni_container_refs":        d.Get("sni_container_refs").([]interface{}),
		"default_tls_container_ref": d.Get("default_tls_container_ref").(string),
		"allowed_cidrs":             d.Get("allowed_cidrs").([]interface{}),
		"x_forwarded_for":           d.Get("x_forwarded_for").(bool),
		"x_forwarded_port":          d.Get("x_forwarded_port").(bool),
		"x_forwarded_proto":         d.Get("x_forwarded_proto").(bool),
		// "timeout_member_data":       d.Get("timeout_member_data").(int),
		// "timeout_member_connect":    d.Get("timeout_member_connect").(int),
		// "timeout_tcp_inspect":       d.Get("timeout_tcp_inspect").(int),
		// "connection_limit":       d.Get("connection_limit").(int),
		// "timeout_client_data":       d.Get("timeout_client_data").(int),
	}

	err := waitUntilELBEditable(d.Get("elb_id").(string), d, meta)
	if err != nil {
		return err
	}

	_, err = getClient(meta).ELB.UpdateListener(d.Id(), params)
	if err != nil {
		return fmt.Errorf("Error update ELB Listener: %s", err)
	}
	_, err = waitUntilELBListenerStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error update ELB Listener: %s", err)
	}
	return resourceELBListenerRead(d, meta)
}

func resourceELBListenerCreate(d *schema.ResourceData, meta interface{}) error {
	params := map[string]interface{}{
		"name":                      d.Get("name").(string),
		"description":               d.Get("description").(string),
		"protocol":                  d.Get("protocol").(string),
		"protocol_port":             d.Get("protocol_port").(int),
		"default_pool_id":           d.Get("default_pool_id").(string),
		"sni_container_refs":        d.Get("sni_container_refs").([]interface{}),
		"default_tls_container_ref": d.Get("default_tls_container_ref").(string),
		"allowed_cidrs":             d.Get("allowed_cidrs").([]interface{}),
		"x_forwarded_for":           d.Get("x_forwarded_for").(bool),
		"x_forwarded_port":          d.Get("x_forwarded_port").(bool),
		"x_forwarded_proto":         d.Get("x_forwarded_proto").(bool),
		// "timeout_client_data":       d.Get("timeout_client_data").(int),
		// "timeout_tcp_inspect":       d.Get("timeout_tcp_inspect").(int),
		// "timeout_member_connect":    d.Get("timeout_member_connect").(int),
		// "timeout_member_data":       d.Get("timeout_member_data").(int),
		// "connection_limit":          d.Get("connection_limit").(int),
	}
	// truoc khi tao health monitor can doi ELB het pending
	err := waitUntilELBEditable(d.Get("elb_id").(string), d, meta)
	if err != nil {
		return err
	}

	elblistener, err := getClient(meta).ELB.CreateListener(d.Get("elb_id").(string), params)
	if v, ok := d.GetOk("connection_limit"); ok {
		params["connection_limit"] = v.(int)
	}

	if v, ok := d.GetOk("timeout_client_data"); ok {
		params["timeout_client_data"] = v.(int)
	}

	if v, ok := d.GetOk("timeout_member_connect"); ok {
		params["timeout_member_connect"] = v.(int)
	}

	if v, ok := d.GetOk("timeout_member_data"); ok {
		params["timeout_member_data"] = v.(int)
	}

	if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
		params["timeout_tcp_inspect"] = v.(int)
	}
	if err != nil {
		return fmt.Errorf("Error creating ELB Listener: %s", err)
	}
	d.SetId(elblistener.ID)
	_, err = waitUntilELBListenerStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating ELB Listener: %s", err)
	}
	return resourceELBListenerRead(d, meta)
}

func resourceELBListenerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	elblistener, err := client.ELB.GetListener(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving ELB Listener %s: %v", d.Id(), err)
	}

	if len(elblistener.Loadbalancers) > 0 {
		_ = d.Set("elb_id", elblistener.Loadbalancers[0].ID)
	}
	_ = d.Set("name", elblistener.Name)
	_ = d.Set("description", elblistener.Description)
	_ = d.Set("protocol", elblistener.Protocol)
	_ = d.Set("protocol_port", elblistener.ProtocolPort)
	_ = d.Set("sni_container_refs", elblistener.SniContainerRefs)
	_ = d.Set("default_tls_container_ref", elblistener.DefaultTLSContainerRef)
	_ = d.Set("allowed_cidrs", elblistener.AllowedCidrs)
	_ = d.Set("operating_status", elblistener.OperatingStatus)
	_ = d.Set("provisioning_status", elblistener.ProvisioningStatus)
	_ = d.Set("default_pool_id", elblistener.DefaultPoolID)

	// kiem tra xem co khac voi gia tri hien tai ko, gia tri hien tai co the chi la gia tri default
	setInt(d, "connection_limit", elblistener.ConnectionLimit)
	setInt(d, "timeout_member_connect", elblistener.TimeoutMemberConnect)
	setInt(d, "timeout_member_data", elblistener.TimeoutMemberData)

	protocol := elblistener.Protocol
	if protocol != "UDP" && protocol != "SCTP" {
		setInt(d, "timeout_client_data", elblistener.TimeoutClientData)
		setInt(d, "timeout_tcp_inspect", elblistener.TimeoutTCPInspect)
	}
	if protocol == "HTTP" || protocol == "TERMINATED_HTTPS" {
		if !elblistener.InsertHeaders.IsArray {
			_ = d.Set("x_forwarded_for", elblistener.InsertHeaders.Headers.XForwardedFor)
			_ = d.Set("x_forwarded_port", elblistener.InsertHeaders.Headers.XForwardedPort)
			_ = d.Set("x_forwarded_proto", elblistener.InsertHeaders.Headers.XForwardedProto)
		}
	}

	return nil
}

func resourceELBListenerDelete(d *schema.ResourceData, meta interface{}) error {
	err := waitUntilELBEditable(d.Get("elb_id").(string), d, meta)
	if err != nil {
		return err
	}

	_, err = getClient(meta).ELB.DeleteListener(d.Id())
	if err != nil {
		return fmt.Errorf("Error delete ELB Listener: %v", err)
	}
	_, err = waitUntilELBListenerDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete ELB Listener: %v", err)
	}
	return nil
}

func resourceELBListenerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceELBListenerRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilELBListenerDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetListener(id)
	})
}

func waitUntilELBListenerStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetListener(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELBListener).ProvisioningStatus
	})
}

func getElbIdFromPool(meta interface{}, pool_id string) (string, error) {
	pool, err := getClient(meta).ELB.GetPool(pool_id)
	elb_id := pool.Loadbalancers[0].ID
	if err != nil {
		return "", fmt.Errorf("Error receiving ELB detail: %v", err)
	}
	return elb_id, nil
}
func waitUntilELBEditable(elb_id string, d *schema.ResourceData, meta interface{}) error {
	_, err := waitUntilResourceStatusChanged(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, WaitConf{
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.Get(elb_id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELB).ProvisioningStatus
	})
	if err != nil {
		return fmt.Errorf("ELB is still immutable and cannot be updated: %v", err)
	}
	return nil
}
