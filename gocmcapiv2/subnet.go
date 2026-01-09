package gocmcapiv2

import (
	"encoding/json"
)

// SubnetService interface
type SubnetService interface {
	Get(id string) (Subnet, error)
	List(params map[string]string) ([]Subnet, error)
	Create(params map[string]interface{}) (Subnet, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}
type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
type HostRoute struct {
	Destination string `json:"destination"`
	NextHop     string `json:"nexthop"`
}

// Subnet object
type Subnet struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	IpVersion       int              `json:"ip_version"`
	EnableDhcp      bool             `json:"enable_dhcp"`
	GatewayIP       string           `json:"gateway_ip"`
	Cidr            string           `json:"cidr"`
	AllocationPools []AllocationPool `json:"allocation_pools"`
	HostRoutes      []HostRoute      `json:"host_routes"`
	DNSNameservers  []string         `json:"dns_nameservers"`
	Tags            []Tag            `json:"tags"`
	CreatedAt       string           `json:"created_at"`
	VpcID           string           `json:"vpc_id"`
	NetworkID       string           `json:"network_id"`
}
type subnet struct {
	client *Client
}

// Get subnet detail
func (v *subnet) Get(id string) (Subnet, error) {
	jsonStr, err := v.client.Get("network/subnet/"+id, map[string]string{})
	var obj Subnet
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (s *subnet) List(params map[string]string) ([]Subnet, error) {
	restext, err := s.client.Get("network/subnet", params)
	items := make([]Subnet, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a subnet
func (v *subnet) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/subnet/" + id)
}
func (v *subnet) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/subnet/"+id, params)
}

func (s *subnet) Create(params map[string]interface{}) (Subnet, error) {
	jsonStr, err := s.client.Post("network/subnet", params)
	var response Subnet
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
