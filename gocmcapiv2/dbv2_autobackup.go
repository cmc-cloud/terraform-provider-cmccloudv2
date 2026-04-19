package gocmcapiv2

import (
	"encoding/json"
)

// DBv2AutoBackupService interface
type DBv2AutoBackupService interface {
	Get(id string) (DBv2AutoBackup, error)
	List(params map[string]string) ([]DBv2AutoBackup, error)
	Create(params map[string]interface{}) (DBv2AutoBackupCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}
type DBv2AutoBackupCreateResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}
type DBv2AutoBackupWrapper struct {
	Data DBv2AutoBackup `json:"data"`
}

type DBv2AutoBackupListWrapper struct {
	Data struct {
		Docs      []DBv2AutoBackup `json:"docs"`
		Page      int              `json:"page"`
		Size      int              `json:"size"`
		Total     int              `json:"total"`
		TotalPage int              `json:"totalPage"`
	} `json:"data"`
}
type DBv2AutoBackup struct {
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

type dbv2autobackup struct {
	client *Client
}

// Get dbv2autobackup detail
func (v *dbv2autobackup) Get(id string) (DBv2AutoBackup, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/backup-schedule/"+id, map[string]string{})
	var response DBv2AutoBackupWrapper
	var nilres DBv2AutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (s *dbv2autobackup) List(params map[string]string) ([]DBv2AutoBackup, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response DBv2AutoBackupListWrapper
	var nilres []DBv2AutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a dbv2autobackup
func (v *dbv2autobackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/backup-schedule", map[string]interface{}{"backupScheduleIds": []string{id}})

}
func (v *dbv2autobackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/backup-schedule/"+id, params)
}

func (s *dbv2autobackup) Create(params map[string]interface{}) (DBv2AutoBackupCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response DBv2AutoBackupCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
