package gocmcapiv2

import (
	"encoding/json"
)

// PortService interface
type PortService interface {
	Get(id string) (Port, error)
	List(params map[string]string) ([]Port, error)
	Create(params map[string]interface{}) (Port, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Patch(id string, params map[string]interface{}) (ActionResponse, error)
}

// Port object
type Port struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NetworkID   string `json:"network_id"`
	TenantID    string `json:"tenant_id"`
	MacAddress  string `json:"mac_address"`
	Status      string `json:"status"`
	DeviceID    string `json:"device_id"`
	DeviceOwner string `json:"device_owner"`
	FixedIps    []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
	AllowedAddressPairs []struct {
		MacAddress string `json:"mac_address"`
		IPAddress  string `json:"ip_address"`
	} `json:"allowed_address_pairs"`
	SecurityGroups      []string `json:"security_groups"`
	PortSecurityEnabled bool     `json:"port_security_enabled"`
	IPAddress           string   `json:"ip_address"`
	// ExtraDhcpOpts       []any  `json:"extra_dhcp_opts"`
	// Description         string `json:"description"`
	// QosPolicyID         any    `json:"qos_policy_id"`
	// QosNetworkPolicyID  any    `json:"qos_network_policy_id"`
	// DNSName             string `json:"dns_name"`
	// DNSAssignment       []struct {
	// 	Hostname  string `json:"hostname"`
	// 	Fqdn      string `json:"fqdn"`
	// } `json:"dns_assignment"`
	// DNSDomain       string    `json:"dns_domain"`
	// IPAllocation    string    `json:"ip_allocation"`
	// Tags            []any     `json:"tags"`
	// CreatedAt       string `json:"created_at"`
	// UpdatedAt       string `json:"updated_at"`
	// ProjectID       string    `json:"project_id"`
}
type port struct {
	client *Client
}

// Get port detail
func (v *port) Get(id string) (Port, error) {
	jsonStr, err := v.client.Get("network/port/"+id, map[string]string{})
	var obj Port
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *port) List(params map[string]string) ([]Port, error) {
	restext, err := v.client.Get("network/port", params)
	items := make([]Port, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a port
func (v *port) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/port/" + id)
}
func (v *port) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/port/"+id, params)
}

// Patch a port
func (v *port) Patch(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformPatch("network/port/"+id, params)
}
func (v *port) Create(params map[string]interface{}) (Port, error) {
	jsonStr, err := v.client.Post("network/port", params)
	var response Port
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
