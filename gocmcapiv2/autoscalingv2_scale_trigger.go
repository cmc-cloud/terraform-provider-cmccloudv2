package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingV2ScaleTriggerService interface
type AutoScalingV2ScaleTriggerService interface {
	Get(group_id string, id string) (AutoScalingV2ScaleTrigger, error)
	List(group_id string, params map[string]string) ([]AutoScalingV2ScaleTrigger, error)
	Create(group_id string, params map[string]interface{}) (AutoScalingV2ScaleTrigger, error)
	Delete(group_id string, id string) (ActionResponse, error)
	Update(group_id string, id string, params map[string]interface{}) (ActionResponse, error)
}

type AutoScalingV2ScaleTrigger struct {
	Annotations struct {
		Description string `json:"description"`
	} `json:"annotations"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Function    string `json:"function"`
	Metric      string `json:"metric"`
	Comparator  string `json:"comparator"`
	Threadhold  int    `json:"threadhold"`
	Interval    int    `json:"interval"`
	Count       int    `json:"count"`
	Action      string `json:"action"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type asv2scaletrigger struct {
	client *Client
}

// Get asv2scaletrigger detail
func (s *asv2scaletrigger) Get(group_id string, id string) (AutoScalingV2ScaleTrigger, error) {
	jsonStr, err := s.client.Get("asv2/group/"+group_id+"/trigger/"+id, map[string]string{"detail": "true"})
	var obj AutoScalingV2ScaleTrigger
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *asv2scaletrigger) List(group_id string, params map[string]string) ([]AutoScalingV2ScaleTrigger, error) {
	restext, err := s.client.Get("asv2/group/"+group_id+"/trigger", params)
	items := make([]AutoScalingV2ScaleTrigger, 0)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	return items, err
}

// Delete a asv2scaletrigger
func (s *asv2scaletrigger) Delete(group_id string, id string) (ActionResponse, error) {
	return s.client.PerformDelete("asv2/group/" + group_id + "/trigger/" + id)
}

func (s *asv2scaletrigger) Update(group_id string, id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("asv2/group/"+group_id+"/trigger/"+id, params)
}

// Create a new asv2scaletrigger
func (s *asv2scaletrigger) Create(group_id string, params map[string]interface{}) (AutoScalingV2ScaleTrigger, error) {
	jsonStr, err := s.client.Post("asv2/group/"+group_id+"/trigger", params)
	var response AutoScalingV2ScaleTrigger
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
