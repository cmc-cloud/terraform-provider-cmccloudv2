package cmccloudv2

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernetesv2NodeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesv2NodeGroupCreate,
		Read:   resourceKubernetesv2NodeGroupRead,
		Update: resourceKubernetesv2NodeGroupUpdate,
		Delete: resourceKubernetesv2NodeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernetesv2NodeGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernetesv2NodeGroupSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			enable_autoscale := diff.Get("enable_autoscale").(bool)
			enable_autohealing := diff.Get("enable_autohealing").(bool)

			if enable_autoscale {
				if !isSet(diff, "min_node") || !isSet(diff, "max_node") { //} || !isSet(diff, "max_pods") || !isSet(diff, "cpu_threshold_percent") || !isSet(diff, "memory_threshold_percent") || !isSet(diff, "disk_threshold_percent") {
					return fmt.Errorf("when `enable_autoscale` is 'true', `min_node, max_node must be set")
				}
				if diff.Get("min_node").(int) >= diff.Get("max_node").(int) {
					return fmt.Errorf("when `enable_autoscale` is 'true', `max_node` must > `min_node`")
				}
			} else {
				if diff.Get("min_node").(int) != diff.Get("max_node").(int) {
					return fmt.Errorf("when `enable_autoscale` is 'false', `max_node` must equals to `min_node`")
				}
			}

			if enable_autohealing {
				if !isSet(diff, "max_unhealthy_percent") || !isSet(diff, "node_startup_timeout_minutes") {
					return fmt.Errorf("when `enable_autohealing` is 'true', `max_unhealthy_percent, node_startup_timeout_minutes must be set")
				}
			} else {
				if isSet(diff, "max_unhealthy_percent") || isSet(diff, "node_startup_timeout_minutes") {
					return fmt.Errorf("when `enable_autohealing` is 'false', `max_unhealthy_percent, node_startup_timeout_minutes must not be set")
				}
			}
			return nil
		},
	}
}

