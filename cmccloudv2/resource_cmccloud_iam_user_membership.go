package cmccloudv2

import (
	"fmt"
	"strings"
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

	splitIdx := strings.LastIndex(id, "-")
	if splitIdx == -1 {
		return fmt.Errorf("invalid ID format: %s", id)
	}

	groupName := id[:splitIdx]
	userName := id[splitIdx+1:]

	if groupName == "" || userName == "" {
		return fmt.Errorf("invalid ID format: %s", id)
	}

	_ = d.Set("group_name", groupName)
	_ = d.Set("user_name", userName)

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
