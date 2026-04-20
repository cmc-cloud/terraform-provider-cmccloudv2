package gocmcapiv2

import (
	"encoding/json"
)

// StorageGatewayService interface
type StorageGatewayService interface {
	Get(id string) (StorageGateway, error)
	List(params map[string]string) ([]StorageGateway, error)
	Create(params map[string]interface{}) (StorageGateway, error)
	Delete(id string) (ActionResponse, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
}

// StorageGateway object
type StorageGateway struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Bucket        string `json:"bucket"`
	ProjectID     string `json:"project_id"`
	ProtocolType  string `json:"protocol_type"`
	VpcID         string `json:"vpc_id"`
	SubnetID      string `json:"subnet_id"`
	CreatedAt     string `json:"created_at"`
	Endpoint      string `json:"endpoint"`
	LastRequestID string `json:"last_request_id"`
	Tags          []Tag  `json:"tags"`
	Status        string `json:"status"`
	CommandLine   string `json:"command_line"`
	SharedPath    string `json:"shared_path"`
	BillingMode   string `json:"billing_mode"`
}
type storagegateway struct {
	client *Client
}

// Get storagegateway detail
func (v *storagegateway) Get(id string) (StorageGateway, error) {
	jsonStr, err := v.client.Get("storagegateway/s3gateway/"+id, map[string]string{})
	var obj StorageGateway
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *storagegateway) List(params map[string]string) ([]StorageGateway, error) {
	restext, err := v.client.Get("storagegateway/s3gateway", params)
	items := make([]StorageGateway, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a storagegateway
func (v *storagegateway) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("storagegateway/s3gateway/" + id)
}
func (v *storagegateway) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("storagegateway/s3gateway/"+id, params)
}

func (v *storagegateway) Create(params map[string]interface{}) (StorageGateway, error) {
	jsonStr, err := v.client.Post("storagegateway/s3gateway", params)
	var response StorageGateway
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