func getAutoScaleConfig(d *schema.ResourceData, meta interface{}) (map[string]interface{}, gocmcapiv2.Flavor, error) {
	client := meta.(*CombinedConfig).goCMCClient()
	flavor, err := client.Flavor.Get(d.Get("flavor_id").(string))
	if err != nil {
		return nil, flavor, fmt.Errorf("error receiving flavor %s: %v", d.Get("flavor_id").(string), err)
	}

	if !flavor.ExtraSpecs.IsK8sFlavor {
		return nil, flavor, fmt.Errorf("flavor %s is not a valid kubernetes flavor", d.Get("flavor_id").(string))
	}

	// cpuThreshold := (float64(d.Get("cpu_threshold_percent").(int)) / 100.0) * float64(flavor.Vcpus)
	// memoryThreshold := (float64(d.Get("memory_threshold_percent").(int)) / 100.0) * float64(flavor.RAM)
	// diskThreshold := (float64(d.Get("disk_threshold_percent").(int)) / 100.0) * float64(flavor.Disk)

	params := map[string]interface{}{
		"minNode": d.Get("min_node").(int),
		"maxNode": d.Get("max_node").(int),
		// "maxPods": d.Get("max_pods").(int),
		// "metaDataAutoScale": map[string]int{
		// 	"percentCpu":    d.Get("cpu_threshold_percent").(int),
		// 	"percentMemory": d.Get("memory_threshold_percent").(int),
		// 	"percentDisk":   d.Get("disk_threshold_percent").(int),
		// },
		// "cpuThreshold":    strconv.FormatFloat(cpuThreshold, 'f', 2, 64),
		// "memoryThreshold": strconv.FormatFloat(memoryThreshold, 'f', 2, 64) + "mb",
		// "diskThreshold":   strconv.FormatFloat(diskThreshold, 'f', 2, 64) + "Gb",
	}

	if !d.Get("enable_autoscale").(bool) {
		params["maxNode"] = params["minNode"]
	}
	return params, flavor, nil
}
func resourceKubernetesv2NodeGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	networkId := ""
	if v, ok := d.GetOk("subnet_id"); ok {
		if v.(string) != "" {
			subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
			if err != nil {
				return fmt.Errorf("error receving subnet with id = %s: %v", d.Get("subnet_id").(string), err)
			}
			networkId = subnet.NetworkID
		}
	}
	if networkId == "" {
		// lay networkId tu cluster
		cluster_id := d.Get("cluster_id").(string)
		cluster, err := client.Kubernetesv2.Get(cluster_id)
		if err != nil {
			return fmt.Errorf("error receving cluster with id = %s: %v", cluster_id, err)
		}
		networkId = cluster.VpcID
	}

	cluster_id := d.Get("cluster_id").(string)
	params, flavor, err := getAutoScaleConfig(d, meta)
	if err != nil {
		return err
	}
	params["region"] = client.Configs.RegionId
	params["project"] = client.Configs.ProjectId
	params["nameNodeGroup"] = d.Get("name").(string)
	params["clusterId"] = cluster_id
	params["flavorId"] = flavor.Name
	params["securityGroups"] = d.Get("security_group_ids").([]interface{})
	params["zone"] = d.Get("zone").(string)
	params["sshKeyName"] = d.Get("key_name").(string)
	params["isAutoscale"] = d.Get("enable_autoscale").(bool)
	params["workerImageGPUTag"] = d.Get("image_gpu_tag").(string)
	params["volumeType"] = d.Get("volume_type").(string)
	params["volumeSize"] = d.Get("volume_size").(int)
	params["billingMode"] = d.Get("billing_mode").(string)
	params["networkID"] = networkId
	// params["isNTP"] = d.Get("ntp_enabled").(bool)
	params["ntpServers"] = flatternNtpServers(d)
	// params["cidrBlockPod"] = d.Get("cidr_block_pod").(string)

	// kiểm tra max_pods được khai báo không
	if v, ok := d.GetOk("max_pods"); ok {
		params["maxPodPerNode"] = v.(int)
	}
	if v, ok := d.GetOk("init_current_node"); ok {
		params["currentNode"] = v.(int)
	}
	if _, ok := d.GetOk("node_metadatas"); ok {
		params["nodeMetadatas"] = d.Get("node_metadatas").([]interface{})
	}

	if d.Get("enable_autoscale").(bool) {
		// kiem tra xem cluster co enable auto scale ko, neu ko enable => ko support
		status, err := client.Kubernetesv2.GetStatus(cluster_id)
		if err != nil {
			return fmt.Errorf("error getting Kubernetesv2 Cluster: %s", err)
		}
		if !status.EnableAutoScale {
			return fmt.Errorf("you need to enable the autoscale on the cluster before creating a node group with the autoscale feature %v", err)
		}
	}

	kubernetesv2nodegroup, err := client.Kubernetesv2.CreateNodeGroup(cluster_id, params)
	if err != nil {
		return fmt.Errorf("error creating Kubernetesv2 NodeGroup: %s", err)
	}
	d.SetId(kubernetesv2nodegroup.ID)

	_, err = waitUntilKubernetesv2NodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating Kubernetesv2 NodeGroup: %v", err)
	}

	if d.Get("enable_autohealing").(bool) {
		params := map[string]interface{}{
			"action":                "enable",
			"maxUnhealthy":          strconv.Itoa(d.Get("max_unhealthy_percent").(int)) + "%",
			"nodeStartupTimeout":    strconv.Itoa(d.Get("node_startup_timeout_minutes").(int)) + "m",
			"externalProviderNames": "auto-healing-node-group",
			"nodeGroupId":           d.Id(),
		}
		_, err := client.Kubernetesv2.UpdateNodeGroup(cluster_id, params)
		if err != nil {
			return fmt.Errorf("error enable auto healing of Kubernetes NodeGroup: %v", err)
		}

		_, err = waitUntilKubernetesv2NodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error creating Kubernetes NodeGroup: %v", err)
		}
	}

	res := resourceKubernetesv2NodeGroupRead(d, meta)

	// kiểm tra có block gpu_configs không, nếu có thì thực hiện ConfigGpu, cần đủ các thuộc tính "gpu_model", "driver", "strategy", "mig_supported", "mig_profile", "time_slicing","gpu_profiles" thuộc gpu_configs
	// if v, ok := d.GetOk("gpu_config"); ok && len(v.([]interface{})) > 0 {
	// 	gpuParams, err := getGpuConfig(d, meta)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	_, err = client.Kubernetesv2.ConfigGpu(cluster_id, d.Id(), gpuParams)
	// 	if err != nil {
	// 		return fmt.Errorf("error configuring GPU for Kubernetes NodeGroup: %v", err)
	// 	}
	// }

	return res
}

