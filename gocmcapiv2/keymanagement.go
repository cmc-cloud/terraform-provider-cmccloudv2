package gocmcapiv2

import (
	"encoding/json"
)

// KeyManagementService interface
type KeyManagementService interface {
	Get(id string) (KeyManagementContainer, error)
	List(params map[string]string) ([]KeyManagementContainer, error)
	Create(params map[string]interface{}) (KeyManagementContainerCreateResponse, error)
	Delete(id string) (ActionResponse, error)

	GetSecrets(id string) ([]KeyManagementSecret, error)
	GetSecret(id string) (KeyManagementSecret, error)
	CreateSecret(params map[string]interface{}) (KeyManagementSecretCreateResponse, error)
	DeleteSecret(secret_id string) (ActionResponse, error)

	GetToken(id string) (KeyManagementToken, error)
	CreateToken(params map[string]interface{}) (KeyManagementTokenCreateResponse, error)
	RenewToken(id string, expiredDate string) (ActionResponse, error)
	DeleteToken(id string) (ActionResponse, error)
}

type KeyManagementContainerWrapper struct {
	Data KeyManagementContainer `json:"data"`
}

type KeyManagementSecretWrapper struct {
	Data KeyManagementSecret `json:"data"`
}

type KeyManagementTokenWrapper struct {
	Data KeyManagementToken `json:"data"`
}

type KeyManagementContainerListWrapper struct {
	Data struct {
		Docs      []KeyManagementContainer `json:"docs"`
		Page      int                      `json:"page"`
		Size      int                      `json:"size"`
		Total     int                      `json:"total"`
		TotalPage int                      `json:"totalPage"`
	} `json:"data"`
}

type KeyManagementSecretListWrapper struct {
	Data struct {
		Docs      []KeyManagementSecret `json:"secretRefs"`
		Page      int                   `json:"page"`
		Size      int                   `json:"size"`
		Total     int                   `json:"total"`
		TotalPage int                   `json:"totalPage"`
	} `json:"data"`
}

type KeyManagementContainer struct {
	ID           string `json:"containerUuid"`
	Name         string `json:"containerName"`
	Type         string `json:"containerType"`
	ContainerRef string `json:"containerRef"`
	Created      string `json:"created"`
}
type KeyManagementToken struct {
	TokenID        int    `json:"tokenId"`
	Token          string `json:"token"`
	ContainerNames []struct {
		ContainerName string `json:"containerName"`
	} `json:"containerNames"`
	ExpireDateTime string `json:"expireDateTime"`
	Description    string `json:"description"`
	CreatedTime    string `json:"createdTime"`
}
type KeyManagementContainerCreateResponse struct {
	Data struct {
		ID           string `json:"containerUuid"`
		ContainerRef string `json:"containerRef"`
	} `json:"data"`
}

type KeyManagementSecretCreateResponse struct {
	Data struct {
		Secrets []struct {
			ID string `json:"secretUuid"`
		} `json:"secrets"`
	} `json:"data"`
}
type KeyManagementTokenCreateResponse struct {
	Data struct {
		Token string `json:"token"`
		ID    string `json:"id"`
	} `json:"data"`
}
type KeyManagementSecret struct {
	Name       string `json:"name"`
	ID         string `json:"secretUuid"`
	SecretType string `json:"secretType"`
	// Algorithm  string `json:"algorithm"`
	// BitLength  string `json:"bitLength"`
	ExpireTime string `json:"expireTime"`
	Created    string `json:"created"`
}
type keymanagement struct {
	client *Client
}

type KeyManagementCreatedResponse struct {
	Data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

// Get keymanagement detail
func (v *keymanagement) Get(id string) (KeyManagementContainer, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/bbc-keys/containers/"+id, map[string]string{})
	var obj KeyManagementContainerWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	if err != nil {
		return KeyManagementContainer{}, err
	}
	return obj.Data, err
}

// Get keymanagement detail
func (v *keymanagement) GetToken(id string) (KeyManagementToken, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/bbc-keys/tokens/"+id, map[string]string{})
	var obj KeyManagementTokenWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	if err != nil {
		return KeyManagementToken{}, err
	}
	return obj.Data, err
}

func (s *keymanagement) List(params map[string]string) ([]KeyManagementContainer, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/bbc-keys/containers", params)
	var obj KeyManagementContainerListWrapper
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &obj)
	}
	if err != nil {
		return []KeyManagementContainer{}, err
	}
	return obj.Data.Docs, err
}
func (s *keymanagement) GetSecrets(id string) ([]KeyManagementSecret, error) {
	jsonStr, err := s.client.Get("cloudops-core/api/v1/bbc-keys/containers/"+id, map[string]string{})
	var obj KeyManagementSecretListWrapper
	if err != nil {
		var nilres []KeyManagementSecret
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return []KeyManagementSecret{}, err
	}
	return obj.Data.Docs, err
}
func (v *keymanagement) GetSecret(secret_id string) (KeyManagementSecret, error) {
	jsonStr, err := v.client.Get("cloudops-core/api/v1/bbc-keys/secrets/"+secret_id, map[string]string{})
	var response KeyManagementSecretWrapper
	if err != nil {
		var nilres KeyManagementSecret
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return KeyManagementSecret{}, err
	}
	return response.Data, nil
}

// Delete a keymanagement
func (v *keymanagement) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("cloudops-core/api/v1/bbc-keys/containers/" + id)
}
func (v *keymanagement) DeleteSecret(secret_id string) (ActionResponse, error) {
	return v.client.PerformDelete("cloudops-core/api/v1/bbc-keys/secrets/" + secret_id)
}
func (s *keymanagement) Create(params map[string]interface{}) (KeyManagementContainerCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/bbc-keys/containers", params)
	var response KeyManagementContainerCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return KeyManagementContainerCreateResponse{}, err
	}
	return response, nil
}

func (s *keymanagement) CreateSecret(params map[string]interface{}) (KeyManagementSecretCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/bbc-keys/secrets", params)
	var response KeyManagementSecretCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (s *keymanagement) CreateToken(params map[string]interface{}) (KeyManagementTokenCreateResponse, error) {
	jsonStr, err := s.client.Post("cloudops-core/api/v1/bbc-keys/tokens", params)
	var response KeyManagementTokenCreateResponse
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *keymanagement) DeleteToken(id string) (ActionResponse, error) {
	return v.client.PerformDelete("cloudops-core/api/v1/bbc-keys/tokens/" + id)
}
func (v *keymanagement) RenewToken(id string, expiredDate string) (ActionResponse, error) {
	return v.client.PerformAction("cloudops-core/api/v1/bbc-keys/tokens/"+id+"/renew", map[string]interface{}{
		"expiredDate": expiredDate,
	})
}
