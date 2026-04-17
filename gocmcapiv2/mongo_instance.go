package gocmcapiv2

import (
	"encoding/json"
)

// MongoInstanceService interface
type MongoInstanceService interface {
	Get(id string) (MongoInstance, error)
	List(params map[string]string) ([]MongoInstance, error)
	ListDatastore(params map[string]string) ([]Datastore, error)
	Create(params map[string]interface{}) (MongoInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	SetConfigurationGroupId(id string, Mongo_configuration_id string) (ActionResponse, error)
	Resize(id string, flavorId string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	GetUser(id string, username string) (MongoUser, error)
	ListUsers(id string) ([]MongoUser, error)
	CreateUser(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateUser(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteUser(id string, username string) (ActionResponse, error)

	GetDatabase(id string, name string) (MongoDatabase, error)
	CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteDatabase(id string, databaseName string) (ActionResponse, error)
	ListDatabases(id string) ([]MongoDatabase, error)

	GetBackup(id string, backup_id string) (MongoBackup, error)
	DeleteBackup(id string) (ActionResponse, error)
	ListBackups(id string) ([]MongoBackup, error)
	CreateBackup(id string, name string) (ActionResponse, error)
}

type MongoUser struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type MongoDatabase struct {
	Name  string `json:"database"`
	Owner string `json:"owner"`
}
type MongoUserWrapper struct {
	Data MongoUser `json:"data"`
}

type MongoUserListWrapper struct {
	Data struct {
		Docs []MongoUser `json:"docs"`
	} `json:"data"`
}

type MongoInstanceWrapper struct {
	Data MongoInstance `json:"data"`
}

type MongoInstance struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	VCpus              int    `json:"vCpus"`
	RAM                int    `json:"ram"`
	Disk               int    `json:"disk"`
	VolumeSize         int    `json:"volumeSize"`
	DatastoreName      string `json:"datastoreName"`
	DatastoreVersion   string `json:"datastoreVersion"`
	DatastoreVersionID string `json:"datastoreVersionId"`
	DatastoreMode      string `json:"datastoreMode"`
	DatastoreModeID    string `json:"datastoreModeId"`
	Status             string `json:"status"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	GroupConfigID      string `json:"groupConfigId"`
	Tags               []any  `json:"tags"`
}
type MongoInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type MongoBackup struct {
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
type MongoInstanceListWrapper struct {
	Data struct {
		Docs      []MongoInstance `json:"docs"`
		Page      int             `json:"page"`
		Size      int             `json:"size"`
		Total     int             `json:"total"`
		TotalPage int             `json:"totalPage"`
	} `json:"data"`
}

type mongoinstance struct {
	client *Client
}

// Get mongoinstance detail
func (v *mongoinstance) Get(id string) (MongoInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
	var obj MongoInstanceWrapper
	if err != nil {
		return MongoInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return MongoInstance{}, err
	}
	return obj.Data, err
}

func (v *mongoinstance) List(params map[string]string) ([]MongoInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
	var obj MongoInstanceListWrapper
	if err != nil {
		return []MongoInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []MongoInstance{}, err
	}
	return obj.Data.Docs, err
}
func (v *mongoinstance) ListDatastore(params map[string]string) ([]Datastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=mongodb", params)
	var obj DatastoreListWrapper

	if err != nil {
		return []Datastore{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []Datastore{}, err
	}
	return obj.Data.Docs, err
}

// Delete a mongoinstance
func (v *mongoinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *mongoinstance) SetConfigurationGroupId(id string, Mongo_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": Mongo_configuration_id,
		},
	})
}
func (v *mongoinstance) Resize(id string, flavorId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_instance",
		"requestData": map[string]interface{}{
			"newFlavorId": flavorId,
		},
	})
}

func (v *mongoinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_volume",
		"requestData": map[string]interface{}{
			"newVolumeSize": volume_size,
		},
	})
}
func (s *mongoinstance) Create(params map[string]interface{}) (MongoInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response MongoInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *mongoinstance) GetUser(id string, username string) (MongoUser, error) {
	users, err := v.ListUsers(id)
	if err != nil {
		return MongoUser{}, err
	}
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return MongoUser{}, nil
}

func (v *mongoinstance) ListUsers(id string) ([]MongoUser, error) {
	params := map[string]interface{}{
		"command": "get_list_user",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []MongoUser `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []MongoUser{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []MongoUser{}, err
	}
	return obj.Data.Docs, nil
}

// executeAction is a helper to reuse for all action-based calls.
func (v *mongoinstance) postAction(id string, action string, params map[string]interface{}) (string, error) {
	bytes, _ := json.Marshal(params)
	return v.client.Post("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}

// executeAction is a helper to reuse for all action-based calls.
func (v *mongoinstance) performAction(id string, action string, params map[string]interface{}) (ActionResponse, error) {
	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}

func (v *mongoinstance) CreateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_user",
		"body":    params,
	})
}

func (v *mongoinstance) UpdateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_user",
		"body":    params,
	})
}

func (v *mongoinstance) CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_database",
		"body":    params,
	})
}

func (v *mongoinstance) DeleteDatabase(id string, databaseName string) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "drop_database",
		"body": map[string]interface{}{
			"databaseName": databaseName,
		},
	})
}

func (v *mongoinstance) GetDatabase(id string, database string) (MongoDatabase, error) {
	databases, err := v.ListDatabases(id)
	if err != nil {
		return MongoDatabase{}, err
	}
	for _, db := range databases {
		if db.Name == database {
			return db, nil
		}
	}
	return MongoDatabase{}, nil
}
func (v *mongoinstance) ListDatabases(id string) ([]MongoDatabase, error) {
	params := map[string]interface{}{
		"command": "get_list_database",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []MongoDatabase `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []MongoDatabase{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []MongoDatabase{}, err
	}
	return obj.Data.Docs, nil
}

func (v *mongoinstance) CreateBackup(id string, name string) (ActionResponse, error) {
	bytes, _ := json.Marshal(map[string]interface{}{"name": name})
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "create_backup",
		"requestData": map[string]interface{}{
			"backupStrategyType": "S3",
			"name":               string(bytes),
		},
	})
}

func (v *mongoinstance) DeleteUser(id string, username string) (ActionResponse, error) {
	params := map[string]interface{}{
		"command": "drop_user",
		"body": map[string]interface{}{
			"username": username,
		},
	}

	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "db_action",
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
		"requestId": genUUID(),
	})
}

func (v *mongoinstance) DeleteBackup(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/backup", map[string]interface{}{"backupIds": []string{id}})
}

func (v *mongoinstance) GetBackup(id string, backup_id string) (MongoBackup, error) {
	databases, err := v.ListBackups(id)
	if err != nil {
		return MongoBackup{}, err
	}
	for _, db := range databases {
		if db.ID == backup_id {
			return db, nil
		}
	}
	return MongoBackup{}, nil
}
func (v *mongoinstance) ListBackups(id string) ([]MongoBackup, error) {
	restext, err := v.client.Get("cloudops-core/api/v1/dbaas/backup-list?page=1&size=15&datastoreCode=mongodb", map[string]string{})
	items := make([]MongoBackup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
