package gocmcapiv2

import (
	"encoding/json"
)

// DnsZoneService interface
type DnsZoneService interface {
	Get(id string) (DnsZone, error)
	List(params map[string]string) ([]DnsZone, error)
	Create(params map[string]interface{}) (DnsZone, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type DnsZoneListWrapper struct {
	PageInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"page_info"`
	Success  bool          `json:"success"`
	Error    interface{}   `json:"error"`
	Messages []interface{} `json:"messages"`
	Result   []DnsZone     `json:"result"`
}

// DnsZone object
type DnsZone struct {
	ID        string `json:"id"`
	Zone      string `json:"zone"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

type DnsZoneWrapper struct {
	Success   bool          `json:"success"`
	Error     interface{}   `json:"error"`
	Messages  []interface{} `json:"messages"`
	RequestID string        `json:"request_id"`
	Zone      DnsZone       `json:"result"`
}

type dns struct {
	client *Client
}

// Get DnsZone detail
func (v *dns) Get(id string) (DnsZone, error) {
	jsonStr, err := v.client.Get("dns/dns/zones/"+id, map[string]string{})
	var dns DnsZoneWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dns)
	}
	return dns.Zone, err
}
func (v *dns) List(params map[string]string) ([]DnsZone, error) {
	jsonStr, err := v.client.Get("dns/dns/zones", map[string]string{})
	var dns DnsZoneListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dns)
	}
	if err != nil {
		return []DnsZone{}, err
	}
	return dns.Result, err
}
func (v *dns) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("dns/dns/zones/" + id)
}
func (v *dns) Create(params map[string]interface{}) (DnsZone, error) {
	jsonStr, err := v.client.Post("dns/dns/zones", params)
	var response DnsZoneWrapper
	if err != nil {
		var resnil DnsZone
		return resnil, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response.Zone, err
}

func (v *dns) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dns/apps/"+id, params)
}
