package gocmcapiv2

import (
	"encoding/json"
	"strings"
)

// KafkaInstanceService interface
type KafkaInstanceService interface {
	Get(id string) (KafkaInstance, error)
	List(params map[string]string) ([]KafkaInstance, error)
	ListDatastore(params map[string]string) ([]KafkaDatastore, error)
	Create(params map[string]interface{}) (KafkaInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	SetPassword(id string, password string) (ActionResponse, error)
	SetConfigurationGroupId(id string, Kafka_configuration_id string) (ActionResponse, error)
	AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)
	DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)

	Resize(id string, flavor_id string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error)
	UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error)

	CreateBackup(id string, params map[string]interface{}) (KafkaBackup, error)
}

type KafkaBackup struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	InstanceID string  `json:"instance_id"`
	ParentID   string  `json:"parent_id"`
	Created    string  `json:"created"`
	Size       float64 `json:"size"`
	Status     string  `json:"status"`
	RealSize   int     `json:"real_size"`
	RealSizeGB float64 `json:"real_size_gb"`
}

type KafkaInstanceWrapper struct {
	Data KafkaInstance `json:"data"`
}

type KafkaInstance struct {
	ID                string                    `json:"id"`
	Name              string                    `json:"name"`
	DatastoreName     string                    `json:"datastoreName"`
	DatastoreVersion  string                    `json:"datastoreVersion"`
	DatastoreMode     string                    `json:"datastoreMode"`
	GroupConfigID     string                    `json:"groupConfigId"`
	SecurityClientIds string                    `json:"securityClientIds"`
	SubnetID          string                    `json:"subnetId"`
	Status            string                    `json:"status"`
	FlavorID          string                    `json:"flavorId"`
	SubnetName        string                    `json:"subnetName"`
	FlavorName        string                    `json:"flavorName"`
	VolumeSize        int                       `json:"volumeSize"`
	Created           string                    `json:"created"`
	Updated           string                    `json:"updated"`
	DataDetail        KafkaDataDetailFromString `json:"dataDetail"`
}

type KafkaDatastore struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	VersionInfos []struct {
		ID          string `json:"id"`
		VersionName string `json:"versionName"`
		ModeInfo    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"modeInfo"`
	} `json:"versionInfos"`
}
type KafkaInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type KafkaDataDetailFromString KafkaDataDetail

func (b *KafkaDataDetailFromString) UnmarshalJSON(data []byte) error {
	var val KafkaDataDetail
	input := string(data)
	input = strings.Trim(input, `"`)
	input = strings.ReplaceAll(input, `\`, ``)
	if err := json.Unmarshal([]byte(input), &val); err != nil {
		Logo("KafkaDataDetailFromString Unmarshal err =", err)
		return err
	}
	// Logo("AutoScalev2Config after Unmarshal = ", val)
	*b = KafkaDataDetailFromString(val)
	return nil
}

type KafkaDataDetail struct {
	MasterInfo struct {
		ID                string `json:"id"`
		OsServerID        string `json:"osServerId"`
		Role              string `json:"role"`
		IPAddress         string `json:"ipAddress"`
		RAM               int    `json:"ram"`
		Disk              int    `json:"disk"`
		VolumeSize        int    `json:"volumeSize"`
		ZoneName          string `json:"zoneName"`
		Status            string `json:"status"`
		MonitorResourceID string `json:"monitorResourceId"`
		Vcpus             int    `json:"vcpus"`
	} `json:"masterInfo"`
	SlavesInfo []struct {
		ID                 string `json:"id"`
		OsServerID         string `json:"osServerId"`
		Role               string `json:"role"`
		IPAddress          string `json:"ipAddress"`
		RAM                int    `json:"ram"`
		Disk               int    `json:"disk"`
		VolumeSize         int    `json:"volumeSize"`
		ZoneName           string `json:"zoneName"`
		Status             string `json:"status"`
		MonitorResourceID  string `json:"monitorResourceId"`
		StatusAgentMonitor string `json:"statusAgentMonitor"`
		Vcpus              int    `json:"vcpus"`
	} `json:"slavesInfo"`
}

type KafkaDatastoreWrapper struct {
	Data KafkaDatastore `json:"data"`
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
type KafkaDatastoreListWrapper struct {
	Data struct {
		Docs      []KafkaDatastore `json:"docs"`
		Page      int              `json:"page"`
		Size      int              `json:"size"`
		Total     int              `json:"total"`
		TotalPage int              `json:"totalPage"`
	} `json:"data"`
}
type kafkainstance struct {
	client *Client
}

// Get kafkainstance detail
func (v *kafkainstance) Get(id string) (KafkaInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
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
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
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
func (v *kafkainstance) ListDatastore(params map[string]string) ([]KafkaDatastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=Kafka", params)
	var obj KafkaDatastoreListWrapper

	if err != nil {
		return []KafkaDatastore{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []KafkaDatastore{}, err
	}
	return obj.Data.Docs, err
}

// Delete a kafkainstance
func (v *kafkainstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *kafkainstance) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id, params)
}
func (v *kafkainstance) SetPassword(id string, password string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "set_password",
		"requestData": map[string]interface{}{
			"password": password,
		},
	})
}
func (v *kafkainstance) SetConfigurationGroupId(id string, Kafka_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": Kafka_configuration_id,
		},
	})
}
func (v *kafkainstance) DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "detach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *kafkainstance) AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "attach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *kafkainstance) Resize(id string, flavor_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize", map[string]interface{}{"flavor_id": flavor_id})
}
func (v *kafkainstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize_volume", map[string]interface{}{"size": volume_size})
}
func (v *kafkainstance) UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/accessbility", params)
}
func (v *kafkainstance) UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/upgrade_datastore_version", map[string]interface{}{"datastore_version": datastore_version})
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

func (s *kafkainstance) CreateBackup(id string, params map[string]interface{}) (KafkaBackup, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance/"+id+"/backup", params)
	var response KafkaBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
