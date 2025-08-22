package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamCustomRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamCustomRoleAssignmentCreate,
		Read:   resourceIamCustomRoleAssignmentRead,
		Importer: &schema.ResourceImporter{
			State: resourceIamCustomRoleAssignmentImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Delete: resourceIamCustomRoleAssignmentDelete,
		Schema: iamCustomRoleAssignmentSchema(),
	}
}

func resourceIamCustomRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamCustomRole.AssignCustomRoleFromGroupOnProject(d.Get("project_id").(string), d.Get("group_id").(string), d.Get("custom_role_id").(string))

	if err != nil {
		return fmt.Errorf("error creating iam custom role assignment: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s_%s", d.Get("project_id").(string), d.Get("group_id").(string), d.Get("custom_role_id").(string)))
	return nil
}
func resourceIamCustomRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
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
		_ = d.Set("custom_role_id", parts[2])
	}
	return nil
}

func resourceIamCustomRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamCustomRole.UnsignCustomRoleFromGroupOnProject(d.Get("project_id").(string), d.Get("group_id").(string), d.Get("custom_role_id").(string))

	if err != nil {
		return fmt.Errorf("error delete iam custom role assignment: %v", err)
	}
	return nil
}

func resourceIamCustomRoleAssignmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamCustomRoleAssignmentRead(d, meta)
	return []*schema.ResourceData{d}, err
}
