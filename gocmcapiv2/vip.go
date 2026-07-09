package gocmcapiv2

import (
	"encoding/json"
)

// VIPService interface
type VIPService interface {
	Get(id string) (VIP, error)
	List(params map[string]string) ([]VIP, error)
	Create(params map[string]interface{}) (VIP, error)
	Delete(id string) (ActionResponse, error)
	BindPort(id string, port_id string) (ActionResponse, error)
	UnbindPort(id string, port_id string) (ActionResponse, error)
}

// VIP object
type VIP struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MacAddress   string `json:"mac_address"`
	AdminStateUp bool   `json:"admin_state_up"`
	Status       string `json:"status"`
	FixedIps     []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
	CreatedAt   string `json:"created_at"`
	DisplayText string `json:"display_text"`
}

type vip struct {
	client *Client
}

// Get VIP detail
func (v *vip) Get(id string) (VIP, error) {
	jsonStr, err := v.client.Get("network/virtual-ip/"+id, map[string]string{})
	var vip VIP
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vip)
	}
	return vip, err
}
func (s *vip) List(params map[string]string) ([]VIP, error) {
	restext, err := s.client.Get("network/virtual-ip", params)
	items := make([]VIP, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a VIP
func (v *vip) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/virtual-ip/" + id)
}
func (v *vip) BindPort(id string, port_id string) (ActionResponse, error) {
	return v.client.PerformAction("network/virtual-ip/"+id+"/bind_port", map[string]interface{}{"port_id": port_id})
}
func (v *vip) UnbindPort(id string, port_id string) (ActionResponse, error) {
	return v.client.PerformAction("network/virtual-ip/"+id+"/unbind_port", map[string]interface{}{"port_id": port_id})
}

func (v *vip) Create(params map[string]interface{}) (VIP, error) {
	jsonStr, err := v.client.Post("network/virtual-ip", params)
	var response VIP
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