func resourceKubernetesv2NodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	nodegroup, err := client.Kubernetesv2.GetNodeGroup(d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving Kubernetesv2 NodeGroup %s: %v", d.Id(), err)
	}

	_ = d.Set("id", nodegroup.ID)
	_ = d.Set("name", nodegroup.Name)
	// _ = d.Set("clusterId", nodegroup.ID)
	// _ = d.Set("flavor_id", nodegroup.MetadataMachineDeployment.FlavorName)
	_ = d.Set("key_name", nodegroup.KeyName)
	// _ = d.Set("security_group_ids", nodegroup.MinNodeCount)
	// _ = d.Set("zone", nodegroup.DockerVolumeSize)
	// _ = d.Set("image_gpu_tag", "")
	_ = d.Set("enable_autoscale", false)
	_ = d.Set("enable_autohealing", false)
	_ = d.Set("status", nodegroup.Status)
	_ = d.Set("ntp_servers", convertNtpServers(nodegroup.NtpServers))
	for _, provider := range nodegroup.ExternalProviders {
		if strings.Contains(provider.Name, "auto-scale") {
			if provider.Status == "active" && provider.Config.MaxNode > provider.Config.MinNode {
				_ = d.Set("enable_autoscale", true)
			}
			setInt(d, "min_node", provider.Config.MinNode)
			setInt(d, "max_node", provider.Config.MaxNode)
			// setInt(d, "max_pods", provider.Config.MaxPods)
			// setInt(d, "current_node", provider.Config.CurrentNode)
			// nodeMetadatas := make([]map[string]interface{}, 0, len(provider.Config.NodeMetadatas))
			// for _, meta := range provider.Config.NodeMetadatas {
			// 	m := map[string]interface{}{
			// 		"key":    meta.Key,
			// 		"value":  meta.Value,
			// 		"type":   meta.Type,
			// 		"effect": meta.Effect,
			// 	}
			// 	nodeMetadatas = append(nodeMetadatas, m)
			// }
			// _ = d.Set("node_metadatas", nodeMetadatas)
			// setInt(d, "cpu_threshold_percent", provider.Config.MetaDataAutoScale.PercentCPU)
			// setInt(d, "memory_threshold_percent", provider.Config.MetaDataAutoScale.PercentMemory)
			// setInt(d, "disk_threshold_percent", provider.Config.MetaDataAutoScale.PercentDisk)
		}
		if strings.Contains(provider.Name, "auto-healing") {
			if provider.Status == "active" {
				_ = d.Set("enable_autohealing", true)
			}
			setInt(d, "max_unhealthy_percent", int(provider.Config.MaxUnhealthy))
			setInt(d, "node_startup_timeout_minutes", int(provider.Config.NodeStartupTimeout))
		}
	}

	params := map[string]string{
		"network_id": nodegroup.VpcID, // VpcId o day la network id cua subnet
	}
	subnets, err := client.Subnet.List(params)
	if err != nil {
		return fmt.Errorf("error retrieving list of subnets %v", err)
	}
	for _, subnet := range subnets {
		if subnet.NetworkID == nodegroup.VpcID {
			_ = d.Set("subnet_id", subnet.ID)
			break
		}
	}
	return nil
}

