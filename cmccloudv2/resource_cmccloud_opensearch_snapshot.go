package cmccloudv2

import (
	"fmt"
	"strings"
	"time"

	"github.com/cmc-cloud/terraform-provider-cmccloudv2/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOpenSearchSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenSearchSnapshotCreate,
		Read:   resourceOpenSearchSnapshotRead,
		Delete: resourceOpenSearchSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOpenSearchSnapshotImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		SchemaVersion: 1,
		Schema:        opensearchSnapshotSchema(),
	}
}

func resourceOpenSearchSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	// #endregion
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.OpenSearch.CreateSnapshot(d.Get("cluster_id").(string), d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("error creating OpenSearch snapshot: %s", err)
	}
	d.SetId(instance.ID)
	_, err = waitUntilOpenSearchSnapshotActive(d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error creating OpenSearch: %s", err)
	}
	return resourceOpenSearchSnapshotRead(d, meta)
}

func resourceOpenSearchSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	instance, err := client.OpenSearch.GetSnapshot(d.Get("cluster_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("error retrieving OpenSearch snapshot %s: %v", d.Id(), err)
	}
	_ = d.Set("name", instance.Name)
	_ = d.Set("status", instance.Status)
	_ = d.Set("created_at", instance.StartTime)
	return nil
}
func resourceOpenSearchSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).goCMCClient()
	_, err := client.OpenSearch.Delete(d.Id())

	if err != nil {
		return fmt.Errorf("error delete opensearch snapshot: %v", err)
	}
	_, err = waitUntilOpenSearchSnapshotDeleted(d, meta)
	if err != nil {
		return fmt.Errorf("error delete opensearch snapshot: %v", err)
	}
	return nil
}

func resourceOpenSearchSnapshotImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceOpenSearchSnapshotRead(d, meta)
	return []*schema.ResourceData{d}, err
}

func waitUntilOpenSearchSnapshotActive(d *schema.ResourceData, meta interface{}, timeout time.Duration) (interface{}, error) {
	return waitUntilResourceStatusChanged(d, meta, []string{"SUCCESS"}, []string{"FAILED", "PARTIAL"}, WaitConf{ // IN_PROGRESS
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).OpenSearch.GetSnapshot(d.Get("cluster_id").(string), d.Id())
	}, func(obj interface{}) string {
		return strings.ToUpper(obj.(gocmcapiv2.OpenSearchSnapshot).Status)
	})
}

func waitUntilOpenSearchSnapshotDeleted(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	return waitUntilResourceDeleted(d, meta, WaitConf{
		Delay:      10 * time.Second,
		MinTimeout: 20 * time.Second,
	}, func(id string) (any, error) {
		return getClient(meta).OpenSearch.Get(id)
	})
}
