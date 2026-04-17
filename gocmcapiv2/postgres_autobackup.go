package gocmcapiv2

import (
	"encoding/json"
)

// PostgresAutoBackupService interface
type PostgresAutoBackupService interface {
	Get(id string) (PostgresAutoBackup, error)
	List(params map[string]string) ([]PostgresAutoBackup, error)
	Create(params map[string]interface{}) (PostgresAutoBackup, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

type PostgresAutoBackupWrapper struct {
	Data PostgresAutoBackup `json:"data"`
}

type PostgresAutoBackupListWrapper struct {
	Data struct {
		Docs      []PostgresAutoBackup `json:"docs"`
		Page      int                  `json:"page"`
		Size      int                  `json:"size"`
		Total     int                  `json:"total"`
		TotalPage int                  `json:"totalPage"`
	} `json:"data"`
}
type PostgresAutoBackup struct {
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

type postgresautobackup struct {
	client *Client
}

// Get postgresautobackup detail
func (v *postgresautobackup) Get(id string) (PostgresAutoBackup, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/backup-schedule/"+id, map[string]string{})
	var response PostgresAutoBackupWrapper
	var nilres PostgresAutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data, nil
}

func (s *postgresautobackup) List(params map[string]string) ([]PostgresAutoBackup, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response PostgresAutoBackupListWrapper
	var nilres []PostgresAutoBackup
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nilres, err
	}
	return response.Data.Docs, nil
}

// Delete a postgresautobackup
func (v *postgresautobackup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/backup-schedule", map[string]interface{}{"backupScheduleIds": []string{id}})

}
func (v *postgresautobackup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/backup-schedule/"+id, params)
}

func (s *postgresautobackup) Create(params map[string]interface{}) (PostgresAutoBackup, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/backup-schedule", params)
	var response PostgresAutoBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
