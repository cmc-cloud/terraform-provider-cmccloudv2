terraform {
  required_providers {
    cmccloudv2 = {
      source = "github.com/terraform-providers/cmccloudv2"
    }
  }
}
provider "cmccloudv2" {
    api_key = "your_api_key"  		# get from IAM Management
    project_id = "your_project_id"  # get from IAM Management
    region_id = "hn-1"
}
  
# Configure the provider for custom api endpoint
provider "cmccloudv2" {
    api_endpoint = "https://apiv2.cloud.cmctelecom.vn"
    api_key = "your_api_key" 
    project_id = "your_project_id"
    region_id = "hn-1"
}