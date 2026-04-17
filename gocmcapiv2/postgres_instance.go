package gocmcapiv2

import (
	"encoding/json"
	"strings"
)

// PostgresInstanceService interface
type PostgresInstanceService interface {
	Get(id string) (PostgresInstance, error)
	List(params map[string]string) ([]PostgresInstance, error)
	ListDatastore(params map[string]string) ([]PostgresDatastore, error)
	Create(params map[string]interface{}) (PostgresInstanceCreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	SetPassword(id string, password string) (ActionResponse, error)
	SetConfigurationGroupId(id string, Postgres_configuration_id string) (ActionResponse, error)
	AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)
	DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error)

	Resize(id string, flavor_id string) (ActionResponse, error)
	ResizeVolume(id string, volume_size int) (ActionResponse, error)
	UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error)
	UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error)

	CreateBackup(id string, params map[string]interface{}) (PostgresBackup, error)
}

type PostgresBackup struct {
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

type PostgresInstanceWrapper struct {
	Data PostgresInstance `json:"data"`
}

type PostgresInstance struct {
	ID                string                       `json:"id"`
	Name              string                       `json:"name"`
	DatastoreName     string                       `json:"datastoreName"`
	DatastoreVersion  string                       `json:"datastoreVersion"`
	DatastoreMode     string                       `json:"datastoreMode"`
	GroupConfigID     string                       `json:"groupConfigId"`
	SecurityClientIds string                       `json:"securityClientIds"`
	SubnetID          string                       `json:"subnetId"`
	Status            string                       `json:"status"`
	FlavorID          string                       `json:"flavorId"`
	SubnetName        string                       `json:"subnetName"`
	FlavorName        string                       `json:"flavorName"`
	VolumeSize        int                          `json:"volumeSize"`
	Created           string                       `json:"created"`
	Updated           string                       `json:"updated"`
	DataDetail        PostgresDataDetailFromString `json:"dataDetail"`
}

type PostgresInstanceCreateResponse struct {
	Data struct {
		InstanceID string `json:"instanceId"`
	} `json:"data"`
}
type PostgresDataDetailFromString PostgresDataDetail

func (b *PostgresDataDetailFromString) UnmarshalJSON(data []byte) error {
	var val PostgresDataDetail
	input := string(data)
	input = strings.Trim(input, `"`)
	input = strings.ReplaceAll(input, `\`, ``)
	if err := json.Unmarshal([]byte(input), &val); err != nil {
		Logo("PostgresDataDetailFromString Unmarshal err =", err)
		return err
	}
	// Logo("AutoScalev2Config after Unmarshal = ", val)
	*b = PostgresDataDetailFromString(val)
	return nil
}

type PostgresDataDetail struct {
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

type PostgresDatastoreWrapper struct {
	Data PostgresDatastore `json:"data"`
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
type PostgresDatastoreListWrapper struct {
	Data struct {
		Docs      []PostgresDatastore `json:"docs"`
		Page      int                 `json:"page"`
		Size      int                 `json:"size"`
		Total     int                 `json:"total"`
		TotalPage int                 `json:"totalPage"`
	} `json:"data"`
}
type Postgresinstance struct {
	client *Client
}

// Get Postgresinstance detail
func (v *Postgresinstance) Get(id string) (PostgresInstance, error) {
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

func (v *Postgresinstance) List(params map[string]string) ([]PostgresInstance, error) {
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
func (v *Postgresinstance) ListDatastore(params map[string]string) ([]PostgresDatastore, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/datastore?datastoreCode=Postgres", params)
	var obj PostgresDatastoreListWrapper

	if err != nil {
		return []PostgresDatastore{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)

	if err != nil {
		return []PostgresDatastore{}, err
	}
	return obj.Data.Docs, err
}

// Delete a Postgresinstance
func (v *Postgresinstance) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}
func (v *Postgresinstance) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id, params)
}
func (v *Postgresinstance) SetPassword(id string, password string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "set_password",
		"requestData": map[string]interface{}{
			"password": password,
		},
	})
}
func (v *Postgresinstance) SetConfigurationGroupId(id string, Postgres_configuration_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "change_group_config",
		"requestData": map[string]interface{}{
			"groupConfigId": Postgres_configuration_id,
		},
	})
}
func (v *Postgresinstance) DetachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "detach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *Postgresinstance) AttachSecurityGroupId(id string, security_group_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "attach_security_group",
		"requestData": map[string]interface{}{
			"securityGroupIds": security_group_id,
		},
	})
}
func (v *Postgresinstance) Resize(id string, flavor_id string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize", map[string]interface{}{"flavor_id": flavor_id})
}
func (v *Postgresinstance) ResizeVolume(id string, volume_size int) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/instance/"+id+"/resize_volume", map[string]interface{}{"size": volume_size})
}
func (v *Postgresinstance) UpdateInstanceAccessbility(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/accessbility", params)
}
func (v *Postgresinstance) UpgradeDatastoreVersion(id string, datastore_version string) (ActionResponse, error) {
	return v.client.PerformUpdate("cloudops-core/api/v1/dbaas/instance/"+id+"/upgrade_datastore_version", map[string]interface{}{"datastore_version": datastore_version})
}

func (s *Postgresinstance) Create(params map[string]interface{}) (PostgresInstanceCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance", params)
	var response PostgresInstanceCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (s *Postgresinstance) CreateBackup(id string, params map[string]interface{}) (PostgresBackup, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/dbaas/instance/"+id+"/backup", params)
	var response PostgresBackup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
