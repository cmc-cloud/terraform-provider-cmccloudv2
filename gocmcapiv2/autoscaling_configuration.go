package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingConfigurationService interface
type AutoScalingConfigurationService interface {
	Get(id string) (AutoScalingConfiguration, error)
	List(params map[string]string) ([]AutoScalingConfiguration, error)
	Create(params map[string]interface{}) (AutoScalingConfiguration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

type AutoScalingConfigurationNetwork struct {
	FloatingNetwork string   `json:"floating_network"`
	FloatingQos     string   `json:"floating_qos"`
	Network         string   `json:"network"`
	SecurityGroups  []string `json:"security_groups"`

	Subnet struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"subnet"`
}
type AutoScalingConfigurationVolume struct {
	BootIndex           int    `json:"boot_index"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	DestinationType     string `json:"destination_type"`
	SourceType          string `json:"source_type"`
	UUID                string `json:"uuid"`
	Size                int    `json:"volume_size"`
	Type                string `json:"volume_type"`
}
type AutoScalingConfiguration struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Spec      struct {
		Properties struct {
			AdminPass string                           `json:"admin_pass"`
			Volumes   []AutoScalingConfigurationVolume `json:"block_device_mapping_v2"`
			Context   struct {
				RegionName string `json:"region_name"`
			} `json:"context"`
			FlavorID              string                            `json:"flavor"`
			Networks              []AutoScalingConfigurationNetwork `json:"networks"`
			DomesticBandwidthMbps int                               `json:"domestic_bandwidth_mbps"`
			InterBandwidthMbps    int                               `json:"inter_bandwidth_mbps"`
			UseEip                bool                              `json:"use_eip"`
			UserData              string                            `json:"user_data"`
			SchedulerHints        struct {
				Group string `json:"group"`
			} `json:"scheduler_hints"`
		} `json:"properties"`
	} `json:"spec"`
}

type asconfiguration struct {
	client *Client
}

// Get asconfiguration detail
func (s *asconfiguration) Get(id string) (AutoScalingConfiguration, error) {
	jsonStr, err := s.client.Get("as/configuration/"+id, map[string]string{"detail": "true"})
	var obj AutoScalingConfiguration
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
		if err != nil {
			return obj, err
		}
	}
	return obj, err
}

func (s *asconfiguration) List(params map[string]string) ([]AutoScalingConfiguration, error) {
	restext, err := s.client.Get("as/configuration", params)
	items := make([]AutoScalingConfiguration, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a asconfiguration
func (s *asconfiguration) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("as/configuration/" + id)
}

func (s *asconfiguration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/configuration/"+id, params)
}

// Create a new asconfiguration
func (s *asconfiguration) Create(params map[string]interface{}) (AutoScalingConfiguration, error) {
	jsonStr, err := s.client.Post("as/configuration", params)
	var response AutoScalingConfiguration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
