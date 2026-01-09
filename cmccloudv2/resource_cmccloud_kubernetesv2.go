package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
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
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			enable_autoscale := diff.Get("enable_autoscale").(bool)
			if enable_autoscale {
				if !isSet(diff, "autoscale_max_node") || !isSet(diff, "autoscale_max_ram_gb") || !isSet(diff, "autoscale_max_core") {
					return fmt.Errorf("when `enable_autoscale` is 'true', `autoscale_max_node, autoscale_max_ram_gb, autoscale_max_core must be set")
				}
			} else {
				if isSet(diff, "autoscale_max_node") || isSet(diff, "autoscale_max_ram_gb") || isSet(diff, "autoscale_max_core") {
					return fmt.Errorf("when `enable_autoscale` is 'false', `autoscale_max_node, autoscale_max_ram_gb, autoscale_max_core must not be set")
				}
			}

			// driver := diff.Get("network_driver").(string)

			// if driver != "cilium" {
			// 	// Nếu không phải cilium mà user cố tình set mode → báo lỗi
			// 	if isSet(diff, "network_driver_mode") {
			// 		return fmt.Errorf("`network_driver_mode` only avaiable when `network_driver = \"cilium\"`")
			// 	}
			// }
			return nil
		},
	}
}

func resourceKubernetesv2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("error receving subnet with id = %s: %v", d.Get("subnet_id").(string), err)
	}

	params := map[string]interface{}{
		"billingMode":                 d.Get("billing_mode").(string),
		"region":                      client.Configs.RegionId,
		"project":                     client.Configs.ProjectId,
		"name":                        d.Get("name").(string),
		"zone":                        d.Get("zone").(string),
		"subnetId":                    d.Get("subnet_id").(string),
		"vpcId":                       subnet.NetworkID,
		"version":                     d.Get("kubernetes_version").(string),
		"replicas":                    d.Get("master_count").(int),
		"masterFlavorId":              d.Get("master_flavor_name").(string),
		"cidrBlockNode":               subnet.Cidr,
		"cidrBlockPod":                d.Get("cidr_block_pod").(string),
		"cidrBlockService":            d.Get("cidr_block_service").(string),
		"clusterNetworkServiceDomain": d.Get("network_driver").(string),
		"nodeMaskCidr":                d.Get("node_mask_cidr").(int),
		// "isNTP":                            d.Get("ntp_enabled").(bool),
		"ntpServers":                       flatternNtpServers(d),
		"clusterNetworkServicesCidrBlocks": "",
		"clusterNetworkApiServerPort":      "",
		"rolloutStrategyType":              "",
		"rolloutStrategyMaxSurge":          "",
		// "workerNumberEstimate":             d.Get("max_node_count").(int),
	}

	if d.Get("network_driver").(string) == "cilium" {
		params["clusterNetworkDriverMode"] = d.Get("network_driver_mode").(string)
	}
	kubernetes, err := client.Kubernetesv2.Create(params)

	if err != nil {
		return fmt.Errorf("error creating Kubernetesv2: %s", err)
	}
	d.SetId(kubernetes.Data.ID)

	_, err = waitUntilKubernetesv2StatusChangedState(d, meta, []string{"HEALTHY", "RUNNING", "active", "Ready", "Running"}, []string{"ERROR", "SHUTDOWN", "FAILURE", "failure", "deleting"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating Kubernetesv2: %s", err)
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
		return fmt.Errorf("error retrieving Kubernetesv2 %s: %v", d.Id(), err)
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

	_ = d.Set("node_mask_cidr", kubernetes.NodeMaskCidr)
	_ = d.Set("network_driver_mode", kubernetes.ClusterNetworkDriverMode)
	// setInt(d, "node_mask_cidr", kubernetes.NodeMaskCidr)
	// setString(d, "network_driver_mode", kubernetes.ClusterNetworkDriverMode)
	_ = d.Set("ntp_servers", convertNtpServers(kubernetes.NtpServers))

	status, err := client.Kubernetesv2.GetStatus(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Kubernetesv2 status %s: %v", d.Id(), err)
	}
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
	_, err := getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	if err != nil {
		return fmt.Errorf("error update auto healing addon: %v", err)
	}
	_, err = waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error update auto healing addon: %v", err)
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
		// "minNodeCluster":        d.Get("autoscale_min_node").(int),
		"maxNodeCluster": d.Get("autoscale_max_node").(int),
		// "minRamCluster":         d.Get("autoscale_min_ram_gb").(int),
		"maxRamCluster": d.Get("autoscale_max_ram_gb").(int),
		// "minCoreCluster":        d.Get("autoscale_min_core").(int),
		"maxCoreCluster": d.Get("autoscale_max_core").(int),
	}
	_, err := getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	if err != nil {
		return fmt.Errorf("error update autoscale addon: %v", err)
	}
	_, err = waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error update autoscale addon: %v", err)
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
	_, err := getClient(meta).Kubernetesv2.UpdateAddon(d.Id(), params)
	if err != nil {
		return fmt.Errorf("error update monitoring addon: %v", err)
	}
	_, err = waitUntilKubernetesv2StatusChangedStateReady(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error update monitoring addon: %v", err)
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

	if d.HasChange("enable_autoscale") || d.HasChange("autoscale_max_node") || d.HasChange("autoscale_max_ram_gb") || d.HasChange("autoscale_max_core") {
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
		return fmt.Errorf("error delete kubernetes [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernetesv2Deleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete kubernetes [%s]: %v", d.Id(), err)
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

func flatternNtpServers(d *schema.ResourceData) []map[string]interface{} {
	servers := d.Get("ntp_servers").([]interface{})
	result := make([]map[string]interface{}, len(servers))
	for i, server := range servers {
		r := server.(map[string]interface{})
		result[i] = map[string]interface{}{
			"host":     r["host"].(string),
			"protocol": r["protocol"].(string),
			"port":     r["port"].(int),
		}
	}
	return result
}

func convertNtpServers(servers []gocmcapiv2.NtpServer) []map[string]interface{} {
	result := make([]map[string]interface{}, len(servers))
	for i, server := range servers {
		result[i] = map[string]interface{}{
			"host":     server.Host,
			"port":     server.Port,
			"protocol": server.Protocol,
		}
	}
	return result
}
