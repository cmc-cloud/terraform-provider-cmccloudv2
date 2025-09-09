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
			"cmccloudv2_server":                   resourceServer(),
			"cmccloudv2_volume":                   resourceVolume(),
			"cmccloudv2_volume_autobackup":        resourceVolumeAutoBackup(),
			"cmccloudv2_volume_attachment":        resourceVolumeAttachment(),
			"cmccloudv2_volume_snapshot":          resourceVolumeSnapshot(),
			"cmccloudv2_volume_backup":            resourceVolumeBackup(),
			"cmccloudv2_vpc":                      resourceVPC(),
			"cmccloudv2_subnet":                   resourceSubnet(),
			"cmccloudv2_eip":                      resourceEIP(),
			"cmccloudv2_eip_port_forwarding_rule": resourceEIPPortForwardingRule(),
			"cmccloudv2_elb":                      resourceELB(),
			"cmccloudv2_elb_pool":                 resourceELBPool(),
			"cmccloudv2_elb_listener":             resourceELBListener(),
			"cmccloudv2_elb_l7policy":             resourceELBL7policy(),
			"cmccloudv2_elb_l7policy_rule":        resourceELBL7policyRule(),
			"cmccloudv2_elb_healthmonitor":        resourceELBHealthMonitor(),
			"cmccloudv2_elb_pool_member":          resourceELBPoolMember(),
			"cmccloudv2_ecs_group":                resourceEcsGroup(),
			"cmccloudv2_eip_port":                 resourceEIPPort(),
			"cmccloudv2_efs":                      resourceEFS(),
			"cmccloudv2_security_group":           resourceSecurityGroup(),
			"cmccloudv2_kubernetes":               resourceKubernetes(),
			"cmccloudv2_kubernetes_nodegroup":     resourceKubernetesNodeGroup(),
			"cmccloudv2_kubernetesv2":             resourceKubernetesv2(),
			"cmccloudv2_kubernetesv2_nodegroup":   resourceKubernetesv2NodeGroup(),
			"cmccloudv2_database_configuration":   resourceDatabaseConfiguration(),
			"cmccloudv2_database_instance":        resourceDatabaseInstance(),
			"cmccloudv2_database_autobackup":      resourceDatabaseAutoBackup(),

			"cmccloudv2_autoscalingv2_group":         resourceAutoScalingV2Group(),
			"cmccloudv2_autoscalingv2_configuration": resourceAutoScalingV2Configuration(),
			"cmccloudv2_autoscalingv2_scale_trigger": resourceAutoScalingV2ScaleTrigger(),

			"cmccloudv2_autoscaling_group":               resourceAutoScalingGroup(),
			"cmccloudv2_autoscaling_configuration":       resourceAutoScalingConfiguration(),
			"cmccloudv2_autoscaling_health_check_policy": resourceAutoScalingHealthCheckPolicy(),
			"cmccloudv2_autoscaling_delete_policy":       resourceAutoScalingDeletePolicy(),
			"cmccloudv2_autoscaling_scale_in_policy":     resourceAutoScalingScaleInPolicy(),
			"cmccloudv2_autoscaling_scale_out_policy":    resourceAutoScalingScaleOutPolicy(),
			"cmccloudv2_autoscaling_az_policy":           resourceAutoScalingAZPolicy(),
			"cmccloudv2_autoscaling_lb_policy":           resourceAutoScalingLBPolicy(),
			"cmccloudv2_server_interface":                resourceServerInterface(),
			"cmccloudv2_redis_instance":                  resourceRedisInstance(),
			"cmccloudv2_redis_configuration":             resourceRedisConfiguration(),
			"cmccloudv2_keymanagement_container":         resourceKeyManagementContainer(),
			"cmccloudv2_keymanagement_secret":            resourceKeyManagementSecret(),
			"cmccloudv2_keymanagement_token":             resourceKeyManagementToken(),
			"cmccloudv2_devops_project":                  resourceDevopsProject(),
			"cmccloudv2_container_registry_repo":         resourceContainerRegistryRepository(),
			"cmccloudv2_va":                              resourceVA(),
			"cmccloudv2_waf":                             resourceWaf(),
			"cmccloudv2_waf_cert":                        resourceWafCert(),
			"cmccloudv2_waf_ip":                          resourceWafIP(),
			"cmccloudv2_waf_rule":                        resourceWafRule(),
			"cmccloudv2_waf_whitelist":                   resourceWafWhitelist(),
			"cmccloudv2_dns_zone":                        resourceDns(),
			"cmccloudv2_dns_record":                      resourceDnsRecord(),
			"cmccloudv2_dns_acl":                         resourceDnsAcl(),
			"cmccloudv2_cdn_cert":                        resourceCDNCert(),
			"cmccloudv2_cdn":                             resourceCDN(),
			"cmccloudv2_iam_group":                       resourceIamGroup(),
			"cmccloudv2_iam_user":                        resourceIamUser(),
			"cmccloudv2_iam_user_membership":             resourceIamUserMembership(),
			"cmccloudv2_iam_role_assignment":             resourceIamRoleAssignment(),
			"cmccloudv2_iam_user_server_permission":      resourceIamUserServerPermission(),
			"cmccloudv2_iam_custom_role":                 resourceIamCustomRole(),
			"cmccloudv2_iam_custom_role_assignment":      resourceIamCustomRoleAssignment(),
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
			"cmccloudv2_redis_configuration":       datasourceRedisConfiguration(),
			"cmccloudv2_security_group":            datasourceSecurityGroup(),
			"cmccloudv2_keymanagement_container":   datasourceKeyManagementContainer(),
			"cmccloudv2_keymanagement_secret":      datasourceKeyManagementSecret(),
			"cmccloudv2_certificate":               datasourceCertificate(),
			"cmccloudv2_iam_project":               datasourceIamProject(),
			"cmccloudv2_iam_group":                 datasourceIamGroup(),
			"cmccloudv2_iam_user":                  datasourceIamUser(),
			"cmccloudv2_iam_role":                  datasourceIamRole(),
			"cmccloudv2_iam_custom_role":           datasourceIamCustomRole(),
			"cmccloudv2_devops_project":            datasourceDevopsProject(),
			"cmccloudv2_container_registry_repo":   datasourceContainerRegistryRepository(),
			// "cmccloudv2_redis_instance":                  datasourceRedisInstance(),
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
