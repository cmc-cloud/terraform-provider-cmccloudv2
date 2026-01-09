package gocmcapiv2

import (
	"encoding/json"
	"fmt"
)

// DatabaseConfigurationService interface
type DatabaseConfigurationService interface {
	Get(id string) (DatabaseConfiguration, error)
	List(params map[string]string) ([]DatabaseConfiguration, error)
	Create(params map[string]interface{}) (DatabaseConfiguration, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error)
}

type ArrayOrMap struct {
	Object   map[string]interface{}
	Array    []interface{}
	IsObject bool
}

func (v *ArrayOrMap) UnmarshalJSON(data []byte) error {
	// Try to unmarshal data into an array
	var arr []interface{}
	if err := json.Unmarshal(data, &arr); err == nil {
		v.Array = arr
		v.IsObject = false
		return nil
	}

	// Try to unmarshal data into an object
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err == nil {
		v.Object = obj
		v.IsObject = true
		return nil
	}

	return fmt.Errorf("ArrayOrMap is neither an array nor an object")
}

// DatabaseConfiguration object
type DatabaseConfiguration struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Created              string     `json:"created"`
	DatastoreVersionID   string     `json:"datastore_version_id"`
	DatastoreName        string     `json:"datastore_name"`
	DatastoreVersionName string     `json:"datastore_version_name"`
	Parameters           ArrayOrMap `json:"values"`
}

type databaseconfiguration struct {
	client *Client
}

// Get databaseconfiguration detail
func (v *databaseconfiguration) Get(id string) (DatabaseConfiguration, error) {
	jsonStr, err := v.client.Get("dbaas/configuration/"+id, map[string]string{})
	var obj DatabaseConfiguration
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *databaseconfiguration) List(params map[string]string) ([]DatabaseConfiguration, error) {
	restext, err := s.client.Get("dbaas/configuration", params)
	items := make([]DatabaseConfiguration, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a databaseconfiguration
func (v *databaseconfiguration) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("dbaas/configuration/" + id)
}
func (v *databaseconfiguration) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/configuration/"+id, params)
}
func (v *databaseconfiguration) UpdateParameters(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/configuration/"+id+"/parameter", params)
}

func (s *databaseconfiguration) Create(params map[string]interface{}) (DatabaseConfiguration, error) {
	jsonStr, err := s.client.Post("dbaas/configuration", params)
	var response DatabaseConfiguration
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
