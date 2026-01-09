package gocmcapiv2

import (
	"encoding/json"
)

// RedisConfigurationService interface
type RedisConfigurationService interface {
	Get(id string) (RedisConfiguration, error)
	GetDefaultConfiguration(id string) (RedisConfiguration, error)
	List(params map[string]string) ([]RedisConfiguration, error)
	Create(params map[string]interface{}) (RedisConfiguration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error)
}

type RedisConfigurationWrapper struct {
	Data RedisConfiguration `json:"data"`
}

type RedisConfigurationListWrapper struct {
	Data struct {
		Docs      []RedisConfiguration `json:"docs"`
		Page      int                  `json:"page"`
		Size      int                  `json:"size"`
		Total     int                  `json:"total"`
		TotalPage int                  `json:"totalPage"`
	} `json:"data"`
}
type RedisDatastore struct {
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
type RedisConfigurationParameter struct {
	Name  string `json:"paramName"`
	Value string `json:"paramValue"`
}

// RedisConfiguration object
type RedisConfiguration struct {
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
	Parameters         []RedisConfigurationParameter `json:"parameters"`
}

type redisconfiguration struct {
	client *Client
}

// Get redisconfiguration detail
func (v *redisconfiguration) Get(id string) (RedisConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/group-configuration/"+id, map[string]string{})
	var response RedisConfigurationWrapper
	var nilres RedisConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (v *redisconfiguration) GetDefaultConfiguration(id string) (RedisConfiguration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/configurations-default/"+id, map[string]string{})
	var response RedisConfigurationWrapper
	var nilres RedisConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}
func (s *redisconfiguration) List(params map[string]string) ([]RedisConfiguration, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response RedisConfigurationListWrapper
	var nilres []RedisConfiguration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a redisconfiguration
func (v *redisconfiguration) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("cloudops-core/api/v1/dbaas/group-configuration/" + id)
}
func (v *redisconfiguration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/group-configuration/"+id, params)
}
func (v *redisconfiguration) UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/group-configuration/"+id+"/parameter", map[string]interface{}{
		"configurations": params,
	})
}

func (s *redisconfiguration) Create(params map[string]interface{}) (RedisConfiguration, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response RedisConfiguration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
