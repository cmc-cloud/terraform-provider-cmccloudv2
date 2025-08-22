package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamRoleAssignmentCreate,
		Read:   resourceIamRoleAssignmentRead,
		Importer: &schema.ResourceImporter{
			State: resourceIamRoleAssignmentImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Delete: resourceIamRoleAssignmentDelete,
		Schema: iamRoleAssignmentSchema(),
	}
}

func resourceIamRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	project_id := d.Get("project_id").(string)
	group_id := d.Get("group_id").(string)
	// group_name := d.Get("group_name").(string)
	// group, err := getClient(meta).IamGroup.GetGroupOfProject(project_id, group_name)
	// group_id := group.ID
	// if err != nil {
	// 	return fmt.Errorf("not found group `%s` in project %s: %s", group_name, project_id, err)
	// }
	_, err := getClient(meta).IamProject.AssignRoleFromGroupOnProject(project_id, group_id, d.Get("role_id").(string))

	if err != nil {
		return fmt.Errorf("error creating iam role assignment: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s_%s", d.Get("project_id").(string), group_id, d.Get("role_id").(string)))
	return nil
}
func resourceIamRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	parts := make([]string, 3)
	splitIdx1 := -1
	splitIdx2 := -1
	// Find the first underscore
	for i := 0; i < len(id); i++ {
		if id[i] == '_' {
			splitIdx1 = i
			break
		}
	}
	// Find the second underscore
	if splitIdx1 != -1 {
		for j := splitIdx1 + 1; j < len(id); j++ {
			if id[j] == '_' {
				splitIdx2 = j
				break
			}
		}
	}
	if splitIdx1 != -1 && splitIdx2 != -1 {
		parts[0] = id[:splitIdx1]
		parts[1] = id[splitIdx1+1 : splitIdx2]
		parts[2] = id[splitIdx2+1:]
		_ = d.Set("project_id", parts[0])
		_ = d.Set("group_id", parts[1])
		_ = d.Set("role_id", parts[2])

		// group_id := parts[1]
		// groups, err := getClient(meta).IamGroup.List(map[string]string{"project_id": parts[0]})
		// if err != nil {
		// 	return fmt.Errorf("not found group `%s` in project %s: %s", parts[1], parts[0], err)
		// }
		// var foundGroup *gocmcapiv2.IamGroup
		// for _, group := range groups {
		// 	if group.ID == group_id {
		// 		foundGroup = &group
		// 		break
		// 	}
		// }
		// if foundGroup == nil {
		// 	return fmt.Errorf("Group with id `%s` not found in project %s", parts[1], parts[0])
		// }
		// d.Set("group_name", foundGroup.Name)
	}
	return nil
}

func resourceIamRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	project_id := d.Get("project_id").(string)
	group_id := d.Get("group_id").(string)
	// group_name := d.Get("group_name").(string)
	// group, err := getClient(meta).IamGroup.GetGroupOfProject(project_id, group_name)
	// group_id := group.ID
	// if err != nil {
	// 	return fmt.Errorf("not found group `%s` in project %s: %s", group_name, project_id, err)
	// }
	_, err := getClient(meta).IamProject.UnsignRoleFromGroupOnProject(project_id, group_id, d.Get("role_id").(string))

	if err != nil {
		return fmt.Errorf("error delete iam role assignment: %v", err)
	}
	return nil
}

func resourceIamRoleAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamRoleAssignmentRead(d, meta)
	return []*schema.ResourceData{d}, err
}
