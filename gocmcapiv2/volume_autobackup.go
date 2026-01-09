package gocmcapiv2

import (
	"encoding/json"
)

// VolumeAutoBackupService interface
type VolumeAutoBackupService interface {
	Get(id string) (VolumeAutoBackup, error)
	Create(params map[string]interface{}) (VolumeAutoBackup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

// VolumeAutoBackup object
type VolumeAutoBackup struct {
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
	LastRun      any         `json:"last_run"`
	Status       string      `json:"status"`
	Error        any         `json:"error"`
	NextRun      string      `json:"next_run"`
}
type volumeautobackup struct {
	client *Client
}

// Get volume detail
func (v *volumeautobackup) Get(id string) (VolumeAutoBackup, error) {
	jsonStr, err := v.client.Get("backup/auto-backup/"+id, map[string]string{})
	var volume VolumeAutoBackup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &volume)
	}
	return volume, err
}

// Delete a volume
func (v *volumeautobackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("backup/auto-backup/" + id)
}
func (v *volumeautobackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("backup/auto-backup/"+id, params)
}
func (s *volumeautobackup) Create(params map[string]interface{}) (VolumeAutoBackup, error) {
	jsonStr, err := s.client.Post("backup/auto-backup", params)
	var response VolumeAutoBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
