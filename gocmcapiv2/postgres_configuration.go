package gocmcapiv2

import (
	"encoding/json"
)

// PostgresConfigurationService interface
type PostgresConfigurationService interface {
	Get(id string) (PostgresConfiguration, error)
	GetDefaultConfiguration(id string) (PostgresConfiguration, error)
	List(params map[string]string) ([]PostgresConfiguration, error)
	Create(params map[string]interface{}) (PostgresConfiguration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error)
}

type PostgresConfigurationWrapper struct {
	Data PostgresConfiguration `json:"data"`
}

type PostgresConfigurationListWrapper struct {
	Data struct {
		Docs      []PostgresConfiguration `json:"docs"`
		Page      int                     `json:"page"`
		Size      int                     `json:"size"`
		Total     int                     `json:"total"`
		TotalPage int                     `json:"totalPage"`
	} `json:"data"`
}
type PostgresDatastore struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	VersionInfos []struct {
		ID          string `json:"id"`
		VersionName string `json:"versionName"`
		ModeInfo    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"modeInfo"`
	} `json:"versionInfos"`
}
type PostgresConfigurationParameter struct {
	ID           string `json:"id"`
	Name         string `json:"paramName"`
	Value        string `json:"paramValue"`
	DefaultValue string `json:"defaultValue"`
	ValueRange   string `json:"valueRange"`
	ValueType    string `json:"valueType"`
	Description  string `json:"description"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// PostgresConfiguration object
type PostgresConfiguration struct {
	ID                 string                           `json:"id"`
	ID2                string                           `json:"groupConfigId"`
	Name               string                           `json:"name"`
	Description        string                           `json:"description"`
	DatastoreName      string                           `json:"datastoreName"`
	DatastoreVersionID string                           `json:"datastoreVersionId"`
	DatastoreVersion   string                           `json:"datastoreVersion"`
	DatastoreMode      string                           `json:"datastoreMode"`
	DatastoreModeID    string                           `json:"datastoreModeId"`
	CreatedAt          string                           `json:"createdAt"`
	IsGroupDefault     bool                             `json:"isGroupDefault"`
	Parameters         []PostgresConfigurationParameter `json:"configurations"`
}

type Postgresconfiguration struct {
	client *Client
}

// Get Postgresconfiguration detail
func (v *Postgresconfiguration) Get(id string) (PostgresConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/group-configuration/"+id, map[string]string{})
	var response PostgresConfigurationWrapper
	var nilres PostgresConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (v *Postgresconfiguration) GetDefaultConfiguration(id string) (PostgresConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/configurations-default/"+id, map[string]string{})
	var response PostgresConfigurationWrapper
	var nilres PostgresConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}
func (s *Postgresconfiguration) List(params map[string]string) ([]PostgresConfiguration, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response PostgresConfigurationListWrapper
	var nilres []PostgresConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a Postgresconfiguration
func (v *Postgresconfiguration) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("cloudops-core/api/v1/dbaas/group-configuration/" + id)
}
func (v *Postgresconfiguration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/group-configuration/"+id, params)
}
func (v *Postgresconfiguration) UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/group-configuration/"+id+"/parameter", map[string]interface{}{
		"configurations": params,
	})
}

func (s *Postgresconfiguration) Create(params map[string]interface{}) (PostgresConfiguration, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response PostgresConfiguration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
