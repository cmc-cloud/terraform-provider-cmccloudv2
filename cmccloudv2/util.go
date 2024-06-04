package cmccloudv2

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/cmc-cloud/gocmcapiv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type WaitConf struct {
	Delay      time.Duration // Wait this time before starting checks
	MinTimeout time.Duration // Smallest time to wait before refreshes
	Timeout    time.Duration // The amount of time to wait before timeout
}

func isSet(diff *schema.ResourceDiff, key string) bool {
	v, ok := diff.GetOkExists(key)
	if !ok {
		gocmcapiv2.Logo("isSet "+key+" not exists ", v)
		return false
	} else {
		gocmcapiv2.Logo("isSet "+key+" exists ", v)
		switch v.(type) {
		case int:
			if v.(int) == 0 {
				return false
			}
			return true
		case string:
			if v.(string) == "" {
				return false
			}
			return true
		case bool:
			if v.(bool) == false {
				return false
			}
			return true
		default:
		}
	}
	return ok
}
func setInt(d *schema.ResourceData, key string, newval int) {
	// kiem tra xem co khac voi gia tri hien tai ko, gia tri hien tai co the chi la gia tri default
	if v, ok := d.GetOk(key); ok && v.(int) != newval {
		_ = d.Set(key, newval)
	}
}
func setString(d *schema.ResourceData, key string, newval string) {
	// kiem tra xem co khac voi gia tri hien tai ko, gia tri hien tai co the chi la gia tri default
	v, ok := d.GetOk(key)
	// gocmcapiv2.Logs("setString old val " + v.(string) + ", newval = " + newval)
	if ok && v.(string) != newval {
		// gocmcapiv2.Logs("setString old val " + v.(string) + ", newval = " + newval)
		_ = d.Set(key, newval)
	}
}

// set các giá trị dạng []string không phân biệt thứ tự (TypeSet)
func setTypeSet(d *schema.ResourceData, key string, newval []string) {
	// kiem tra xem co khac voi gia tri hien tai ko, gia tri hien tai co the chi la gia tri default
	v, ok := d.GetOk(key)
	items := v.(*schema.Set).List()
	itemStrings := make([]string, len(items))
	index := 0
	for _, val := range items {
		itemStrings[index] = val.(string)
		index++
	}
	if ok && !areTypeSetEqual(itemStrings, newval) {
		_ = d.Set(key, newval)
	}
}

// areSlicesEqual so sánh hai slice mà không phân biệt thứ tự phần tử
func areTypeSetEqual(a, b []string) bool {
	// Nếu độ dài khác nhau, chắc chắn hai slice không giống nhau
	if len(a) != len(b) {
		return false
	}

	// Tạo map để đếm số lần xuất hiện của từng phần tử trong slice a
	counts := make(map[string]int)
	for _, item := range a {
		counts[item]++
	}

	// Giảm đếm số lần xuất hiện của từng phần tử trong slice b
	for _, item := range b {
		if _, found := counts[item]; !found {
			return false
		}
		counts[item]--
		if counts[item] == 0 {
			delete(counts, item)
		}
	}

	// Kiểm tra xem map counts có trống không
	return len(counts) == 0
}

func arrayContains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
func getClient(meta interface{}) *gocmcapiv2.Client {
	return meta.(*CombinedConfig).goCMCClient()
}
func _checkDeletedRefreshFunc(d *schema.ResourceData, meta interface{}, getResourceFunc func(id string) (interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resource, err := getResourceFunc(d.Id())
		if errors.Is(err, gocmcapiv2.ErrNotFound) {
			return resource, "true", nil
		}
		return resource, "false", nil
	}
}
func waitUntilResourceDeleted(d *schema.ResourceData, meta interface{}, timeout WaitConf, getResourceFunc func(id string) (interface{}, error)) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"false"},
		Target:         []string{"true"},
		Refresh:        _checkDeletedRefreshFunc(d, meta, getResourceFunc),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          timeout.Delay,
		MinTimeout:     timeout.MinTimeout,
		NotFoundChecks: 3,
	}
	return stateConf.WaitForState()
}

func _checkStatusRefreshFunc(d *schema.ResourceData, meta interface{}, errorStatus []string,
	getResourceFunc func(id string) (interface{}, error),
	getStatusFunc func(obj interface{}) string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resource, err := getResourceFunc(d.Id())
		if err != nil {
			fmt.Errorf("Error retrieving resource %s: %v", d.Id(), err)
			return nil, "", err
		}
		newStatus := getStatusFunc(resource)
		for _, v := range errorStatus {
			if v == newStatus {
				return resource, newStatus, fmt.Errorf("got the status " + newStatus)
			}
		}
		return resource, newStatus, nil
	}
}
func waitUntilResourceStatusChanged(d *schema.ResourceData, meta interface{}, targetStatus []string, errorStatus []string, timeout WaitConf,
	getResourceFunc func(id string) (interface{}, error),
	getStatusFunc func(obj interface{}) string) (interface{}, error) {
	_timeout := d.Timeout(schema.TimeoutCreate) // neu ko set gia tri timeout thi mac dinh lay TimeoutCreate
	if timeout.Timeout > 0 {                    // co set gia tri timeout, mot vai truong hop timeout nay se set = TimeoutUpdate
		_timeout = timeout.Timeout
	}
	stateConf := &resource.StateChangeConf{
		// Pending:        []string{""},
		Target:         targetStatus,
		Refresh:        _checkStatusRefreshFunc(d, meta, errorStatus, getResourceFunc, getStatusFunc),
		Timeout:        _timeout,
		Delay:          timeout.Delay,
		MinTimeout:     timeout.MinTimeout,
		NotFoundChecks: 3,
	}
	return stateConf.WaitForState()
}

