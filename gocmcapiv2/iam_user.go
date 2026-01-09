package gocmcapiv2

import (
	"encoding/json"
)

// IamUserService interface
type IamUserService interface {
	Get(id string) (IamUser, error)
	GetServerPermission(username string) ([]UserServerPermission, error)
	List(params map[string]string) ([]IamUser, error)
	Create(params map[string]interface{}) (IamUser, error)
	Update(username string, params map[string]interface{}) (ActionResponse, error)
	UpdateEmail(username string, email string) (ActionResponse, error)
	SetPassword(username string, params map[string]interface{}) (ActionResponse, error)
	SetServerPermission(username string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
	Enable(id string) (ActionResponse, error)
	Disable(id string) (ActionResponse, error)
}

// IamUser object
type IamUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Enabled   bool   `json:"enabled"`
	Totp      bool   `json:"totp"`
	// APIKeyEnabled    bool   `json:"api_key_enabled"`
	ShortName string `json:"short_name"`
	// S3RootEnabled    bool   `json:"s3_root_enabled"`
	// S3Enabled        bool   `json:"s3_enabled"`
	// S3Access         string `json:"s3_access"`
	// S3Endpoint       string `json:"s3_endpoint"`
	// S3ArnID          string `json:"s3_arn_id"`
	// S3SecretKey      string `json:"s3_secret_key"`
	// S3Key            string `json:"s3_key"`
}
type UserServerPermission struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	DisplayName string `json:"display_name"`
	ServerID    string `json:"server_id"`
	Created     string `json:"created"`
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
	Blocked     bool   `json:"blocked"`
	AllowView   bool   `json:"allow_view"`
	AllowEdit   bool   `json:"allow_edit"`
	AllowCreate bool   `json:"allow_create"`
	AllowDelete bool   `json:"allow_delete"`
	RegionID    string `json:"region_id"`
}
type iamuser struct {
	client *Client
}

// Get IamUser detail
func (v *iamuser) Get(id string) (IamUser, error) {
	jsonStr, err := v.client.Get("iam/user/"+id, map[string]string{})
	var obj IamUser
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *iamuser) GetServerPermission(username string) ([]UserServerPermission, error) {
	restext, err := v.client.Get("iam/user/"+username+"/server_permissions", map[string]string{})
	items := make([]UserServerPermission, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *iamuser) List(params map[string]string) ([]IamUser, error) {
	restext, err := v.client.Get("iam/user", params)
	items := make([]IamUser, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *iamuser) Create(params map[string]interface{}) (IamUser, error) {
	jsonStr, err := v.client.Post("iam/user", params)
	var response IamUser
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
func (v *iamuser) Update(username string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/user/"+username, params)
}
func (v *iamuser) UpdateEmail(username string, email string) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/user/"+username+"/email", map[string]interface{}{"email": email})
}

func (v *iamuser) SetPassword(username string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/user/"+username, params)
}

func (v *iamuser) SetServerPermission(username string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/user/"+username+"/set_server_permissions", params)
}

// Delete a IamUser
func (v *iamuser) Delete(username string) (ActionResponse, error) {
	return v.client.PerformDelete("iam/user/" + username)
}

func (v *iamuser) Enable(username string) (ActionResponse, error) {
	return v.client.PerformAction("iam/user/"+username+"/set_user_status", map[string]interface{}{"enabled": true})
}

func (v *iamuser) Disable(username string) (ActionResponse, error) {
	return v.client.PerformAction("iam/user/"+username+"/set_user_status", map[string]interface{}{"enabled": false})
}
