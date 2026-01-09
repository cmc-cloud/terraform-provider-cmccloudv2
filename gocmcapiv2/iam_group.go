package gocmcapiv2

import (
	"encoding/json"
)

// IamGroupService interface
type IamGroupService interface {
	Get(id string) (IamGroup, error)
	GetGroupOfProject(project_id string, group_name string) (IamGroup, error)
	List(params map[string]string) ([]IamGroup, error)
	Create(params map[string]interface{}) (IamGroup, error)
	Delete(group_name string) (ActionResponse, error)
	Update(group_name string, params map[string]interface{}) (ActionResponse, error)
	AddUserToGroup(userID string, groupID string) (string, error)
	RemoveUserFromGroup(userID string, groupID string) (ActionResponse, error)
}

// IamGroup object
type IamGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type iamgroup struct {
	client *Client
}

// Get IamGroup detail
func (v *iamgroup) Get(id string) (IamGroup, error) {
	// jsonStr, err := v.client.Get("iam/group", map[string]string{"group_name": group_name})
	jsonStr, err := v.client.Get("iam/group/"+id, map[string]string{})
	var obj IamGroup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *iamgroup) GetGroupOfProject(project_id string, group_name string) (IamGroup, error) {
	jsonStr, err := v.client.Get("iam/group", map[string]string{"project_id": project_id, "group_name": group_name})
	var obj IamGroup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}
func (v *iamgroup) List(params map[string]string) ([]IamGroup, error) {
	restext, err := v.client.Get("iam/group", params)
	items := make([]IamGroup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

func (v *iamgroup) Create(params map[string]interface{}) (IamGroup, error) {
	jsonStr, err := v.client.Post("iam/group", params)
	var response IamGroup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

// Delete a IamGroup
func (v *iamgroup) Delete(group_name string) (ActionResponse, error) {
	return v.client.PerformDelete("iam/group/" + group_name)
}
func (v *iamgroup) Update(group_name string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("iam/group/"+group_name, params)
}

func (v *iamgroup) AddUserToGroup(userID string, groupID string) (string, error) {
	params := map[string]interface{}{}
	return v.client.Post("iam/group/"+groupID+"/user/"+userID, params)
}

func (v *iamgroup) RemoveUserFromGroup(userID string, groupID string) (ActionResponse, error) {
	return v.client.PerformDelete("iam/group/" + groupID + "/user/" + userID)
}
