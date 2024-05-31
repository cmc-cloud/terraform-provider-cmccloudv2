package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernatesNodeGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernatesNodeGroupCreate,
		Read:   resourceKubernatesNodeGroupRead,
		Update: resourceKubernatesNodeGroupUpdate,
		Delete: resourceKubernatesNodeGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKubernatesNodeGroupImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernatesNodeGroupSchema(),
	}
}

func resourceKubernatesNodeGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	params := map[string]interface{}{
		"name":           d.Get("name").(string),
		"flavor_id":      d.Get("flavor_id").(string),
		"node_count":     d.Get("node_count").(int),
		"min_node_count": d.Get("min_node_count").(int),
		"max_node_count": d.Get("max_node_count").(int),
		"billing_mode":   d.Get("billing_mode").(string),
		"zone":           d.Get("zone").(string),
		"labels": map[string]interface{}{
			"docker_volume_size": d.Get("docker_volume_size").(int),
			"docker_volume_type": d.Get("docker_volume_type").(string),
			"availability_zone":  d.Get("zone").(string),
		},
	}

	kubernatesnodegroup, err := client.Kubernates.CreateNodeGroup(d.Get("cluster_id").(string), params)
	if err != nil {
		return fmt.Errorf("Error creating Kubernates NodeGroup: %s", err)
	}
	d.SetId(kubernatesnodegroup.ID)

	_, err = waitUntilKubernatesNodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Kubernates NodeGroup: %v", err)
	}

	return resourceKubernatesNodeGroupRead(d, meta)
}

func resourceKubernatesNodeGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	nodegroup, err := client.Kubernates.GetNodeGroup(d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving Kubernates NodeGroup %s: %v", d.Id(), err)
	}

	_ = d.Set("id", nodegroup.ID)
	_ = d.Set("cluster_id", nodegroup.ClusterID)
	_ = d.Set("name", nodegroup.Name)
	_ = d.Set("flavor_id", nodegroup.FlavorID)
	_ = d.Set("node_count", nodegroup.NodeCount)
	_ = d.Set("min_node_count", nodegroup.MinNodeCount)
	_ = d.Set("max_node_count", nodegroup.MaxNodeCount)
	_ = d.Set("docker_volume_size", nodegroup.DockerVolumeSize)
	_ = d.Set("docker_volume_type", nodegroup.Labels.DockerVolumeType)
	_ = d.Set("zone", nodegroup.Labels.AvailabilityZone)
	_ = d.Set("created_at", nodegroup.CreatedAt)
	_ = d.Set("billing_mode", nodegroup.BillingMode)

	return nil
}

func resourceKubernatesNodeGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") {
		return fmt.Errorf("Can't change name of nodegroup %s after created", id)
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetKubernateBilingMode(d.Get("cluster_id").(string), d.Get("billing_mode").(string), "worker")
		if err != nil {
			return fmt.Errorf("Error when update billing mode of Kubernates NodeGroup [%s]: %v", id, err)
		}
	}
	if d.HasChange("node_count") {
		_, err := client.Kubernates.ResizeNodeGroup(d.Get("cluster_id").(string), map[string]interface{}{
			"nodegroup":  d.Get("name").(string),
			"node_count": d.Get("node_count").(int),
		})
		if err != nil {
			return fmt.Errorf("Error when change Kubernates NodeGroup [%s] node count: %v", id, err)
		}
		_, err = waitUntilKubernatesNodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when change Kubernates NodeGroup [%s] node count: %v", id, err)
		}
	}

	if d.HasChange("min_node_count") || d.HasChange("max_node_count") {
		_, err := client.Kubernates.UpdateNodeGroup(d.Get("cluster_id").(string), id, d.Get("min_node_count").(int), d.Get("max_node_count").(int))
		if err != nil {
			return fmt.Errorf("Error when change min_node_count/max_node_count of Kubernates NodeGroup [%s]: %v", id, err)
		}
		_, err = waitUntilKubernatesNodeGroupStatusChangedState(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error when change min_node_count/max_node_count of Kubernates NodeGroup [%s]: %v", id, err)
		}
	}
	return resourceKubernatesNodeGroupRead(d, meta)
}

func resourceKubernatesNodeGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernates.DeleteNodeGroup(d.Get("cluster_id").(string), d.Id())

	if err != nil {
		return fmt.Errorf("Error delete kubernates nodegroup [%s]: %v", d.Id(), err)
	}
	_, err = waitUntilKubernatesNodeGroupDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete kubernates nodegroup [%s]: %v", d.Id(), err)
	}
	return nil
}

func resourceKubernatesNodeGroupImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceKubernatesNodeGroupRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilKubernatesNodeGroupDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      30 * time.Second,
		MinTimeout: 5 * 60 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernates.GetNodeGroup(d.Get("cluster_id").(string), id)
	})
}

func waitUntilKubernatesNodeGroupStatusChangedState(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"CREATE_COMPLETE", "UPDATE_COMPLETE", "HEALTHY"}, []string{"CREATE_FAILED", "UPDATE_FAILED"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Kubernates.GetNodeGroup(d.Get("cluster_id").(string), id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.KubernatesNodeGroup).Status
	})
}
