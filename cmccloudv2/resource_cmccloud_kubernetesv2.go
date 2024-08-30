package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernetesv2() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesv2Create,
		Read:   resourceKubernetesv2Read,
		Update: resourceKubernetesv2Update,
		Delete: resourceKubernetesv2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernetesv2Import,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernetesv2Schema(),
	}
}

func resourceKubernetesv2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("Error receving subnet with id = %s: %v", d.Get("subnet_id").(string), err)
	}
	kubernetes, err := client.Kubernetesv2.Create(map[string]interface{}{
		"region":         client.Configs.RegionId,
		"project":        client.Configs.ProjectId,
		"name":           d.Get("name").(string),
		"zone":           d.Get("zone").(string),
		"subnetId":       d.Get("subnet_id").(string),
		"vpcId":          subnet.NetworkID,
		"version":        d.Get("kubernetes_version").(string),
		"replicas":       d.Get("master_count").(int),
		"masterFlavorId": d.Get("master_flavor_name").(string),
		// "workerNumberEstimate":             d.Get("max_node_count").(int),
		"cidrBlockPod":                     d.Get("cidr_block_pod").(string),
		"cidrBlockService":                 d.Get("cidr_block_service").(string),
		"clusterNetworkServiceDomain":      d.Get("network_driver").(string),
		"clusterNetworkServicesCidrBlocks": "",
		"clusterNetworkApiServerPort":      "",
		"rolloutStrategyType":              "",
		"rolloutStrategyMaxSurge":          "",
	})

	if err != nil {
		return fmt.Errorf("Error creating Kubernetesv2: %s", err)
	}
	d.SetId(kubernetes.Data.ID)

	_, err = waitUntilKubernetesv2StatusChangedState(d, meta, []string{"HEALTHY", "RUNNING", "active", "Ready", "Running"}, []string{"ERROR", "SHUTDOWN", "FAILURE", "failure", "deleting"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Kubernetesv2: %s", err)
	}
	if d.Get("enable_autohealing").(bool) {
		err := updateAutoHealingAddon(d, meta)
		if err != nil {
			return err
		}
	}
	if d.Get("enable_autoscale").(bool) {
		err := updateAutoScaleAddon(d, meta)
		if err != nil {
			return err
		}
	}
	if d.Get("enable_monitoring").(bool) {
		err := updateMonitoringAddon(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceKubernetesv2Read(d, meta)
}

func resourceKubernetesv2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	kubernetes, err := client.Kubernetesv2.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernetesv2 %s: %v", d.Id(), err)
	}

	_ = d.Set("id", kubernetes.ClusterID)
	_ = d.Set("name", kubernetes.ClusterName)
	_ = d.Set("subnetId", kubernetes.SubnetID)
	_ = d.Set("vpcId", kubernetes.VpcID)
	_ = d.Set("cidr_block_pod", kubernetes.CidrBlockPod)
	_ = d.Set("cidr_block_service", kubernetes.CidrBlockService)
	_ = d.Set("kubernetes_version", kubernetes.KubeletVersion)
	_ = d.Set("master_count", kubernetes.NumberMasterNode)
	_ = d.Set("network_driver", kubernetes.ServiceDomain)
	_ = d.Set("created_at", kubernetes.CreatedAt)
	_ = d.Set("state", kubernetes.State)

	status, err := client.Kubernetesv2.GetStatus(d.Id())
	_ = d.Set("enable_autohealing", status.EnableAutoHealing)
	_ = d.Set("enable_autoscale", status.EnableAutoScale)
	_ = d.Set("enable_monitoring", status.EnableMonitor)

	// "zone":                        d.Get("zone").(string),
	// 	"workerNumberEstimate":        d.Get("max_node_count").(int),

	return nil
}

func updateAutoHealingAddon(d *schema.ResourceData, meta interface{}) error {
	action := "disable"
	if d.Get("enable_autohealing").(bool) {
		action = "enable"
	}
	params := map[string]interface{}{
		"action":                action,
		"externalProviderNames": "auto-healing-control-plane",
	}
	getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	_, err := waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error update auto healing addon: %v", err)
	}
	return nil
}

func updateAutoScaleAddon(d *schema.ResourceData, meta interface{}) error {
	action := "disable"
	if d.Get("enable_autoscale").(bool) {
		action = "enable"
	}
	params := map[string]interface{}{
		"action":                action,
		"externalProviderNames": "auto-scale",
		"minNodeCluster":        d.Get("autoscale_min_node").(int),
		"maxNodeCluster":        d.Get("autoscale_max_node").(int),
		"minRamCluster":         d.Get("autoscale_min_ram_gb").(int),
		"maxRamCluster":         d.Get("autoscale_max_ram_gb").(int),
		"minCoreCluster":        d.Get("autoscale_min_core").(int),
		"maxCoreCluster":        d.Get("autoscale_max_core").(int),
	}
	getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	_, err := waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error update autoscale addon: %v", err)
	}
	return nil
}

func updateMonitoringAddon(d *schema.ResourceData, meta interface{}) error {
	action := "disable"
	if d.Get("enable_autohealing").(bool) {
		action = "enable"
	}
	params := map[string]interface{}{
		"action":                action,
		"externalProviderNames": "monitoring",
	}
	getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	_, err := waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error update monitoring addon: %v", err)
	}
	return nil
}
func resourceKubernetesv2Update(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("enable_autohealing") {
		err := updateAutoHealingAddon(d, meta)
		if err != nil {
			return err
		}
	}
	if d.HasChange("enable_autoscale") {
		err := updateAutoScaleAddon(d, meta)
		if err != nil {
			return err
		}
	}
	if d.HasChange("enable_monitoring") {
		err := updateMonitoringAddon(d, meta)
		if err != nil {
			return err
		}
	}
	return resourceKubernetesv2Read(d, meta)
}

func resourceKubernetesv2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernetesv2.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernetesv2Deleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceKubernetesv2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKubernetesv2Read(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKubernetesv2Deleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      20 * time.Second,
		MinTimeout: 3 * 60 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetesv2.Get(id)
	})
}

func waitUntilKubernetesv2StatusChangedStateReady(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilKubernetesv2StatusChangedState(d, meta, []string{"HEALTHY", "RUNNING", "active", "Ready", "Running"}, []string{"ERROR", "SHUTDOWN", "FAILURE", "failure", "deleting"}, d.Timeout(schema.TimeoutCreate))
}
func waitUntilKubernetesv2StatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetesv2.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Kubernetesv2).State
	})
}
