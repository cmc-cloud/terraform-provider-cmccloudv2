package gocmcapiv2

import (
	"encoding/json"
)

type DBv2Service interface {
	DeleteUser(id string, userName string, host string) (ActionResponse, error)

	GetMysqlUser(id string, userName string, host string) (DBv2MysqlUser, error)
	ListMysqlUsers(id string, params map[string]string) ([]DBv2MysqlUser, error)

	GetKafkaUser(id string, userName string) (DBv2KafkaUser, error)
	ListKafkaUsers(id string, params map[string]string) ([]DBv2KafkaUser, error)

	GetPostgresUser(id string, userName string) (DBv2PostgresUser, error)
	ListPostgresUsers(id string, params map[string]string) ([]DBv2PostgresUser, error)

	GetMongoUser(id string, userName string) (DBv2MongoUser, error)
	ListMongoUsers(id string, params map[string]string) ([]DBv2MongoUser, error)

	GetDatabase(id string, userName string) (DBv2Database, error)
	ListDatabases(id string, params map[string]string) ([]DBv2Database, error)
}

type DBv2MysqlUser struct {
	Host            string `json:"allowHost"`
	UserPermissions []struct {
		DatabaseName string   `json:"databaseName"`
		Permissions  []string `json:"permissions"`
	} `json:"userPermissions"`
	Name string `json:"username"`
}
type DBv2PostgresUser struct {
	OwnedDatabases []string `json:"owned_databases"`
	Permissions    []string `json:"permissions"`
	Name           string   `json:"username"`
}
type DBv2MongoUser struct {
	Dbs struct {
		Admin []string `json:"admin"`
	} `json:"dbs"`
	Name string `json:"user"`
}
type DBv2KafkaUser struct {
	Name string `json:"username"`
}

type DBv2Database struct {
	Name string `json:"name"`
}
type dbv2 struct {
	client *Client
}

func (v *dbv2) GetMysqlUser(id string, userName string, host string) (DBv2MysqlUser, error) {
	users, err := ListUsers[DBv2MysqlUser](v, id, map[string]string{})
	if err != nil {
		return DBv2MysqlUser{}, err
	}
	for _, t := range users {
		if t.Name == userName && t.Host == host {
			return t, nil
		}
	}
	return DBv2MysqlUser{}, err
}

func (v *dbv2) ListMysqlUsers(id string, params map[string]string) ([]DBv2MysqlUser, error) {
	return ListUsers[DBv2MysqlUser](v, id, map[string]string{})
}

func (v *dbv2) GetKafkaUser(id string, userName string) (DBv2KafkaUser, error) {
	users, err := ListUsers[DBv2KafkaUser](v, id, map[string]string{})
	if err != nil {
		return DBv2KafkaUser{}, err
	}
	for _, t := range users {
		if t.Name == userName {
			return t, nil
		}
	}
	return DBv2KafkaUser{}, err
}

func (v *dbv2) ListKafkaUsers(id string, params map[string]string) ([]DBv2KafkaUser, error) {
	return ListUsers[DBv2KafkaUser](v, id, map[string]string{})
}

func (v *dbv2) GetPostgresUser(id string, userName string) (DBv2PostgresUser, error) {
	users, err := ListUsers[DBv2PostgresUser](v, id, map[string]string{})
	if err != nil {
		return DBv2PostgresUser{}, err
	}
	for _, t := range users {
		if t.Name == userName {
			return t, nil
		}
	}
	return DBv2PostgresUser{}, err
}

func (v *dbv2) ListPostgresUsers(id string, params map[string]string) ([]DBv2PostgresUser, error) {
	return ListUsers[DBv2PostgresUser](v, id, map[string]string{})
}

func (v *dbv2) GetMongoUser(id string, userName string) (DBv2MongoUser, error) {
	users, err := ListUsers[DBv2MongoUser](v, id, map[string]string{})
	if err != nil {
		return DBv2MongoUser{}, err
	}
	for _, t := range users {
		if t.Name == userName {
			return t, nil
		}
	}
	return DBv2MongoUser{}, err
}

func (v *dbv2) ListMongoUsers(id string, params map[string]string) ([]DBv2MongoUser, error) {
	return ListUsers[DBv2MongoUser](v, id, map[string]string{})
}
func ListUsers[T any](v *dbv2, id string, params map[string]string) ([]T, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/users", params)
	var obj DBActionResponse
	if err != nil {
		return []T{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []T{}, err
	}
	actionID := obj.Data.ActionID
	var users []T
	users, err = WaitForActionResult[T](v.client, "cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, actionID, 2)
	if err != nil {
		return []T{}, err
	}
	return users, err
}

func (v *dbv2) GetDatabase(id string, userName string) (DBv2Database, error) {
	databases, err := DBv2ListDatabases[DBv2Database](v, id, map[string]string{})
	if err != nil {
		return DBv2Database{}, err
	}
	for _, t := range databases {
		if t.Name == userName {
			return t, nil
		}
	}
	return DBv2Database{}, err
}

func (v *dbv2) ListDatabases(id string, params map[string]string) ([]DBv2Database, error) {
	return DBv2ListDatabases[DBv2Database](v, id, map[string]string{})
}

func DBv2ListDatabases[T any](v *dbv2, id string, params map[string]string) ([]T, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/databases", params)
	var obj DBActionResponse
	if err != nil {
		return []T{}, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []T{}, err
	}
	actionID := obj.Data.ActionID
	var users []T
	users, err = WaitForActionResult[T](v.client, "cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, actionID, 2)
	if err != nil {
		return []T{}, err
	}
	return users, err
}

// Get dbv2 action result
func (v *dbv2) GetAction(id string, actionID string) (string, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/dbaas/instance/"+id+"/actions/"+actionID, map[string]string{})
	if err != nil {
		return "", err
	}
	return jsonStr, err
}

func (v *dbv2) DeleteUser(id string, userName string, host string) (ActionResponse, error) {
	body := map[string]interface{}{
		"username": userName,
	}
	if host != "" {
		body["allowHost"] = host
	}
	params := map[string]interface{}{
		"command": "drop_user",
		"body":    body,
	}

	bytes, _ := json.Marshal(params)
	return v.client.PerformAction("cloudops-core/api/v1/dbaas/execute-action", map[string]interface{}{
		"instanceId": id,
		"action":     "db_action",
		"requestData": map[string]interface{}{
			"requestDbAction": string(bytes),
		},
		// "requestId": genUUID(),
	})
}
