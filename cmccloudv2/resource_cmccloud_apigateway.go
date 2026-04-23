package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApiGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiGatewayCreate,
		Read:   resourceApiGatewayRead,
		Update: resourceApiGatewayUpdate,
		Delete: resourceApiGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: resourceApiGatewayImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        apigatewaySchema(),
		// Validate that if public_access is true, bandwidth must be set
		CustomizeDiff: func(d *schema.ResourceDiff, meta interface{}) error {
			publicAccess := d.Get("public_access").(bool)
			bandwidth, ok := d.GetOk("bandwidth")
			if publicAccess && (!ok || bandwidth.(int) == 0) {
				return fmt.Errorf("when 'public_access' is true, 'bandwidth' must be set")
			}
			return nil
		},
	}
}

func resourceApiGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("subnet id is not valid %v", err)
	}

	zones := getStringArrayFromTypeSet(d.Get("zones").(*schema.Set))
	params := map[string]interface{}{
		"name":        d.Get("name").(string),
		"billingMode": "monthly",
		"zoneName":    zones[0],
		"zoneNames":   zones,
		"flavorId":    d.Get("flavor_id").(string),
		"volumeSize":  d.Get("volume_size").(int),
		"volumeType":  d.Get("volume_type").(string),
		"networkId":   subnet.NetworkID,
		"subnetId":    subnet.ID,
		"vpcId":       subnet.VpcID,
		"mode":        "standalone",
		"projectId":   client.Configs.ProjectId,
		"regionId":    client.Configs.RegionId,
		"createType":  "new",
		"backupId":    "", //d.Get("backup_id").(string),
	}

	if d.Get("public_access").(bool) {
		params["publicAccess"] = 1
		params["publicAccessBandwidth"] = d.Get("bandwidth").(int)
	} else {
		params["publicAccess"] = 0
		params["publicAccessBandwidth"] = 0
	}

	if backupId, ok := d.GetOk("backup_id"); ok {
		if backupStr, ok := backupId.(string); ok && backupStr != "" {
			params["createType"] = "restore"
		}
	}

	instance, err := client.ApiGateway.Create(params)
	if err != nil {
		return fmt.Errorf("error creating ApiGateway Instance: %s", err)
	}
	d.SetId(instance.Data.ID)

	_, err = client.Tag.UpdateTag(instance.Data.ID, "APIG", d)
	if err != nil {
		fmt.Printf("error updating ApiGateway tags: %s\n", err)
	}

	_, err = waitUntilApiGatewayJobFinished(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating ApiGateway Instance: %s", err)
	}
	return resourceApiGatewayRead(d, meta)
}

func resourceApiGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.ApiGateway.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving ApiGateway Instance %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("flavor_id", instance.Flavor.ID)
	_ = d.Set("volume_size", instance.VolumeSize)
	_ = d.Set("public_access", instance.PublicAccess)
	if !d.Get("public_access").(bool) {
		_ = d.Set("bandwith", nil)
	} else {
		_ = d.Set("bandwith", instance.PublicAccessBandwidth)
	}
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.CreatedAt)
	return nil
}

func resourceApiGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("tags") {
		_, err := client.Tag.UpdateTag(id, "APIG", d)
		if err != nil {
			return fmt.Errorf("error when set Api Gateway tags [%s]: %v", id, err)
		}
	}
	if d.HasChange("name") {
		_, err := client.ApiGateway.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("error when update Api Gateway name [%s]: %v", id, err)
		}
	}
	return resourceApiGatewayRead(d, meta)
}

func resourceApiGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.ApiGateway.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete Api Gateway instance: %v", err)
	}
	_, err = waitUntilApiGatewayDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete Api Gateway instance: %v", err)
	}
	return nil
}

func resourceApiGatewayImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceApiGatewayRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilApiGatewayJobFinished(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"COMPLETED", "RUNNING"}, []string{"ERROR", "FAILED", "STOPPED"}, WaitConf{
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ApiGateway.Get(id)
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.ApiGateway).Status)
	})
}

func waitUntilApiGatewayDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).ApiGateway.Get(id)
	})
}
