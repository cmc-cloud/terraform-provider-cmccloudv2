package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceKubernetesv2NodeGroupGpuConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesv2NodeGroupGpuConfigCreate,
		Read:   resourceKubernetesv2NodeGroupGpuConfigRead,
		Update: resourceKubernetesv2NodeGroupGpuConfigUpdate,
		Delete: resourceKubernetesv2NodeGroupGpuConfigDelete,
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        kubernetesv2NodeGroupGpuConfigSchema(),
	}
}

func resourceKubernetesv2NodeGroupGpuConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	gpuParams, err := getGpuConfig(d, meta)
	if err != nil {
		return err
	}
	_, err = client.Kubernetesv2.ConfigGpu(d.Get("cluster_id").(string), d.Get("nodegroup_id").(string), gpuParams)
	if err != nil {
		return fmt.Errorf("error configuring GPU for Kubernetes NodeGroup: %v", err)
	}
	return resourceKubernetesv2NodeGroupGpuConfigRead(d, meta)
}

func getGpuConfig(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*CombinedConfig).goCMCClient()
	driver := d.Get("driver").(string)
	gpu_model := d.Get("gpu_model").(string)
	gpuDrivers, err := client.Kubernetesv2.GetGpuDrivers(d.Get("cluster_id").(string), d.Get("nodegroup_id").(string), gpu_model)
	if err != nil {
		return nil, fmt.Errorf("error getting Kubernetesv2 GpuDrivers: %s", err)
	}
	if len(gpuDrivers) == 0 {
		return nil, fmt.Errorf("no GPU drivers found for model %s", gpu_model)
	}
	for _, gpuDriver := range gpuDrivers {
		if gpuDriver.Name == driver {
			gpuParams := map[string]interface{}{
				"gpuModel":     gpu_model, // config["gpu_model"].(string),
				"driver":       driver,
				"strategy":     d.Get("strategy").(string),
				"migProfile":   d.Get("mig_profile").(string),
				"migSupported": gpuDriver.MigSupported,
				"timeSlicing":  gpuDriver.TimeSlicing, //map[bool]int{true: 1, false: 0}[config["time_slicing"].(bool)],
				"gpuProfiles":  d.Get("gpu_profiles").([]interface{}),
				"nodeGroupId":  d.Id(),
				"isGpu":        1,
			}
			return gpuParams, nil
		}
	}
	return nil, fmt.Errorf("not valid driver name/model: %s/%s", driver, gpu_model)
}
func resourceKubernetesv2NodeGroupGpuConfigRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("cluster_id").(string) + "_" + d.Get("nodegroup_id").(string)
	d.SetId(id)
	return nil
}

func resourceKubernetesv2NodeGroupGpuConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	gpuParams, err := getGpuConfig(d, meta)
	if err != nil {
		return err
	}
	_, err = client.Kubernetesv2.ConfigGpu(d.Get("cluster_id").(string), d.Get("nodegroup_id").(string), gpuParams)
	if err != nil {
		return fmt.Errorf("error configuring GPU for Kubernetes NodeGroup: %v", err)
	}
	return resourceKubernetesv2NodeGroupGpuConfigRead(d, meta)
}

func resourceKubernetesv2NodeGroupGpuConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Kubernetesv2.DisableGpu(d.Get("cluster_id").(string), d.Get("nodegroup_id").(string))

	if err != nil {
		return fmt.Errorf("error disabling GPU for Kubernetes NodeGroup [%s]: %v", d.Id(), err)
	}
	return nil
}
