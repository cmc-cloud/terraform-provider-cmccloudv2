package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVolumeSnapshotSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"snapshot_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the snapshot",
		},
		"volume_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the volume",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of snapshot (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "Filter by status of snapshot (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"is_latest": {
			Type:        schema.TypeBool,
			Description: "true if you want to get the latest snapshot that match other filter",
			Optional:    true,
			ForceNew:    true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
		"real_size_gb": {
			Type:     schema.TypeFloat,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceVolumeSnapshotRead,
		Schema: datasourceVolumeSnapshotSchema(),
	}
}

func dataSourceVolumeSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allSnapshots []gocmcapiv2.Snapshot
	if snapshot_id := d.Get("snapshot_id").(string); snapshot_id != "" {
		volume, err := client.Snapshot.Get(snapshot_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("Unable to retrieve volume [%s]: %s", snapshot_id, err)
			}
		}
		allSnapshots = append(allSnapshots, volume)
	} else {
		params := map[string]string{
			"name":      d.Get("name").(string),
			"status":    d.Get("status").(string),
			"volume_id": d.Get("volume_id").(string),
		}
		snapshots, err := client.Snapshot.List(params)
		if err != nil {
			return fmt.Errorf("Error when get snapshots %v", err)
		}
		allSnapshots = append(allSnapshots, snapshots...)
	}
	if len(allSnapshots) > 0 {
		var filteredSnapshots []gocmcapiv2.Snapshot
		for _, snapshot := range allSnapshots {
			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(snapshot.Name), strings.ToLower(v)) {
					continue
				}
			}

			if v := d.Get("status").(string); v != "" {
				if !strings.Contains(strings.ToLower(snapshot.Status), strings.ToLower(v)) {
					continue
				}
			}
			filteredSnapshots = append(filteredSnapshots, snapshot)
		}
		allSnapshots = filteredSnapshots
	}
	if len(allSnapshots) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	if len(allSnapshots) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allSnapshots)

		if v, ok := d.GetOkExists("is_latest"); ok {
			if v.(bool) {
				// lay ban snapshot dau tien vi snapshot duoc list theo thu tu tao gan nhat truoc
				return dataSourceComputeVolumeSnapshotAttributes(d, allSnapshots[0])
			}
		}

		return fmt.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeVolumeSnapshotAttributes(d, allSnapshots[0])
}

func dataSourceComputeVolumeSnapshotAttributes(d *schema.ResourceData, snapshot gocmcapiv2.Snapshot) error {
	log.Printf("[DEBUG] Retrieved volume %s: %#v", snapshot.ID, snapshot)
	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("status", snapshot.Status)
	d.Set("volume_id", snapshot.VolumeID)
	d.Set("real_size_gb", snapshot.RealSizeGB)
	d.Set("created_at", snapshot.CreatedAt)
	return nil
}
