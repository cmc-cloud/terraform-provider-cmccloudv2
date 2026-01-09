package gocmcapiv2

import (
	"encoding/json"
)

// SecurityGroupService interface
type SecurityGroupService interface {
	Get(id string) (SecurityGroup, error)
	List(params map[string]string) ([]SecurityGroup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	SaveRules(id string, rules string) (ActionResponse, error)
	GetRule(id string) (SecurityGroupRule, error)
	CreateRule(id string, params map[string]interface{}) (SecurityGroupRule, error)
	DeleteRule(id string) (ActionResponse, error)
	Create(params map[string]interface{}) (SecurityGroup, error)
}

// SecurityGroupRule object
type SecurityGroupRule struct {
	ID                   string `json:"id"`
	EtherType            string `json:"ethertype"`
	Direction            string `json:"direction"`
	Protocol             string `json:"protocol"`
	PortRangeMin         int    `json:"port_range_min"`
	PortRangeMax         int    `json:"port_range_max"`
	CIDR                 any    `json:"remote_ip_prefix"`
	RemoteAddressGroupID any    `json:"remote_address_group_id"`
	RemoteGroupID        any    `json:"remote_group_id"`
	Description          string `json:"description"`
	CreatedAt            string `json:"created_at"`
}

type SecurityGroup struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Stateful    bool                `json:"stateful"`
	Description string              `json:"description"`
	CreatedAt   string              `json:"created_at"`
	Rules       []SecurityGroupRule `json:"security_group_rules"`
}

type securitygroup struct {
	client *Client
}

// Get SecurityGroup detail
func (v *securitygroup) Get(id string) (SecurityGroup, error) {
	jsonStr, err := v.client.Get("network/securitygroup/"+id, map[string]string{})
	var firewall SecurityGroup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &firewall)
	}
	return firewall, err
}

func (s *securitygroup) List(params map[string]string) ([]SecurityGroup, error) {
	restext, err := s.client.Get("network/securitygroup", params)
	items := make([]SecurityGroup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a SecurityGroup
func (v *securitygroup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/securitygroup/" + id)
}
func (v *securitygroup) DeleteRule(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/securitygroup/rule/" + id)
}

func (v *securitygroup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/securitygroup/"+id, params)
}

func (v *securitygroup) SaveRules(id string, rules string) (ActionResponse, error) {
	return v.client.PerformUpdate("network/securitygroup/"+id, map[string]interface{}{"rules": rules})
}

func (v *securitygroup) GetRule(id string) (SecurityGroupRule, error) {
	jsonStr, err := v.client.Get("network/securitygroup/rule/"+id, map[string]string{})
	var response SecurityGroupRule
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (v *securitygroup) CreateRule(id string, params map[string]interface{}) (SecurityGroupRule, error) {
	jsonStr, err := v.client.Post("network/securitygroup/"+id+"/rule", params)
	var response SecurityGroupRule
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *securitygroup) Create(params map[string]interface{}) (SecurityGroup, error) {
	jsonStr, err := s.client.Post("network/securitygroup", params)
	var response SecurityGroup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
