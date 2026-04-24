package gocmcapiv2

import (
	"encoding/json"
)

// OpenSearchService interface
type OpenSearchService interface {
	Get(id string) (OpenSearch, error)
	List(params map[string]string) ([]OpenSearch, error)
	Create(params map[string]interface{}) (CreateResponse, error)
	Delete(id string) (ActionResponse, error)
	Resize(id string, params map[string]interface{}) (ActionResponse, error)
	ChangePassword(id string, password string) (ActionResponse, error)
	UpdateStorageAutoscaling(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateSnapshotPolicy(id string, params map[string]interface{}) (ActionResponse, error)
	GetSnapshot(id string, snapshotName string) (OpenSearchSnapshot, error)
	DeleteSnapshot(id string, snapshotName string) (ActionResponse, error)
	ListSnapshots(id string) ([]OpenSearchSnapshot, error)
	CreateSnapshot(id string, name string) (CreateResponse, error)
	ListFlavors() ([]OpenSearchFlavor, error)
	ListDashboardFlavors() ([]OpenSearchFlavor, error)
}

type OpenSearch struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	APILbURL                    any    `json:"api_lb_url"`
	DashboardLbURL              any    `json:"dashboard_lb_url"`
	VolumeSize                  int    `json:"volume_size"`
	Version                     string `json:"version"`
	NodeCount                   int    `json:"node_count"`
	MasterCount                 int    `json:"master_count"`
	EnableDrainNodes            bool   `json:"enable_drain_nodes"`
	LbSubnetID                  string `json:"lb_subnet_id"`
	EnableLbInternal            bool   `json:"enable_lb_internal"`
	FlavorID                    string `json:"flavor_id"`
	DashboardFlavorID           string `json:"dashboard_flavor_id"`
	EnableSnapshot              bool   `json:"enable_snapshot"`
	SnapshotCreationCron        string `json:"snapshot_creation_cron"`
	SnapshotDeletionCron        string `json:"snapshot_deletion_cron"`
	SnapshotTimezone            string `json:"snapshot_timezone"`
	DashboardReplicas           int    `json:"dashboard_replicas"`
	RentationMaxAge             int    `json:"rentation_max_age"`
	RentationMinCount           int    `json:"rentation_min_count"`
	RentationMaxCount           int    `json:"rentation_max_count"`
	EnableStorageAutoscaling    bool   `json:"enable_storage_autoscaling"`
	StorageAutoscalingThreshold int    `json:"storage_autoscaling_threshold"`
	StorageAutoscalingIncrement int    `json:"storage_autoscaling_increment"`
	StorageAutoscalingMax       int    `json:"storage_autoscaling_max"`
	ApiDomain                   string `json:"api_domain"`
	DashboardDomain             string `json:"dashboard_domain"`
	Status                      string `json:"status"`
	Created                     string `json:"created"`
	VpcID                       string `json:"vpc_id"`
	Tags                        []Tag  `json:"tags"`
	BillingMode                 string `json:"billing_mode"`
	TaskState                   string `json:"task_state"`
}
type OpenSearchFlavor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Vcpus       int    `json:"vcpus"`
	RAM         int    `json:"ram"`
	Description string `json:"description"`
}

type OpenSearchSnapshot struct {
	Name         string   `json:"name"`
	Status       string   `json:"status"`
	StartTime    string   `json:"start_time"`
	EndTime      string   `json:"end_time"`
	DurationMs   int      `json:"duration_ms"`
	Indices      []string `json:"indices"`
	IndicesCount int      `json:"indices_count"`
	Id           string   `json:"snapshot_id"`
	// ID           string   `json:"id"`
}

type opensearch struct {
	client *Client
}

// Get opensearch detail
func (v *opensearch) Get(id string) (OpenSearch, error) {
	jsonStr, err := v.client.Get("opensearch/cluster/"+id, map[string]string{})
	var obj OpenSearch
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (v *opensearch) List(params map[string]string) ([]OpenSearch, error) {
	restext, err := v.client.Get("opensearch/cluster", params)
	items := make([]OpenSearch, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *opensearch) ListFlavors() ([]OpenSearchFlavor, error) {
	restext, err := v.client.Get("opensearch/flavor", map[string]string{})
	items := make([]OpenSearchFlavor, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *opensearch) ListDashboardFlavors() ([]OpenSearchFlavor, error) {
	restext, err := v.client.Get("opensearch/dashboard_flavor", map[string]string{})
	items := make([]OpenSearchFlavor, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a opensearch
func (v *opensearch) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("opensearch/cluster/" + id)
}
func (v *opensearch) Resize(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformAction("opensearch/cluster/"+id+"/resize", params)
}
func (v *opensearch) ChangePassword(id string, password string) (ActionResponse, error) {
	return v.client.PerformAction("opensearch/cluster/"+id+"/change_password", map[string]interface{}{"password": password})
}
func (v *opensearch) UpdateStorageAutoscaling(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("opensearch/cluster/"+id+"/storage_autoscaling", params)
}
func (v *opensearch) UpdateSnapshotPolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("opensearch/cluster/"+id+"/snapshot_policy", params)
}

func (v *opensearch) Create(params map[string]interface{}) (CreateResponse, error) {
	return v.client.PerformCreate("opensearch/cluster", params)
}

func (v *opensearch) CreateSnapshot(id string, name string) (CreateResponse, error) {
	return v.client.PerformCreate("opensearch/cluster/"+id+"/snapshot", map[string]interface{}{"name": name})
}

func (v *opensearch) DeleteSnapshot(id string, snapshotName string) (ActionResponse, error) {
	return v.client.PerformDelete("opensearch/cluster/" + id + "/snapshot/" + snapshotName)
}

func (v *opensearch) GetSnapshot(id string, snapshotNameOrId string) (OpenSearchSnapshot, error) {
	jsonStr, err := v.client.Get("opensearch/cluster/"+id+"/snapshot/"+snapshotNameOrId, map[string]string{})
	var obj OpenSearchSnapshot
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *opensearch) ListSnapshots(id string) ([]OpenSearchSnapshot, error) {
	restext, err := v.client.Get("opensearch/cluster/"+id+"/snapshot", map[string]string{})
	items := make([]OpenSearchSnapshot, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
