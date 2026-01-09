package gocmcapiv2

import (
	"encoding/json"
)

// WafRuleService interface
type WafRuleService interface {
	Get(id string) (WafRule, error)
	List(waf_id string, params map[string]string) ([]WafRule, error)
	Create(params map[string]interface{}) (WafRule, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type WafRule struct {
	ID          string `json:"id"`
	RuleType    string `json:"rule_type"`
	Msg         string `json:"msg"`
	Detection   string `json:"detection"`
	Mz          string `json:"mz"`
	Score       string `json:"score"`
	RulesetID   string `json:"ruleset_id"`
	Rmks        string `json:"rmks"`
	Active      bool   `json:"active"`
	Negative    bool   `json:"negative"`
	Timestamp   int    `json:"timestamp"`
	WafID       string `json:"waf_id"`
	Sid         int    `json:"sid"`
	Description string `json:"description"`
}
type WafRuleListWrapper struct {
	Items []WafRule `json:"items"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
}

type wafrule struct {
	client *Client
}

// Get WafRule detail
func (v *wafrule) Get(id string) (WafRule, error) {
	jsonStr, err := v.client.Get("waf/rules/"+id, map[string]string{})
	var wafrule WafRule
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafrule)
	}
	return wafrule, err
}
func (v *wafrule) List(waf_id string, params map[string]string) ([]WafRule, error) {
	jsonStr, err := v.client.Get("waf/list/rules/"+waf_id, map[string]string{})
	var wafrule WafRuleListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafrule)
	}
	if err != nil {
		return []WafRule{}, err
	}
	return wafrule.Items, err
}
func (v *wafrule) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("waf/rules/"+id, params)
}
func (v *wafrule) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("waf/rules/" + id)
}
func (v *wafrule) Create(params map[string]interface{}) (WafRule, error) {
	jsonStr, err := v.client.Post("waf/rules/", params)
	var response WafRule
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
