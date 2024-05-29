package cmccloudv2

import (
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
