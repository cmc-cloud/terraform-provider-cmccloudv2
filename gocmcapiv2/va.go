package gocmcapiv2

import (
	"encoding/json"
)

// VAService interface
type VAService interface {
	Get(id string) (VA, error)
	List(params map[string]string) ([]VA, error)
	Create(params map[string]interface{}) (VA, error)
	Delete(id string) (ActionResponse, error)
}

// VA object
type VA struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Target          []string `json:"target"`
	Engine          string   `json:"engine"`
	Status          string   `json:"status"`
	Schedule        string   `json:"schedule"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Description     string   `json:"description"`
	EngineScanID    string   `json:"engine_scan_id"`
	TargetID        string   `json:"target_id"`
	ReportID        string   `json:"report_id"`
	UserID          string   `json:"user_id"`
	ResourceGroupID string   `json:"resource_group_id"`
}

type VAListWrapper struct {
	Items []VA `json:"items"`
	Page  int  `json:"page"`
	Size  int  `json:"size"`
	Total int  `json:"total"`
}
type va struct {
	client *Client
}

// Get VA detail
func (v *va) Get(id string) (VA, error) {
	jsonStr, err := v.client.Get("vulner/"+id, map[string]string{})
	var va VA
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &va)
	}
	return va, err
}
func (v *va) List(params map[string]string) ([]VA, error) {
	jsonStr, err := v.client.Get("vulner", map[string]string{})
	var va VAListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &va)
	}
	if err != nil {
		return []VA{}, err
	}
	return va.Items, err
}
func (v *va) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("vulner/" + id)
}
func (v *va) Create(params map[string]interface{}) (VA, error) {
	jsonStr, err := v.client.Post("vulner/scan", params)
	var response VA
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
