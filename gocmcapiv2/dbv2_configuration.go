package gocmcapiv2

import (
	"encoding/json"
	"strings"
)

// DBv2ConfigurationService interface
type DBv2ConfigurationService interface {
	Get(id string) (DBv2Configuration, error)
	GetDefaultConfiguration(id string) (DBv2Configuration, error)
	ListDatastore(dbType string, params map[string]string) ([]Datastore, error)
	List(params map[string]string) ([]DBv2Configuration, error)
	Create(params map[string]interface{}) (DBv2Configuration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error)
}

type DBv2Configuration struct {
	ID                 string                       `json:"id"`
	ID2                string                       `json:"groupConfigId"`
	Name               string                       `json:"name"`
	Description        string                       `json:"description"`
	DatastoreName      string                       `json:"datastoreName"`
	DatastoreVersionID string                       `json:"datastoreVersionId"`
	DatastoreVersion   string                       `json:"datastoreVersion"`
	DatastoreMode      string                       `json:"datastoreMode"`
	DatastoreModeID    string                       `json:"datastoreModeId"`
	CreatedAt          string                       `json:"createdAt"`
	IsGroupDefault     bool                         `json:"isGroupDefault"`
	Parameters         []DBv2ConfigurationParameter `json:"configurations"`
}

type DBv2ConfigurationParameter struct {
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

type DBv2ConfigurationWrapper struct {
	Data DBv2Configuration `json:"data"`
}

type DBv2ConfigurationListWrapper struct {
	Data struct {
		Docs      []DBv2Configuration `json:"docs"`
		Page      int                 `json:"page"`
		Size      int                 `json:"size"`
		Total     int                 `json:"total"`
		TotalPage int                 `json:"totalPage"`
	} `json:"data"`
}

type dbv2configuration struct {
	client *Client
}

// Get dbaasconfiguration detail
func (v *dbv2configuration) Get(id string) (DBv2Configuration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/group-configuration/"+id, map[string]string{})
	var response DBv2ConfigurationWrapper
	var nilres DBv2Configuration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (v *dbv2configuration) GetDefaultConfiguration(id string) (DBv2Configuration, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/configurations-default/"+id, map[string]string{})
	var response DBv2ConfigurationWrapper
	var nilres DBv2Configuration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (s *dbv2configuration) ListDatastore(dbType string, params map[string]string) ([]Datastore, error) {
	switch {
	case strings.Contains(strings.ToLower(dbType), "mongo"):
		dbType = "mongodb"
	case strings.Contains(strings.ToLower(dbType), "postgres"):
		dbType = "postgresql"
	case strings.Contains(strings.ToLower(dbType), "redis"):
		dbType = "redis"
	case strings.Contains(strings.ToLower(dbType), "mysql"):
		dbType = "mysql"
	}
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode="+dbType, params)
	var obj DatastoreListWrapper

	if err != nil {
		return []Datastore{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []Datastore{}, err
	}
	return obj.Data.Docs, err
}

func (s *dbv2configuration) List(params map[string]string) ([]DBv2Configuration, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response DBv2ConfigurationListWrapper
	var nilres []DBv2Configuration
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a dbaasconfiguration
func (v *dbv2configuration) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/group-configuration", map[string]interface{}{"groupConfigIds": []string{id}})

}
func (v *dbv2configuration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/configurations/"+id, params)
}
func (v *dbv2configuration) UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/configurations/"+id, map[string]interface{}{
		"configurations": params,
	})
}

func (s *dbv2configuration) Create(params map[string]interface{}) (DBv2Configuration, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/group-configuration", params)
	var response DBv2Configuration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
