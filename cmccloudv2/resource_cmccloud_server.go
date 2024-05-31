package cmccloudv2

import (
	"fmt"

	// "strconv"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
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
			// gocmcapiv2.Logs("CustomizeDiff")
			old, new := d.GetChange("volume_size")
			if old.(int) > new.(int) {
				return fmt.Errorf("Can't shrink volume_size, new `volume_size` must be > %d", old.(int))
			}
			// bỏ qua các thay đổi của trường nics
			if d.HasChange("nics") {
				d.Clear("nics")
			}
			return nil
		},
	}
}

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
		"subnets":              d.Get("nics").([]interface{}),
		// "subnets":              d.Get("nics").(*schema.Set).List(),
		// "eip_id":               d.Get("eip_id").(string),
		// "domestic_bandwidth":   d.Get("domestic_bandwidth").(int),
		// "inter_bandwidth":      d.Get("inter_bandwidth").(int),
	}
	res, err := client.Server.Create(datas)

	if err != nil {
		return fmt.Errorf("Error creating server: %v", err.Error())
	}
	d.SetId(res.Server.ID)
	// waitUntilServerChangeState(d, meta, res.Server.ID, []string{"building"}, []string{"active"})
	_, err = waitUntilServerStatusChangedState(d, meta, []string{"active"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("create server failed: %v", err)
	}
	return resourceServerRead(d, meta)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server, err := client.Server.Get(d.Id(), true)
	if err != nil {
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
	_ = d.Set("volumes", convertVolumeAttachs(server.VolumesAttached))
	_ = d.Set("nics", convertNics(server.Nics))

	// neu server chi co 1 nics / chua set id cho nics nao thi moi set gia tri nics, neu khong order co the thay doi
	// if curr_nics, ok := d.GetOkExists("nics"); !ok {
	// 	gocmcapiv2.Logo("curr_nics =", curr_nics)
	// 	if len(server.Nics) == 1 {
	// 		// chi co 1 nic
	// 	} else {
	// 		// kiem tra xem truoc do nics da set chua
	// 		if len(curr_nics.([]gocmcapiv2.Nic)) > 0 {
	// 			da_set_id := false
	// 			for _, nic := range curr_nics.([]gocmcapiv2.Nic) {
	// 				if nic.Id != "" {
	// 					da_set_id = true
	// 				}
	// 			}
	// 			if !da_set_id {
	// 				// chua set id cho nics nao thi moi set gia tri nics
	// 				_ = d.Set("nics", convertNics(server.Nics))
	// 			}
	// 		}
	// 	}
	// } else {
	// 	_ = d.Set("nics", convertNics(server.Nics))
	// }
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
			return fmt.Errorf("Error when set server tags [%s]: %v", id, err)
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
				return fmt.Errorf("Add security group [%s] to server [%s] error: %v", add.(string), d.Id(), err)
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
		// _, err = waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "stopped", "active"}, []string{"resized"})
		_, err = waitUntilServerStatusChangedState(d, meta, []string{"resized"}, []string{"error"})
		if err != nil {
			return fmt.Errorf("Resize server failed: %v", err)
		}

		_, err = client.Server.ConfirmResize(id)
		if err != nil {
			return fmt.Errorf("Error when resize server [%s]: %v", id, err)
		}
		// _, err = waitUntilServerChangeState(d, meta, d.Id(), []string{"resized"}, []string{"active", "stopped"})
		_, err = waitUntilServerStatusChangedState(d, meta, []string{"stopped", "active"}, []string{"error"})
		if err != nil {
			return fmt.Errorf("Resize server failed: %v", err)
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
				return fmt.Errorf("Error when start server: %v", err)
			}
			// waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "stopped"}, []string{"active"})
			_, err = waitUntilServerStatusChangedState(d, meta, []string{"active"}, []string{"error"})
			if err != nil {
				return fmt.Errorf("Start server failed: %v", err)
			}
		} else if newState.(string) == "stopped" {
			_, err := client.Server.Stop(d.Id())
			if err != nil {
				return fmt.Errorf("Error when stop server: %v", err)
			}
			// waitUntilServerChangeState(d, meta, d.Id(), []string{"building", "active"}, []string{"stopped"})
			_, err = waitUntilServerStatusChangedState(d, meta, []string{"stopped"}, []string{"error"})
			if err != nil {
				return fmt.Errorf("Stop server failed: %v", err)
			}
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
		return fmt.Errorf("Error delete server: %v", err)
	}
	_, err = waitUntilServerDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("Error delete server: %v", err)
	}
	return nil
}

func resourceServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceServerRead(d, meta)
	return []*schema.ResourceData{d}, err
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

func waitUntilServerDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Server.Get(id, false)
	})
}

func waitUntilServerStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 30 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Server.Get(id, false)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Server).VMState
	})
}
