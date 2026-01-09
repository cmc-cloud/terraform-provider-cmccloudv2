package gocmcapiv2

import (
	"encoding/json"
	"strconv"
)

type ContainerRegistryService interface {
	Get(devops_project_id string, id string) (ContainerRegistryRepository, error)
	List(devops_project_id string, params map[string]string) ([]ContainerRegistryRepository, error)
	Create(devops_project_id string, params map[string]interface{}) (ContainerRegistryRepository, error)
	Delete(devops_project_id string, id string) (ActionResponse, error)
}
type ContainerRegistryRepositoryListWrapper struct {
	Data struct {
		Page      int                           `json:"page"`
		Size      int                           `json:"size"`
		Total     int                           `json:"total"`
		TotalPage int                           `json:"totalPage"`
		Docs      []ContainerRegistryRepository `json:"docs"`
	} `json:"data"`
}

type ContainerRegistryRepositoryWrapper struct {
	Data struct {
		Repository ContainerRegistryRepository `json:"repository"`
	} `json:"data"`
}
type ContainerRegistryRepository struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URI       string `json:"uri"`
	ProjectId int    `json:"projectId"`
	CreatedAt string `json:"createdAt"`
}
type containerregistry struct {
	client *Client
}

// Get containerregistry detail
func (s *containerregistry) Get(devops_project_id string, id string) (ContainerRegistryRepository, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/projects/"+devops_project_id+"/repositories/"+id, map[string]string{})
	var response ContainerRegistryRepositoryWrapper
	var nilres ContainerRegistryRepository
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Repository, nil
}

func (s *containerregistry) List(devops_project_id string, params map[string]string) ([]ContainerRegistryRepository, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/projects/"+devops_project_id+"/repositories?projectId="+devops_project_id+"&page=1&size=1000", params)
	var response ContainerRegistryRepositoryListWrapper
	var nilres []ContainerRegistryRepository
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a containerregistry
func (s *containerregistry) Delete(devops_project_id string, id string) (ActionResponse, error) {
	devops_project_id_int, _ := strconv.Atoi(devops_project_id)
	id_int, _ := strconv.Atoi(id)
	return s.client.PerformDeleteWithBody("cloudops-core/api/v1/repositories", map[string]interface{}{
		"projectId": devops_project_id_int,
		"repoIds":   []int{id_int},
	})
}

// Create a new containerregistry
func (s *containerregistry) Create(devops_project_id string, params map[string]interface{}) (ContainerRegistryRepository, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/projects/"+devops_project_id+"/repositories", params)
	var response ContainerRegistryRepositoryWrapper
	var nilres ContainerRegistryRepository
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Repository, nil
}
