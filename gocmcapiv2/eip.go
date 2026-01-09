package gocmcapiv2

import (
	"encoding/json"
)

// EIPService interface
type EIPService interface {
	Get(id string) (EIP, error)
	List(params map[string]string) ([]EIP, error)
	Create(params map[string]interface{}) (EIP, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	AttachPort(id string, port_id string, fix_ip_address string) (ActionResponse, error)
	DetachPort(id string) (ActionResponse, error)
	GetPortForwardingRule(id string, rule_id string) (PortForwardingRule, error)
	CreatePortForwardingRule(id string, params map[string]interface{}) (PortForwardingRule, error)
	UpdatePortForwardingRule(id string, rule_id string, params map[string]interface{}) (ActionResponse, error)
	DeletePortForwardingRule(id string, rule_id string) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type PortForwardingRule struct {
	ID                string `json:"id"`
	Protocol          string `json:"protocol"`
	InternalPortID    string `json:"internal_port_id"`
	InternalIPAddress string `json:"internal_ip_address"`
	InternalPort      int    `json:"internal_port"`
	ExternalPort      int    `json:"external_port"`
	InternalPortRange string `json:"internal_port_range"`
	ExternalPortRange string `json:"external_port_range"`
	Description       string `json:"description"`
}

// EIP object
type EIP struct {
	ID                string `json:"id"`
	FloatingIPAddress string `json:"floating_ip_address"`
	FloatingNetworkID string `json:"floating_network_id"`
	PortID            string `json:"port_id"`
	FixedIPAddress    string `json:"fixed_ip_address"`
	Status            string `json:"status"`
	Description       string `json:"description"`
	QosPolicyID       string `json:"qos_policy_id"`
	PortDetails       struct {
		Name         string `json:"name"`
		NetworkID    string `json:"network_id"`
		MacAddress   string `json:"mac_address"`
		AdminStateUp bool   `json:"admin_state_up"`
		Status       string `json:"status"`
		DeviceID     string `json:"device_id"`
		DeviceOwner  string `json:"device_owner"`
	} `json:"port_details"`
	DNSDomain             string `json:"dns_domain"`
	DNSName               string `json:"dns_name"`
	Tags                  []Tag  `json:"tags"`
	CreatedAt             string `json:"created_at"`
	DomesticBandwidthMbps int    `json:"domestic_bandwidth_mbps"`
	InterBandwidthMbps    int    `json:"inter_bandwidth_mbps"`
	BillingMode           string `json:"billing_mode"`
	// PortForwardings       []PortForwarding `json:"port_forwardings"`
}

type eip struct {
	client *Client
}

// Get EIP detail
func (v *eip) Get(id string) (EIP, error) {
	jsonStr, err := v.client.Get("network/eip/"+id, map[string]string{})
	var eip EIP
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &eip)
	}
	return eip, err
}
func (s *eip) List(params map[string]string) ([]EIP, error) {
	restext, err := s.client.Get("network/eip", params)
	items := make([]EIP, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a EIP
func (v *eip) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/eip/" + id)
}
func (v *eip) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/eip/"+id, params)
}
func (v *eip) Create(params map[string]interface{}) (EIP, error) {
	jsonStr, err := v.client.Post("network/eip", params)
	var response EIP
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *eip) GetPortForwardingRule(id string, rule_id string) (PortForwardingRule, error) {
	jsonStr, err := v.client.Get("network/eip/"+id+"/port-forwarding-rule/"+rule_id, map[string]string{})
	var eip PortForwardingRule
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &eip)
	}
	return eip, err
}

func (v *eip) CreatePortForwardingRule(id string, params map[string]interface{}) (PortForwardingRule, error) {
	jsonStr, err := v.client.Post("network/eip/"+id+"/port-forwarding-rule", params)
	var response PortForwardingRule
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (v *eip) UpdatePortForwardingRule(id string, rule_id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/eip/"+id+"/port-forwarding-rule/"+rule_id, params)
}
func (v *eip) DeletePortForwardingRule(id string, rule_id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/eip/" + id + "/port-forwarding-rule/" + rule_id)
}

func (v *eip) AttachPort(id string, port_id string, fix_ip_address string) (ActionResponse, error) {
	return v.client.PerformAction("network/eip/"+id+"/associate", map[string]interface{}{"port_id": port_id})
}

func (v *eip) DetachPort(id string) (ActionResponse, error) {
	return v.client.PerformAction("network/eip/"+id+"/disassociate", map[string]interface{}{})
}
