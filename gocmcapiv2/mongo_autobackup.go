package gocmcapiv2

import (
	"encoding/json"
)

// MongoAutoBackupService interface
type MongoAutoBackupService interface {
	Get(id string) (MongoAutoBackup, error)
	List(params map[string]string) ([]MongoAutoBackup, error)
	Create(params map[string]interface{}) (MongoAutoBackup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

type MongoAutoBackupWrapper struct {
	Data MongoAutoBackup `json:"data"`
}

type MongoAutoBackupListWrapper struct {
	Data struct {
		Docs      []MongoAutoBackup `json:"docs"`
		Page      int               `json:"page"`
		Size      int               `json:"size"`
		Total     int               `json:"total"`
		TotalPage int               `json:"totalPage"`
	} `json:"data"`
}
type MongoAutoBackup struct {
	ID               string `json:"id"`
	InstanceID       string `json:"instanceId"`
	InstanceName     string `json:"instanceName"`
	Hour             int    `json:"hour"`
	Minute           int    `json:"minute"`
	TimeZone         string `json:"timeZone"`
	IntervalNum      int    `json:"intervalNum"`
	IntervalType     string `json:"intervalType"`
	NextBackupTime   string `json:"nextBackupTime"`
	KeepRecordBackup int    `json:"keepRecordBackup"`
	Created          string `json:"created"`
	Updated          string `json:"updated"`
}

type mongoautobackup struct {
	client *Client
}

// Get mongoautobackup detail
func (v *mongoautobackup) Get(id string) (MongoAutoBackup, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/backup-schedule/"+id, map[string]string{})
	var response MongoAutoBackupWrapper
	var nilres MongoAutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (s *mongoautobackup) List(params map[string]string) ([]MongoAutoBackup, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response MongoAutoBackupListWrapper
	var nilres []MongoAutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a mongoautobackup
func (v *mongoautobackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/backup-schedule", map[string]interface{}{"backupScheduleIds": []string{id}})

}
func (v *mongoautobackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/backup-schedule/"+id, params)
}

func (s *mongoautobackup) Create(params map[string]interface{}) (MongoAutoBackup, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response MongoAutoBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
