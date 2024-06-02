package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for CMC Cloud.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL use for the CMC Cloud API",
				DefaultFunc: schema.EnvDefaultFunc("CMC_CLOUD_API_ENDPOINT", "https://apiv2.cloud.cmctelecom.vn"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API key get from account settings in https://portalv2.cloud.cmctelecom.vn",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of project",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of region, eg hn-1,hcm-1",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cmccloudv2_server":                          resourceServer(),
			"cmccloudv2_volume":                          resourceVolume(),
			"cmccloudv2_volume_autobackup":               resourceVolumeAutoBackup(),
			"cmccloudv2_volume_attachment":               resourceVolumeAttachment(),
			"cmccloudv2_volume_snapshot":                 resourceVolumeSnapshot(),
			"cmccloudv2_volume_backup":                   resourceVolumeBackup(),
			"cmccloudv2_vpc":                             resourceVPC(),
			"cmccloudv2_subnet":                          resourceSubnet(),
			"cmccloudv2_eip":                             resourceEIP(),
			"cmccloudv2_eip_port_forwarding_rule":        resourceEIPPortForwardingRule(),
			"cmccloudv2_elb":                             resourceELB(),
			"cmccloudv2_elb_pool":                        resourceELBPool(),
			"cmccloudv2_elb_listener":                    resourceELBListener(),
			"cmccloudv2_elb_healthmonitor":               resourceELBHealthMonitor(),
			"cmccloudv2_elb_pool_member":                 resourceELBPoolMember(),
			"cmccloudv2_ecs_group":                       resourceEcsGroup(),
			"cmccloudv2_eip_port":                        resourceEIPPort(),
			"cmccloudv2_efs":                             resourceEFS(),
			"cmccloudv2_security_group":                  resourceSecurityGroup(),
			"cmccloudv2_kubernetes":                      resourceKubernates(),
			"cmccloudv2_kubernetes_nodegroup":            resourceKubernatesNodeGroup(),
			"cmccloudv2_database_configuration":          resourceDatabaseConfiguration(),
			"cmccloudv2_database_instance":               resourceDatabaseInstance(),
			"cmccloudv2_database_autobackup":             resourceDatabaseAutoBackup(),
			"cmccloudv2_autoscaling_group":               resourceAutoScalingGroup(),
			"cmccloudv2_autoscaling_configuration":       resourceAutoScalingConfiguration(),
			"cmccloudv2_autoscaling_health_check_policy": resourceAutoScalingHealthCheckPolicy(),
			"cmccloudv2_autoscaling_delete_policy":       resourceAutoScalingDeletePolicy(),
			"cmccloudv2_autoscaling_scale_in_policy":     resourceAutoScalingScaleInPolicy(),
			"cmccloudv2_autoscaling_scale_out_policy":    resourceAutoScalingScaleOutPolicy(),
			"cmccloudv2_autoscaling_az_policy":           resourceAutoScalingAZPolicy(),
			"cmccloudv2_autoscaling_lb_policy":           resourceAutoScalingLBPolicy(),
			"cmccloudv2_server_interface":                resourceServerInterface(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cmccloudv2_image":                     datasourceImage(),
			"cmccloudv2_flavor_ec":                 datasourceFlavorForEC(),
			"cmccloudv2_flavor_dbaas":              datasourceFlavorForDB(),
			"cmccloudv2_flavor_k8s":                datasourceFlavorForK8s(),
			"cmccloudv2_flavor_elb":                datasourceFlavorForELB(),
			"cmccloudv2_elb":                       datasourceELB(),
			"cmccloudv2_eip":                       datasourceEIP(),
			"cmccloudv2_ecs_group":                 datasourceEcsGroup(),
			"cmccloudv2_efs":                       datasourceEFS(),
			"cmccloudv2_vpc":                       datasourceVPC(),
			"cmccloudv2_subnet":                    datasourceSubnet(),
			"cmccloudv2_volume":                    datasourceVolume(),
			"cmccloudv2_volume_type":               datasourceVolumeType(),
			"cmccloudv2_volume_type_database":      datasourceVolumeTypeDatabase(),
			"cmccloudv2_server":                    datasourceServer(),
			"cmccloudv2_keypair":                   datasourceKeypair(),
			"cmccloudv2_backup":                    datasourceVolumeBackup(),
			"cmccloudv2_snapshot":                  datasourceVolumeSnapshot(),
			"cmccloudv2_autoscaling_configuration": datasourceAutoScalingConfiguration(),
			"cmccloudv2_autoscaling_group":         datasourceAutoScalingGroup(),
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	} /**/
	return p
}

/**/
func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		APIEndpoint: d.Get("api_endpoint").(string),
		APIKey:      d.Get("api_key").(string),
		ProjectId:   d.Get("project_id").(string),
		RegionId:    d.Get("region_id").(string),
		// TerraformVersion: terraformVersion,
	}
	return config.Client()
}