func resourceKubernetesv2NodeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	cluster_id := d.Get("cluster_id").(string)
	if d.HasChange("enable_autoscale") || d.HasChange("min_node") || d.HasChange("max_node") || d.HasChange("max_pods") || d.HasChange("cpu_threshold_percent") || d.HasChange("memory_threshold_percent") || d.HasChange("disk_threshold_percent") {
		if d.Get("enable_autoscale").(bool) {
			// kiem tra xem cluster co enable auto scale ko, neu ko enable => ko support
			status, err := client.Kubernetesv2.GetStatus(cluster_id)
			if err != nil {
				return fmt.Errorf("error getting Kubernetesv2 Cluster: %s", err)
			}
			if !status.EnableAutoScale {
				return fmt.Errorf("you need to enable the autoscale on the cluster before creating a node group with the autoscale feature %v", err)
			}
		}

		params, _, err := getAutoScaleConfig(d, meta)
		if err != nil {
			return err
		}
		action := "disable"
		if d.Get("enable_autoscale").(bool) {
			action = "enable"
		}
		params["action"] = action
		params["externalProviderNames"] = "auto-scale-node-group"
		params["nodeGroupId"] = d.Id()
		_, err = client.Kubernetesv2.UpdateNodeGroup(cluster_id, params)
		if err != nil {
			return fmt.Errorf("error updating Kubernetes NodeGroup: %v", err)
		}

		_, err = waitUntilKubernetesv2NodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error creating Kubernetes NodeGroup: %v", err)
		}
	}
	if d.HasChange("enable_autohealing") || d.HasChange("max_unhealthy_percent") || d.HasChange("node_startup_timeout_minutes") {
		action := "disable"
		if d.Get("enable_autohealing").(bool) {
			action = "enable"
		}
		params := map[string]interface{}{
			"action":                action,
			"maxUnhealthy":          strconv.Itoa(d.Get("max_unhealthy_percent").(int)) + "%",
			"nodeStartupTimeout":    strconv.Itoa(d.Get("node_startup_timeout_minutes").(int)) + "m",
			"externalProviderNames": "auto-healing-node-group",
			"nodeGroupId":           d.Id(),
		}
		_, err := client.Kubernetesv2.UpdateNodeGroup(cluster_id, params)
		if err != nil {
			return fmt.Errorf("error updating Kubernetes NodeGroup: %v", err)
		}

		_, err = waitUntilKubernetesv2NodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error creating Kubernetes NodeGroup: %v", err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetKubernateNodeGroupBilingMode(cluster_id, d.Id(), d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("error when update billing mode of Nodegroup [%s]: %v", d.Id(), err)
		}
	}

	// if d.HasChange("gpu_config") {
	// 	if v, ok := d.GetOk("gpu_config"); ok && len(v.([]interface{})) > 0 {
	// 		gpuParams, err := getGpuConfig(d, meta)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		_, err = client.Kubernetesv2.ConfigGpu(cluster_id, d.Id(), gpuParams)
	// 		if err != nil {
	// 			return fmt.Errorf("error configuring GPU for Kubernetes NodeGroup: %v", err)
	// 		}
	// 	} else {
	// 		// gpu_config block removed => disable GPU
	// 		_, err := client.Kubernetesv2.DisableGpu(cluster_id, d.Id())
	// 		if err != nil {
	// 			return fmt.Errorf("error disabling GPU for Kubernetes NodeGroup: %v", err)
	// 		}
	// 	}
	// }
	return resourceKubernetesv2NodeGroupRead(d, meta)
}

func resourceKubernetesv2NodeGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernetesv2.DeleteNodeGroup(d.Get("cluster_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("error delete kubernetesv2 nodegroup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernetesv2NodeGroupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete kubernetesv2 nodegroup [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceKubernetesv2NodeGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	// Tách ID thành hai tham số
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected import ID to be in the format 'cluster_id/nodegroup_id'")
	}

	cluster_id := parts[0]
	nodegroup_id := parts[1]

	_ = d.Set("cluster_id", cluster_id)

	// Thiết lập ID cho resource
	d.SetId(nodegroup_id)
	err := resourceKubernetesv2NodeGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKubernetesv2NodeGroupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      30 * time.Second,
		MinTimeout: 5 * 60 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetesv2.GetNodeGroup(d.Get("cluster_id").(string), id)
	})
}

func waitUntilKubernetesv2NodeGroupStatusChangedState(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"HEALTHY", "RUNNING", "active", "Ready", "Running"}, []string{"ERROR", "SHUTDOWN", "FAILURE", "failure", "deleting"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernetesv2.GetNodeGroup(d.Get("cluster_id").(string), id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Kubernetesv2NodeGroup).Status
	})
}
