package gocmcapiv2

import (
	"encoding/json"
)

// WafWhitelistService interface
type WafWhitelistService interface {
	Get(id string) (WafWhitelist, error)
	List(waf_id string, params map[string]string) ([]WafWhitelist, error)
	Create(params map[string]interface{}) (WafWhitelist, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type WafWhitelist struct {
	ID             string `json:"id"`
	WhitelistType  string `json:"whitelist_type"`
	Wl             string `json:"wl"`
	Mz             string `json:"mz"`
	WafID          string `json:"waf_id"`
	Description    string `json:"description"`
	WhitelistSetID string `json:"whitelist_set_id"`
	Rmks           string `json:"rmks"`
	Active         bool   `json:"active"`
	Negative       bool   `json:"negative"`
	Timestamp      int    `json:"timestamp"`
}
type WafWhitelistListWrapper struct {
	Items []WafWhitelist `json:"items"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
}

type wafwhitelist struct {
	client *Client
}

// Get WafWhitelist detail
func (v *wafwhitelist) Get(id string) (WafWhitelist, error) {
	jsonStr, err := v.client.Get("waf/whitelist/"+id, map[string]string{})
	var wafwhitelist WafWhitelist
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafwhitelist)
	}
	return wafwhitelist, err
}
func (v *wafwhitelist) List(waf_id string, params map[string]string) ([]WafWhitelist, error) {
	jsonStr, err := v.client.Get("waf/list/whitelist/"+waf_id, map[string]string{})
	var wafwhitelist WafWhitelistListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafwhitelist)
	}
	if err != nil {
		return []WafWhitelist{}, err
	}
	return wafwhitelist.Items, err
}
func (v *wafwhitelist) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("waf/whitelist/"+id, params)
}
func (v *wafwhitelist) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("waf/whitelist/" + id)
}
func (v *wafwhitelist) Create(params map[string]interface{}) (WafWhitelist, error) {
	jsonStr, err := v.client.Post("waf/whitelist/", params)
	var response WafWhitelist
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
