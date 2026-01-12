package cmccloudv2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func cdnSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"vod": {
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
			Description: "If true, it will be the VOD site",
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The name of the CDN",
		},
		"origin_type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"host", "s3"}, false),
			Description:  "The origin type of the CDN: Host Origin or S3 Origin",
		},
		"domain_or_ip": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.Any(validation.IsIPAddress, validateDomainName),
			Description:  "The domain or IP of the CDN. Only valid for Host Origin type",
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"http", "https"}, false),
			Description:  "The protocol of the CDN, `http`, `https`. Only valid for Host Origin type",
		},
		"port": {
			Type:         schema.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsPortNumber,
			Description:  "The port of the CDN. Only valid for Host Origin type",
		},
		"origin_path": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The origin path of the CDN, example `/content/`. Only valid for Host Origin type & vod is false",
		},

		"s3_access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The access key of the S3 Origin. Only valid for S3 Origin type",
		},
		"s3_secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The secret key of the S3 Origin. Only valid for S3 Origin type",
		},
		"s3_bucket_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The bucket name of the S3 Origin. Only valid for S3 Origin type",
		},
		"s3_region": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The region of the S3 Origin. Only valid for S3 Origin type",
		},
		"s3_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The endpoint of the S3 Origin. Only valid for S3 Origin type",
		},

		"status": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The status of the CDN",
		},
		"cdn_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The CDN URL of the CDN",
		},
		"multi_cdn_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The multi CDN URL of the CDN",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The updated time of the CDN",
		},
	}
}
