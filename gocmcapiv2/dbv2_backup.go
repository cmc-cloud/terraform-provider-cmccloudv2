package gocmcapiv2

import (
	"encoding/json"
	"strings"
)

// DBv2InstanceService interface
type DBv2BackupService interface {
	Get(dbType string, instanceId string, backupId string) (DBv2Backup, error)
	Delete(id string) (ActionResponse, error)
	List(dbType string, instanceId string) ([]DBv2Backup, error)
	Create(instanceId string, name string) (CreateBackupResponse, error)
}
type CreateBackupResponse struct {
	Status int    `json:"status"`
	Error  any    `json:"error"`
	Msg    string `json:"msg"`
	Data   struct {
		BackupID string `json:"backupId"`
	} `json:"data"`
	Success bool `json:"success"`
}
type DBv2BackupListWrapper struct {
	Status      int    `json:"status"`
	Error       any    `json:"error"`
	Msg         string `json:"msg"`
	CurrentTime string `json:"currentTime"`
	Data        struct {
		Page      int          `json:"page"`
		Size      int          `json:"size"`
		Total     int          `json:"total"`
		TotalPage int          `json:"totalPage"`
		Docs      []DBv2Backup `json:"docs"`
	} `json:"data"`
	Success bool `json:"success"`
}
type DBv2Backup struct {
	ID                 string `json:"backupId"`
	BackupName         string `json:"backupName"`
	DatastoreName      string `json:"datastoreName"`
	DatastoreVersion   string `json:"datastoreVersion"`
	InstanceID         string `json:"instanceId"`
	InstanceName       string `json:"instanceName"`
	Size               int    `json:"size"`
	BackupStrategyType string `json:"backupStrategyType"`
	Status             string `json:"status"`
	Created            string `json:"created"`
	DatastoreCode      string `json:"datastoreCode"`
}
type dbv2backup struct {
	client *Client
}

func (v *dbv2backup) Create(instanceId string, name string) (CreateBackupResponse, error) {
	jsonStr, err := v.client.Post("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": instanceId,
		"action":     "create_backup",
		"requestData": map[string]interface{}{
			"backupStrategyType": "S3",
			"name":               name, //              string(bytes),
		},
	})
	var res CreateBackupResponse
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		return res, err
	}
	return res, err
}
func (v *dbv2backup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/backup", map[string]interface{}{"backupIds": []string{id}})
}

func (v *dbv2backup) Get(dbType string, id string, backupId string) (DBv2Backup, error) {
	databases, err := v.List(dbType, id)
	if err != nil {
		return DBv2Backup{}, err
	}
	for _, db := range databases {
		if db.ID == backupId {
			return db, nil
		}
	}
	return DBv2Backup{}, &NotFoundError{}
}
func (v *dbv2backup) List(dbType string, id string) ([]DBv2Backup, error) {
	// fix dbType, if contains mongo => mongodb, if contains postgres => postgresql, if contains redis => redis, if contains mysql => mysql
	switch {
	case strings.Contains(strings.ToLower(dbType), "mongo"):
		dbType = "mongodb"
	case strings.Contains(strings.ToLower(dbType), "postgres"):
		dbType = "postgresql"
	case strings.Contains(strings.ToLower(dbType), "redis"):
		dbType = "redis"
	case strings.Contains(strings.ToLower(dbType), "mysql"):
		dbType = "mysql"
	}
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/backup-list?page=1&size=20&instanceId="+id+"&datastoreCode="+dbType, map[string]string{})
	var obj DBv2BackupListWrapper

	if err != nil {
		return []DBv2Backup{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []DBv2Backup{}, err
	}
	return obj.Data.Docs, err
}
