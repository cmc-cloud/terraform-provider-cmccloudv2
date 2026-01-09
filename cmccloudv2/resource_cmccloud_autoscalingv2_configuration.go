package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAutoScalingV2Configuration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoScalingV2ConfigurationCreate,
		Read:   resourceAutoScalingV2ConfigurationRead,
		Update: resourceAutoScalingV2ConfigurationUpdate,
		Delete: resourceAutoScalingV2ConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAutoScalingV2ConfigurationImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        autoscalingV2ConfigurationSchema(),
	}
}

func resourceAutoScalingV2ConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	datas := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"source_type":          d.Get("source_type").(string),
		"source_id":            d.Get("source_id").(string),
		"flavor_id":            d.Get("flavor_id").(string),
		"subnet_ids":           d.Get("subnet_ids").([]interface{}),
		"use_eip":              d.Get("use_eip").(bool),
		"domestic_bandwidth":   d.Get("domestic_bandwidth").(int),
		"inter_bandwidth":      d.Get("inter_bandwidth").(int),
		"volumes":              d.Get("volumes").([]interface{}),
		"security_group_names": d.Get("security_group_names").(*schema.Set).List(),
		"key_name":             d.Get("key_name").(string),
		"user_data":            d.Get("user_data").(string),
		"password":             d.Get("password").(string),
		"ecs_group_id":         d.Get("ecs_group_id").(string),
	}
	res, err := client.AutoScalingV2Configuration.Create(datas)

	if err != nil {
		return fmt.Errorf("error creating configuration: %v", err.Error())
	}
	d.SetId(res.ID)
	return resourceAutoScalingV2ConfigurationRead(d, meta)
}

func resourceAutoScalingV2ConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	configuration, err := client.AutoScalingV2Configuration.Get(d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving configuration %s: %v", d.Id(), err)
	}
	_ = d.Set("name", configuration.Name)
	_ = d.Set("source_type", configuration.SourceType)
	_ = d.Set("source_id", configuration.SourceID)
	_ = d.Set("flavor_id", configuration.FlavorID)
	_ = d.Set("subnet_ids", configuration.SubnetIds)
	_ = d.Set("use_eip", configuration.UseEip)
	setBool(d, "use_eip", configuration.UseEip)
	setInt(d, "domestic_bandwidth", configuration.DomesticBandwidth)
	setInt(d, "inter_bandwidth", configuration.InterBandwidth)
	_ = d.Set("volumes", v2ConvertConfigVolumes(configuration.Volumes))
	// _ = d.Set("security_group_names", configuration.SecurityGroupNames)
	setTypeSet(d, "security_group_names", configuration.SecurityGroupNames)
	setString(d, "ecs_group_id", configuration.EcsGroupID)
	setString(d, "user_data", configuration.UserData)
	setString(d, "key_name", configuration.KeyName)
	_ = d.Set("created", configuration.Created)

	return nil
}

func v2ConvertConfigVolumes(volumes []gocmcapiv2.AutoScalingV2ConfigurationVolume) []map[string]interface{} {
	result := make([]map[string]interface{}, len(volumes))
	for i, volume := range volumes {
		result[i] = map[string]interface{}{
			"type":                  volume.Type,
			"size":                  volume.Size,
			"delete_on_termination": volume.DeleteOnTermination,
		}
	}
	return result
}
func resourceAutoScalingV2ConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	datas := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"source_type":          d.Get("source_type").(string),
		"source_id":            d.Get("source_id").(string),
		"flavor_id":            d.Get("flavor_id").(string),
		"subnet_ids":           d.Get("subnet_ids").([]interface{}),
		"use_eip":              d.Get("use_eip").(bool),
		"domestic_bandwidth":   d.Get("domestic_bandwidth").(int),
		"inter_bandwidth":      d.Get("inter_bandwidth").(int),
		"volumes":              d.Get("volumes").([]interface{}),
		"security_group_names": d.Get("security_group_names").(*schema.Set).List(),
		"key_name":             d.Get("key_name").(string),
		"user_data":            d.Get("user_data").(string),
		"password":             d.Get("password").(string),
		"ecs_group_id":         d.Get("ecs_group_id").(string),
	}

	_, err := client.AutoScalingV2Configuration.Update(id, datas)
	if err != nil {
		return fmt.Errorf("error when rename autoscale configuration [%s]: %v", id, err)
	}
	return resourceAutoScalingV2ConfigurationRead(d, meta)
}

func resourceAutoScalingV2ConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.AutoScalingV2Configuration.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete autoscale configuration: %v", err)
	}
	_, err = waitUntilAutoScalingV2ConfigurationDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete autoscale configuration: %v", err)
	}
	return nil
}

func resourceAutoScalingV2ConfigurationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceAutoScalingV2ConfigurationRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilAutoScalingV2ConfigurationDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      3 * time.Second,
		MinTimeout: 10 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).AutoScalingV2Configuration.Get(id)
	})
}
