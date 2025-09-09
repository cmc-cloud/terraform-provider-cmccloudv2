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
			old, new := d.GetChange("volume_size")
			if old.(int) > new.(int) {
				return fmt.Errorf("can't shrink volume_size, new `volume_size` must be > %d", old.(int))
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
		"type":                  d.Get("volume_type").(string),
		"size":                  d.Get("volume_size").(int),
		"delete_on_termination": d.Get("delete_on_termination").(bool),
	}

	subnets := make([]map[string]interface{}, 1)
	subnets[0] = map[string]interface{}{
		"subnet_id": d.Get("subnet_id").(string),
	}
	if d.Get("ip_address").(string) != "" {
		subnet, err := client.Subnet.Get(d.Get("subnet_id").(string))

		if err != nil {
			return fmt.Errorf("error when getting subnet info: %v", err)
		}
		_, err = isIpBelongToCidr(d.Get("ip_address").(string), subnet.Cidr)
		if err != nil {
			return err
		}
		subnets[0]["ip_address"] = d.Get("ip_address").(string)
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
		"subnets":              subnets, // d.Get("nics").([]interface{}),
		// "subnets":              d.Get("nics").(*schema.Set).List(),
		// "eip_id":               d.Get("eip_id").(string),
		// "domestic_bandwidth":   d.Get("domestic_bandwidth").(int),
		// "inter_bandwidth":      d.Get("inter_bandwidth").(int),
	}
	res, err := client.Server.Create(datas)

	if err != nil {
		return fmt.Errorf("error creating server: %v", err.Error())
	}
	d.SetId(res.Server.ID)
	_, err = waitUntilServerStatusChangedState(d, meta, []string{"active"}, []string{"error"})
	if err != nil {
		return fmt.Errorf("create server failed: %v", err)
	}
	return readOrImport(d, meta, false)
}

