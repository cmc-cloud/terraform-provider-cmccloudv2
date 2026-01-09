package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingV2GroupService interface
type AutoScalingV2GroupService interface {
	Get(id string) (AutoScalingV2Group, error)
	List(params map[string]string) ([]AutoScalingV2Group, error)
	Create(params map[string]interface{}) (AutoScalingV2Group, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	// ChangeConfiguration(id string, configuration_id string) (ActionResponse, error)
	// UpdateCapacity(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

// AutoScalingV2Group object
type AutoScalingV2Group struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	CreationTime string `json:"creation_time"`
	Parameters   struct {
		AvailabilityZone        string `json:"availability_zone"`
		DesiredCapacity         int    `json:"desired_capacity"`
		MinSize                 int    `json:"min_size"`
		MaxSize                 int    `json:"max_size"`
		LbPool                  string `json:"lb_pool"`
		LbMemberPort            int    `json:"lb_member_port"`
		Cooldown                int    `json:"cooldown"`
		ScaleUpAdjustmentType   string `json:"scale_up_adjustment_type"`
		ScaleUpCooldown         int    `json:"scale_up_cooldown"`
		ScaleUpAdjustment       int    `json:"scale_up_adjustment"`
		ScaleDownAdjustmentType string `json:"scale_down_adjustment_type"`
		ScaleDownCooldown       int    `json:"scale_down_cooldown"`
		ScaleDownAdjustment     int    `json:"scale_down_adjustment"`
		// FlavorID                   string `json:"flavor"`
		// ImageID                    string `json:"image"`
		// KeyName                    string `json:"key_name"`
		// UserData                   string `json:"user_data"`
		// QosPolicyID      string `json:"qos_policy_id"`
		// EnableFloatingIP bool   `json:"enable_floating_ip"`
		// VolumeSize1          int    `json:"volume_size_1"`
		// VolumeType1          string `json:"volume_type_1"`
		// DeleteOnTermination1 bool   `json:"delete_on_termination_1"`
		// NetworkID1           string `json:"network_id_1"`
		// SubnetID1            string `json:"subnet_id_1"`
		// SecurityGroups1      []any  `json:"security_groups_1"`
		// RollingUpdatesMaxBatchSize int    `json:"rolling_updates_max_batch_size"`
		// RollingUpdatesMinInService int    `json:"rolling_updates_min_in_service"`
		// RollingUpdatesPauseTime    int    `json:"rolling_updates_pause_time"`
		// GroupPolicies              string `json:"group_policies"`
		// FloatingNetwork            string `json:"floating_network"`
	} `json:"parameters"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
}
type autoscalingv2group struct {
	client *Client
}

// Get autoscalingv2group detail
func (s *autoscalingv2group) Get(id string) (AutoScalingV2Group, error) {
	jsonStr, err := s.client.Get("asv2/group/"+id, map[string]string{})
	var obj AutoScalingV2Group
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *autoscalingv2group) List(params map[string]string) ([]AutoScalingV2Group, error) {
	restext, err := s.client.Get("asv2/group", params)
	items := make([]AutoScalingV2Group, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a autoscalingv2group
func (s *autoscalingv2group) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("asv2/group/" + id)
}
func (s *autoscalingv2group) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("asv2/group/"+id, params)
}

// func (s *autoscalingv2group) ChangeConfiguration(id string, configuration_id string) (ActionResponse, error) {
// 	return s.client.PerformUpdate("asv2/group/"+id+"/change_configuration", map[string]interface{}{
// 		"configuration_id": configuration_id,
// 	})
// }

// func (s *autoscalingv2group) UpdateCapacity(id string, params map[string]interface{}) (ActionResponse, error) {
// 	return s.client.PerformUpdate("asv2/group/"+id+"/capacity", params)
// }

// Create a new autoscalingv2group
func (s *autoscalingv2group) Create(params map[string]interface{}) (AutoScalingV2Group, error) {
	jsonStr, err := s.client.Post("asv2/group", params)
	var response AutoScalingV2Group
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
