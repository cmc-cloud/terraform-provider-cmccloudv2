package gocmcapiv2

import (
	"encoding/json"
)

// WafIPService interface
type WafIPService interface {
	Get(id string) (WafIP, error)
	List(waf_id string, params map[string]string) ([]WafIP, error)
	Create(params map[string]interface{}) (WafIP, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type WafIP struct {
	ID          string `json:"id"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Expired     int    `json:"expired"`
	Created     int    `json:"created"`
	Description string `json:"description"`
	WafID       string `json:"waf_id"`
}
type WafIPListWrapper struct {
	Items []WafIP `json:"items"`
	Page  int     `json:"page"`
	Size  int     `json:"size"`
	Total int     `json:"total"`
}

type wafip struct {
	client *Client
}

// Get WafIP detail
func (v *wafip) Get(id string) (WafIP, error) {
	jsonStr, err := v.client.Get("waf/iplist/"+id, map[string]string{})
	var wafip WafIP
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafip)
	}
	return wafip, err
}
func (v *wafip) List(waf_id string, params map[string]string) ([]WafIP, error) {
	jsonStr, err := v.client.Get("waf/list/iplist/"+waf_id, map[string]string{})
	var wafip WafIPListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafip)
	}
	if err != nil {
		return []WafIP{}, err
	}
	return wafip.Items, err
}
func (v *wafip) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("waf/iplist/"+id, params)
}
func (v *wafip) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("waf/iplist/" + id)
}
func (v *wafip) Create(params map[string]interface{}) (WafIP, error) {
	jsonStr, err := v.client.Post("waf/iplist/", params)
	var response WafIP
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
