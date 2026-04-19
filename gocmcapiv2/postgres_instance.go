package gocmcapiv2

import (
	"encoding/json"
)

// PostgresInstanceService interface
type PostgresInstanceService interface {
	Get(id string) (PostgresInstance, error)
	List(params map[string]string) ([]PostgresInstance, error)
	ListDatastore(params map[string]string) ([]Datastore, error)
	Create(params map[string]interface{}) (PostgresInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	SetConfigurationGroupId(id string, Postgres_configuration_id string) (ActionResponse, error)
	Resize(id string, flavorId string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	GetUser(id string, username string) (PostgresUser, error)
	ListUsers(id string) ([]PostgresUser, error)
	CreateUser(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateUser(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteUser(id string, username string) (ActionResponse, error)

	GetDatabase(id string, name string) (PostgresDatabase, error)
	CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteDatabase(id string, databaseName string) (ActionResponse, error)
	ListDatabases(id string) ([]PostgresDatabase, error)
}

type PostgresUser struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type PostgresDatabase struct {
	Name  string `json:"database"`
	Owner string `json:"owner"`
}
type PostgresUserWrapper struct {
	Data PostgresUser `json:"data"`
}

type PostgresUserListWrapper struct {
	Data struct {
		Docs []PostgresUser `json:"docs"`
	} `json:"data"`
}

type PostgresInstanceWrapper struct {
	Data PostgresInstance `json:"data"`
}

type PostgresInstance struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	DatastoreName     string `json:"datastoreName"`
	DatastoreVersion  string `json:"datastoreVersion"`
	DatastoreMode     string `json:"datastoreMode"`
	GroupConfigID     string `json:"groupConfigId"`
	SecurityClientIds string `json:"securityClientIds"`
	VpcID             string `json:"vpcId"`
	SubnetID          string `json:"subnetId"`
	Status            string `json:"status"`
	// DataDetail        string `json:"dataDetail"`
	VolumeSize int `json:"volumeSize"`
	FlavorInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flavorInfo"`
	Created         string `json:"created"`
	QuantityOfNodes int    `json:"quantityOfNodes"`
	// QuantityOfSlave int    `json:"quantityOfSlave"`
	// ProxyQuantity int    `json:"proxyQuantity"`
	// ProxyFlavorID string `json:"proxyFlavorId"`
	Connections []struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Type string `json:"type"`
	} `json:"connections"`
}
type PostgresInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type PostgresInstanceListWrapper struct {
	Data struct {
		Docs      []PostgresInstance `json:"docs"`
		Page      int                `json:"page"`
		Size      int                `json:"size"`
		Total     int                `json:"total"`
		TotalPage int                `json:"totalPage"`
	} `json:"data"`
}
type postgresinstance struct {
	client *Client
}

// Get postgresinstance detail
func (v *postgresinstance) Get(id string) (PostgresInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
	var obj PostgresInstanceWrapper
	if err != nil {
		return PostgresInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return PostgresInstance{}, err
	}
	return obj.Data, err
}

func (v *postgresinstance) List(params map[string]string) ([]PostgresInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
	var obj PostgresInstanceListWrapper
	if err != nil {
		return []PostgresInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []PostgresInstance{}, err
	}
	return obj.Data.Docs, err
}
func (v *postgresinstance) ListDatastore(params map[string]string) ([]Datastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=postgresql", params)
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

// Delete a postgresinstance
func (v *postgresinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *postgresinstance) SetConfigurationGroupId(id string, Postgres_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": Postgres_configuration_id,
		},
	})
}
func (v *postgresinstance) Resize(id string, flavorId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_instance",
		"requestData": map[string]interface{}{
			"newFlavorId": flavorId,
		},
	})
}

func (v *postgresinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_volume",
		"requestData": map[string]interface{}{
			"newVolumeSize": volume_size,
		},
	})
}
func (s *postgresinstance) Create(params map[string]interface{}) (PostgresInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response PostgresInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *postgresinstance) GetUser(id string, username string) (PostgresUser, error) {
	users, err := v.ListUsers(id)
	if err != nil {
		return PostgresUser{}, err
	}
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return PostgresUser{}, nil
}

func (v *postgresinstance) ListUsers(id string) ([]PostgresUser, error) {
	params := map[string]interface{}{
		"command": "get_list_user",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []PostgresUser `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []PostgresUser{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []PostgresUser{}, err
	}
	return obj.Data.Docs, nil
}

// executeAction is a helper to reuse for all action-based calls.
func (v *postgresinstance) postAction(id string, action string, params map[string]interface{}) (string, error) {
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
func (v *postgresinstance) performAction(id string, action string, params map[string]interface{}) (ActionResponse, error) {
	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}

func (v *postgresinstance) CreateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_user",
		"body":    params,
	})
}

func (v *postgresinstance) UpdateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_user",
		"body":    params,
	})
}

func (v *postgresinstance) CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_database",
		"body":    params,
	})
}

func (v *postgresinstance) DeleteDatabase(id string, databaseName string) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "drop_database",
		"body": map[string]interface{}{
			"databaseName": databaseName,
		},
	})
}

func (v *postgresinstance) UpdateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_database",
		"body":    params,
	})
}

func (v *postgresinstance) GetDatabase(id string, database string) (PostgresDatabase, error) {
	databases, err := v.ListDatabases(id)
	if err != nil {
		return PostgresDatabase{}, err
	}
	for _, db := range databases {
		if db.Name == database {
			return db, nil
		}
	}
	return PostgresDatabase{}, nil
}
func (v *postgresinstance) ListDatabases(id string) ([]PostgresDatabase, error) {
	params := map[string]interface{}{
		"command": "get_list_database",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []PostgresDatabase `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []PostgresDatabase{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []PostgresDatabase{}, err
	}
	return obj.Data.Docs, nil
}

func (v *postgresinstance) DeleteUser(id string, username string) (ActionResponse, error) {
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
