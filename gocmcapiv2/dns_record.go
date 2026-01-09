package gocmcapiv2

import (
	"encoding/json"
)

// DnsRecordService interface
type DnsRecordService interface {
	Get(zone_id string, id string) (DnsRecord, error)
	List(zone_id string, params map[string]string) ([]DnsRecord, error)
	Create(zone_id string, params map[string]interface{}) (DnsRecord, error)
	Update(zone_id string, id string, params map[string]interface{}) (ActionResponse, error)
	Delete(zone_id string, id string) (ActionResponse, error)
}

type DnsRecordListWrapper struct {
	PageInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"page_info"`
	Success  bool          `json:"success"`
	Error    interface{}   `json:"error"`
	Messages []interface{} `json:"messages"`
	Result   []DnsRecord   `json:"result"`
}

// DnsRecord object
type DnsRecord struct {
	ID                string        `json:"id"`
	ZoneID            string        `json:"zone_id"`
	Domain            string        `json:"domain"`
	Type              string        `json:"type"`
	TTL               int           `json:"ttl"`
	Detail            []DnsRecordIP `json:"detail"`
	LoadbalanceType   string        `json:"loadbalance_type"`
	LoadbalanceStatus string        `json:"loadbalance_status"`
	CreatedAt         int           `json:"created_at"`
	UpdatedAt         int           `json:"updated_at"`
}
type DnsRecordIP struct {
	Content string `json:"content"`
	Weight  int    `json:"weight"`
}
type DnsRecordWrapper struct {
	Success   bool          `json:"success"`
	Error     interface{}   `json:"error"`
	Messages  []interface{} `json:"messages"`
	RequestID string        `json:"request_id"`
	Record    DnsRecord     `json:"result"`
}

type dnsrecord struct {
	client *Client
}

// Get DnsRecord detail
func (v *dnsrecord) Get(zone_id string, id string) (DnsRecord, error) {
	jsonStr, err := v.client.Get("dns/dns/zones/"+zone_id+"/rrsets/"+id, map[string]string{})
	var dnsrecord DnsRecordWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dnsrecord)
	}

	if err != nil {
		return DnsRecord{}, err
	}
	return dnsrecord.Record, err
}
func (v *dnsrecord) List(zone_id string, params map[string]string) ([]DnsRecord, error) {
	jsonStr, err := v.client.Get("dns/dns/zones/"+zone_id+"/rrsets", map[string]string{})
	var dnsrecord DnsRecordListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &dnsrecord)
	}
	if err != nil {
		return []DnsRecord{}, err
	}
	return dnsrecord.Result, err
}
func (v *dnsrecord) Delete(zone_id string, id string) (ActionResponse, error) {
	return v.client.PerformDelete("dns/dns/zones/" + zone_id + "/rrsets/" + id)
}
func (v *dnsrecord) Create(zone_id string, params map[string]interface{}) (DnsRecord, error) {
	jsonStr, err := v.client.Post("dns/dns/zones/"+zone_id+"/rrsets", params)
	var response DnsRecordWrapper
	if err != nil {
		var nilres DnsRecord
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return DnsRecord{}, err
	}
	return response.Record, err
}

func (v *dnsrecord) Update(zone_id string, id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dns/dns/zones/"+zone_id+"/rrsets/"+id, params)
}
