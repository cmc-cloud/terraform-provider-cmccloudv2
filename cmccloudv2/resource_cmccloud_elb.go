package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceELB() *schema.Resource {
	return &schema.Resource{
		Create: resourceELBCreate,
		Read:   resourceELBRead,
		Update: resourceELBUpdate,
		Delete: resourceELBDelete,
		Importer: &schema.ResourceImporter{
			State: resourceELBImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        elbSchema(),
		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			networkType := diff.Get("network_type").(string)
			subnetID, subnetIDSet := diff.GetOk("subnet_id")
			bandwidthMbps, bandwidthMbpsSet := diff.GetOk("bandwidth_mbps")

			if networkType == "public" {
				if !bandwidthMbpsSet || bandwidthMbps.(int) <= 0 {
					return fmt.Errorf("when network_type is 'public', 'bandwidth_mbps' must be set and greater than 0")
				}
				if subnetIDSet && len(subnetID.(string)) > 0 {
					return fmt.Errorf("when network_type is 'public', 'subnet_id' must be not be set")
				}
			}

			if networkType == "private" {
				if !subnetIDSet || subnetID.(string) == "" {
					return fmt.Errorf("when network_type is 'private', 'subnet_id' must be set and not empty")
				}
				if bandwidthMbpsSet && bandwidthMbps.(int) > 0 {
					return fmt.Errorf("when network_type is 'private', 'bandwidth_mbps' must be not be set")
				}
			}

			return nil
		},
	}
}

func resourceELBUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("tags") {
		_, err := client.ELB.Update(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"tags":        d.Get("tags").(*schema.Set).List(),
		})
		if err != nil {
			return fmt.Errorf("error when update ELB [%s]: %v", id, err)
		}
	}

	if d.HasChange("flavor_id") {
		_, err := client.ELB.Resize(id, map[string]interface{}{
			"flavor_id": d.Get("flavor_id").(string),
		})
		if err != nil {
			return fmt.Errorf("error when change flavor of ELB [%s]: %v", id, err)
		}
		_, err = waitUntilELBStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error when change flavor of ELB [%s]: %v", id, err)
		}
	}
	if d.HasChange("bandwidth_mbps") {
		_, err := client.ELB.Resize(id, map[string]interface{}{
			"bandwidth_mbps": d.Get("bandwidth_mbps").(int),
		})
		if err != nil {
			return fmt.Errorf("error when change Internet bandwidth of ELB [%s]: %v", id, err)
		}
		_, err = waitUntilELBStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("error when change bandwidth_mbps of ELB [%s]: %v", id, err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetLoadBalancerBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("error when update billing mode of LoadBalancer [%s]: %v", id, err)
		}
	}
	return resourceELBRead(d, meta)
}

func resourceELBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	bandwidth_mbps := d.Get("bandwidth_mbps").(int)
	if bandwidth_mbps > 0 && d.Get("network_type").(string) == "private" {
		return fmt.Errorf("bandwidth_mbps is not avaiable when network_type = private")
	}

	if len(d.Get("subnet_id").(string)) > 0 && d.Get("network_type").(string) == "public" {
		return fmt.Errorf("subnet_id is not avaiable when network_type = public")
	}
	elb, err := client.ELB.Create(map[string]interface{}{
		"name":           d.Get("name").(string),
		"description":    d.Get("description").(string),
		"flavor_id":      d.Get("flavor_id").(string),
		"zone":           d.Get("zone").(string),
		"network_type":   d.Get("network_type").(string),
		"subnet_id":      d.Get("subnet_id").(string),
		"tags":           d.Get("tags").(*schema.Set).List(),
		"billing_mode":   d.Get("billing_mode").(string),
		"bandwidth_mbps": d.Get("bandwidth_mbps").(int),
	})
	if err != nil {
		return fmt.Errorf("error creating ELB: %s", err)
	}
	d.SetId(elb.ID)
	_, err = waitUntilELBStatusChangedState(d, meta, []string{"ONLINE", "ACTIVE", "OFFLINE", "NO_MONITOR"}, []string{"ERROR", "DELETED", "DEGRADED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ELB: %s", err)
	}
	return resourceELBRead(d, meta)
}

func resourceELBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	elb, err := client.ELB.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ELB %s: %v", d.Id(), err)
	}

	network_type := "private"
	if !isPrivateIP(elb.VipAddress) {
		network_type = "public"
	}
	_ = d.Set("billing_mode", elb.BillingMode)
	_ = d.Set("zone", elb.AvailabilityZone)
	_ = d.Set("flavor_id", elb.FlavorID)
	_ = d.Set("name", elb.Name)
	_ = d.Set("network_type", network_type)
	_ = d.Set("created_at", elb.CreatedAt)
	_ = d.Set("tags", convertTagsToSet(elb.Tags))
	_ = d.Set("description", elb.Description)
	_ = d.Set("operating_status", elb.OperatingStatus)
	_ = d.Set("provisioning_status", elb.ProvisioningStatus)
	if network_type == "public" {
		_ = d.Set("bandwidth_mbps", elb.DomesticBandwidthMbps)
	} else {
		_ = d.Set("subnet_id", elb.VipSubnetID)
	}
	return nil
}

func resourceELBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ELB.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete ELB: %v", err)
	}
	_, err = waitUntilELBDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete ELB: %v", err)
	}
	return nil
}

func resourceELBImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceELBRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilELBDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.Get(id)
	})
}

func waitUntilELBStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ELB.Get(id)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.ELB).ProvisioningStatus
	})
}
