package gocmcapiv2

import (
	"encoding/json"
)

// MysqlInstanceService interface
type MysqlInstanceService interface {
	Get(id string) (MysqlInstance, error)
	List(params map[string]string) ([]MysqlInstance, error)
	ListDatastore(params map[string]string) ([]Datastore, error)
	Create(params map[string]interface{}) (MysqlInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	SetConfigurationGroupId(id string, Mysql_configuration_id string) (ActionResponse, error)
	Resize(id string, flavorId string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	GetUser(id string, username string) (MysqlUser, error)
	ListUsers(id string) ([]MysqlUser, error)
	CreateUser(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateUser(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteUser(id string, username string) (ActionResponse, error)

	GetDatabase(id string, name string) (MysqlDatabase, error)
	CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteDatabase(id string, databaseName string) (ActionResponse, error)
	ListDatabases(id string) ([]MysqlDatabase, error)
}

type MysqlUser struct {
	Name      string   `json:"name"`
	Host      string   `json:"host"`
	Databases []string `json:"databases"`
}

type MysqlDatabase struct {
	Name string `json:"databaseName"`
}
type MysqlUserWrapper struct {
	Data MysqlUser `json:"data"`
}

type MysqlUserListWrapper struct {
	Data struct {
		Docs []MysqlUser `json:"docs"`
	} `json:"data"`
}

type MysqlInstanceWrapper struct {
	Data MysqlInstance `json:"data"`
}

type MysqlInstance struct {
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
	Tags               []Tag  `json:"tags"`
}
type MysqlInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}

type MysqlInstanceListWrapper struct {
	Data struct {
		Docs      []MysqlInstance `json:"docs"`
		Page      int             `json:"page"`
		Size      int             `json:"size"`
		Total     int             `json:"total"`
		TotalPage int             `json:"totalPage"`
	} `json:"data"`
}

type mysqlinstance struct {
	client *Client
}

// Get mysqlinstance detail
func (v *mysqlinstance) Get(id string) (MysqlInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
	var obj MysqlInstanceWrapper
	if err != nil {
		return MysqlInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return MysqlInstance{}, err
	}
	return obj.Data, err
}

func (v *mysqlinstance) List(params map[string]string) ([]MysqlInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
	var obj MysqlInstanceListWrapper
	if err != nil {
		return []MysqlInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []MysqlInstance{}, err
	}
	return obj.Data.Docs, err
}
func (v *mysqlinstance) ListDatastore(params map[string]string) ([]Datastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=mysql", params)
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

// Delete a mysqlinstance
func (v *mysqlinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *mysqlinstance) SetConfigurationGroupId(id string, Mysql_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": Mysql_configuration_id,
		},
	})
}
func (v *mysqlinstance) Resize(id string, flavorId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_instance",
		"requestData": map[string]interface{}{
			"newFlavorId": flavorId,
		},
	})
}

func (v *mysqlinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_volume",
		"requestData": map[string]interface{}{
			"newVolumeSize": volume_size,
		},
	})
}
func (s *mysqlinstance) Create(params map[string]interface{}) (MysqlInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response MysqlInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *mysqlinstance) GetUser(id string, username string) (MysqlUser, error) {
	users, err := v.ListUsers(id)
	if err != nil {
		return MysqlUser{}, err
	}
	for _, user := range users {
		if user.Name == username {
			return user, nil
		}
	}
	return MysqlUser{}, nil
}

func (v *mysqlinstance) ListUsers(id string) ([]MysqlUser, error) {
	params := map[string]interface{}{
		"command": "get_list_user",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []MysqlUser `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []MysqlUser{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []MysqlUser{}, err
	}
	return obj.Data.Docs, nil
}

// executeAction is a helper to reuse for all action-based calls.
func (v *mysqlinstance) postAction(id string, action string, params map[string]interface{}) (string, error) {
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
func (v *mysqlinstance) performAction(id string, action string, params map[string]interface{}) (ActionResponse, error) {
	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}

func (v *mysqlinstance) CreateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_user",
		"body":    params,
	})
}

func (v *mysqlinstance) UpdateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_user",
		"body":    params,
	})
}

func (v *mysqlinstance) CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_database",
		"body":    params,
	})
}

func (v *mysqlinstance) DeleteDatabase(id string, databaseName string) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "drop_database",
		"body": map[string]interface{}{
			"databaseName": databaseName,
		},
	})
}

func (v *mysqlinstance) GetDatabase(id string, database string) (MysqlDatabase, error) {
	databases, err := v.ListDatabases(id)
	if err != nil {
		return MysqlDatabase{}, err
	}
	for _, db := range databases {
		if db.Name == database {
			return db, nil
		}
	}
	return MysqlDatabase{}, nil
}
func (v *mysqlinstance) ListDatabases(id string) ([]MysqlDatabase, error) {
	params := map[string]interface{}{
		"command": "get_list_database",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []MysqlDatabase `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []MysqlDatabase{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []MysqlDatabase{}, err
	}
	return obj.Data.Docs, nil
}

func (v *mysqlinstance) DeleteUser(id string, username string) (ActionResponse, error) {
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
