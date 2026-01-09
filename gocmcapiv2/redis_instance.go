package gocmcapiv2

import (
	"encoding/json"
	"strings"
)

// RedisInstanceService interface
type RedisInstanceService interface {
	Get(id string) (RedisInstance, error)
	List(params map[string]string) ([]RedisInstance, error)
	ListDatastore(params map[string]string) ([]RedisDatastore, error)
	Create(params map[string]interface{}) (RedisInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	SetPassword(id string, password string) (ActionResponse, error)
	SetConfigurationGroupId(id string, redis_configuration_id string) (ActionResponse, error)
	AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)
	DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)

	Resize(id string, flavor_id string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error)
	UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error)

	CreateBackup(id string, params map[string]interface{}) (RedisBackup, error)
	CreateSnapshot(id string, params map[string]interface{}) (RedisSnapshot, error)
}

type RedisSnapshot struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	InstanceID string  `json:"instance_id"`
	Created    string  `json:"created"`
	Size       float64 `json:"size"`
	Status     string  `json:"status"`
	RealSize   int     `json:"real_size"`
}
type RedisBackup struct {
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

type RedisInstanceWrapper struct {
	Data RedisInstance `json:"data"`
}

type RedisInstance struct {
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
	DataDetail        RedisDataDetailFromString `json:"dataDetail"`
}

type RedisInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type RedisDataDetailFromString RedisDataDetail

func (b *RedisDataDetailFromString) UnmarshalJSON(data []byte) error {
	var val RedisDataDetail
	input := string(data)
	input = strings.Trim(input, `"`)
	input = strings.ReplaceAll(input, `\`, ``)
	if err := json.Unmarshal([]byte(input), &val); err != nil {
		Logo("RedisDataDetailFromString Unmarshal err =", err)
		return err
	}
	// Logo("AutoScalev2Config after Unmarshal = ", val)
	*b = RedisDataDetailFromString(val)
	return nil
}

type RedisDataDetail struct {
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

type RedisDatastoreWrapper struct {
	Data RedisDatastore `json:"data"`
}
type RedisInstanceListWrapper struct {
	Data struct {
		Docs      []RedisInstance `json:"docs"`
		Page      int             `json:"page"`
		Size      int             `json:"size"`
		Total     int             `json:"total"`
		TotalPage int             `json:"totalPage"`
	} `json:"data"`
}
type RedisDatastoreListWrapper struct {
	Data struct {
		Docs      []RedisDatastore `json:"docs"`
		Page      int              `json:"page"`
		Size      int              `json:"size"`
		Total     int              `json:"total"`
		TotalPage int              `json:"totalPage"`
	} `json:"data"`
}
type redisinstance struct {
	client *Client
}

// Get redisinstance detail
func (v *redisinstance) Get(id string) (RedisInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id, map[string]string{})
	var obj RedisInstanceWrapper
	if err != nil {
		return RedisInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return RedisInstance{}, err
	}
	return obj.Data, err
}

func (v *redisinstance) List(params map[string]string) ([]RedisInstance, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance", params)
	var obj RedisInstanceListWrapper
	if err != nil {
		return []RedisInstance{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []RedisInstance{}, err
	}
	return obj.Data.Docs, err
}
func (v *redisinstance) ListDatastore(params map[string]string) ([]RedisDatastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=redis", params)
	var obj RedisDatastoreListWrapper

	if err != nil {
		return []RedisDatastore{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []RedisDatastore{}, err
	}
	return obj.Data.Docs, err
}

// Delete a redisinstance
func (v *redisinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *redisinstance) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id, params)
}
func (v *redisinstance) SetPassword(id string, password string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "set_password",
		"requestData": map[string]interface{}{
			"password": password,
		},
	})
}
func (v *redisinstance) SetConfigurationGroupId(id string, redis_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": redis_configuration_id,
		},
	})
}
func (v *redisinstance) DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "detach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *redisinstance) AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "attach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *redisinstance) Resize(id string, flavor_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize", map[string]interface{}{"flavor_id": flavor_id})
}
func (v *redisinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize_volume", map[string]interface{}{"size": volume_size})
}
func (v *redisinstance) UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/accessbility", params)
}
func (v *redisinstance) UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/upgrade_datastore_version", map[string]interface{}{"datastore_version": datastore_version})
}

func (s *redisinstance) Create(params map[string]interface{}) (RedisInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response RedisInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (s *redisinstance) CreateBackup(id string, params map[string]interface{}) (RedisBackup, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance/"+id+"/backup", params)
	var response RedisBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (s *redisinstance) CreateSnapshot(id string, params map[string]interface{}) (RedisSnapshot, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance/"+id+"/snapshot", params)
	var response RedisSnapshot
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
