package gocmcapiv2

import (
	"encoding/json"
)

// KafkaInstanceService interface
type KafkaInstanceService interface {
	Get(id string) (KafkaInstance, error)
	List(params map[string]string) ([]KafkaInstance, error)
	ListDatastore() ([]Datastore, error)
	Create(params map[string]interface{}) (KafkaInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	AttachSecurityGroupId(id string, securityGroupId string) (ActionResponse, error)
	DetachSecurityGroupId(id string, securityGroupId string) (ActionResponse, error)
	Resize(id string, flavorId string) (ActionResponse, error)
	ResizeVolume(id string, volumeSize int) (ActionResponse, error)
}

type KafkaInstanceWrapper struct {
	Data KafkaInstance `json:"data"`
}

type KafkaInstance struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	DatastoreName     string `json:"datastoreName"`
	DatastoreVersion  string `json:"datastoreVersion"`
	DatastoreMode     string `json:"datastoreMode"`
	SecurityClientIds string `json:"securityClientIds"`
	VpcID             string `json:"vpcId"`
	SubnetID          string `json:"subnetId"`
	Status            string `json:"status"`
	// DataDetail        KafkaDataDetail `json:"dataDetail"`
	VpcName    string `json:"vpcName"`
	SubnetName string `json:"subnetName"`
	VolumeSize int    `json:"volumeSize"`
	FlavorInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flavorInfo"`
	Created         string `json:"created"`
	Updated         string `json:"updated"`
	QuantityOfNodes int    `json:"quantityOfNodes"`
}

type KafkaInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}

type KafkaInstanceListWrapper struct {
	Data struct {
		Docs      []KafkaInstance `json:"docs"`
		Page      int             `json:"page"`
		Size      int             `json:"size"`
		Total     int             `json:"total"`
		TotalPage int             `json:"totalPage"`
	} `json:"data"`
}
type kafkainstance struct {
	client *Client
}

// Get kafkainstance detail
func (v *kafkainstance) Get(id string) (KafkaInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{"v": "2"})
	var obj KafkaInstanceWrapper
	if err != nil {
		return KafkaInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return KafkaInstance{}, err
	}
	return obj.Data, err
}

func (v *kafkainstance) List(params map[string]string) ([]KafkaInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instances", params)
	var obj KafkaInstanceListWrapper
	if err != nil {
		return []KafkaInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []KafkaInstance{}, err
	}
	return obj.Data.Docs, err
}
func (v *kafkainstance) ListDatastore() ([]Datastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=kafka", map[string]string{})
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

// Delete a kafkainstance
func (v *kafkainstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *kafkainstance) DetachSecurityGroupId(id string, securityGroupId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "detach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": securityGroupId,
		},
	})
}
func (v *kafkainstance) AttachSecurityGroupId(id string, securityGroupId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "attach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": securityGroupId,
		},
	})
}
func (v *kafkainstance) Resize(id string, flavorId string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_instance",
		"requestData": map[string]interface{}{
			"newFlavorId": flavorId,
		},
	})
}
func (v *kafkainstance) ResizeVolume(id string, volumeSize int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "resize_volume",
		"requestData": map[string]interface{}{
			"newVolumeSize": volumeSize,
		},
	})
}

func (s *kafkainstance) Create(params map[string]interface{}) (KafkaInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response KafkaInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
