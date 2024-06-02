package cmccloudv2

import (
	"fmt"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the certificate",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of certificate, match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func datasourceCertificate() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceCertificateRead,
		Schema: datasourceCertificateSchema(),
	}
}

func dataSourceCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allCertificates []gocmcapiv2.Certificate
	if certificate_id := d.Get("certificate_id").(string); certificate_id != "" {
		certificate, err := client.Certificate.Get(certificate_id)
		if err != nil {
			return fmt.Errorf("Unable to retrieve certificate [%s]: %s", certificate_id, err)
		}
		allCertificates = append(allCertificates, certificate)
	} else {
		params := map[string]string{}
		certificates, err := client.Certificate.List(params)
		if err != nil {
			return fmt.Errorf("Error when get certificates %v", err)
		}
		allCertificates = append(allCertificates, certificates...)
	}
	if len(allCertificates) > 0 {
		var filteredCertificates []gocmcapiv2.Certificate
		for _, certificate := range allCertificates {
			if v := d.Get("name").(string); v != "" {
				if v != certificate.Name {
					continue
				}
			}
			filteredCertificates = append(filteredCertificates, certificate)
		}
		allCertificates = filteredCertificates
	}
	if len(allCertificates) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allCertificates) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allCertificates)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeCertificateAttributes(d, allCertificates[0])
}

func dataSourceComputeCertificateAttributes(d *schema.ResourceData, certificate gocmcapiv2.Certificate) error {
	d.SetId(certificate.Name)
	d.Set("name", certificate.Name)
	d.Set("certificate_id", certificate.Name)
	d.Set("created_at", certificate.Created)
	return nil
}
