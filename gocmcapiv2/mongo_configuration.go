package gocmcapiv2

import (
	"encoding/json"
)

// MongoConfigurationService interface
type MongoConfigurationService interface {
	Get(id string) (MongoConfiguration, error)
	GetDefaultConfiguration(id string) (MongoConfiguration, error)
	List(params map[string]string) ([]MongoConfiguration, error)
	Create(params map[string]interface{}) (MongoConfiguration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error)
}

type MongoConfigurationWrapper struct {
	Data MongoConfiguration `json:"data"`
}

type MongoConfigurationListWrapper struct {
	Data struct {
		Docs      []MongoConfiguration `json:"docs"`
		Page      int                  `json:"page"`
		Size      int                  `json:"size"`
		Total     int                  `json:"total"`
		TotalPage int                  `json:"totalPage"`
	} `json:"data"`
}
type MongoConfigurationParameter struct {
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

// MongoConfiguration object
type MongoConfiguration struct {
	ID                 string                        `json:"id"`
	ID2                string                        `json:"groupConfigId"`
	Name               string                        `json:"name"`
	Description        string                        `json:"description"`
	DatastoreName      string                        `json:"datastoreName"`
	DatastoreVersionID string                        `json:"datastoreVersionId"`
	DatastoreVersion   string                        `json:"datastoreVersion"`
	DatastoreMode      string                        `json:"datastoreMode"`
	DatastoreModeID    string                        `json:"datastoreModeId"`
	CreatedAt          string                        `json:"createdAt"`
	IsGroupDefault     bool                          `json:"isGroupDefault"`
	Parameters         []MongoConfigurationParameter `json:"configurations"`
}

type mongoconfiguration struct {
	client *Client
}

// Get mongoconfiguration detail
func (v *mongoconfiguration) Get(id string) (MongoConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/group-configuration/"+id, map[string]string{})
	var response MongoConfigurationWrapper
	var nilres MongoConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (v *mongoconfiguration) GetDefaultConfiguration(id string) (MongoConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/configurations-default/"+id, map[string]string{})
	var response MongoConfigurationWrapper
	var nilres MongoConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}
func (s *mongoconfiguration) List(params map[string]string) ([]MongoConfiguration, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response MongoConfigurationListWrapper
	var nilres []MongoConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a mongoconfiguration
func (v *mongoconfiguration) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/group-configuration", map[string]interface{}{"groupConfigIds": []string{id}})

}
func (v *mongoconfiguration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/configurations/"+id, params)
}
func (v *mongoconfiguration) UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/configurations/"+id, map[string]interface{}{
		"configurations": params,
	})
}

func (s *mongoconfiguration) Create(params map[string]interface{}) (MongoConfiguration, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response MongoConfiguration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
