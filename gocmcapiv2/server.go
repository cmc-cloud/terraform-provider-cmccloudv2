package gocmcapiv2

import (
	"encoding/json"
	"strconv"
)

// ServerService interface
type ServerService interface {
	Get(id string, more_info bool) (Server, error)
	List(params map[string]string) ([]Server, error)
	ListInterfaces(id string) ([]NetworkInterface, error)
	GetVolumeAttachments(id string) ([]VolumeAttachment, error)
	GetVolumeAttachmentDetail(id string, volume_id string) (VolumeAttachmentDetail, error)
	Create(params map[string]interface{}) (ServerCreatedResponse, error)
	Rename(id, newName string) (ActionResponse, error)
	Resize(id string, flavor_id string) (ActionResponse, error)
	ConfirmResize(id string) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
	ChangePassword(id string, password string) (ActionResponse, error)
	SetTags(id string, tags []interface{}) (ActionResponse, error)
	Stop(id string) (ActionResponse, error)
	Start(id string) (ActionResponse, error)
	RemoveSecurityGroup(id string, security_group_name string) (ActionResponse, error)
	AddSecurityGroup(id string, security_group_name string) (ActionResponse, error)
}

type IpAddress struct {
	IpAddress  string `json:"addr"`
	MacAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
}

// Nic object

type Nic struct {
	ID        string `json:"id"`
	NetworkID string `json:"net_id"`
	// VpcID      string `json:"vpc_id"`
	MacAddress string `json:"mac_addr"`
	FixedIps   []struct {
		SubnetID  string `json:"subnet_id"`
		IPAddress string `json:"ip_address"`
	} `json:"fixed_ips"`
	SecurityGroups []string `json:"security_groups"`
}
type VolumeAttach struct {
	ID                  string `json:"id"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}
type VolumeAttachmentDetail struct {
	AttachmentID        string `json:"attachment_id"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
	Device              string `json:"device"`
	ServerID            string `json:"serverId"`
	VolumeID            string `json:"volumeId"`
}
type VolumeAttachment struct {
	ID          string `json:"id"`
	Attachments []struct {
		AttachmentID string `json:"attachment_id"`
		Device       string `json:"device"`
		ServerID     string `json:"server_id"`
		VolumeID     string `json:"volume_id"`
	} `json:"attachments"`
	DeleteOnTermination bool `json:"delete_on_termination"`
	// Status             string `json:"status"`
	// Size               int    `json:"size"`
	// AvailabilityZone   string `json:"availability_zone"`
	// CreatedAt          string `json:"created_at"`
	// UpdatedAt          string `json:"updated_at"`
	// Name               string `json:"name"`
	// Description        any    `json:"description"`
	// VolumeType         string `json:"volume_type"`
	// SnapshotID         any    `json:"snapshot_id"`
	// SourceVolid        any    `json:"source_volid"`
	// Metadata           []any  `json:"metadata"`
	// UserID             string `json:"user_id"`
	// Bootable           string `json:"bootable"`
	// Encrypted          bool   `json:"encrypted"`
}
type ServerCreatedResponse struct {
	Server struct {
		ID        string `json:"id"`
		AdminPass string `json:"adminPass"`
	} `json:"server"`
	Success bool `json:"success"`
}
type ServerSecurityGroup struct {
	Name string `json:"name"`
}

// Server object
type Server struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"` // BUILD, ACTIVE, SHUTOFF, VERIFY_RESIZE, ERROR
	TaskState  any    `json:"OS-EXT-STS:task_state"`
	VMState    string `json:"OS-EXT-STS:vm_state"` // building, active, stopped, resized, error,
	PowerState int    `json:"OS-EXT-STS:power_state"`
	Flavor     struct {
		CPU          int    `json:"vcpus"`
		RAM          int    `json:"ram"`
		Disk         int    `json:"disk"`
		ID           string `json:"id"`
		OriginalName string `json:"original_name"`
	} `json:"flavor"`
	Created          string                `json:"created"`
	Addresses        interface{}           `json:"addresses"`
	AvailabilityZone string                `json:"OS-EXT-AZ:availability_zone"`
	KeyName          string                `json:"key_name"`
	SecurityGroups   []ServerSecurityGroup `json:"security_groups"`
	VolumesAttached  []VolumeAttach        `json:"os-extended-volumes:volumes_attached"`
	ServerGroups     []string              `json:"server_groups"`
	BillingMode      string                `json:"billing_mode"`
	Volumes          []Volume              `json:"volumes"`
	Description      string                `json:"description"`
	Tags             []Tag                 `json:"tags"`
	Nics             []Nic                 `json:"nics"`
	// Metadata   []any  `json:"metadata"`
	// Image      string `json:"image"`
	// DisplayText  string   `json:"display_text"`
	// Locked           bool                  `json:"locked"`
	// LockedReason     any                   `json:"locked_reason"`
}

type server struct {
	client *Client
}

// Get server detail
func (s *server) Get(id string, more_info bool) (Server, error) {
	jsonStr, err := s.client.Get("server/"+id, map[string]string{"more_info": strconv.FormatBool(more_info)})
	var obj Server
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *server) List(params map[string]string) ([]Server, error) {
	restext, err := s.client.Get("server", params)
	items := make([]Server, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

func (s *server) ListInterfaces(id string) ([]NetworkInterface, error) {
	restext, err := s.client.Get("server/"+id+"/interface", map[string]string{})
	items := make([]NetworkInterface, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Get server detail
func (s *server) GetVolumeAttachments(id string) ([]VolumeAttachment, error) {
	jsonStr, err := s.client.Get("server/"+id+"/volume_attachment", map[string]string{})
	var obj []VolumeAttachment
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *server) GetVolumeAttachmentDetail(id string, volume_id string) (VolumeAttachmentDetail, error) {
	jsonStr, err := s.client.Get("server/"+id+"/volume_attachment/"+volume_id, map[string]string{})
	var obj VolumeAttachmentDetail
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *server) Rename(id, newName string) (ActionResponse, error) {
	return s.client.PerformUpdate("server/"+id, map[string]interface{}{"name": newName})
}

func (s *server) SetTags(id string, tags []interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("server/"+id+"/tags", map[string]interface{}{"tags": tags})
}

// Delete a server
func (s *server) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("server/" + id)
}
func (s *server) Resize(id string, flavor_id string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/resize", map[string]interface{}{"flavor_id": flavor_id})
}
func (s *server) ConfirmResize(id string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/confirm_resize", map[string]interface{}{})
}
func (s *server) Stop(id string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/stop", map[string]interface{}{})
}
func (s *server) Start(id string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/start", map[string]interface{}{})
}
func (s *server) RemoveSecurityGroup(id string, security_group_name string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/detach_security_group", map[string]interface{}{"security_group_name": security_group_name})
}
func (s *server) AddSecurityGroup(id string, security_group_name string) (ActionResponse, error) {
	return s.client.PerformAction("server/"+id+"/attach_security_group", map[string]interface{}{"security_group_name": security_group_name})
}
func (s *server) ChangePassword(id string, password string) (ActionResponse, error) {
	jsonStr, err := s.client.Post("server/"+id+"/change_pass", map[string]interface{}{"password": password})
	var res ActionResponse
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(jsonStr), &res)
	return res, err
}

// Create a new server
func (s *server) Create(params map[string]interface{}) (ServerCreatedResponse, error) {
	jsonStr, err := s.client.Post("server", params)
	var response ServerCreatedResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
