package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceVolumeBackupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"backup_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the backup",
		},
		"volume_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of the volume",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Filter by name of backup (case-insenitive)",
			Optional:    true,
			ForceNew:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "Filter by status of backup (case-insenitive), match exactly",
			Optional:    true,
			ForceNew:    true,
		},
		"is_latest": {
			Type:        schema.TypeBool,
			Description: "true if you want to get the latest backup that match other filter",
			Optional:    true,
			ForceNew:    true,
		},
		"is_incremental": {
			Type:        schema.TypeBool,
			Description: "filter by volume bootable",
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

func datasourceVolumeBackup() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceVolumeBackupRead,
		Schema: datasourceVolumeBackupSchema(),
	}
}

func dataSourceVolumeBackupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allBackups []gocmcapiv2.Backup
	if backup_id := d.Get("backup_id").(string); backup_id != "" {
		volume, err := client.Backup.Get(backup_id)
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve volume [%s]: %s", backup_id, err)
			}
		}
		allBackups = append(allBackups, volume)
	} else {
		params := map[string]string{
			"name":      d.Get("name").(string),
			"status":    d.Get("status").(string),
			"volume_id": d.Get("volume_id").(string),
		}
		backups, err := client.Backup.List(params)
		if err != nil {
			return fmt.Errorf("error when get backups %v", err)
		}
		allBackups = append(allBackups, backups...)
	}
	if len(allBackups) > 0 {
		var filteredBackups []gocmcapiv2.Backup
		for _, backup := range allBackups {
			if v := d.Get("name").(string); v != "" {
				if !strings.Contains(strings.ToLower(backup.Name), strings.ToLower(v)) {
					continue
				}
			}

			if v := d.Get("status").(string); v != "" {
				if !strings.Contains(strings.ToLower(backup.Status), strings.ToLower(v)) {
					continue
				}
			}
			if v, ok := d.GetOk("is_incremental"); ok {
				// user explicitly set trong .tf (true hoặc false)
				if backup.IsIncremental != v.(bool) {
					continue
				}
			} //else {
			// user không set, đang dùng default
			//}

			// if v, ok := d.GetOkExists("is_incremental"); ok {
			// 	if backup.IsIncremental != v.(bool) {
			// 		continue
			// 	}
			// }
			filteredBackups = append(filteredBackups, backup)
		}
		allBackups = filteredBackups
	}
	if len(allBackups) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allBackups) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allBackups)

		// if v, ok := d.GetOkExists("is_latest"); ok {
		// 	if v.(bool) {
		// 		// lay ban backup dau tien vi backup duoc list theo thu tu tao gan nhat truoc
		// 		return dataSourceComputeVolumeBackupAttributes(d, allBackups[0])
		// 	}
		// }
		if v, ok := d.GetOk("is_latest"); ok {
			// user explicitly set trong .tf (true hoặc false)
			if v.(bool) {
				return dataSourceComputeVolumeBackupAttributes(d, allBackups[0])
			}
		} //else {
		// user không set, đang dùng default
		//}

		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeVolumeBackupAttributes(d, allBackups[0])
}

func dataSourceComputeVolumeBackupAttributes(d *schema.ResourceData, backup gocmcapiv2.Backup) error {
	log.Printf("[DEBUG] Retrieved volume %s: %#v", backup.ID, backup)
	d.SetId(backup.ID)
	_ = d.Set("name", backup.Name)
	_ = d.Set("status", backup.Status)
	_ = d.Set("is_incremental", backup.IsIncremental)
	_ = d.Set("volume_id", backup.VolumeID)
	_ = d.Set("real_size_gb", backup.RealSizeGB)
	_ = d.Set("created_at", backup.CreatedAt)
	return nil
}
