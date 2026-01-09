package gocmcapiv2

import (
	"encoding/json"
)

// SnapshotService interface
type SnapshotService interface {
	Get(id string) (Snapshot, error)
	List(params map[string]string) ([]Snapshot, error)
	Delete(id string) (ActionResponse, error)
	Rename(id string, newName string) (ActionResponse, error)
}

// Snapshot object
type Snapshot struct {
	ID            string  `json:"id"`
	Status        string  `json:"status"`
	Size          int     `json:"size"`
	CreatedAt     string  `json:"created_at"`
	Name          string  `json:"name"`
	VolumeID      string  `json:"volume_id"`
	IsIncremental bool    `json:"is_incremental"`
	SnapshotID    any     `json:"snapshot_id"`
	RealSizeGB    float64 `json:"real_size_gb"`
	Volume        struct {
		ID        string `json:"id"`
		DeletedAt any    `json:"deleted_at"`
		Name      string `json:"name"`
	} `json:"volume"`
}

type snapshot struct {
	client *Client
}

// Get snapshot detail
func (v *snapshot) Get(id string) (Snapshot, error) {
	jsonStr, err := v.client.Get("snapshot/"+id, map[string]string{})
	var snapshot Snapshot
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &snapshot)
	}
	return snapshot, err
}

func (s *snapshot) List(params map[string]string) ([]Snapshot, error) {
	restext, err := s.client.Get("snapshot", params)
	items := make([]Snapshot, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a snapshot
func (v *snapshot) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("snapshot/" + id)
}
func (v *snapshot) Rename(id string, newName string) (ActionResponse, error) {
	return v.client.PerformUpdate("snapshot/"+id, map[string]interface{}{"id": id, "name": newName})
}
