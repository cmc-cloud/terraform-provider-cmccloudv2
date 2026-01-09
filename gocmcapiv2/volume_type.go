package gocmcapiv2

import (
	"encoding/json"
)

// VolumeTypeService interface
type VolumeTypeService interface {
	Get(id string) (VolumeType, error)
	List(params map[string]string) ([]VolumeType, error)
}

// VolumeType object
type VolumeType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsPublic    bool   `json:"is_public"`
	Description string `json:"description"`
}
type volumetype struct {
	client *Client
}

// Get volumetype detail
func (v *volumetype) Get(id string) (VolumeType, error) {
	jsonStr, err := v.client.Get("volume/type/"+id, map[string]string{})
	var volumetype VolumeType
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &volumetype)
	}
	return volumetype, err
}

func (s *volumetype) List(params map[string]string) ([]VolumeType, error) {
	restext, err := s.client.Get("volume/type", params)
	items := make([]VolumeType, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
