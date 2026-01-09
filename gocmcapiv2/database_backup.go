package gocmcapiv2

import (
	"encoding/json"
)

// DatabaseBackupService interface
type DatabaseBackupService interface {
	Get(id string) (DatabaseBackup, error)
	List(params map[string]string) ([]DatabaseBackup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

type DatabaseSnapshot struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	InstanceID string  `json:"instance_id"`
	Created    string  `json:"created"`
	Size       float64 `json:"size"`
	Status     string  `json:"status"`
	RealSize   int     `json:"real_size"`
}
type DatabaseBackup struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	InstanceID string  `json:"instance_id"`
	ParentID   string  `json:"parent_id"`
	Created    string  `json:"created"`
	Size       float64 `json:"size"`
	Status     string  `json:"status"`
	RealSize   int     `json:"real_size"`
	RealSizeGB float64 `json:"real_size_gb"`
}

type databasebackup struct {
	client *Client
}

// Get databasebackup detail
func (v *databasebackup) Get(id string) (DatabaseBackup, error) {
	jsonStr, err := v.client.Get("dbaas/instance/"+id, map[string]string{})
	var obj DatabaseBackup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *databasebackup) List(params map[string]string) ([]DatabaseBackup, error) {
	restext, err := s.client.Get("dbaas/backup", params)
	items := make([]DatabaseBackup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a databasebackup
func (v *databasebackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("dbaas/instance/" + id)
}
func (v *databasebackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/instance/"+id, params)
}

func (s *databasebackup) Create(params map[string]interface{}) (DatabaseBackup, error) {
	jsonStr, err := s.client.Post("dbaas/instance", params)
	var response DatabaseBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
