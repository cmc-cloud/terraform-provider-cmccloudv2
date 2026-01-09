package gocmcapiv2

import (
	"encoding/json"
)

// CertificateService interface
type CertificateService interface {
	Get(id string) (Certificate, error)
	List(params map[string]string) ([]Certificate, error)
	Create(params map[string]interface{}) (Certificate, error)
	Delete(id string) (ActionResponse, error)
}

// Certificate object
type Certificate struct {
	ID         string `json:"id"`
	Created    string `json:"created"`
	Status     string `json:"status"`
	Name       string `json:"name"`
	SecretType string `json:"secret_type"`
	SecretRef  string `json:"secret_ref"`
}
type certificate struct {
	client *Client
}

func (v *certificate) Get(name string) (Certificate, error) {
	jsonStr, err := v.client.Get("certificate/"+name, map[string]string{})
	var vpc Certificate
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vpc)
	}
	return vpc, err
}

func (s *certificate) List(params map[string]string) ([]Certificate, error) {
	restext, err := s.client.Get("certificate", params)
	certificates := make([]Certificate, 0)
	if err != nil {
		return certificates, err
	}
	err = json.Unmarshal([]byte(restext), &certificates)
	return certificates, err
}

func (v *certificate) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("certificate/" + id)
}
func (v *certificate) Create(params map[string]interface{}) (Certificate, error) {
	jsonStr, err := v.client.Post("certificate", params)
	var response Certificate
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
