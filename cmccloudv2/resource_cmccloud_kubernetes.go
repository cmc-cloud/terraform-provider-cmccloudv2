package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernetes() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesCreate,
		Read:   resourceKubernetesRead,
		Update: resourceKubernetesUpdate,
		Delete: resourceKubernetesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernetesImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernetesSchema(),
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			default_worker := d.Get("default_worker").([]interface{})[0].(map[string]interface{})
			if default_worker["min_node_count"].(int) > default_worker["node_count"].(int) {
				return fmt.Errorf("default_worker: min_node_count must be <= node_count")
			}
			if default_worker["max_node_count"].(int) < default_worker["node_count"].(int) {
				return fmt.Errorf("default_worker: max_node_count must be >= node_count")
			}
			// if v, ok := d.GetOk("labels"); ok {
			// 	blockList := v.([]interface{})
			// 	if len(blockList) > 0 {
			// 		labels := blockList[0].(map[string]interface{})
			// 		auto_scaling_enabled := labels["auto_scaling_enabled"].(bool)

			// 		default_worker := d.Get("default_worker").([]interface{})[0].(map[string]interface{})
			// 		if auto_scaling_enabled {
			// 			// co enable => phai set 2 truong nay, khong set => thong bao
			// 			if default_worker["max_node_count"].(int) <= 0 { // khong duoc set max_node_count
			// 				return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enabled is 'true'")
			// 			}
			// 			if default_worker["min_node_count"].(int) <= 0 { // khong duoc set min_node_count
			// 				return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enabled is 'true'")
			// 			}
			// 		} else {
			// 			// khong enable => ko set 2 truong nay
			// 			if default_worker["max_node_count"].(int) > 0 {
			// 				return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enabled is 'false'")
			// 			}
			// 			if default_worker["min_node_count"].(int) > 0 {
			// 				return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enabled is 'false'")
			// 			}
			// 		}
			// 	}
			// }
			return nil
		},
	}
}

func resourceKubernetesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	metas := getFirstBlock(d, "labels")
	default_worker := getFirstBlock(d, "default_worker")
	default_master := getFirstBlock(d, "default_master")
	labels := map[string]interface{}{
		"kube_dashboard_enabled": metas["kube_dashboard_enabled"].(bool),
		"metrics_server_enabled": metas["metrics_server_enabled"].(bool),
		"npd_enabled":            metas["npd_enabled"].(bool),
		"auto_scaling_enabled":   metas["auto_scaling_enabled"].(bool),
		"auto_healing_enabled":   metas["auto_healing_enabled"].(bool),
		"max_node_count":         default_worker["max_node_count"].(int),
		"min_node_count":         default_worker["min_node_count"].(int),
		"kube_tag":               metas["kube_tag"].(string),
		"network-driver":         metas["network_driver"].(string),
		"calico_ipv4pool":        metas["calico_ipv4pool"].(string),
		"docker_volume_type":     d.Get("docker_volume_type").(string),
		"zone":                   d.Get("zone").(string),
	}

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

	kubernetes, err := client.Kubernetes.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes: %s", err)
	}
	d.SetId(kubernetes.ID)

	_, err = waitUntilKubernetesStatusChangedState(d, meta, []string{"CREATE_COMPLETE", "HEALTHY"}, []string{"CREATE_FAILED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Kubernetes: %s", err)
	}
	return resourceKubernetesRead(d, meta)
}

func resourceKubernetesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	kubernetes, err := client.Kubernetes.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernetes %s: %v", d.Id(), err)
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
		// "zone": kubernetes.Labels.AvailabilityZone,
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
		"flavor_id":    kubernetes.NodeFlavorID,
		"billing_mode": kubernetes.NodeBillingMode,
	}
	nodegroups, _ := client.Kubernetes.GetNodeGroups(d.Id(), false) // pass
	for _, nodegroup := range nodegroups {
		if nodegroup.Name == "default-worker" {
			// kubernetes.NodeCount = tong so node cua tat ca nodegroup loai worker
			// khong lay gia tri min_node_count va max_node_count tu labels, vi gia tri nay chi la gia tri init tu khi tao cluster
			// khi doi min_node_count, max_node_count cua nodegroup default worker thi gia tri nay van khong thay doi
			default_worker["min_node_count"] = nodegroup.MinNodeCount
			default_worker["max_node_count"] = nodegroup.MaxNodeCount
			default_worker["node_count"] = nodegroup.NodeCount
		}
	}

	// if kubernetes.Labels.AutoScalingEnabled {
	// 	default_worker_block := getFirstBlock(d, "default_worker")
	// 	if default_worker_block["min_node_count"].(int) != 0 {
	// 		default_worker["min_node_count"] = kubernetes.Labels.MinNodeCount
	// 	}
	// 	if default_worker_block["max_node_count"] != 0 {
	// 		default_worker["max_node_count"] = kubernetes.Labels.MaxNodeCount
	// 	}
	// }

	d.Set("default_worker", []interface{}{default_worker})

	_ = d.Set("created_at", kubernetes.CreatedAt)
	_ = d.Set("labels", labels)

	return nil
}

func resourceKubernetesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	master_billing_mode_changed, new_master_billing_mode := isSubBlockFieldChanged(d, "default_master", "billing_mode")
	worker_billing_mode_changed, new_worker_billing_mode := isSubBlockFieldChanged(d, "default_worker", "billing_mode")
	worker_node_count_changed, new_worker_node_count := isSubBlockFieldChanged(d, "default_worker", "node_count")
	min_node_count_changed, _ := isSubBlockFieldChanged(d, "default_worker", "min_node_count")
	max_node_count_changed, _ := isSubBlockFieldChanged(d, "default_worker", "max_node_count")
	default_worker_block := getFirstBlock(d, "default_worker")

	if worker_node_count_changed || min_node_count_changed || max_node_count_changed {
		// _, err := client.Kubernetes.UpdateNodeCount(id, new_worker_node_count.(int))
		found := false
		nodegroups, _ := client.Kubernetes.GetNodeGroups(id, false)
		for _, nodegroup := range nodegroups {
			// gocmcapiv2.Logs("nodegroup " + nodegroup.Name)
			if nodegroup.Name == "default-worker" {
				found = true
				change_minmax_first := false
				// nếu  Node count hiện tại nằm trong khoảng giá trị min_node_count mới & max_node_count
				// thì update min,max trước
				if default_worker_block["min_node_count"].(int) <= nodegroup.NodeCount && nodegroup.NodeCount <= default_worker_block["max_node_count"].(int) {
					change_minmax_first = true
				}

				if change_minmax_first {
					if min_node_count_changed || max_node_count_changed {
						_, err := client.Kubernetes.UpdateNodeGroup(id, nodegroup.ID, default_worker_block["min_node_count"].(int), default_worker_block["max_node_count"].(int))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker min/max node count [%s]: %v", id, err)
						}
						_, err = waitUntilKubernetesStatusChangedState(d, meta, []string{"UPDATE_COMPLETE", "HEALTHY"}, []string{"UPDATE_FAILED"}, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker min/max node count [%s]: %v", id, err)
						}
					}

					if worker_node_count_changed {
						_, err := client.Kubernetes.ResizeNodeGroup(id, map[string]interface{}{
							"node_count": new_worker_node_count,
							"nodegroup":  nodegroup.ID,
						})
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker node count [%s]: %v", id, err)
						}
						_, err = waitUntilKubernetesStatusChangedState(d, meta, []string{"UPDATE_COMPLETE", "HEALTHY"}, []string{"UPDATE_FAILED"}, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker node count [%s]: %v", id, err)
						}
					}
				} else {
					if worker_node_count_changed {
						_, err := client.Kubernetes.ResizeNodeGroup(id, map[string]interface{}{
							"node_count": new_worker_node_count,
							"nodegroup":  nodegroup.ID,
						})
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker node count [%s]: %v", id, err)
						}
						_, err = waitUntilKubernetesStatusChangedState(d, meta, []string{"UPDATE_COMPLETE", "HEALTHY"}, []string{"UPDATE_FAILED"}, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker node count [%s]: %v", id, err)
						}
					}

					if min_node_count_changed || max_node_count_changed {
						_, err := client.Kubernetes.UpdateNodeGroup(id, nodegroup.ID, default_worker_block["min_node_count"].(int), default_worker_block["max_node_count"].(int))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker min/max node count [%s]: %v", id, err)
						}
						_, err = waitUntilKubernetesStatusChangedState(d, meta, []string{"UPDATE_COMPLETE", "HEALTHY"}, []string{"UPDATE_FAILED"}, d.Timeout(schema.TimeoutUpdate))
						if err != nil {
							return fmt.Errorf("Error when update Kubernetes worker min/max node count [%s]: %v", id, err)
						}
					}
				}
			}
		}
		if !found {
			return fmt.Errorf("Not found default_worker nodegroup of kubernetes [%s]", id)
		}
	}

	if master_billing_mode_changed {
		_, err := client.BillingMode.SetKubernateBilingMode(id, new_master_billing_mode.(string), "master")
		if err != nil {
			return fmt.Errorf("Error when change default master biling mode of Kubernetes cluster [%s]: %v", id, err)
		}
	}

	if worker_billing_mode_changed {
		_, err := client.BillingMode.SetKubernateBilingMode(id, new_worker_billing_mode.(string), "worker")
		if err != nil {
			return fmt.Errorf("Error when change default worker biling mode of Kubernetes cluster [%s]: %v", id, err)
		}
	}

	return resourceKubernetesRead(d, meta)
}

func resourceKubernetesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernetes.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernetesDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete kubernetes [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceKubernetesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKubernetesRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKubernetesDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      20 * time.Second,
		MinTimeout: 3 * 60 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetes.Get(id)
	})
}

func waitUntilKubernetesStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetes.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Kubernetes).Status
	})
}
