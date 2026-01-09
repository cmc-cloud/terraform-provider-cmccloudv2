package gocmcapiv2

import (
	"encoding/json"
)

type DevopsProjectService interface {
	Get(id string) (DevopsProject, error)
	List(params map[string]string) ([]DevopsProject, error)
	Create(params map[string]interface{}) (DevopsProject, error)
}
type DevopsProjectListWrapper struct {
	Data struct {
		Page      int             `json:"page"`
		Size      int             `json:"size"`
		Total     int             `json:"total"`
		TotalPage int             `json:"totalPage"`
		Docs      []DevopsProject `json:"docs"`
	} `json:"data"`
}

type DevopsProjectWrapper struct {
	Data DevopsProject `json:"data"`
}
type DevopsProject struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   string `json:"isDefault"`
	CreatedAt   string `json:"createdAt"`
}
type devopsproject struct {
	client *Client
}

// Get devopsproject detail
func (s *devopsproject) Get(id string) (DevopsProject, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/teams/projects/"+id, map[string]string{})
	var response DevopsProjectWrapper
	var nilres DevopsProject
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (s *devopsproject) List(params map[string]string) ([]DevopsProject, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/teams/projects", params)
	var response DevopsProjectListWrapper
	var nilres []DevopsProject
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Create a new devopsproject
func (s *devopsproject) Create(params map[string]interface{}) (DevopsProject, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/teams/projects", params)
	var response DevopsProject
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
