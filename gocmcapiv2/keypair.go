package gocmcapiv2

import (
	"encoding/json"
)

// KeypairService interface
type KeypairService interface {
	Get(id string) (Keypair, error)
	List(params map[string]string) ([]Keypair, error)
}

// Keypair object
type Keypair struct {
	// Keypair struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	// } `json:"keypair"`
}

// type Keypairs []Keypair

type keypair struct {
	client *Client
}

func (v *keypair) Get(name string) (Keypair, error) {
	jsonStr, err := v.client.Get("keypair/"+name, map[string]string{})
	var vpc Keypair
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vpc)
	}
	return vpc, err
}

func (s *keypair) List(params map[string]string) ([]Keypair, error) {
	restext, err := s.client.Get("keypair", params)
	keypairs := make([]Keypair, 0)
	if err != nil {
		return keypairs, err
	}
	err = json.Unmarshal([]byte(restext), &keypairs)
	return keypairs, err
}
