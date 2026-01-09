package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingV2ConfigurationService interface
type AutoScalingV2ConfigurationService interface {
	Get(id string) (AutoScalingV2Configuration, error)
	List(params map[string]string) ([]AutoScalingV2Configuration, error)
	Create(params map[string]interface{}) (AutoScalingV2Configuration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

type AutoScalingV2ConfigurationVolume struct {
	Type                string `json:"type"`
	Size                int    `json:"size"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}
type AutoScalingV2Configuration struct {
	ID                 string                             `json:"id"`
	Name               string                             `json:"name"`
	SourceType         string                             `json:"source_type"`
	SourceID           string                             `json:"source_id"`
	FlavorID           string                             `json:"flavor_id"`
	VpcID              string                             `json:"vpc_id"`
	SubnetIds          []string                           `json:"subnet_ids"`
	Volumes            []AutoScalingV2ConfigurationVolume `json:"volumes"`
	SecurityGroupNames []string                           `json:"security_group_names"`
	UseEip             bool                               `json:"use_eip"`
	DomesticBandwidth  int                                `json:"domestic_bandwidth"`
	InterBandwidth     int                                `json:"inter_bandwidth"`
	KeyName            string                             `json:"key_name"`
	UserData           string                             `json:"user_data"`
	EcsGroupID         string                             `json:"ecs_group_id"`
	Created            string                             `json:"created"`
	PasswordEncrypted  string                             `json:"password_encrypted"`
}

type asv2configuration struct {
	client *Client
}

// Get asv2configuration detail
func (s *asv2configuration) Get(id string) (AutoScalingV2Configuration, error) {
	jsonStr, err := s.client.Get("asv2/configuration/"+id, map[string]string{"detail": "true"})
	var obj AutoScalingV2Configuration
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *asv2configuration) List(params map[string]string) ([]AutoScalingV2Configuration, error) {
	restext, err := s.client.Get("asv2/configuration", params)
	items := make([]AutoScalingV2Configuration, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a asv2configuration
func (s *asv2configuration) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("asv2/configuration/" + id)
}

func (s *asv2configuration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("asv2/configuration/"+id, params)
}

// Create a new asv2configuration
func (s *asv2configuration) Create(params map[string]interface{}) (AutoScalingV2Configuration, error) {
	jsonStr, err := s.client.Post("asv2/configuration", params)
	var response AutoScalingV2Configuration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
