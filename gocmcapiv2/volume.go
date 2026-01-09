package gocmcapiv2

import (
	"encoding/json"
)

// VolumeService interface
type VolumeService interface {
	Get(id string) (Volume, error)
	List(params map[string]string) ([]Volume, error)
	Create(params map[string]interface{}) (Volume, error)
	CreateBackup(id string, params map[string]interface{}) (Backup, error)
	CreateSnapshot(id string, params map[string]interface{}) (Snapshot, error)
	Delete(id string) (ActionResponse, error)
	Resize(id string, new_size int) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Attach(id string, params map[string]interface{}) (ActionResponse, error)
	Detach(id string, server_id string) (ActionResponse, error)
	Rename(id string, name string) (ActionResponse, error)
}

// Volume object
type Volume struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	Size             int    `json:"size"`
	AvailabilityZone string `json:"availability_zone"`
	CreatedAt        string `json:"created_at"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	VolumeType       string `json:"volume_type"`
	Bootable         string `json:"bootable"`
	Encrypted        bool   `json:"encrypted"`
	Multiattach      bool   `json:"multiattach"`
	Attachments      []struct {
		ID           string `json:"id"`
		AttachmentID string `json:"attachment_id"`
		VolumeID     string `json:"volume_id"`
		ServerID     string `json:"server_id"`
		Device       string `json:"device"`
	} `json:"attachments"`
	Tags                []Tag  `json:"tags"`
	BillingMode         string `json:"billing_mode"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	EncryptionKeyID     string `json:"encryption_key_id"`
}
type volume struct {
	client *Client
}

// Get volume detail
func (v *volume) Get(id string) (Volume, error) {
	jsonStr, err := v.client.Get("volume/"+id, map[string]string{})
	var volume Volume
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &volume)
	}
	return volume, err
}

func (s *volume) List(params map[string]string) ([]Volume, error) {
	restext, err := s.client.Get("volume", params)
	items := make([]Volume, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a volume
func (v *volume) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("volume/" + id)
}
func (v *volume) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("volume/"+id, params)
}
func (v *volume) Resize(id string, new_size int) (ActionResponse, error) {
	return v.client.PerformAction("volume/"+id+"/resize", map[string]interface{}{"size": new_size})
}
func (v *volume) Attach(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformAction("volume/"+id+"/attach", params)
}
func (v *volume) Detach(id string, server_id string) (ActionResponse, error) {
	return v.client.PerformAction("volume/"+id+"/detach", map[string]interface{}{"server_id": server_id})
}
func (v *volume) Rename(id string, name string) (ActionResponse, error) {
	return v.client.PerformAction("volume/"+id+"/change_name", map[string]interface{}{"name": name})
}
func (s *volume) CreateSnapshot(id string, params map[string]interface{}) (Snapshot, error) {
	jsonStr, err := s.client.Post("volume/"+id+"/snapshot", params)
	var response Snapshot
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *volume) CreateBackup(id string, params map[string]interface{}) (Backup, error) {
	jsonStr, err := s.client.Post("volume/"+id+"/backup", params)
	var response Backup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *volume) Create(params map[string]interface{}) (Volume, error) {
	jsonStr, err := s.client.Post("volume", params)
	var response Volume
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
