package gocmcapiv2

import (
	"encoding/json"
)

// IamRoleService interface
type IamRoleService interface {
	Get(id string) (IamRole, error)
	List(params map[string]string) ([]IamRole, error)
}

// IamRole object
type IamRole struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Group       string `json:"group"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type iamrole struct {
	client *Client
}

// Get IamRole detail
func (v *iamrole) Get(id string) (IamRole, error) {
	jsonStr, err := v.client.Get("iam/role/"+id, map[string]string{})
	var obj IamRole
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (s *iamrole) List(params map[string]string) ([]IamRole, error) {
	restext, err := s.client.Get("iam/role", params)
	items := make([]IamRole, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
