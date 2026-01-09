package gocmcapiv2

import (
	"encoding/json"
)

// EcsGroupService interface
type EcsGroupService interface {
	Get(id string) (EcsGroup, error)
	List(params map[string]string) ([]EcsGroup, error)
	Create(params map[string]interface{}) (EcsGroup, error)
	Delete(id string) (ActionResponse, error)
}

// EcsGroup object
type EcsGroup struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Policy string `json:"policy"`
}
type ecsgroup struct {
	client *Client
}

// Get EcsGroup detail
func (v *ecsgroup) Get(id string) (EcsGroup, error) {
	jsonStr, err := v.client.Get("server/ecs-group/"+id, map[string]string{})
	var ecsgroup EcsGroup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &ecsgroup)
	}
	return ecsgroup, err
}
func (s *ecsgroup) List(params map[string]string) ([]EcsGroup, error) {
	restext, err := s.client.Get("server/ecs-group", params)
	items := make([]EcsGroup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a EcsGroup
func (v *ecsgroup) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("server/ecs-group/" + id)
}
func (v *ecsgroup) Create(params map[string]interface{}) (EcsGroup, error) {
	jsonStr, err := v.client.Post("server/ecs-group", params)
	var response EcsGroup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
