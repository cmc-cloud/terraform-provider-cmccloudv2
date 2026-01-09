package gocmcapiv2

import (
	"encoding/json"
)

// ImageService interface
type ImageService interface {
	Get(id string) (Image, error)
	List(params map[string]string) ([]Image, error)
}

// Image object
type Image struct {
	Architecture string `json:"architecture"`
	OsDistro     string `json:"os_distro"`
	Name         string `json:"name"`
	Visibility   string `json:"visibility"`
	Status       string `json:"status"`
	ID           string `json:"id"`
	CreatedAt    string `json:"created_at"`
	Tags         []Tag  `json:"tags"`
	Os           string `json:"os"`
	// DiskFormat   string   `json:"disk_format"`
	// Protected    bool     `json:"protected"`
	// MinDisk      int      `json:"min_disk"`
	// OsType       string   `json:"os_type"`
}

// type Images []Image

type image struct {
	client *Client
}

func (v *image) Get(id string) (Image, error) {
	jsonStr, err := v.client.Get("image/"+id, map[string]string{})
	var vpc Image
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &vpc)
	}
	return vpc, err
}

func (s *image) List(params map[string]string) ([]Image, error) {
	restext, err := s.client.Get("image", params)
	images := make([]Image, 0)
	if err != nil {
		return images, err
	}
	err = json.Unmarshal([]byte(restext), &images)
	return images, err
}
