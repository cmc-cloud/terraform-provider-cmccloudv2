package gocmcapiv2

import (
	"encoding/json"
)

// DatabaseInstanceService interface
type DatabaseInstanceService interface {
	Get(id string) (DatabaseInstance, error)
	List(params map[string]string) ([]DatabaseInstance, error)
	Create(params map[string]interface{}) (DatabaseInstance, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Resize(id string, flavor_id string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error)
	UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error)

	CreateBackup(id string, params map[string]interface{}) (DatabaseBackup, error)
	CreateSnapshot(id string, params map[string]interface{}) (DatabaseSnapshot, error)
}

// DatabaseInstance object
type DatabaseInstance struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	BillingMode     string `json:"billing_mode"`
	Status          string `json:"status"`
	OperatingStatus string `json:"operating_status"`
	Flavor          struct {
		ID string `json:"id"`
	} `json:"flavor"`
	Datastore struct {
		Type          string `json:"type"`
		Version       string `json:"version"`
		VersionNumber string `json:"version_number"`
	} `json:"datastore"`
	Region string `json:"region"`
	Access struct {
		IsPublic bool `json:"is_public"`
	} `json:"access"`
	Volume struct {
		Size int `json:"size"`
	} `json:"volume"`
	Created           string `json:"created"`
	ComputeInstanceID string `json:"compute_instance_id"`
}

type databaseinstance struct {
	client *Client
}

// Get databaseinstance detail
func (v *databaseinstance) Get(id string) (DatabaseInstance, error) {
	jsonStr, err := v.client.Get("dbaas/instance/"+id, map[string]string{})
	var obj DatabaseInstance
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *databaseinstance) List(params map[string]string) ([]DatabaseInstance, error) {
	restext, err := s.client.Get("dbaas/instance", params)
	items := make([]DatabaseInstance, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a databaseinstance
func (v *databaseinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("dbaas/instance/" + id)
}
func (v *databaseinstance) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/instance/"+id, params)
}
func (v *databaseinstance) Resize(id string, flavor_id string) (ActionResponse, error) {
	return v.client.PerformAction("dbaas/instance/"+id+"/resize", map[string]interface{}{"flavor_id": flavor_id})
}
func (v *databaseinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("dbaas/instance/"+id+"/resize_volume", map[string]interface{}{"size": volume_size})
}
func (v *databaseinstance) UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/instance/"+id+"/accessbility", params)
}
func (v *databaseinstance) UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/instance/"+id+"/upgrade_datastore_version", map[string]interface{}{"datastore_version": datastore_version})
}

func (s *databaseinstance) Create(params map[string]interface{}) (DatabaseInstance, error) {
	jsonStr, err := s.client.Post("dbaas/instance", params)
	var response DatabaseInstance
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (s *databaseinstance) CreateBackup(id string, params map[string]interface{}) (DatabaseBackup, error) {
	jsonStr, err := s.client.Post("dbaas/instance/"+id+"/backup", params)
	var response DatabaseBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *databaseinstance) CreateSnapshot(id string, params map[string]interface{}) (DatabaseSnapshot, error) {
	jsonStr, err := s.client.Post("dbaas/instance/"+id+"/snapshot", params)
	var response DatabaseSnapshot
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
