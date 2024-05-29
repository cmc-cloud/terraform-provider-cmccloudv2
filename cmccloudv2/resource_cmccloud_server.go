package cmccloudv2

import (
	"errors"
	"fmt"
	"log"

	// "strconv"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceServerImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        serverSchema(),

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			old, new := d.GetChange("volume_size")
			if old.(int) > new.(int) {
				return fmt.Errorf("Can't shrink volume_size, new `volume_size` must be > %d", old.(int))
			}
			return nil
		},
	}
}

// type VolumeResourceData struct {
// 	// SourceId            string `json:"source_id"`
// 	// SourceType          string `json:"source_type"` // image,snapshot,volume
// 	Size                int    `json:"size"`
// 	Type                string `json:"type"`
// 	DeleteOnTermination bool   `json:"delete_on_termination"`
// }

// type NicResourceData struct {
// 	SubnetId  string `json:"subnet_id"` // image,snapshot,volume
// 	IpAddress string `json:"ip_address"`
// }

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	flavor_id := d.Get("flavor_id").(string)
	volumes := make([]map[string]interface{}, 1)
	volumes[0] = map[string]interface{}{
		"type": d.Get("volume_type").(string),
		"size": d.Get("volume_size").(int),
	}
	datas := map[string]interface{}{
		"project":              client.Configs.ProjectId,
		"server_name":          d.Get("name").(string),
		"zone":                 d.Get("zone").(string),
		"flavor_id":            flavor_id,
		"subnets":              d.Get("nics").(*schema.Set).List(),
		"volumes":              volumes,
		"security_group_names": d.Get("security_group_names").(*schema.Set).List(),
		"ecs_group_id":         d.Get("ecs_group_id").(string),
		"key_name":             d.Get("key_name").(string),
		"user_data":            d.Get("user_data").(string),
		"password":             d.Get("password").(string),
		"billing_mode":         d.Get("billing_mode").(string),
		"source_type":          d.Get("source_type").(string),
		"source_id":            d.Get("source_id").(string),
		"tags":                 d.Get("tags").(*schema.Set).List(),
		// "eip_id":               d.Get("eip_id").(string),
		// "domestic_bandwidth":   d.Get("domestic_bandwidth").(int),
		// "inter_bandwidth":      d.Get("inter_bandwidth").(int),
	}
	res, err := client.Server.Create(datas)

	if err != nil {
		return fmt.Errorf("Error creating server: %v", err.Error())
	}
	// gocmcapiv2.Logs("set id " + res.Server.ID)
	d.SetId(res.Server.ID)
	// _ = d.Set("password", res.Server.AdminPass)
	// _ = d.Set("flavor_id", flavor_id)
	// _ = d.Set("source_type", d.Get("source_type").(string))
	// _ = d.Set("source_id", d.Get("source_id").(string))

	waitUntilServerChangeState(d, meta, res.Server.ID, []string{"building"}, []string{"active"})
	return resourceServerRead(d, meta)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server, err := client.Server.Get(d.Id(), true)
	if err != nil {
		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			log.Printf("[WARN] CMC Cloud Server with id = (%s) is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving server %s: %v", d.Id(), err)
	}
	_ = d.Set("name", server.Name)
	_ = d.Set("zone", server.AvailabilityZone)

	if server.Flavor.ID != "" {
		_ = d.Set("flavor_id", server.Flavor.ID)
	}
	_ = d.Set("security_group_names", convertSecurityGroups(server.SecurityGroups))
	_ = d.Set("ecs_group_id", strings.Join(server.ServerGroups, ","))
	_ = d.Set("created", server.Created)
	_ = d.Set("tags", server.Tags)
	_ = d.Set("description", server.Description)
	_ = d.Set("billing_mode", server.BillingMode)
	_ = d.Set("vm_state", server.VMState)
	_ = d.Set("nics", convertNics(server.Nics))
	_ = d.Set("volumes", convertVolumeAttachs(server.VolumesAttached))
	// _ = d.Set("volumes", convertVolumes(server.Volumes))
	gocmcapiv2.Logo("volumes is = ", convertVolumeAttachs(server.VolumesAttached))
	if server.KeyName != "" {
		_ = d.Set("key_name", server.KeyName)
	}
	return nil
}
func convertSecurityGroups(groups []gocmcapiv2.ServerSecurityGroup) []string {
	seen := make(map[string]bool)
	var result []string

	for _, group := range groups {
		// If the element is not in the map, add it to the result slice
		if !seen[group.Name] {
			seen[group.Name] = true
			result = append(result, group.Name)
		}
	}

	return result
}

func resourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	id := d.Id()

	if d.HasChange("name") {
		_, err := client.Server.Rename(id, d.Get("name").(string))
		if err != nil {
			return fmt.Errorf("Error when rename server [%s]: %v", id, err)
		}
	}

	if d.HasChange("tags") {
		_, err := client.Server.SetTags(id, d.Get("tags").(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("Error when rename server [%s]: %v", id, err)
		}
	}

	if d.HasChange("security_group_names") {
		removed, added := getDiffSet(d.GetChange("security_group_names"))
		for _, remove := range removed.List() {
			// Logic xử lý phần tử bị xóa
			_, err := client.Server.RemoveSecurityGroup(d.Id(), remove.(string))
			if err != nil {
				return fmt.Errorf("Remove security group [%s] from server [%s] error: %v", remove.(string), d.Id(), err)
			}
		}
		for _, add := range added.List() {
			// Logic xử lý phần tử add them
			_, err := client.Server.AddSecurityGroup(d.Id(), add.(string))
			if err != nil {
				return fmt.Errorf("Add security group [%s] from server [%s] error: %v", add.(string), d.Id(), err)
			}
		}
	}

	if d.HasChange("volume_size") {
		vol, _ := client.Server.Get(d.Id(), false)
		volume_id := vol.VolumesAttached[0].ID
		_, err := client.Volume.Resize(volume_id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("Error when resize volume of server (%s): %v", d.Id(), err)
		}
	}

	if d.HasChange("password") {
		_, err := client.Server.ChangePassword(id, d.Get("password").(string))
		if err != nil {
			return fmt.Errorf("Error when reseting password of server (%s): %v", d.Id(), err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetServerBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("Error when update billing mode of Server [%s]: %v", id, err)
		}
	}

	if d.HasChange("flavor_id") {
		// Resize server to new flavor
		_, err := client.Server.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		_, err = waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "stopped", "active"}, []string{"resized"})
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		_, err = client.Server.ConfirmResize(id)
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		_, err = waitUntilServerChangeState(d, meta, d.Id(), []string{"resized"}, []string{"active", "stopped"})
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
	}

	if d.HasChange("vm_state") {
		oldState, newState := d.GetChange("vm_state")
		if oldState.(string) == "error" {
			return fmt.Errorf("You cannot change server state because old server state is %s", oldState.(string))
		}
		if newState.(string) == "active" {
			_, err := client.Server.Start(d.Id())
			if err != nil {
				return fmt.Errorf("Error when changing state of server: %v", err)
			}

			log.Printf("Wait until server state to be active", d.Id())
			waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "stopped"}, []string{"active"})
		} else if newState.(string) == "stopped" {
			_, err := client.Server.Stop(d.Id())
			if err != nil {
				return fmt.Errorf("Error when changing state of server: %v", err)
			}
			log.Printf("Wait until server state to be stopped", d.Id())
			waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "active"}, []string{"stopped"})
		} else {
			return fmt.Errorf("New state of server must be 'active' or 'stopped'")
		}
	}

	return resourceServerRead(d, meta)
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Server.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("Error delete cloud server: %v", err)
	}
	return nil
}

func resourceServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceServerRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilServerChangeState(d *schema.ResourceData, meta interface{}, id string, pendingStatus []string, targetStatus []string) (interface{}, error) {
	log.Printf("[INFO] Waiting for server with id (%s) to be "+strings.Join(targetStatus, ","), id)
	stateConf := &resource.StateChangeConf{
		Pending:        pendingStatus,
		Target:         targetStatus,
		Refresh:        serverStateRefreshfunc(d, meta, id),
		Timeout:        1200 * time.Second,
		Delay:          20 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
	// return stateConf.WaitForStateContext(context.Background())
}

func serverStateRefreshfunc(d *schema.ResourceData, meta interface{}, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := meta.(*CombinedConfig).goCMCClient()
		server, err := client.Server.Get(d.Id(), false)
		if err != nil {
			fmt.Errorf("Error retrieving server %s: %v", id, err)
			return nil, "", err
		}
		log.Println("[DEBUG] Server status = " + server.VMState)
		gocmcapiv2.Logs("[DEBUG] Server " + d.Id() + " status = " + server.VMState)
		return server, server.VMState, nil
	}
}
func convertVolumeAttachs(vols []gocmcapiv2.VolumeAttach) []map[string]interface{} {
	result := make([]map[string]interface{}, len(vols))
	for i, vol := range vols {
		result[i] = map[string]interface{}{
			"id":                    vol.ID,
			"delete_on_termination": vol.DeleteOnTermination,
		}
	}
	return result
}

func convertNics(nics []gocmcapiv2.Nic) []map[string]interface{} {
	result := make([]map[string]interface{}, len(nics))
	for i, nic := range nics {
		result[i] = map[string]interface{}{
			// "id":                 nic.Id,
			"subnet_id":  nic.FixedIps[0].SubnetID,
			"ip_address": nic.FixedIps[0].IPAddress,
			// "security_group_ids": nic.SecurityGroups,
			// "mac_address": nic.MacAddress,
		}
	}
	return result
}
