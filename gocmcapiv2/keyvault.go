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
func (s *keyvault) Create(params map[string]interface{}) (KeyVaultCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response KeyVaultCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
