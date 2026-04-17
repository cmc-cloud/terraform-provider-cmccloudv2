package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRdsCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRdsClusterCreate,
		Read:   resourceRdsClusterRead,
		Update: resourceRdsClusterUpdate,
		Delete: resourceRdsClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRdsClusterImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        rdsClusterSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			if diff.Get("enable_backup").(bool) {
				requiredFields := []string{"enable_pitr", "backup_schedule", "backup_retention"}
				for _, field := range requiredFields {
					val, ok := diff.GetOk(field)
					if !ok || val == nil || (val == 0 && diff.Get(field) == 0) {
						return fmt.Errorf("%s is required when enable_backup is true", field)
					}
				}
			}
			if diff.Get("enable_storage_autoscaling").(bool) {
				requiredFields := []string{"storage_autoscaling_threshold", "storage_autoscaling_increment"}
				for _, field := range requiredFields {
					val, ok := diff.GetOk(field)
					if !ok || val == nil || (val == 0 && diff.Get(field) == 0) {
						return fmt.Errorf("%s is required when enable_storage_autoscaling is true", field)
					}
				}
			}
			if !diff.Get("enable_backup").(bool) {
				if diff.Get("enable_pitr").(bool) {
					return fmt.Errorf("enable_pitr must be false when enable_backup is false")
				}
			}
			return nil
		},
	}
}

func resourceRdsClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	params := map[string]interface{}{
		"name":                          d.Get("name").(string),
		"billing_mode":                  d.Get("billing_mode").(string),
		"flavor_id":                     d.Get("flavor_id").(string),
		"volume_size":                   d.Get("volume_size").(int),
		"db_engine":                     d.Get("db_engine").(string),
		"db_version":                    d.Get("db_version").(string),
		"subnet_id":                     d.Get("subnet_id").(string),
		"mode":                          d.Get("mode").(string),
		"cluster_size":                  d.Get("cluster_size").(int),
		"proxy_size":                    d.Get("proxy_size").(int),
		"enable_backup":                 d.Get("enable_backup").(bool),
		"enable_pitr":                   d.Get("enable_pitr").(bool),
		"backup_schedule":               d.Get("backup_schedule").(string),
		"backup_retention":              d.Get("backup_retention").(int),
		"enable_storage_autoscaling":    d.Get("enable_storage_autoscaling").(bool),
		"storage_autoscaling_threshold": d.Get("storage_autoscaling_threshold").(int),
		"storage_autoscaling_increment": d.Get("storage_autoscaling_increment").(int),
	}
	instance, err := client.RdsCluster.Create(params)
	if err != nil {
		return fmt.Errorf("error creating RDS Cluster Instance: %s", err)
	}
	d.SetId(instance.ID)
	_, err = waitUntilRdsClusterActive(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating RDS Cluster Instance: %s", err)
	}
	return resourceRdsClusterRead(d, meta)
}

func resourceRdsClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.RdsCluster.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving RdsCluster Instance %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("billing_mode", instance.BillingMode)
	_ = d.Set("db_engine", instance.DbEngine)
	_ = d.Set("db_version", instance.DbVersion)
	_ = d.Set("flavor_id", instance.FlavorID)
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("db_engine", instance.DbEngine)
	_ = d.Set("db_version", instance.DbVersion)
	_ = d.Set("subnet_id", instance.SubnetID)
	_ = d.Set("mode", instance.Mode)
	_ = d.Set("cluster_size", instance.ClusterSize)
	_ = d.Set("proxy_size", instance.ProxySize)
	_ = d.Set("enable_backup", instance.EnableBackup)
	_ = d.Set("enable_pitr", instance.EnablePitr)
	_ = d.Set("backup_schedule", instance.BackupSchedule)
	_ = d.Set("backup_retention", instance.BackupRetention)
	_ = d.Set("enable_storage_autoscaling", instance.EnableStorageAutoscaling)
	_ = d.Set("storage_autoscaling_threshold", instance.StorageAutoscalingThreshold)
	_ = d.Set("storage_autoscaling_increment", instance.StorageAutoscalingIncrement)

	_ = d.Set("lb_vip_ipaddress", instance.LbVipIPAddress)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.Created)
	return nil
}

func resourceRdsClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("enable_storage_autoscaling") ||
		d.HasChange("storage_autoscaling_threshold") ||
		d.HasChange("storage_autoscaling_increment") {
		params := map[string]interface{}{
			"enable_storage_autoscaling":    d.Get("enable_storage_autoscaling").(bool),
			"storage_autoscaling_threshold": d.Get("storage_autoscaling_threshold").(int),
			"storage_autoscaling_increment": d.Get("storage_autoscaling_increment").(int),
		}
		_, err := client.RdsCluster.UpdateStorageAutoscaling(id, params)
		if err != nil {
			return fmt.Errorf("error updating storage autoscaling for RdsCluster %s: %v", id, err)
		}
		_, err = waitUntilRdsClusterJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for RdsCluster storage autoscaling update job %s: %v", id, err)
		}
	}

	if d.HasChange("enable_backup") ||
		d.HasChange("enable_pitr") ||
		d.HasChange("backup_schedule") ||
		d.HasChange("backup_retention") {
		params := map[string]interface{}{
			"enable_backup":    d.Get("enable_backup").(bool),
			"enable_pitr":      d.Get("enable_pitr").(bool),
			"backup_schedule":  d.Get("backup_schedule").(string),
			"backup_retention": d.Get("backup_retention").(int),
		}
		_, err := client.RdsCluster.UpdateBackupPolicy(id, params)
		if err != nil {
			return fmt.Errorf("error updating backup policy for RdsCluster %s: %v", id, err)
		}
		_, err = waitUntilRdsClusterJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for RdsCluster backup policy update job %s: %v", id, err)
		}
	}

	if d.HasChange("flavor_id") ||
		d.HasChange("volume_size") ||
		d.HasChange("cluster_size") ||
		d.HasChange("proxy_size") {
		params := map[string]interface{}{
			"flavor_id":    d.Get("flavor_id").(string),
			"volume_size":  d.Get("volume_size").(int),
			"cluster_size": d.Get("cluster_size").(int),
			"proxy_size":   d.Get("proxy_size").(int),
		}
		_, err := client.RdsCluster.Resize(id, params)
		if err != nil {
			return fmt.Errorf("error resizing RdsCluster %s: %v", id, err)
		}
		_, err = waitUntilRdsClusterJobFinished(d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for RdsCluster resize job %s: %v", id, err)
		}
	}
	return resourceRdsClusterRead(d, meta)
}

func resourceRdsClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.RdsCluster.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete rds cluster: %v", err)
	}
	_, err = waitUntilRdsClusterDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete rds cluster: %v", err)
	}
	return nil
}

func resourceRdsClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRdsClusterRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilRdsClusterActive(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"ACTIVE"}, []string{"ERROR"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RdsCluster.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.RdsCluster).Status)
	})
}
func waitUntilRdsClusterJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"done"}, []string{"processing"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RdsCluster.Get(id)
	}, func(obj interface{}) string {
		taskState := strings.ToLower(obj.(gocmcapiv2.RdsCluster).TaskState)
		if taskState != "" {
			return "processing"
		}
		return "done"
	})
}

func waitUntilRdsClusterDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).RdsCluster.Get(id)
	})
}
