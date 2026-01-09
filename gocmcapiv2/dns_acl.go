package gocmcapiv2

import (
	"encoding/json"
)

// DnsAclService interface
type DnsAclService interface {
	Get(zone_id string, id string) (DnsAcl, error)
	List(zone_id string, params map[string]string) ([]DnsAcl, error)
	Create(zone_id string, params map[string]interface{}) (DnsAcl, error)
	Update(zone_id string, id string, params map[string]interface{}) (ActionResponse, error)
	Delete(zone_id string, id string) (ActionResponse, error)
}

type DnsAclListWrapper struct {
	PageInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"page_info"`
	Success  bool          `json:"success"`
	Error    interface{}   `json:"error"`
	Messages []interface{} `json:"messages"`
	Result   []DnsAcl      `json:"result"`
}

// DnsAcl object
type DnsAcl struct {
	ID         string `json:"_id"`
	ZoneID     string `json:"zone_id"`
	Domain     string `json:"domain"`
	RecordType string `json:"record_type"`
	Action     string `json:"action"`
	SourceIP   string `json:"source_ip"`
	Status     string `json:"status"`
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at"`
}
type DnsAclWrapper struct {
	Success   bool          `json:"success"`
	Error     interface{}   `json:"error"`
	Messages  []interface{} `json:"messages"`
	RequestID string        `json:"request_id"`
	Acl       DnsAcl        `json:"result"`
}

type dnsacl struct {
	client *Client
}

// Get DnsAcl detail
func (v *dnsacl) Get(zone_id string, id string) (DnsAcl, error) {
	jsonStr, err := v.client.Get("dns/dns/zones/"+zone_id+"/acls/"+id, map[string]string{})
	var dnsacl DnsAclWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dnsacl)
	}

	if err != nil {
		return DnsAcl{}, err
	}
	return dnsacl.Acl, err
}
func (v *dnsacl) List(zone_id string, params map[string]string) ([]DnsAcl, error) {
	jsonStr, err := v.client.Get("dns/dns/zones/"+zone_id+"/acls", map[string]string{})
	var dnsacl DnsAclListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dnsacl)
	}

	if err != nil {
		return []DnsAcl{}, err
	}
	return dnsacl.Result, err
}
func (v *dnsacl) Delete(zone_id string, id string) (ActionResponse, error) {
	return v.client.PerformDelete("dns/dns/zones/" + zone_id + "/acls/" + id)
}
func (v *dnsacl) Create(zone_id string, params map[string]interface{}) (DnsAcl, error) {
	jsonStr, err := v.client.Post("dns/dns/zones/"+zone_id+"/acls", params)
	var response DnsAclWrapper
	if err != nil {
		return DnsAcl{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response.Acl, err
}

func (v *dnsacl) Update(zone_id string, id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformPatch("dns/dns/zones/"+zone_id+"/acls/"+id, params)
}
