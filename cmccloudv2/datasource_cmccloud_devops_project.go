package cmccloudv2

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func datasourceDevopsProjectSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"devops_project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Id of the devops project",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"owner", "share"}, false),
			ForceNew:     true,
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "search by name, match exactly",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "search by description",
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			ForceNew: true,
		},
	}
}

func datasourceDevopsProject() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceDevopsProjectRead,
		Schema: datasourceDevopsProjectSchema(),
	}
}

func dataSourceDevopsProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()

	var allDevopsProjects []gocmcapiv2.DevopsProject
	if project_id := d.Get("devops_project_id").(int); project_id != 0 {
		project, err := client.DevopsProject.Get(strconv.Itoa(project_id))
		if err != nil {
			if errors.Is(err, gocmcapiv2.ErrNotFound) {
				d.SetId("")
				return fmt.Errorf("unable to retrieve devops project [%v]: %s", project_id, err)
			}
		}
		allDevopsProjects = append(allDevopsProjects, project)
	} else {
		params := map[string]string{
			"q":    d.Get("name").(string),
			"type": d.Get("type").(string),
			"page": "1",
			"size": "1000",
		}
		projects, err := client.DevopsProject.List(params)
		if err != nil {
			return fmt.Errorf("error when get devops project %v", err)
		}
		allDevopsProjects = append(allDevopsProjects, projects...)
	}
	if len(allDevopsProjects) > 0 {
		var filteredDevopsProjects []gocmcapiv2.DevopsProject
		for _, project := range allDevopsProjects {
			if v := d.Get("name").(string); v != "" {
				if !strings.EqualFold(project.Name, v) {
					continue
				}
			}
			if v := d.Get("description").(string); v != "" {
				if !strings.Contains(strings.ToLower(project.Description), strings.ToLower(v)) {
					continue
				}
			}
			filteredDevopsProjects = append(filteredDevopsProjects, project)
		}
		allDevopsProjects = filteredDevopsProjects
	}
	if len(allDevopsProjects) < 1 {
		return fmt.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(allDevopsProjects) > 1 {
		gocmcapiv2.Logo("[DEBUG] Multiple results found: %#v", allDevopsProjects)
		return fmt.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	return dataSourceComputeDevopsProjectAttributes(d, allDevopsProjects[0])
}

func dataSourceComputeDevopsProjectAttributes(d *schema.ResourceData, project gocmcapiv2.DevopsProject) error {
	log.Printf("[DEBUG] Retrieved devops project %v: %#v", project.ID, project)
	d.SetId(strconv.Itoa(project.ID))
	_ = d.Set("name", project.Name)
	_ = d.Set("description", project.Description)
	_ = d.Set("created_at", project.CreatedAt)
	return nil
}
