package gocmcapiv2

import (
	"encoding/json"
)

// KeyVaultService interface
type KeyVaultService interface {
	Get(id string) (KeyVault, error)
	List(params map[string]string) ([]KeyVault, error)
	ListDatastore(params map[string]string) ([]Datastore, error)
	Create(params map[string]interface{}) (KeyVaultCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	SetConfigurationGroupId(id string, KeyVault_configuration_id string) (ActionResponse, error)
	Resize(id string, flavorId string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	GetUser(id string, username string) (KeyVaultUser, error)
	ListUsers(id string) ([]KeyVaultUser, error)
	CreateUser(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateUser(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteUser(id string, username string) (ActionResponse, error)

	GetDatabase(id string, name string) (KeyVaultDatabase, error)
	CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateDatabase(id string, params map[string]interface{}) (ActionResponse, error)
	DeleteDatabase(id string, databaseName string) (ActionResponse, error)
	ListDatabases(id string) ([]KeyVaultDatabase, error)
}

type KeyVaultUser struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
}

type KeyVaultDatabase struct {
	Name  string `json:"database"`
	Owner string `json:"owner"`
}
type KeyVaultUserWrapper struct {
	Data KeyVaultUser `json:"data"`
}

type KeyVaultUserListWrapper struct {
	Data struct {
		Docs []KeyVaultUser `json:"docs"`
	} `json:"data"`
}

type KeyVaultWrapper struct {
	Data KeyVault `json:"data"`
}

type KeyVault struct {
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
}
type KeyVaultCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type KeyVaultListWrapper struct {
	Data struct {
		Docs      []KeyVault `json:"docs"`
		Page      int        `json:"page"`
		Size      int        `json:"size"`
		Total     int        `json:"total"`
		TotalPage int        `json:"totalPage"`
	} `json:"data"`
}
type keyvault struct {
	client *Client
}

// Get keyvault detail
func (v *keyvault) Get(id string) (KeyVault, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
	var obj KeyVaultWrapper
	if err != nil {
		return KeyVault{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return KeyVault{}, err
	}
	return obj.Data, err
}

func (v *keyvault) List(params map[string]string) ([]KeyVault, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
	var obj KeyVaultListWrapper
	if err != nil {
		return []KeyVault{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []KeyVault{}, err
	}
	return obj.Data.Docs, err
}
func (v *keyvault) ListDatastore(params map[string]string) ([]Datastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=keyvault", params)
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

// Delete a keyvault
func (v *keyvault) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *keyvault) SetConfigurationGroupId(id string, KeyVault_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": KeyVault_configuration_id,
		},
	})
}
func (v *keyvault) Resize(id string, flavorId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_instance",
		"requestData": map[string]interface{}{
			"newFlavorId": flavorId,
		},
	})
}

func (v *keyvault) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_volume",
		"requestData": map[string]interface{}{
			"newVolumeSize": volume_size,
		},
	})
}
func (s *keyvault) Create(params map[string]interface{}) (KeyVaultCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response KeyVaultCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *keyvault) GetUser(id string, username string) (KeyVaultUser, error) {
	users, err := v.ListUsers(id)
	if err != nil {
		return KeyVaultUser{}, err
	}
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}
	return KeyVaultUser{}, nil
}

func (v *keyvault) ListUsers(id string) ([]KeyVaultUser, error) {
	params := map[string]interface{}{
		"command": "get_list_user",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []KeyVaultUser `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []KeyVaultUser{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []KeyVaultUser{}, err
	}
	return obj.Data.Docs, nil
}

// executeAction is a helper to reuse for all action-based calls.
func (v *keyvault) postAction(id string, action string, params map[string]interface{}) (string, error) {
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
func (v *keyvault) performAction(id string, action string, params map[string]interface{}) (ActionResponse, error) {
	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}

func (v *keyvault) CreateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_user",
		"body":    params,
	})
}

func (v *keyvault) UpdateUser(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_user",
		"body":    params,
	})
}

func (v *keyvault) CreateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "create_database",
		"body":    params,
	})
}

func (v *keyvault) DeleteDatabase(id string, databaseName string) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "drop_database",
		"body": map[string]interface{}{
			"databaseName": databaseName,
		},
	})
}

func (v *keyvault) UpdateDatabase(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.performAction(id, "db_action", map[string]interface{}{
		"command": "update_database",
		"body":    params,
	})
}

func (v *keyvault) GetDatabase(id string, database string) (KeyVaultDatabase, error) {
	databases, err := v.ListDatabases(id)
	if err != nil {
		return KeyVaultDatabase{}, err
	}
	for _, db := range databases {
		if db.Name == database {
			return db, nil
		}
	}
	return KeyVaultDatabase{}, nil
}
func (v *keyvault) ListDatabases(id string) ([]KeyVaultDatabase, error) {
	params := map[string]interface{}{
		"command": "get_list_database",
		"body":    map[string]interface{}{},
	}
	jsonStr, err := v.postAction(id, "db_action", params)
	var obj struct {
		Data struct {
			Docs []KeyVaultDatabase `json:"docs"`
		} `json:"data"`
	}
	if err != nil {
		return []KeyVaultDatabase{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []KeyVaultDatabase{}, err
	}
	return obj.Data.Docs, nil
}

func (v *keyvault) DeleteUser(id string, username string) (ActionResponse, error) {
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
