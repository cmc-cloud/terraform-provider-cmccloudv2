package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingGroupService interface
type AutoScalingGroupService interface {
	Get(id string) (AutoScalingGroup, error)
	List(params map[string]string) ([]AutoScalingGroup, error)
	GetAction(id string) (AutoScalingAction, error)
	Create(params map[string]interface{}) (AutoScalingGroup, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateCapacity(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type AutoScalingAction struct {
	Action       string  `json:"action"`
	Cause        string  `json:"cause"`
	CreatedAt    string  `json:"created_at"`
	EndTime      float64 `json:"end_time"`
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	StartTime    float64 `json:"start_time"`
	Status       string  `json:"status"`
	StatusReason string  `json:"status_reason"`
	Target       string  `json:"target"`
	Timeout      int     `json:"timeout"`
}

// AutoScalingGroup object
type AutoScalingGroup struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	CreatedAt       string   `json:"created_at"`
	DesiredCapacity int      `json:"desired_capacity"`
	MaxSize         int      `json:"max_size"`
	MinSize         int      `json:"min_size"`
	Nodes           []string `json:"nodes"`
	Policies        []string `json:"policies"`
	ProfileID       string   `json:"profile_id"`
	ProfileName     string   `json:"profile_name"`
	Status          string   `json:"status"`
	StatusReason    string   `json:"status_reason"`
	Timeout         int      `json:"timeout"`
	// Config          []string    `json:"config"`
	// Domain          string      `json:"domain"`
	// Metadata        []any     `json:"metadata"`
}
type autoscalinggroup struct {
	client *Client
}

// Get autoscalinggroup detail
func (s *autoscalinggroup) Get(id string) (AutoScalingGroup, error) {
	jsonStr, err := s.client.Get("as/group/"+id, map[string]string{})
	var obj AutoScalingGroup
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *autoscalinggroup) List(params map[string]string) ([]AutoScalingGroup, error) {
	restext, err := s.client.Get("as/group", params)
	items := make([]AutoScalingGroup, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	if err != nil {
		return items, err
	}
	return items, err
}

func (s *autoscalinggroup) GetAction(id string) (AutoScalingAction, error) {
	jsonStr, err := s.client.Get("as/group/action/"+id, map[string]string{})
	var obj AutoScalingAction
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

// Delete a autoscalinggroup
func (s *autoscalinggroup) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("as/group/" + id)
}
func (s *autoscalinggroup) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/group/"+id, params)
}
func (s *autoscalinggroup) UpdateCapacity(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/group/"+id+"/capacity", params)
}

// Create a new autoscalinggroup
func (s *autoscalinggroup) Create(params map[string]interface{}) (AutoScalingGroup, error) {
	jsonStr, err := s.client.Post("as/group", params)
	var response AutoScalingGroup
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