// func waitUntilServerChangeState(d *schema.ResourceData, meta interface{}, id string, pendingStatus []string, targetStatus []string) (interface{}, error) {
// 	log.Printf("[INFO] Waiting for server with id (%s) to be "+strings.Join(targetStatus, ","), id)
// 	stateConf := &resource.StateChangeConf{
// 		Pending:        pendingStatus,
// 		Target:         targetStatus,
// 		Refresh:        serverStateRefreshfunc(d, meta, id),
// 		Timeout:        60 * 20 * time.Second,
// 		Delay:          10 * time.Second,
// 		MinTimeout:     10 * time.Second,
// 		NotFoundChecks: 5,
// 	}
// 	return stateConf.WaitForState()
// }

//	func serverStateRefreshfunc(d *schema.ResourceData, meta interface{}, id string) resource.StateRefreshFunc {
//		return func() (interface{}, string, error) {
//			client := meta.(*CombinedConfig).goCMCClient()
//			server, err := client.Server.Get(d.Id(), false)
//			if err != nil {
//				fmt.Errorf("Error retrieving server %s: %v", id, err)
//				return nil, "", err
//			}
//			return server, server.VMState, nil
//		}
//	}
//
// kiem tra xem 1 truong trong sub block co thay doi hay khong
func isSubBlockFieldChanged(d *schema.ResourceData, block_name string, field_name string) (bool, interface{}) {
	if d.HasChange(block_name) {
		// Get the old and new values
		old, new := d.GetChange(block_name)

		oldSubBlocks := old.([]interface{})
		newSubBlocks := new.([]interface{})

		for i := range oldSubBlocks {
			oldSubBlock := oldSubBlocks[i].(map[string]interface{})
			newSubBlock := newSubBlocks[i].(map[string]interface{})

			// Check if field_name has changed
			if oldSubBlock[field_name] != newSubBlock[field_name] {
				return true, newSubBlock[field_name]
			}
		}
	}
	return false, nil
}
func getFirstBlock(d *schema.ResourceData, key string) map[string]interface{} {
	// var block map[string]interface{}
	if v, ok := d.GetOk(key); ok {
		blockList := v.([]interface{})
		if len(blockList) > 0 {
			return blockList[0].(map[string]interface{})
		}
	}
	return nil
}
func getDiffSet(olds interface{}, news interface{}) (*schema.Set, *schema.Set) {
	oldSet := olds.(*schema.Set)
	newSet := news.(*schema.Set)

	// Tìm các phần tử bị xóa
	removed := oldSet.Difference(newSet)

	// Tìm các phần tử mới
	added := newSet.Difference(oldSet)
	return removed, added
}
func getStringArrayFromTypeSet(set *schema.Set) []string {
	items := set.List()
	stringArray := make([]string, 0)
	for i := 0; i < len(items); i++ {
		networkInterface := items[i].(string)
		stringArray = append(stringArray, networkInterface)
	}
	return stringArray
}

// isPrivateIP checks if the given IP address is a private IP address.
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	privateIPv4Ranges := []string{
		"10.0.0.0/8",     // 10.0.0.0 - 10.255.255.255
		"172.16.0.0/12",  // 172.16.0.0 - 172.31.255.255
		"192.168.0.0/16", // 192.168.0.0 - 192.168.255.255
	}

	privateIPv6Ranges := []string{
		"fc00::/7", // fc00::/7 (fc00:: - fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff)
	}

	// Check if it's an IPv4 private address
	if ip.To4() != nil {
		for _, cidr := range privateIPv4Ranges {
			_, ipnet, _ := net.ParseCIDR(cidr)
			if ipnet.Contains(ip) {
				return true
			}
		}
	}

	// Check if it's an IPv6 private address
	if ip.To16() != nil && strings.Contains(ipStr, ":") {
		for _, cidr := range privateIPv6Ranges {
			_, ipnet, _ := net.ParseCIDR(cidr)
			if ipnet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// IfThenElse one line if else condition: IfThenElse(1 == 1, "Yes", false) => "Yes"
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func stringArrayToSet(items []string) *schema.Set {
	set := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range items {
		set.Add(v)
	}
	return set
}

func setToStringArray(items *schema.Set) []string {
	flatten := make([]string, items.Len())

	for i, v := range items.List() {
		flatten[i] = v.(string)
	}
	return flatten
}

/*
func interfaceToString(items []interface{}) []string {
	flatten := make([]string, len(items))

	for i, v := range items {
		flatten[i] = fmt.Sprint(v)
	}
	return flatten
}
*/
