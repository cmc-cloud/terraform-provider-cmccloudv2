package gocmcapiv2

import (
	"encoding/json"
)

// EFSService interface
type EFSService interface {
	Get(id string) (EFS, error)
	List(params map[string]string) ([]EFS, error)
	Create(params map[string]interface{}) (EFS, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Resize(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

// EFS object
type EFS struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Used         int    `json:"used"`
	Capacity     int    `json:"capacity"`
	ProtocolType string `json:"protocol_type"`
	VpcID        string `json:"vpc_id"`
	SubnetID     string `json:"subnet_id"`
	CreatedAt    string `json:"created_at"`
	Endpoint     string `json:"endpoint"`
	Tags         []Tag  `json:"tags"`
	Status       string `json:"status"`
	CommandLine  string `json:"command_line"`
	SharedPath   string `json:"shared_path"`
	BillingMode  string `json:"billing_mode"`
}

type efs struct {
	client *Client
}

// Get EFS detail
func (v *efs) Get(id string) (EFS, error) {
	jsonStr, err := v.client.Get("efs/"+id, map[string]string{})
	var efs EFS
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &efs)
	}
	return efs, err
}
func (s *efs) List(params map[string]string) ([]EFS, error) {
	restext, err := s.client.Get("efs", params)
	items := make([]EFS, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *efs) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("efs/" + id)
}
func (v *efs) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("efs/"+id, params)
}
func (v *efs) Resize(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformAction("efs/"+id+"/resize", params)
}
func (v *efs) Create(params map[string]interface{}) (EFS, error) {
	jsonStr, err := v.client.Post("efs", params)
	var response EFS
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
