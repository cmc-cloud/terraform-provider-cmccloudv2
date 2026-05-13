package gocmcapiv2

import (
	"encoding/json"
)

// ImageService interface
type ImageService interface {
	Get(id string) (Image, error)
	CreateFromVolume(volumeId string, params map[string]interface{}) (Image, error)
	List(params map[string]string) ([]Image, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
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

func (v *image) List(params map[string]string) ([]Image, error) {
	restext, err := v.client.Get("image", params)
	images := make([]Image, 0)
	if err != nil {
		return images, err
	}
	err = json.Unmarshal([]byte(restext), &images)
	return images, err
}

// Create image from volume
func (v *image) CreateFromVolume(volumeId string, params map[string]interface{}) (Image, error) {
	jsonStr, err := v.client.Post("volume/"+volumeId+"/upload_to_image", params)
	var response Image
	if err != nil {
		return response, err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	return response, err
}

func (v *image) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return v.client.PerformUpdate("image/"+id, params)
}

// Delete a server
func (v *image) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("image/" + id)
}
