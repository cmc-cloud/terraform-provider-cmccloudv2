package gocmcapiv2

import (
	"encoding/json"
)

// ApiGatewayService interface
type ApiGatewayService interface {
	Get(id string) (ApiGateway, error)
	List(params map[string]string) ([]ApiGateway, error)
	Create(params map[string]interface{}) (ApiGatewayCreateResponse, error)
	Rename(id string, name string) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)

	GetBackup(id string) (ApiGatewayBackup, error)
	DeleteBackup(id string) (ActionResponse, error)
	CreateBackup(instanceId string, name string) (ApiGatewayCreateResponse, error)
}

type ApiGatewayWrapper struct {
	Data ApiGateway `json:"data"`
}
type ApiGatewayBackupWrapper struct {
	Data ApiGatewayBackup `json:"data"`
}

type ApiGateway struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	PortalProjectID       string   `json:"portalProjectId"`
	RegionID              string   `json:"regionId"`
	TeamCodeID            string   `json:"teamCodeId"`
	Status                string   `json:"status"`
	IP                    []string `json:"ip"`
	VolumeSize            int      `json:"volumeSize"`
	PublicAccess          int      `json:"publicAccess"`
	PublicAccessBandwidth int      `json:"publicAccessBandwidth"`
	CreatedAt             string   `json:"createdAt"`
	UpdatedAt             string   `json:"updatedAt"`
	Flavor                struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		RAM         int    `json:"ram"`
		Disk        int    `json:"disk"`
		Vcpus       int    `json:"vcpus"`
		GpuName     string `json:"gpuName"`
		GpuQuantity int    `json:"gpuQuantity"`
	} `json:"flavor"`
}
type ApiGatewayCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}
type ApiGatewayListWrapper struct {
	Data struct {
		Docs      []ApiGateway `json:"docs"`
		Page      int          `json:"page"`
		Size      int          `json:"size"`
		Total     int          `json:"total"`
		TotalPage int          `json:"totalPage"`
	} `json:"data"`
}
type ApiGatewayBackup struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Size     int    `json:"size"`
	Instance struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		PortalProjectID string `json:"portalProjectId"`
		RegionID        string `json:"regionId"`
		TeamCodeID      string `json:"teamCodeId"`
		Status          string `json:"status"`
	} `json:"instance"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
type apigateway struct {
	client *Client
}

// Get apigateway detail
func (v *apigateway) Get(id string) (ApiGateway, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/api-gateways/instances/"+id, map[string]string{})
	var obj ApiGatewayWrapper
	if err != nil {
		return ApiGateway{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return ApiGateway{}, err
	}
	return obj.Data, err
}

func (v *apigateway) List(params map[string]string) ([]ApiGateway, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/api-gateways/instances", params)
	var obj ApiGatewayListWrapper
	if err != nil {
		return []ApiGateway{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []ApiGateway{}, err
	}
	return obj.Data.Docs, err
}

// Delete a apigateway
func (v *apigateway) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/api-gateways/instances/"+id, map[string]interface{}{})
}
func (s *apigateway) Create(params map[string]interface{}) (ApiGatewayCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/api-gateways/instances", params)
	var response ApiGatewayCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *apigateway) Rename(id string, name string) (ActionResponse, error) {
	return s.client.PerformUpdate("cloudops-core/api/v1/api-gateways/instances/"+id, map[string]interface{}{"name": name})
}

func (s *apigateway) GetBackup(id string) (ApiGatewayBackup, error) {
	var response ApiGatewayBackupWrapper
	jsonStr, err := s.client.Get("cloudops-core/api/v1/api-gateways/backups/"+id, map[string]string{})
	if err != nil {
		return ApiGatewayBackup{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response.Data, err
}

func (s *apigateway) DeleteBackup(id string) (ActionResponse, error) {
	return s.client.PerformDelete("cloudops-core/api/v1/api-gateways/backups/" + id)
}

func (s *apigateway) CreateBackup(instanceId string, name string) (ApiGatewayCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/api-gateways/backups", map[string]interface{}{
		"instanceId": instanceId,
		"name":       name,
	})
	var response ApiGatewayCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
