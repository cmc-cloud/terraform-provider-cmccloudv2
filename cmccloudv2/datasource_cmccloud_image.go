package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceImageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"image_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the image",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of image (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"visibility": {
			Type:         schema.TypeString,
			Description:  "Visibility of image, accept public/private/shared",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"public", "private", "shared"}, false),
			ForceNew:     true,
		},
		"os": {
			Type:        schema.TypeString,
			Description: "Filter by os name of image (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
	}
}

func datasourceImage() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceImageRead,
		Schema: datasourceImageSchema(),
	}
}

func dataSourceImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allImages []gocmcapiv2.Image
	if image_id := d.Get("image_id").(string); image_id != "" {
		image, err := client.Image.Get(image_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve flavor [%s]: %s", image_id, err)
			}
		}
		allImages = append(allImages, image)
	} else {
		params := map[string]string{
			"visibility": d.Get("visibility").(string),
			"status":     "active",
		}
		images, err := client.Image.List(params)
		if err != nil {
			return fmt.Errorf("Error when get images %v", err)
		}
		allImages = append(allImages, images...)
	}
	if len(allImages) > 0 {
		var filteredImages []gocmcapiv2.Image
		for _, image := range allImages {
			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(image.Name), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("os").(string); v != "" {
				if !strings.Contains(strings.ToLower(image.OsDistro), strings.ToLower(v)) {
					continue
				}
			}
			filteredImages = append(filteredImages, image)
		}
		allImages = filteredImages
	}
	if len(allImages) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allImages) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allImages)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeImageAttributes(d, allImages[0])
}

func dataSourceComputeImageAttributes(d *schema.ResourceData, image gocmcapiv2.Image) error {
	log.Printf("[DEBUG] Retrieved image %s: %#v", image.ID, image)
	d.SetId(image.ID)
	d.Set("name", image.Name)
	d.Set("os", image.Os)
	d.Set("image_id", image.ID)
	return nil
}
