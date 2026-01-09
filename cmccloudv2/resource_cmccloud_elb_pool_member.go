package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELBPoolMember() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBPoolMemberCreate,
		Read:   resourceELBPoolMemberRead,
		Update: resourceELBPoolMemberUpdate,
		Delete: resourceELBPoolMemberDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBPoolMemberImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        elbPoolMemberSchema(),
	}
}

func resourceELBPoolMemberCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":            d.Get("name").(string),
		"address":         d.Get("address").(string),
		"protocol_port":   d.Get("protocol_port").(int),
		"weight":          d.Get("weight").(int),
		"subnet_id":       d.Get("subnet_id").(string),
		"monitor_address": d.Get("monitor_address").(string),
		"backup":          d.Get("backup").(bool),
	}
	// truong nay lay rieng, neu ko gia tri default = 0 se bi loi ko thuoc rang 0 65535
	if monitor_port, ok := d.GetOk("monitor_port"); ok && monitor_port.(int) > 0 {
		params["monitor_port"] = monitor_port.(int)
	}

	elb_id, err := getElbIdFromPool(meta, d.Get("pool_id").(string))
	if err != nil {
		return err
	}

	// doi cho den khi elb het pending update
	err = waitUntilELBEditable(elb_id, d, meta)
	if err != nil {
		return err
	}

	member, err := client.ELB.CreatePoolMember(d.Get("pool_id").(string), params)
	if err != nil {
		return fmt.Errorf("error creating ELB Pool Member: %s", err)
	}
	d.SetId(member.ID)
	// _, err = waitUntilELBPoolMemberStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	// if err != nil {
	// 	return fmt.Errorf("error creating ELB Pool Member: %s", err)
	// }
	return resourceELBPoolMemberRead(d, meta)
}

func resourceELBPoolMemberUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":            d.Get("name").(string),
		"weight":          d.Get("weight").(int),
		"subnet_id":       d.Get("subnet_id").(string),
		"monitor_address": d.Get("monitor_address").(string),
		"backup":          d.Get("backup").(bool),
	}
	// truong nay lay rieng, neu ko gia tri default = 0 se bi loi ko thuoc rang 0 65535
	if monitor_port, ok := d.GetOk("monitor_port"); ok && monitor_port.(int) > 0 {
		params["monitor_port"] = monitor_port.(int)
	}
	_, err := client.ELB.UpdatePoolMember(d.Get("pool_id").(string), d.Id(), params)
	if err != nil {
		return fmt.Errorf("error updating ELB Pool Member: %s", err)
	}
	_, err = waitUntilELBPoolMemberStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error updating ELB Pool Member: %s", err)
	}
	return resourceELBPoolMemberRead(d, meta)
}

func resourceELBPoolMemberRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	member, err := client.ELB.GetPoolMember(d.Get("pool_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ELB Pool Member %s: %v", d.Id(), err)
	}

	_ = d.Set("name", member.Name)
	_ = d.Set("address", member.Address)
	_ = d.Set("protocol_port", member.ProtocolPort)
	_ = d.Set("weight", member.Weight)
	_ = d.Set("subnet_id", member.SubnetID)
	_ = d.Set("monitor_address", member.MonitorAddress)
	_ = d.Set("monitor_port", member.MonitorPort)
	_ = d.Set("backup", member.Backup)
	_ = d.Set("created_at", member.CreatedAt)
	_ = d.Set("operating_status", member.OperatingStatus)
	_ = d.Set("provisioning_status", member.ProvisioningStatus)

	return nil
}

func resourceELBPoolMemberDelete(d *schema.ResourceData, meta interface{}) error {
	elb_id, err := getElbIdFromPool(meta, d.Get("pool_id").(string))
	if err != nil {
		return err
	}

	err = waitUntilELBEditable(elb_id, d, meta)
	if err != nil {
		return err
	}

	_, err = getClient(meta).ELB.DeletePoolMember(d.Get("pool_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete ELB Pool Member: %v", err)
	}
	_, err = waitUntilELBPoolMemberDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete ELB Pool Member: %v", err)
	}
	return nil
}

func resourceELBPoolMemberImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceELBPoolMemberRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilELBPoolMemberDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetPoolMember(d.Get("pool_id").(string), d.Id())
	})
}

func waitUntilELBPoolMemberStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.GetPoolMember(d.Get("pool_id").(string), d.Id())
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELBPoolMember).ProvisioningStatus
	})
}
