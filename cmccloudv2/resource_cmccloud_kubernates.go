package cmccloudv2

import (
	"fmt"
	"time"

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
		Schema:        kubernatesSchema(),
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if v, ok := d.GetOk("labels"); ok {
				blockList := v.([]interface{})
				if len(blockList) > 0 {
					labels := blockList[0].(map[string]interface{})
					auto_scaling_enable := labels["auto_scaling_enable"].(bool)
					if auto_scaling_enable {
						// co enable => phai set 2 truong nay, khong set => thong bao
						if labels["max_node_count"].(int) <= 0 { // khong duoc set max_node_count
							return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enable is 'true'")
						}
						if labels["min_node_count"].(int) <= 0 { // khong duoc set min_node_count
							return fmt.Errorf("min_node_count & max_node_count must be set > 0 when auto_scaling_enable is 'true'")
						}
					} else {
						// khong enable => ko set 2 truong nay
						if labels["max_node_count"].(int) > 0 {
							return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enable is 'false'")
						}
						if labels["min_node_count"].(int) > 0 {
							return fmt.Errorf("min_node_count & max_node_count must not be set when auto_scaling_enable is 'false'")
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
		"kube_dashboard_enable": metas["kube_dashboard_enable"].(bool),
		"metrics_server_enable": metas["metrics_server_enable"].(bool),
		"npd_enable":            metas["npd_enable"].(bool),
		"auto_scaling_enable":   metas["auto_scaling_enable"].(bool),
		"auto_healing_enable":   metas["auto_healing_enable"].(bool),
		"max_node_count":        metas["max_node_count"].(int),
		"min_node_count":        metas["min_node_count"].(int),
		"kube_tag":              metas["kube_tag"].(string),
		"network-driver":        metas["network_driver"].(string),
		"calico_ipv4pool":       metas["calico_ipv4pool"].(string),
		"docker_volume_type":    d.Get("docker_volume_type").(string),
		"zone":                  d.Get("zone").(string),
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

	kubernates, err := client.Kubernates.Create(params)
	if err != nil {
		return fmt.Errorf("Error creating Kubernates: %s", err)
	}
	d.SetId(kubernates.ID)
	return resourceKubernatesRead(d, meta)
}

func resourceKubernatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	kubernates, err := client.Kubernates.Get(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernates %s: %v", d.Id(), err)
	}

	labels := make([]map[string]interface{}, 1)
	labels[0] = map[string]interface{}{
		"kube_dashboard_enable": kubernates.Labels.KubeDashboardEnabled,
		"metrics_server_enable": kubernates.Labels.MetricsServerEnabled,
		"npd_enable":            kubernates.Labels.NpdEnabled,
		"auto_scaling_enable":   kubernates.Labels.AutoScalingEnabled,
		"auto_healing_enable":   kubernates.Labels.AutoHealingEnabled,

		"kube_tag":           kubernates.Labels.KubeTag,
		"network_driver":     kubernates.Labels.NetworkDriver,
		"calico_ipv4pool":    kubernates.Labels.CalicoIpv4Pool,
		"docker_volume_type": kubernates.Labels.DockerVolumeType,
		// "create_timeout":        kubernates.Labels.CreateTimeout,
		"zone": kubernates.Labels.AvailabilityZone,
	}

	if kubernates.Labels.AutoScalingEnabled {
		labels[0]["max_node_count"] = kubernates.Labels.MaxNodeCount
		labels[0]["min_node_count"] = kubernates.Labels.MinNodeCount
	}

	// gocmcapiv2.Logo("labels = ", labels)

	_ = d.Set("id", kubernates.ID)
	_ = d.Set("name", kubernates.Name)
	_ = d.Set("zone", kubernates.Labels.AvailabilityZone)
	_ = d.Set("subnet_id", kubernates.SubnetID)
	_ = d.Set("docker_volume_size", kubernates.DockerVolumeSize)
	_ = d.Set("docker_volume_type", kubernates.Labels.DockerVolumeType)
	_ = d.Set("keypair", kubernates.Keypair)
	_ = d.Set("create_timeout", kubernates.CreateTimeout)

	default_master := map[string]interface{}{
		"node_count":   kubernates.MasterCount,
		"flavor_id":    kubernates.MasterFlavorID,
		"billing_mode": kubernates.MasterBillingMode,
	}
	d.Set("default_master", []interface{}{default_master})

	default_worker := map[string]interface{}{
		"node_count":   kubernates.NodeCount,
		"flavor_id":    kubernates.NodeFlavorID,
		"billing_mode": kubernates.NodeBillingMode,
	}
	d.Set("default_worker", []interface{}{default_worker})

	_ = d.Set("created_at", kubernates.CreatedAt)
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
			return fmt.Errorf("Error when rename Kubernates [%s]: %v", id, err)
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
		return fmt.Errorf("Error delete cloud kubernates: %v", err)
	}
	return nil
}

func resourceKubernatesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKubernatesRead(d, meta)
	return []*schema.ResourceData{d}, err
}
