package gocmcapiv2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// TagService interface
type TagService interface {
	UpdateTag(resource_id string, resource_type string, d *schema.ResourceData) (ActionResponse, error)
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type tag struct {
	client *Client
}

func (s *tag) UpdateTag(resourceId string, resourceType string, d *schema.ResourceData) (ActionResponse, error) {
	tags := d.Get("tags").(*schema.Set).List()
	if len(tags) == 0 {
		return ActionResponse{Success: true}, nil
	}
	return s.client.PerformUpdate("tag/mapping", map[string]interface{}{
		"resource_id":   resourceId,
		"resource_type": resourceType,
		"tags":          tags,
	})
}
