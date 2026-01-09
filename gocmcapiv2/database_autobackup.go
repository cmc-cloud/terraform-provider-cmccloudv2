package gocmcapiv2

import (
	"encoding/json"
)

// DatabaseAutoBackupService interface
type DatabaseAutoBackupService interface {
	Get(id string) (DatabaseAutoBackup, error)
	List(params map[string]string) ([]DatabaseAutoBackup, error)
	Create(params map[string]interface{}) (DatabaseAutoBackup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

// DatabaseAutoBackup object
type DatabaseAutoBackup struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	VolumeID     string      `json:"volume_id"`
	VolumeName   string      `json:"volume_name"`
	VolumeSize   int         `json:"volume_size"`
	Time         string      `json:"time"`
	Interval     int         `json:"interval"`
	MaxKeep      int         `json:"max_keep"`
	Created      string      `json:"created"`
	IsFullBackup BoolFromInt `json:"is_full_backup"`
	LastRun      string      `json:"last_run"`
	Status       string      `json:"status"`
	Error        string      `json:"error"`
}

type databaseautobackup struct {
	client *Client
}

// Get volume detail
func (v *databaseautobackup) Get(id string) (DatabaseAutoBackup, error) {
	jsonStr, err := v.client.Get("dbaas/auto-backup/"+id, map[string]string{})
	var volume DatabaseAutoBackup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &volume)
	}
	return volume, err
}

func (s *databaseautobackup) List(params map[string]string) ([]DatabaseAutoBackup, error) {
	restext, err := s.client.Get("dbaas/auto-backup", params)
	items := make([]DatabaseAutoBackup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a volume
func (v *databaseautobackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("dbaas/auto-backup/" + id)
}
func (v *databaseautobackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("dbaas/auto-backup/"+id, params)
}
func (s *databaseautobackup) Create(params map[string]interface{}) (DatabaseAutoBackup, error) {
	jsonStr, err := s.client.Post("dbaas/auto-backup", params)
	var response DatabaseAutoBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
