package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOpenSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenSearchCreate,
		Read:   resourceOpenSearchRead,
		Update: resourceOpenSearchUpdate,
		Delete: resourceOpenSearchDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOpenSearchImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        opensearchSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			// If enable_snapshot is true, require the rentation fields
			if diff.Get("enable_snapshot").(bool) {
				requiredFields := []string{"snapshot_creation_cron", "snapshot_timezone", "rentation_max_age", "rentation_min_count", "rentation_max_count"}
				for _, field := range requiredFields {
					val, ok := diff.GetOk(field)
					if !ok || val == nil || (val == 0 && diff.Get(field) == 0) {
						return fmt.Errorf("%s is required when enable_snapshot is true", field)
					}
				}
				maxCount, maxOk := diff.GetOk("rentation_max_count")
				minCount, minOk := diff.GetOk("rentation_min_count")
				if maxOk && minOk {
					maxCountInt, ok1 := maxCount.(int)
					minCountInt, ok2 := minCount.(int)
					if ok1 && ok2 {
						if maxCountInt <= minCountInt {
							return fmt.Errorf("rentation_max_count must be greater than rentation_min_count")
						}
					}
				}
			}
			if diff.Get("enable_storage_autoscaling").(bool) {
				requiredFields := []string{"storage_autoscaling_threshold", "storage_autoscaling_increment", "storage_autoscaling_max"}
				for _, field := range requiredFields {
					val, ok := diff.GetOk(field)
					if !ok || val == nil || (val == 0 && diff.Get(field) == 0) {
						return fmt.Errorf("%s is required when enable_storage_autoscaling is true", field)
					}
				}
			}

			return nil
		},
	}
}

func resourceOpenSearchCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"version":                       d.Get("version").(string),
		"flavor_id":                     d.Get("flavor_id").(string),
		"dashboard_flavor_id":           d.Get("dashboard_flavor_id").(string),
		"volume_size":                   d.Get("volume_size").(int),
		"admin_password":                d.Get("admin_password").(string),
		"node_count":                    d.Get("node_count").(int),
		"master_count":                  d.Get("master_count").(int),
		"enable_drain_nodes":            true, //d.Get("enable_drain_nodes").(bool),
		"dashboard_replicas":            d.Get("dashboard_replicas").(int),
		"enable_snapshot":               d.Get("enable_snapshot").(bool),
		"snapshot_creation_cron":        d.Get("snapshot_creation_cron").(string),
		"snapshot_timezone":             d.Get("snapshot_timezone").(string),
		"rentation_max_age":             d.Get("rentation_max_age").(int),
		"rentation_min_count":           d.Get("rentation_min_count").(int),
		"rentation_max_count":           d.Get("rentation_max_count").(int),
		"lb_subnet_id":                  d.Get("lb_subnet_id").(string),
		"enable_storage_autoscaling":    d.Get("enable_storage_autoscaling").(bool),
		"storage_autoscaling_threshold": d.Get("storage_autoscaling_threshold").(int),
		"storage_autoscaling_increment": d.Get("storage_autoscaling_increment").(int),
		"storage_autoscaling_max":       d.Get("storage_autoscaling_max").(int),
		"billing_mode":                  d.Get("billing_mode").(string),
		"api_domain":                    d.Get("api_domain").(string),
		"dashboard_domain":              d.Get("dashboard_domain").(string),
		"enable_lb_internal":            true, // d.Get("enable_lb_internal").(bool),
		"tags":                          d.Get("tags").(*schema.Set).List(),
		// "snapshot_deletion_cron":       d.Get("snapshot_deletion_cron").(string),
	}
	instance, err := client.OpenSearch.Create(params)
	if err != nil {
		return fmt.Errorf("error creating OpenSearch: %s", err)
	}
	d.SetId(instance.ID)
	_, err = waitUntilOpenSearchActive(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating OpenSearch: %s", err)
	}
	return resourceOpenSearchRead(d, meta)
}

func resourceOpenSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.OpenSearch.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving OpenSearch %s: %v", d.Id(), err)
	}
	_ = d.Set("dashboard_flavor_id", instance.DashboardFlavorID)
	_ = d.Set("node_count", instance.NodeCount)
	_ = d.Set("master_count", instance.MasterCount)
	_ = d.Set("dashboard_replicas", instance.DashboardReplicas)
	_ = d.Set("enable_snapshot", instance.EnableSnapshot)
	_ = d.Set("snapshot_creation_cron", instance.SnapshotCreationCron)
	_ = d.Set("snapshot_timezone", instance.SnapshotTimezone)
	_ = d.Set("rentation_max_age", instance.RentationMaxAge)
	_ = d.Set("rentation_min_count", instance.RentationMinCount)
	_ = d.Set("rentation_max_count", instance.RentationMaxCount)
	_ = d.Set("lb_subnet_id", instance.LbSubnetID)
	_ = d.Set("enable_storage_autoscaling", instance.EnableStorageAutoscaling)
	_ = d.Set("storage_autoscaling_threshold", instance.StorageAutoscalingThreshold)
	_ = d.Set("storage_autoscaling_increment", instance.StorageAutoscalingIncrement)
	_ = d.Set("storage_autoscaling_max", instance.StorageAutoscalingMax)
	_ = d.Set("name", instance.Name)
	_ = d.Set("version", instance.Version)
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("api_domain", instance.ApiDomain)
	_ = d.Set("dashboard_domain", instance.DashboardDomain)
	_ = d.Set("tags", convertTagsToSet(instance.Tags))
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	// _ = d.Set("enable_lb_internal", instance.EnableLbInternal)
	// _ = d.Set("enable_drain_nodes", instance.EnableDrainNodes)
	return nil
}

func resourceOpenSearchUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("tags") {
		_, err := client.Tag.UpdateTag(id, "OpenSearch", d)
		if err != nil {
			return fmt.Errorf("error when set opensearch tags [%s]: %v", id, err)
		}
	}

	if d.HasChange("flavor_id") ||
		d.HasChange("dashboard_flavor_id") ||
		d.HasChange("node_count") ||
		d.HasChange("dashboard_replicas") ||
		d.HasChange("volume_size") {
		params := map[string]interface{}{
			"flavor_id":           d.Get("flavor_id").(string),
			"dashboard_flavor_id": d.Get("dashboard_flavor_id").(string),
			"node_count":          d.Get("node_count").(int),
			"dashboard_replicas":  d.Get("dashboard_replicas").(int),
			"volume_size":         d.Get("volume_size").(int),
		}
		_, err := client.OpenSearch.Resize(id, params)
		if err != nil {
			return fmt.Errorf("error resizing OpenSearch %s: %v", id, err)
		}
		_, err = waitUntilOpenSearchJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for OpenSearch resize job %s: %v", id, err)
		}
	}

	if d.HasChange("enable_storage_autoscaling") ||
		d.HasChange("storage_autoscaling_threshold") ||
		d.HasChange("storage_autoscaling_increment") ||
		d.HasChange("storage_autoscaling_max") {
		params := map[string]interface{}{
			"enable_storage_autoscaling":    d.Get("enable_storage_autoscaling").(bool),
			"storage_autoscaling_threshold": d.Get("storage_autoscaling_threshold").(int),
			"storage_autoscaling_increment": d.Get("storage_autoscaling_increment").(int),
			"storage_autoscaling_max":       d.Get("storage_autoscaling_max").(int),
		}
		_, err := client.OpenSearch.UpdateStorageAutoscaling(id, params)
		if err != nil {
			return fmt.Errorf("error updating storage autoscaling for OpenSearch %s: %v", id, err)
		}
		_, err = waitUntilOpenSearchJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for OpenSearch storage autoscaling update job %s: %v", id, err)
		}
	}

	if d.HasChange("enable_snapshot") ||
		d.HasChange("snapshot_creation_cron") ||
		d.HasChange("snapshot_timezone") ||
		d.HasChange("rentation_max_age") ||
		d.HasChange("rentation_min_count") ||
		d.HasChange("rentation_max_count") {
		params := map[string]interface{}{
			"enable_snapshot":        d.Get("enable_snapshot").(bool),
			"snapshot_creation_cron": d.Get("snapshot_creation_cron").(string),
			"snapshot_timezone":      d.Get("snapshot_timezone").(string),
			"rentation_max_age":      d.Get("rentation_max_age").(int),
			"rentation_min_count":    d.Get("rentation_min_count").(int),
			"rentation_max_count":    d.Get("rentation_max_count").(int),
		}
		_, err := client.OpenSearch.UpdateSnapshotPolicy(id, params)
		if err != nil {
			return fmt.Errorf("error updating snapshot policy for OpenSearch %s: %v", id, err)
		}
		_, err = waitUntilOpenSearchJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for OpenSearch snapshot policy update job %s: %v", id, err)
		}
	}
	return resourceOpenSearchRead(d, meta)
}

func resourceOpenSearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.OpenSearch.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete opensearch: %v", err)
	}
	_, err = waitUntilOpenSearchDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete opensearch: %v", err)
	}
	return nil
}

func resourceOpenSearchImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceOpenSearchRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilOpenSearchActive(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"active"}, []string{"error"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).OpenSearch.Get(id)
	}, func(obj interface{}) string {
		return strings.ToLower(obj.(gocmcapiv2.OpenSearch).Status)
	})
}
func waitUntilOpenSearchJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"done"}, []string{"processing"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).OpenSearch.Get(id)
	}, func(obj interface{}) string {
		taskState := strings.ToLower(obj.(gocmcapiv2.OpenSearch).TaskState)
		if taskState != "" {
			return "processing"
		}
		return "done"
	})
}

func waitUntilOpenSearchDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).OpenSearch.Get(id)
	})
}
