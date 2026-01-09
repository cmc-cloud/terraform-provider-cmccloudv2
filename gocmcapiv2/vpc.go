package gocmcapiv2

import (
	"encoding/json"
)

// VPCService interface
type VPCService interface {
	Get(id string) (VPC, error)
	List(params map[string]string) ([]VPC, error)
	Create(params map[string]interface{}) (VPC, error)
	CreateSubnet(id string, params map[string]interface{}) (Subnet, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

// VPC object
type VPC struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cidr        string `json:"cidr"`
	CreatedAt   string `json:"created_at"`
	ProjectID   string `json:"project_id"`
	RouterID    string `json:"router_id"`
	Tags        []Tag  `json:"tags"`
	BillingMode string `json:"billing_mode"`
}
type vpc struct {
	client *Client
}

// Get vpc detail
func (v *vpc) Get(id string) (VPC, error) {
	jsonStr, err := v.client.Get("network/vpc/"+id, map[string]string{})
	var obj VPC
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (s *vpc) List(params map[string]string) ([]VPC, error) {
	restext, err := s.client.Get("network/vpc", params)
	items := make([]VPC, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a vpc
func (v *vpc) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("network/vpc/" + id)
}
func (v *vpc) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("network/vpc/"+id, params)
}

func (s *vpc) Create(params map[string]interface{}) (VPC, error) {
	jsonStr, err := s.client.Post("network/vpc", params)
	var response VPC
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (s *vpc) CreateSubnet(id string, params map[string]interface{}) (Subnet, error) {
	jsonStr, err := s.client.Post("network/vpc/"+id+"/subnet", params)
	var response Subnet
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
