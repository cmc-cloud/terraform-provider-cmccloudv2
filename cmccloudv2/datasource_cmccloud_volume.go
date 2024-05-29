package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVolumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"volume_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the volume",
		},
		"server_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the server that volume attached, match exactly",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of volume (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "Filter by status of volume (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"zone": {
			Type:        schema.TypeString,
			Description: "Filter by zone of that volume (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"bootable": {
			Type:        schema.TypeBool,
			Description: "filter by volume bootable",
			Optional:    true,
			ForceNew:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "filter by volume that contains this text (case-insensitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceVolume() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceVolumeRead,
		Schema: datasourceVolumeSchema(),
	}
}

func dataSourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allVolumes []gocmcapiv2.Volume
	if volume_id := d.Get("volume_id").(string); volume_id != "" {
		volume, err := client.Volume.Get(volume_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve volume [%s]: %s", volume_id, err)
			}
		}
		allVolumes = append(allVolumes, volume)
	} else {
		params := map[string]string{
			"name":      d.Get("name").(string),
			"status":    d.Get("status").(string),
			"server_id": d.Get("server_id").(string),
			"zone":      d.Get("zone").(string),
			"bootable":  strconv.FormatBool(d.Get("bootable").(bool)),
		}
		volumes, err := client.Volume.List(params)
		if err != nil {
			return fmt.Errorf("Error when get volumes %v", err)
		}
		allVolumes = append(allVolumes, volumes...)
	}
	if len(allVolumes) > 0 {
		var filteredVolumes []gocmcapiv2.Volume
		for _, volume := range allVolumes {
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(volume.Description), strings.ToLower(v)) {
					continue
				}
			}

			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(volume.Name), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("status").(string); v != "" {
				if !strings.Contains(strings.ToLower(volume.Status), strings.ToLower(v)) {
					continue
				}
			}
			if v := d.Get("zone").(string); v != "" {
				if !strings.Contains(strings.ToLower(volume.AvailabilityZone), strings.ToLower(v)) {
					continue
				}
			}
			if v, ok := d.GetOkExists("bootable"); ok {
				boot, _ := strconv.ParseBool(volume.Bootable)
				if boot != v.(bool) {
					continue
				}
			}

			filteredVolumes = append(filteredVolumes, volume)
		}
		allVolumes = filteredVolumes
	}
	if len(allVolumes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allVolumes) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allVolumes)
		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeVolumeAttributes(d, allVolumes[0])
}

func dataSourceComputeVolumeAttributes(d *schema.ResourceData, volume gocmcapiv2.Volume) error {
	log.Printf("[DEBUG] Retrieved volume %s: %#v", volume.ID, volume)
	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("bootable", volume.Bootable)
	d.Set("zone", volume.AvailabilityZone)
	d.Set("description", volume.Description)
	d.Set("created_at", volume.CreatedAt)
	server_id := ""
	if len(volume.Attachments) > 0 {
		server_id = volume.Attachments[0].ServerID
	}
	d.Set("server_id", server_id)
	return nil
}
