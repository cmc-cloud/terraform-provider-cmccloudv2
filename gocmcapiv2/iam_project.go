package gocmcapiv2

import (
	"encoding/json"
	"fmt"
)

// IamProjectService interface
type IamProjectService interface {
	Get(id string) (IamProject, error)
	List(params map[string]string) ([]IamProject, error)
	Create(params map[string]interface{}) (IamProject, error)
	Delete(id string) (ActionResponse, error)
	AssignRoleFromGroupOnProject(projectID string, groupID string, roleID string) (ActionResponse, error)
	UnsignRoleFromGroupOnProject(projectID string, groupID string, roleID string) (ActionResponse, error)
}

// IamProject object
type IamProject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	DomainId string `json:"domain_id"`
	RegionId string `json:"region_id"`
}
type iamproject struct {
	client *Client
}

// Get IamProject detail
func (v *iamproject) Get(id string) (IamProject, error) {
	jsonStr, err := v.client.Get("iam/project/"+id, map[string]string{})
	var obj IamProject
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *iamproject) List(params map[string]string) ([]IamProject, error) {
	restext, err := v.client.Get("iam/project", params)
	items := make([]IamProject, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a IamProject
func (v *iamproject) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("iam/project/" + id)
}
func (v *iamproject) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/project/"+id, params)
}
func (v *iamproject) AssignRoleFromGroupOnProject(projectID string, groupID string, roleID string) (ActionResponse, error) {
	return v.client.PerformAction(fmt.Sprintf("iam/project/%s/group/%s/role/%s/assign_role_from_group_on_project", projectID, groupID, roleID), map[string]interface{}{})
}
func (v *iamproject) UnsignRoleFromGroupOnProject(projectID string, groupID string, roleID string) (ActionResponse, error) {
	return v.client.PerformDelete(fmt.Sprintf("iam/project/%s/group/%s/role/%s/unassign_role_from_group_on_project", projectID, groupID, roleID))
}
func (v *iamproject) Create(params map[string]interface{}) (IamProject, error) {
	jsonStr, err := v.client.Post("iam/project", params)
	var response IamProject
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
