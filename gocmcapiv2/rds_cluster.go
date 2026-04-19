package gocmcapiv2

import (
	"encoding/json"
	"net/url"
)

// RdsClusterService interface
type RdsClusterService interface {
	Get(id string) (RdsCluster, error)
	List(params map[string]string) ([]RdsCluster, error)
	Create(params map[string]interface{}) (CreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Resize(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateStorageAutoscaling(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateBackupPolicy(id string, params map[string]interface{}) (ActionResponse, error)
	GetBackup(id string, snapshotName string) (RdsClusterBackup, error)
	DeleteBackup(id string, snapshotName string) (ActionResponse, error)
	ListBackups(id string) ([]RdsClusterBackup, error)
	CreateBackup(id string, name string) (CreateResponse, error)
	ListFlavors() ([]RdsClusterFlavor, error)

	GetUser(id string, username string) (RdsClusterUser, error)
	ListUsers(id string) ([]RdsClusterUser, error)
	CreateUser(id string, params map[string]interface{}) (CreateResponse, error)
	UpdateUser(id string, username string, host string, params map[string]interface{}) (ActionResponse, error)
	DeleteUser(id string, username string, host string) (ActionResponse, error)

	GetDatabase(id string, name string) (RdsClusterDatabase, error)
	CreateDatabase(id string, params map[string]interface{}) (CreateResponse, error)
	DeleteDatabase(id string, databaseName string) (ActionResponse, error)
	ListDatabases(id string) ([]RdsClusterDatabase, error)
}

type RdsCluster struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	Mode                        string `json:"mode"`
	VolumeSize                  int    `json:"volume_size"`
	DbEngine                    string `json:"db_engine"`
	DbVersion                   string `json:"db_version"`
	SubnetID                    string `json:"subnet_id"`
	FlavorID                    string `json:"flavor_id"`
	BackupSchedule              string `json:"backup_schedule"`
	BackupRetention             int    `json:"backup_retention"`
	ClusterSize                 int    `json:"cluster_size"`
	ProxySize                   int    `json:"proxy_size"`
	Status                      string `json:"status"`
	Created                     string `json:"created"`
	RegionID                    string `json:"region_id"`
	EnableBackup                bool   `json:"enable_backup"`
	EnablePitr                  bool   `json:"enable_pitr"`
	EnableStorageAutoscaling    bool   `json:"enable_storage_autoscaling"`
	StorageAutoscalingThreshold int    `json:"storage_autoscaling_threshold"`
	StorageAutoscalingIncrement int    `json:"storage_autoscaling_increment"`
	LbVipIPAddress              string `json:"lb_vip_ip_address"`
	SubnetName                  string `json:"subnet_name"`
	VpcName                     string `json:"vpc_name"`
	VpcID                       string `json:"vpc_id"`
	BillingMode                 string `json:"billing_mode"`
	TaskState                   string `json:"task_state"`
	Tags                        []Tag  `json:"tags"`
}
type RdsClusterFlavor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Vcpus       int    `json:"vcpus"`
	RAM         int    `json:"ram"`
	Description string `json:"description"`
}
type RdsClusterCreateResponse struct {
	ID string `json:"id"`
}
type RdsClusterBackup struct {
	ID            string `json:"id"`
	CloudID       string `json:"cloud_id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	StorageStatus string `json:"storage_status"`
	Size          string `json:"size"`
	Created       string `json:"created"`
	S3Location    string `json:"s3_location"`
	EnablePitr    bool   `json:"enable_pitr"`
	PitrStartTime string `json:"pitr_start_time"`
	PitrEndTime   string `json:"pitr_end_time"`
}

type RdsClusterUser struct {
	Name      string   `json:"name"`
	Host      string   `json:"host"`
	Databases []string `json:"databases"`
}
type RdsClusterDatabase struct {
	Name string `json:"name"`
}
type rdscluster struct {
	client *Client
}

// Get rdscluster detail
func (v *rdscluster) Get(id string) (RdsCluster, error) {
	jsonStr, err := v.client.Get("rds-cluster/cluster/"+id, map[string]string{})
	var obj RdsCluster
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (v *rdscluster) List(params map[string]string) ([]RdsCluster, error) {
	restext, err := v.client.Get("rds-cluster/cluster", params)
	items := make([]RdsCluster, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *rdscluster) ListFlavors() ([]RdsClusterFlavor, error) {
	restext, err := v.client.Get("rds-cluster/flavor", map[string]string{})
	items := make([]RdsClusterFlavor, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a rdscluster
func (v *rdscluster) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("rds-cluster/cluster/" + id)
}
func (v *rdscluster) Resize(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformAction("rds-cluster/cluster/"+id+"/resize", params)
}
func (v *rdscluster) UpdateStorageAutoscaling(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("rds-cluster/cluster/"+id+"/storage_autoscaling", params)
}
func (v *rdscluster) UpdateBackupPolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("rds-cluster/cluster/"+id+"/backup_policy", params)
}

func (v *rdscluster) Create(params map[string]interface{}) (CreateResponse, error) {
	return v.client.PerformCreate("rds-cluster/cluster", params)
}

func (v *rdscluster) CreateBackup(id string, name string) (CreateResponse, error) {
	return v.client.PerformCreate("rds-cluster/cluster/"+id+"/backup", map[string]interface{}{"name": name})
}

func (v *rdscluster) DeleteBackup(id string, backupId string) (ActionResponse, error) {
	return v.client.PerformDelete("rds-cluster/cluster/" + id + "/backup/" + backupId)
}

func (v *rdscluster) GetBackup(id string, backupId string) (RdsClusterBackup, error) {
	jsonStr, err := v.client.Get("rds-cluster/cluster/"+id+"/backup/"+backupId, map[string]string{})
	var obj RdsClusterBackup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *rdscluster) ListBackups(id string) ([]RdsClusterBackup, error) {
	restext, err := v.client.Get("rds-cluster/cluster/"+id+"/backup", map[string]string{})
	items := make([]RdsClusterBackup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

func (v *rdscluster) GetUser(id string, username string) (RdsClusterUser, error) {
	users, err := v.ListUsers(id)
	if err != nil {
		return RdsClusterUser{}, err
	}
	for _, user := range users {
		if user.Name == username {
			return user, nil
		}
	}
	return RdsClusterUser{}, nil
}

func (v *rdscluster) ListUsers(id string) ([]RdsClusterUser, error) {
	restext, err := v.client.Get("rds-cluster/cluster/"+id+"/backup", map[string]string{})
	items := make([]RdsClusterUser, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

func (v *rdscluster) CreateUser(id string, params map[string]interface{}) (CreateResponse, error) {
	return v.client.PerformCreate("rds-cluster/cluster/"+id+"/user", params)
}

func (v *rdscluster) UpdateUser(id string, username string, host string, params map[string]interface{}) (ActionResponse, error) {
	// username và host trong url có thể có chứa kí tự đặc biệt
	return v.client.PerformUpdate("rds-cluster/cluster/"+id+"/user/"+username+"/"+url.PathEscape(host), params)
}

func (v *rdscluster) DeleteUser(id string, username string, host string) (ActionResponse, error) {
	return v.client.PerformDelete("rds-cluster/cluster/" + id + "/user/" + username + "/" + url.PathEscape(host))
}

func (v *rdscluster) CreateDatabase(id string, params map[string]interface{}) (CreateResponse, error) {
	return v.client.PerformCreate("rds-cluster/cluster/"+id+"/database", params)
}

func (v *rdscluster) DeleteDatabase(id string, databaseName string) (ActionResponse, error) {
	return v.client.PerformDelete("rds-cluster/cluster/" + id + "/database/" + databaseName)
}

func (v *rdscluster) GetDatabase(id string, database string) (RdsClusterDatabase, error) {
	databases, err := v.ListDatabases(id)
	if err != nil {
		return RdsClusterDatabase{}, err
	}
	for _, db := range databases {
		if db.Name == database {
			return db, nil
		}
	}
	return RdsClusterDatabase{}, nil
}
func (v *rdscluster) ListDatabases(id string) ([]RdsClusterDatabase, error) {
	restext, err := v.client.Get("rds-cluster/cluster/"+id+"/database", map[string]string{})
	items := make([]RdsClusterDatabase, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
