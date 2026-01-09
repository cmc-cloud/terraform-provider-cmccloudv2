package gocmcapiv2

import (
	"encoding/json"
)

// AutoScalingPolicyService interface
type AutoScalingPolicyService interface {
	Get(id string) (AutoScalingPolicy, error)
	List(params map[string]string) ([]AutoScalingPolicy, error)
	Create(params map[string]interface{}) (AutoScalingPolicy, error)
	DetachFromASGroup(policy_id string, as_group_id string) (PolicyActionResponse, error)
	AttachToASGroup(policy_id string, as_group_id string) (PolicyActionResponse, error)
	CreateHealthCheckPolicy(params map[string]interface{}) (AutoScalingPolicy, error)
	CreateDeletePolicy(params map[string]interface{}) (AutoScalingPolicy, error)
	CreateAZPolicy(params map[string]interface{}) (AutoScalingPolicy, error)
	CreateLBPolicy(params map[string]interface{}) (AutoScalingPolicy, error)
	CreateScalePolicy(params map[string]interface{}) (AutoScalingPolicy, error)

	UpdateHealthCheckPolicy(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateDeletePolicy(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateAZPolicy(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateLBPolicy(id string, params map[string]interface{}) (ActionResponse, error)
	UpdateScalePolicy(id string, params map[string]interface{}) (ActionResponse, error)

	Delete(id string) (ActionResponse, error)
}

type PolicyActionResponse struct {
	Success  bool   `json:"success"`
	ActionID string `json:"action"`
}

type AutoScalingPolicy struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Spec      struct {
		Description string `json:"description"`
		Properties  struct {
			// delete policy
			Criteria             string `json:"criteria"`
			DestroyAfterDeletion bool   `json:"destroy_after_deletion"`
			GracePeriod          int    `json:"grace_period"`
			Hooks                struct {
				Params struct {
					URL string `json:"url"`
				} `json:"params"`
				Timeout string `json:"timeout"`
				Type    string `json:"type"`
			} `json:"hooks"`
			ReduceDesiredCapacity bool `json:"reduce_desired_capacity"`

			// az policy
			Zones []struct {
				Name   string `json:"name"`
				Weight int    `json:"weight"`
			} `json:"zones"`

			// scale in/out policy
			Adjustment struct {
				BestEffort bool   `json:"best_effort"`
				Cooldown   int    `json:"cooldown"`
				MinStep    int    `json:"min_step"`
				Number     int    `json:"number"`
				Type       string `json:"type"`
			} `json:"adjustment"`
			Event string `json:"event"`

			// lb policy
			HealthMonitor struct {
				ID string `json:"id"`
			} `json:"health_monitor"`
			LbStatusTimeout int    `json:"lb_status_timeout"`
			Loadbalancer    string `json:"loadbalancer"`
			Pool            struct {
				ID           string `json:"id"`
				ProtocolPort int    `json:"protocol_port"`
				Subnet       string `json:"subnet"`
			} `json:"pool"`
			Vip struct {
				Subnet string `json:"subnet"`
			} `json:"vip"`

			// heath check policy
			Detection struct {
				DetectionModes []struct {
					Type string `json:"type"`
				} `json:"detection_modes"`
				Interval          int `json:"interval"`
				NodeUpdateTimeout int `json:"node_update_timeout"`
			} `json:"detection"`
			Recovery struct {
				Actions []struct {
					Name string `json:"name"`
				} `json:"actions"`
				NodeDeleteTimeout int  `json:"node_delete_timeout"`
				NodeForceRecreate bool `json:"node_force_recreate"`
			} `json:"recovery"`
		} `json:"properties"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"spec"`
	// Data      []string  `json:"data"`
	// Domain    string    `json:"domain"`
	// Project   string `json:"project"`
}

type autoscalingpolicy struct {
	client *Client
}

func (s *autoscalingpolicy) PerformAttachDetachAction(path string, params map[string]interface{}) (PolicyActionResponse, error) {
	jsonStr, err := s.client.Post(path, params)
	var res PolicyActionResponse
	if err != nil {
		return res, err
	}
	err = json.Unmarshal([]byte(jsonStr), &res)
	if err != nil {
		return res, err
	}
	return res, err
}

func (s *autoscalingpolicy) DetachFromASGroup(policy_id string, as_group_id string) (PolicyActionResponse, error) {
	return s.PerformAttachDetachAction("as/policy/"+policy_id+"/detach_from_as_group", map[string]interface{}{"as_group_id": as_group_id})
}
func (s *autoscalingpolicy) AttachToASGroup(policy_id string, as_group_id string) (PolicyActionResponse, error) {
	return s.PerformAttachDetachAction("as/policy/"+policy_id+"/attach_to_as_group", map[string]interface{}{"as_group_id": as_group_id})
}

func (s *autoscalingpolicy) CreateDeletePolicy(params map[string]interface{}) (AutoScalingPolicy, error) {
	return s.CreatePolicy("as/policy/delete-policy", params)
}

func (s *autoscalingpolicy) CreateAZPolicy(params map[string]interface{}) (AutoScalingPolicy, error) {
	return s.CreatePolicy("as/policy/az-policy", params)
}

func (s *autoscalingpolicy) CreateLBPolicy(params map[string]interface{}) (AutoScalingPolicy, error) {
	return s.CreatePolicy("as/policy/lb-policy", params)
}

func (s *autoscalingpolicy) CreateScalePolicy(params map[string]interface{}) (AutoScalingPolicy, error) {
	return s.CreatePolicy("as/policy/scale-policy", params)
}

func (s *autoscalingpolicy) CreateHealthCheckPolicy(params map[string]interface{}) (AutoScalingPolicy, error) {
	return s.CreatePolicy("as/policy/monitor-policy", params)
}

func (s *autoscalingpolicy) CreatePolicy(url string, params map[string]interface{}) (AutoScalingPolicy, error) {
	jsonStr, err := s.client.Post(url, params)
	var obj AutoScalingPolicy
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *autoscalingpolicy) UpdateDeletePolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/policy/"+id+"/delete-policy", params)
}

func (s *autoscalingpolicy) UpdateAZPolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/policy/"+id+"/az-policy", params)
}

func (s *autoscalingpolicy) UpdateLBPolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/policy/"+id+"/lb-policy", params)
}

func (s *autoscalingpolicy) UpdateScalePolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/policy/"+id+"/scale-policy", params)
}

func (s *autoscalingpolicy) UpdateHealthCheckPolicy(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("as/policy/"+id+"/monitor-policy", params)
}

func (s *autoscalingpolicy) Get(id string) (AutoScalingPolicy, error) {
	jsonStr, err := s.client.Get("as/policy/"+id, map[string]string{})
	var obj AutoScalingPolicy
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	return obj, err
}

func (s *autoscalingpolicy) List(params map[string]string) ([]AutoScalingPolicy, error) {
	items := make([]AutoScalingPolicy, 0)
	restext, err := s.client.Get("as/policy", params)
	if err != nil {
		return items, err
	}
	err = json.Unmarshal([]byte(restext), &items)
	if err != nil {
		return nil, err
	}
	return items, err
}

// Delete a autoscaling policy
func (s *autoscalingpolicy) Delete(id string) (ActionResponse, error) {
	return s.client.PerformDelete("as/policy/" + id)
}

// Create a new autoscaling policy
func (s *autoscalingpolicy) Create(params map[string]interface{}) (AutoScalingPolicy, error) {
	jsonStr, err := s.client.Post("as/policy", params)
	var response AutoScalingPolicy
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}
