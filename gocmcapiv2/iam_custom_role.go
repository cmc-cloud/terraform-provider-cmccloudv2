package gocmcapiv2

import (
	"encoding/json"
	"fmt"
)

// IamCustomRoleService interface
type IamCustomRoleService interface {
	Get(id string) (IamCustomRole, error)
	List(params map[string]string) ([]IamCustomRole, error)
	Create(params map[string]interface{}) (IamCustomRole, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
	AssignCustomRoleFromGroupOnProject(projectID string, groupID string, customRoleID string) (ActionResponse, error)
	UnsignCustomRoleFromGroupOnProject(projectID string, groupID string, customRoleID string) (ActionResponse, error)
}

type IamCustomRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Version     string `json:"version"`
	Content string `json:"content"`
	Created string `json:"created"`
}
type iamcustomrole struct {
	client *Client
}

// Get IamCustomRole detail
func (v *iamcustomrole) Get(id string) (IamCustomRole, error) {
	jsonStr, err := v.client.Get("iam/customrole/"+id, map[string]string{})
	var obj IamCustomRole
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *iamcustomrole) List(params map[string]string) ([]IamCustomRole, error) {
	restext, err := v.client.Get("iam/customrole", params)
	items := make([]IamCustomRole, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}
func (v *iamcustomrole) Create(params map[string]interface{}) (IamCustomRole, error) {
	jsonStr, err := v.client.Post("iam/customrole", params)
	var response IamCustomRole
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *iamcustomrole) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/customrole/"+id, params)
}

// Delete a IamCustomRole
func (v *iamcustomrole) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("iam/customrole/" + id)
}

func (v *iamcustomrole) AssignCustomRoleFromGroupOnProject(projectID string, groupID string, customRoleID string) (ActionResponse, error) {
	return v.client.PerformAction(fmt.Sprintf("iam/customrole/%s/project/%s/group/%s", customRoleID, projectID, groupID), map[string]interface{}{})
}
func (v *iamcustomrole) UnsignCustomRoleFromGroupOnProject(projectID string, groupID string, customRoleID string) (ActionResponse, error) {
	return v.client.PerformDelete(fmt.Sprintf("iam/customrole/%s/project/%s/group/%s", customRoleID, projectID, groupID))
}
