package gocmcapiv2

import (
	"encoding/json"
)

// WafCertService interface
type WafCertService interface {
	Get(id string) (WafCert, error)
	List(waf_id string, params map[string]string) ([]WafCert, error)
	Create(params map[string]interface{}) (WafCert, error)
	Delete(id string) (ActionResponse, error)
}

type WafCert struct {
	ID          string `json:"id"`
	Created     int    `json:"created"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	CertName    string `json:"cert_name"`
	CertData    string `json:"cert_data"`
	KeyName     string `json:"key_name"`
	KeyData     string `json:"key_data"`
}
type WafCertListWrapper struct {
	Items []WafCert `json:"items"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Total int       `json:"total"`
}

type wafcert struct {
	client *Client
}

// Get WafCert detail
func (v *wafcert) Get(id string) (WafCert, error) {
	jsonStr, err := v.client.Get("waf/certs/"+id, map[string]string{})
	var wafcert WafCert
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafcert)
	}
	return wafcert, err
}
func (v *wafcert) List(waf_id string, params map[string]string) ([]WafCert, error) {
	jsonStr, err := v.client.Get("waf/certs/"+waf_id, map[string]string{})
	var wafcert WafCertListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &wafcert)
	}
	return wafcert.Items, err
}
func (v *wafcert) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("waf/certs/" + id)
}
func (v *wafcert) Create(params map[string]interface{}) (WafCert, error) {
	jsonStr, err := v.client.Post("waf/certs/", params)
	var response WafCert
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
