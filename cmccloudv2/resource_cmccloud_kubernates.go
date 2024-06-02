package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernates() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernatesCreate,
		Read:   resourceKubernatesRead,
		Update: resourceKubernatesUpdate,
		Delete: resourceKubernatesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernatesImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernetesSchema(),
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if v, ok := d.GetOk("labels"); ok {
				blockList := v.([]interface{})
				if len(blockList) > 0 {
					labels := blockList[0].(map[string]interface{})
					auto_scaling_enabled := labels["auto_scaling_enabled"].(bool)
					if auto_scaling_enabled {
						// co enable => phai set 2 truong nay, khong set => thong bao
						if labels["max_node_count"].(int) <= 0 { // khong duoc set max_node_count
							return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enabled is 'true'")
						}
						if labels["min_node_count"].(int) <= 0 { // khong duoc set min_node_count
							return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enabled is 'true'")
						}
					} else {
						// khong enable => ko set 2 truong nay
						if labels["max_node_count"].(int) > 0 {
							return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enabled is 'false'")
						}
						if labels["min_node_count"].(int) > 0 {
							return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enabled is 'false'")
						}
					}
				}
			}
			return nil
		},
	}
}

func resourceKubernatesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	metas := getFirstBlock(d, "labels")
	labels := map[string]interface{}{
		"kube_dashboard_enabled": metas["kube_dashboard_enabled"].(bool),
		"metrics_server_enabled": metas["metrics_server_enabled"].(bool),
		"npd_enabled":            metas["npd_enabled"].(bool),
		"auto_scaling_enabled":   metas["auto_scaling_enabled"].(bool),
		"auto_healing_enabled":   metas["auto_healing_enabled"].(bool),
		"max_node_count":         metas["max_node_count"].(int),
		"min_node_count":         metas["min_node_count"].(int),
		"kube_tag":               metas["kube_tag"].(string),
		"network-driver":         metas["network_driver"].(string),
		"calico_ipv4pool":        metas["calico_ipv4pool"].(string),
		"docker_volume_type":     d.Get("docker_volume_type").(string),
		"zone":                   d.Get("zone").(string),
	}

	default_master := getFirstBlock(d, "default_master")
	default_worker := getFirstBlock(d, "default_worker")
	params := map[string]interface{}{
		"name": d.Get("name").(string),

		"master_count":        default_master["node_count"].(int),
		"master_flavor_id":    default_master["flavor_id"].(string),
		"master_billing_mode": default_master["billing_mode"].(string),

		"node_count":          default_worker["node_count"].(int),
		"node_flavor_id":      default_worker["flavor_id"].(string),
		"worker_billing_mode": default_worker["billing_mode"].(string),

		"keypair":            d.Get("keypair").(string),
		"docker_volume_size": d.Get("docker_volume_size").(int),
		"subnet_id":          d.Get("subnet_id").(string),
		"create_timeout":     d.Get("create_timeout").(int),
		"zone":               d.Get("zone").(string),
		"labels":             labels,
	}

	kubernetes, err := client.Kubernates.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating Kubernates: %s", err)
	}
	d.SetId(kubernetes.ID)

	_, err = waitUntilKubernatesStatusChangedState(d, meta, []string{"CREATE_COMPLETE", "HEALTHY"}, []string{"CREATE_FAILED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Kubernates: %s", err)
	}
	return resourceKubernatesRead(d, meta)
}

func resourceKubernatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	kubernetes, err := client.Kubernates.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernates %s: %v", d.Id(), err)
	}

	labels := make([]map[string]interface{}, 1)
	labels[0] = map[string]interface{}{
		"kube_dashboard_enabled": kubernetes.Labels.KubeDashboardEnabled,
		"metrics_server_enabled": kubernetes.Labels.MetricsServerEnabled,
		"npd_enabled":            kubernetes.Labels.NpdEnabled,
		"auto_scaling_enabled":   kubernetes.Labels.AutoScalingEnabled,
		"auto_healing_enabled":   kubernetes.Labels.AutoHealingEnabled,

		"kube_tag":           kubernetes.Labels.KubeTag,
		"network_driver":     kubernetes.Labels.NetworkDriver,
		"calico_ipv4pool":    kubernetes.Labels.CalicoIpv4Pool,
		"docker_volume_type": kubernetes.Labels.DockerVolumeType,
		// "create_timeout":        kubernetes.Labels.CreateTimeout,
		"zone": kubernetes.Labels.AvailabilityZone,
	}

	if kubernetes.Labels.AutoScalingEnabled {
		labels[0]["max_node_count"] = kubernetes.Labels.MaxNodeCount
		labels[0]["min_node_count"] = kubernetes.Labels.MinNodeCount
	}

	_ = d.Set("id", kubernetes.ID)
	_ = d.Set("name", kubernetes.Name)
	_ = d.Set("zone", kubernetes.Labels.AvailabilityZone)
	_ = d.Set("subnet_id", kubernetes.SubnetID)
	_ = d.Set("docker_volume_size", kubernetes.DockerVolumeSize)
	_ = d.Set("docker_volume_type", kubernetes.Labels.DockerVolumeType)
	_ = d.Set("keypair", kubernetes.Keypair)
	_ = d.Set("create_timeout", kubernetes.CreateTimeout)

	default_master := map[string]interface{}{
		"node_count":   kubernetes.MasterCount,
		"flavor_id":    kubernetes.MasterFlavorID,
		"billing_mode": kubernetes.MasterBillingMode,
	}
	d.Set("default_master", []interface{}{default_master})

	default_worker := map[string]interface{}{
		"node_count":   kubernetes.NodeCount,
		"flavor_id":    kubernetes.NodeFlavorID,
		"billing_mode": kubernetes.NodeBillingMode,
	}
	d.Set("default_worker", []interface{}{default_worker})

	_ = d.Set("created_at", kubernetes.CreatedAt)
	_ = d.Set("labels", labels)

	return nil
}

func resourceKubernatesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	master_billing_mode_changed, new_master_billing_mode := isSubBlockFieldChanged(d, "default_master", "billing_mode")
	worker_billing_mode_changed, new_worker_billing_mode := isSubBlockFieldChanged(d, "default_worker", "billing_mode")

	if d.HasChange("node_count") {
		_, err := client.Kubernates.UpdateNodeCount(id, d.Get("node_count").(int))
		if err != nil {
			return fmt.Errorf("Error when update Kubernates node count [%s]: %v", id, err)
		}
		_, err = waitUntilKubernatesStatusChangedState(d, meta, []string{"UPDATE_COMPLETE", "HEALTHY"}, []string{"UPDATE_FAILED"}, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when update Kubernates node count [%s]: %v", id, err)
		}
	} else if master_billing_mode_changed {
		_, err := client.BillingMode.SetKubernateBilingMode(id, new_master_billing_mode.(string), "master")
		if err != nil {
			return fmt.Errorf("Error when change default master biling mode of Kubernates cluster [%s]: %v", id, err)
		}
	} else if worker_billing_mode_changed {
		_, err := client.BillingMode.SetKubernateBilingMode(id, new_worker_billing_mode.(string), "worker")
		if err != nil {
			return fmt.Errorf("Error when change default worker biling mode of Kubernates cluster [%s]: %v", id, err)
		}
	} else {
		return fmt.Errorf("Only `node_count`, `billing_mode` fields can be updated after Kubernates cluster created")
	}
	return resourceKubernatesRead(d, meta)
}

func resourceKubernatesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernates.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernatesDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceKubernatesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKubernatesRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKubernatesDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      20 * time.Second,
		MinTimeout: 3 * 60 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernates.Get(id)
	})
}

func waitUntilKubernatesStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernates.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Kubernates).Status
	})
}
