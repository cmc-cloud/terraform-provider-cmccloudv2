package gocmcapiv2

import (
	"encoding/json"
)

// WafService interface
type WafService interface {
	Get(id string) (Waf, error)
	List(params map[string]string) ([]Waf, error)
	Create(params map[string]interface{}) (Waf, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateLoadBalance(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type WafListWrapper struct {
	Items []Waf `json:"items"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int   `json:"total"`
}

// Waf object
type Waf struct {
	ID                string `json:"id"`
	Domain            string `json:"domain"`
	Port              string `json:"port"`
	Realserver        string `json:"realserver"`
	Description       string `json:"description"`
	Created           int    `json:"created"`
	Updated           int    `json:"updated"`
	Ssl               bool   `json:"ssl"`
	Certificate       string `json:"certificate"`
	Certificatekey    string `json:"certificatekey"`
	UserID            string `json:"user_id"`
	LbEnable          bool   `json:"lb_enable"`
	LbMethod          string `json:"lb_method"`
	LbKeepalive       int    `json:"lb_keepalive"`
	CertificateID     string `json:"certificate_id"`
	Sendfile          bool   `json:"sendfile"`
	ClientMaxBodySize int    `json:"client_max_body_size"`
	Validated         bool   `json:"validated"`
	Mode              string `json:"mode"`
	Status            string `json:"status"`
	Protocol          string `json:"protocol"`
	DomainDNS         string `json:"domain_dns"`
	PolicyID          string `json:"policy_id"`
	ProjectID         string `json:"project_id"`
}

type waf struct {
	client *Client
}

// Get Waf detail
func (v *waf) Get(id string) (Waf, error) {
	jsonStr, err := v.client.Get("waf/apps/"+id, map[string]string{})
	var waf Waf
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &waf)
	}
	return waf, err
}
func (v *waf) List(params map[string]string) ([]Waf, error) {
	jsonStr, err := v.client.Get("waf/apps", map[string]string{})
	var waf WafListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &waf)
	}
	if err != nil {
		return []Waf{}, err
	}
	return waf.Items, err
}
func (v *waf) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("waf/apps/" + id)
}
func (v *waf) Create(params map[string]interface{}) (Waf, error) {
	jsonStr, err := v.client.Post("waf/apps/", params)
	var response Waf
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *waf) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("waf/apps/"+id, params)
}
func (v *waf) UpdateLoadBalance(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("waf/loadbalancer/config/"+id, params)
}
