package gocmcapiv2

import (
	"encoding/json"
)

// FlavorService interface
type FlavorService interface {
	Get(id string) (Flavor, error)
	List(params map[string]string) ([]Flavor, error)
}

// Flavor object
type Flavor struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	RAM                    int    `json:"ram"`
	Disk                   int    `json:"disk"`
	Vcpus                  int    `json:"vcpus"`
	OsFlavorAccessIsPublic bool   `json:"os-flavor-access:is_public"`
	Description            string `json:"description"`
	ExtraSpecs             struct {
		IsDatabaseFlavor BoolFromString `json:"aggregate_instance_extra_specs:database"`
		IsK8sFlavor      BoolFromString `json:"aggregate_instance_extra_specs:k8s"`
	} `json:"extra_specs"`
}

// type Flavors []Flavor

type flavor struct {
	client *Client
}

func (v *flavor) Get(id string) (Flavor, error) {
	jsonStr, err := v.client.Get("server/flavor/"+id, map[string]string{})
	var vpc Flavor
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vpc)
	}
	return vpc, err
}

// List lists all flavors. []*Flavor
func (s *flavor) List(params map[string]string) ([]Flavor, error) {
	restext, err := s.client.Get("server/flavors", params)
	flavors := make([]Flavor, 0)
	if err != nil {
		return flavors, err
	}
	err = json.Unmarshal([]byte(restext), &flavors)
	return flavors, err
}