func readOrImport(d *schema.ResourceData, meta interface{}, isImport bool) error {
	client := meta.(*CombinedConfig).goCMCClient()
	server, err := client.Server.Get(d.Id(), true)
	if err != nil {
		return fmt.Errorf("error retrieving server %s: %v", d.Id(), err)
	}
	_ = d.Set("name", server.Name)
	_ = d.Set("zone", server.AvailabilityZone)

	if server.Flavor.ID != "" {
		_ = d.Set("flavor_id", server.Flavor.ID)
	}
	_ = d.Set("security_group_names", convertSecurityGroups(server.SecurityGroups))
	_ = d.Set("ecs_group_id", strings.Join(server.ServerGroups, ","))
	_ = d.Set("created", server.Created)
	_ = d.Set("tags", convertTagsToSet(server.Tags))
	_ = d.Set("description", server.Description)
	_ = d.Set("billing_mode", server.BillingMode)
	_ = d.Set("vm_state", server.VMState)
	_ = d.Set("volumes", convertVolumeAttachs(server.VolumesAttached))
	if len(server.VolumesAttached) > 0 {
		_ = d.Set("delete_on_termination", server.VolumesAttached[0].DeleteOnTermination)
	}
	if len(server.Nics) > 0 {
		if isImport {
			// khong set, neu set co the bi sai neu co >= 2 interfaces
			_ = d.Set("subnet_id", server.Nics[0].FixedIps[0].SubnetID)
		}

		for _, nic := range server.Nics {
			if nic.FixedIps[0].SubnetID == d.Get("subnet_id").(string) {
				// chi set ngay sau khi tao server, vi khi do chua add them interface nao
				// neu add >= 2 interfaces thi thu tu interface co the thay doi => bi doi interface_id
				if d.Get("interface_id").(string) == "" {
					_ = d.Set("interface_id", nic.ID)
				}
				// chua set thi moi set
				if d.Get("ip_address").(string) == "" || d.Get("ip_address").(string) == nic.FixedIps[0].IPAddress {
					setString(d, "ip_address", nic.FixedIps[0].IPAddress)
				}
				break
			}
		}
	}
	// if isImport {
	// 	_ = d.Set("nics", convertNics(server.Nics))
	// }
	if server.KeyName != "" {
		_ = d.Set("key_name", server.KeyName)
	}
	return nil
}
func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	return readOrImport(d, meta, false)
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
			return fmt.Errorf("error when rename server [%s]: %v", id, err)
		}
	}

	if d.HasChange("tags") {
		_, err := client.Server.SetTags(id, d.Get("tags").(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("error when set server tags [%s]: %v", id, err)
		}
	}

	if d.HasChange("security_group_names") {
		removed, added := getDiffSet(d.GetChange("security_group_names"))
		for _, remove := range removed.List() {
			// Logic xử lý phần tử bị xóa
			_, err := client.Server.RemoveSecurityGroup(d.Id(), remove.(string))
			if err != nil {
				return fmt.Errorf("remove security group [%s] from server [%s] error: %v", remove.(string), d.Id(), err)
			}
		}
		for _, add := range added.List() {
			// Logic xử lý phần tử add them
			_, err := client.Server.AddSecurityGroup(d.Id(), add.(string))
			if err != nil {
				return fmt.Errorf("add security group [%s] to server [%s] error: %v", add.(string), d.Id(), err)
			}
		}
	}

	if d.HasChange("volume_size") {
		vol, _ := client.Server.Get(d.Id(), false)
		volume_id := vol.VolumesAttached[0].ID
		_, err := client.Volume.Resize(volume_id, d.Get("volume_size").(int))
		if err != nil {
			return fmt.Errorf("error when resize volume of server (%s): %v", d.Id(), err)
		}
	}

	if d.HasChange("password") {
		_, err := client.Server.ChangePassword(id, d.Get("password").(string))
		if err != nil {
			return fmt.Errorf("error when reseting password of server (%s): %v", d.Id(), err)
		}
	}

	if d.HasChange("billing_mode") {
		_, err := client.BillingMode.SetServerBilingMode(id, d.Get("billing_mode").(string))
		if err != nil {
			return fmt.Errorf("error when update billing mode of Server [%s]: %v", id, err)
		}
	}

	if d.HasChange("flavor_id") {
		// Resize server to new flavor
		_, err := client.Server.Resize(id, d.Get("flavor_id").(string))
		if err != nil {
			return fmt.Errorf("error when resize server [%s]: %v", id, err)
		}
		time.Sleep(5 * time.Second)
		// _, err = waitUntilServerStatusChangedState(d, meta, []string{"resized"}, []string{"error"})
		// if err != nil {
		// 	return fmt.Errorf("resize server failed: %v", err)
		// }

		// _, err = client.Server.ConfirmResize(id)
		// if err != nil {
		// 	return fmt.Errorf("error when resize server [%s]: %v", id, err)
		// }
		_, err = waitUntilServerStatusChangedState(d, meta, []string{"stopped", "active"}, []string{"error"})
		if err != nil {
			return fmt.Errorf("resize server failed: %v", err)
		}
	}

	if d.HasChange("vm_state") {
		oldState, newState := d.GetChange("vm_state")
		if oldState.(string) == "error" {
			return fmt.Errorf("you cannot change server state because old server state is %s", oldState.(string))
		}

		switch newState.(string) {
		case "active":
			_, err := client.Server.Start(d.Id())
			if err != nil {
				return fmt.Errorf("error when start server: %v", err)
			}
			_, err = waitUntilServerStatusChangedState(d, meta, []string{"active"}, []string{"error"})
			if err != nil {
				return fmt.Errorf("start server failed: %v", err)
			}

		case "stopped":
			_, err := client.Server.Stop(d.Id())
			if err != nil {
				return fmt.Errorf("error when stop server: %v", err)
			}
			_, err = waitUntilServerStatusChangedState(d, meta, []string{"stopped"}, []string{"error"})
			if err != nil {
				return fmt.Errorf("stop server failed: %v", err)
			}

		default:
			return fmt.Errorf("new state of server must be 'active' or 'stopped'")
		}
	}

	return readOrImport(d, meta, false)
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.Server.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete server: %v", err)
	}
	_, err = waitUntilServerDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete server: %v", err)
	}
	return nil
}

func resourceServerImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := readOrImport(d, meta, true)
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

// func convertNics(nics []gocmcapiv2.Nic) []map[string]interface{} {
// 	result := make([]map[string]interface{}, len(nics))
// 	for i, nic := range nics {
// 		result[i] = map[string]interface{}{
// 			// "id":                 nic.Id,
// 			"subnet_id":  nic.FixedIps[0].SubnetID,
// 			"ip_address": nic.FixedIps[0].IPAddress,
// 			// "security_group_ids": nic.SecurityGroups,
// 			// "mac_address": nic.MacAddress,
// 		}
// 	}
// 	return result
// }

func waitUntilServerDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Server.Get(id, false)
	})
}

func waitUntilServerStatusChangedState(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, targetStatus, errorStatus, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).Server.Get(id, false)
	}, func(obj interface{}) string {
		return obj.(gocmcapiv2.Server).VMState
	})
}
