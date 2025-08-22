package cmccloudv2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIamUserMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceIamUserMembershipCreate,
		Read:   resourceIamUserMembershipRead,
		Importer: &schema.ResourceImporter{
			State: resourceIamUserMembershipImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Delete: resourceIamUserMembershipDelete,
		Schema: iamUserMembershipSchema(),
	}
}

func resourceIamUserMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamGroup.AddUserToGroup(d.Get("user_name").(string), d.Get("group_name").(string))

	if err != nil {
		return fmt.Errorf("error creating iam user membership: %s", err)
	}

	d.SetId(fmt.Sprintf("%s-%s", d.Get("group_name").(string), d.Get("user_name").(string)))
	return nil
}

func resourceIamUserMembershipRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	parts := make([]string, 2)
	splitIdx := -1
	for i := 0; i < len(id); i++ {
		if id[i] == '-' {
			splitIdx = i
			break
		}
	}
	if splitIdx != -1 {
		parts[0] = id[:splitIdx]
		parts[1] = id[splitIdx+1:]
		_ = d.Set("group_name", parts[0])
		_ = d.Set("user_name", parts[1])
	}

	return nil
}

func resourceIamUserMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := getClient(meta).IamGroup.RemoveUserFromGroup(d.Get("user_name").(string), d.Get("group_name").(string))

	if err != nil {
		return fmt.Errorf("error delete iam user membership: %v", err)
	}
	return nil
}

func resourceIamUserMembershipImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceIamUserMembershipRead(d, meta)
	return []*schema.ResourceData{d}, err
}
