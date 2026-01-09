package gocmcapiv2

import (
	"encoding/json"
)

type NetworkInterfaceService interface {
	Get(id string) (NetworkInterface, error)
	List(params map[string]string) ([]NetworkInterface, error)
	Create(server_id string, params map[string]interface{}) (NetworkInterface, error)
	Delete(id string, server_id string) (ActionResponse, error)
}

// NetworkInterface object
type NetworkInterface struct {
	ID        string `json:"port_id"`
	NetID     string `json:"net_id"`
	MacAddr   string `json:"mac_addr"`
	PortState string `json:"port_state"`
	FixedIps  []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
	Tag any `json:"tag"`
}
type networkinterface struct {
	client *Client
}

// Get networkinterface detail
func (v *networkinterface) Get(id string) (NetworkInterface, error) {
	jsonStr, err := v.client.Get("network/port/"+id, map[string]string{})
	var networkinterface NetworkInterface
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &networkinterface)
	}
	return networkinterface, err
}

func (s *networkinterface) List(params map[string]string) ([]NetworkInterface, error) {
	restext, err := s.client.Get("network/port", params)
	items := make([]NetworkInterface, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

func (s *networkinterface) Create(server_id string, params map[string]interface{}) (NetworkInterface, error) {
	jsonStr, err := s.client.Post("server/"+server_id+"/interface", params)
	var response NetworkInterface
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

// Delete a networkinterface
func (v *networkinterface) Delete(id string, server_id string) (ActionResponse, error) {
	return v.client.PerformDelete("server/" + server_id + "/interface?port_id=" + id)
}
