package gocmcapiv2

import (
	"encoding/json"
)

type DBv2Service interface {
	CreateUser(id string, topicId string, numPartitions int, replicationFactor int) (ActionResponse, error)
	GetUser(id string, topicId string) (DBv2User, error)
	ListUsers(id string, params map[string]string) ([]DBv2User, error)
	UpdateUser(id string, topicId string, partitions int, retentionDay int) (ActionResponse, error)
	DeleteUser(id string, topicId string) (ActionResponse, error)

	// CreateDatabase(id string, topicId string, numPartitions int, replicationFactor int) (ActionResponse, error)
	// GetDatabase(id string, topicId string) (DBv2Database, error)
	// ListDatabases(id string, params map[string]string) ([]DBv2Database, error)
	// UpdateDatabase(id string, topicId string, partitions int, retentionDay int) (ActionResponse, error)
	// DeleteDatabase(id string, topicId string) (ActionResponse, error)
}

type DBv2User struct {
	Name string `json:"name"`
}
type DBv2Database struct {
	Name string `json:"name"`
}
type dbv2 struct {
	client *Client
}

func (v *dbv2) GetUser(id string, topicId string) (DBv2User, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/users", map[string]string{})
	var obj DBActionResponse
	if err != nil {
		return DBv2User{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return DBv2User{}, err
	}
	actionID := obj.Data.ActionID
	var topics []DBv2User
	topics, err = WaitForActionResult[DBv2User](v.client, "cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, actionID, 2)
	if err != nil {
		return DBv2User{}, err
	}
	for _, t := range topics {
		if t.Name == topicId {
			return t, nil
		}
	}
	return DBv2User{}, err
}

// Get dbv2 detail
func (v *dbv2) GetAction(id string, actionID string) (string, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, map[string]string{})
	if err != nil {
		return "", err
	}
	return jsonStr, err
}

func (v *dbv2) ListUsers(id string, params map[string]string) ([]DBv2User, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/topics", params)
	var obj DBActionResponse
	if err != nil {
		return []DBv2User{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []DBv2User{}, err
	}
	actionID := obj.Data.ActionID
	var topics []DBv2User
	topics, err = WaitForActionResult[DBv2User](v.client, "cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, actionID, 2)
	if err != nil {
		return []DBv2User{}, err
	}
	return topics, err
}
func (v *dbv2) CreateUser(id string, topicId string, numPartitions int, replicationFactor int) (ActionResponse, error) {
	params := map[string]interface{}{
		"command": "create_topic",
		"body": map[string]interface{}{
			"topicName":         topicId,
			"numPartitions":     numPartitions,
			"replicationFactor": replicationFactor,
		},
	}
	return v.performAction(id, "db_action", params)
}

func (v *dbv2) DeleteUser(id string, topicId string) (ActionResponse, error) {
	return v.client.PerformDeleteWithBody("cloudops-core/api/v1/dbaas/instances", map[string]interface{}{"instanceIds": []string{id}})
}

func (v *dbv2) UpdateUser(id string, topicId string, partitions int, retentionDay int) (ActionResponse, error) {
	params := map[string]interface{}{
		"command": "edit_topic_config",
		"body": map[string]interface{}{
			"topicName":   topicId,
			"retentionMs": retentionDay * 86400000,
			"partitions":  partitions,
		},
	}
	return v.performAction(id, "db_action", params)
}
func (v *dbv2) performAction(id string, action string, params map[string]interface{}) (ActionResponse, error) {
	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     action,
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
	})
}
