package gocmcapiv2

import (
	"encoding/json"
)

// BackupService interface
type BackupService interface {
	Get(id string) (Backup, error)
	List(params map[string]string) ([]Backup, error)
	Delete(id string) (ActionResponse, error)
	Rename(id string, newName string) (ActionResponse, error)
}

// Backup object
type Backup struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	Size          int     `json:"size"`
	CreatedAt     string  `json:"created_at"`
	Name          string  `json:"name"`
	VolumeID      string  `json:"volume_id"`
	IsIncremental bool    `json:"is_incremental"`
	RealSize      int     `json:"real_size"`
	RealSizeGB    float64 `json:"real_size_gb"`
	Volume        struct {
		ID        string `json:"id"`
		DeletedAt string `json:"deleted_at"`
		Name      string `json:"name"`
	} `json:"volume"`
	// SnapshotID    string  `json:"snapshot_id"`
}
type backup struct {
	client *Client
}

// Get backup detail
func (v *backup) Get(id string) (Backup, error) {
	jsonStr, err := v.client.Get("backup/"+id, map[string]string{})
	var backup Backup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &backup)
	}
	return backup, err
}

func (s *backup) List(params map[string]string) ([]Backup, error) {
	restext, err := s.client.Get("backup", params)
	items := make([]Backup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a backup
func (v *backup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("backup/" + id)
}
func (v *backup) Rename(id string, newName string) (ActionResponse, error) {
	return v.client.PerformUpdate("backup/"+id, map[string]interface{}{"id": id, "name": newName})
}
