package gocmcapiv2

import "encoding/json"

// EIPService interface
type BillingModeService interface {
	SetServerBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetDatabaseInstanceBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetVPCBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetEIPBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetVolumeBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetLoadBalancerBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetKubernateBilingMode(id string, billing_mode string, node_group_role string) (ActionResponse, error)
	SetVPNBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetEFSBilingMode(id string, billing_mode string) (ActionResponse, error)
	SetRedisInstanceBilingMode(id string, billing_mode string) (ActionResponse, error)
	GetBilingMode(id string, resource_type string) (string, error)
	SetKubernateNodeGroupBilingMode(cluster_id string, nodegroup_id string, billing_mode string) (ActionResponse, error)
}

type BillingMode struct {
	BillingMode string `json:"billing_mode"`
}

type billingmode struct {
	client *Client
}

func (b *billingmode) GetBilingMode(id string, resource_type string) (string, error) {
	jsonStr, err := b.client.Get("billing/get_billing_mode", map[string]string{
		"id":            id,
		"resource_type": resource_type,
	})
	var response BillingMode
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", err
	}
	return response.BillingMode, err
}

func (b *billingmode) SetBilingMode(params map[string]interface{}) (ActionResponse, error) {
	return b.client.PerformUpdate("billing/update_billing_mode", params)
}

func (b *billingmode) SetServerBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EC", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetDatabaseInstanceBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EC", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetVPCBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "VPC", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetEIPBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EIP", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetVolumeBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EV", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetLoadBalancerBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "ELB", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetKubernateBilingMode(id string, billing_mode string, node_group_role string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "CKE", "resource_id": id, "billing_mode": billing_mode, "node_group_role": node_group_role})
}
func (b *billingmode) SetKubernateNodeGroupBilingMode(cluster_id string, node_group_id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "Nodegroup", "resource_id": cluster_id, "billing_mode": billing_mode, "node_group_id": node_group_id})
}
func (b *billingmode) SetVPNBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "VPN", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetEFSBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EFS", "resource_id": id, "billing_mode": billing_mode})
}
func (b *billingmode) SetRedisInstanceBilingMode(id string, billing_mode string) (ActionResponse, error) {
	return b.SetBilingMode(map[string]interface{}{"resource_type": "EC", "resource_id": id, "billing_mode": billing_mode})
}
