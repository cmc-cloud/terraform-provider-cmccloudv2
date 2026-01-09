package gocmcapiv2

import (
	"encoding/json"
)

// CDNService interface
type CDNService interface {
	Get(id string) (CDN, error)
	List(cdn_id string, params map[string]string) ([]CDN, error)
	Create(params map[string]interface{}) (string, error)
	Update(id string, params map[string]interface{}) (ActionResponse, error)
	Delete(id string) (ActionResponse, error)
}

type CDN struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	OriginServerURL    string `json:"origin_server_url"`
	Type               string `json:"type"`
	CdnURL             string `json:"cdn_url"`
	MultiCdnURL        string `json:"multi_cdn_url"`
	OriginInfor        string `json:"origin_Infor"`
	Status             string `json:"status"`
	AppID              string `json:"app_Id"`
	LoadBalancerMethod string `json:"load_balancer_method"`
	Zone               string `json:"zone"`
	ZoneActive         string `json:"zone_active"`
	ZoneBackup         string `json:"zone_backup"`
	Vod                string `json:"vod"`
	MultiCdnStatus     string `json:"multi_cdn_status"`
	OriginSetting      struct {
		LoadBalancerMethod string        `json:"load_balancer_method"`
		HostHeader         string        `json:"host_header"`
		Protocol           string        `json:"protocol"`
		Port               IntFromString `json:"port"`
		OriginServerUrls   []struct {
			OriginServerURL string `json:"origin_server_url"`
			Domain          string `json:"domain"`
			FailTimeout     string `json:"fail_timeout"`
			MaxFails        string `json:"max_fails"`
		} `json:"origin_server_urls"`
	} `json:"origin_setting"`
	EdgeSettings struct {
		BrowserCacheTTL       int    `json:"browser_cache_ttl"`
		CachingLevel          string `json:"caching_level"`
		DevelopmentMode       string `json:"development_mode"`
		Thump                 string `json:"thump"`
		ImageResizing         string `json:"image_resizing"`
		Polish                string `json:"polish"`
		ImageVirtualizing     string `json:"image_virtualizing"`
		GzipLevel             int    `json:"gzip_level"`
		BrotliCompression     string `json:"brotli_compression"`
		AutoMinify            string `json:"auto_minify"`
		HTTP2                 string `json:"http2"`
		HTTP2Origin           string `json:"http2_origin"`
		Websocket             string `json:"websocket"`
		BrowserIntegrityCheck string `json:"browser_integrity_check"`
		BotFightMode          string `json:"bot_fight_mode"`
		AlwaysUseHTTPS        string `json:"always_use_https"`
		TLS13                 string `json:"tls13"`
		Hsts                  string `json:"hsts"`
	} `json:"edge_settings"`
	SiteCname []interface{} `json:"site_cname"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`

	ErrorMessage string `json:"message"`
}
type CDNCreateResponse struct {
	Message string `json:"message"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
	StatusCode int `json:"status_code"`
}
type CDNWrapper struct {
	Data CDN `json:"data"`
}
type CDNListWrapper struct {
	Data     []CDN `json:"data"`
	PageInfo struct {
		Page       string `json:"page"`
		PerPage    int    `json:"per_page"`
		TotalCount int    `json:"total_count"`
		TotalPages int    `json:"total_pages"`
	} `json:"page_info"`
}

type cdn struct {
	client *Client
}

// Get CDN detail
func (v *cdn) Get(id string) (CDN, error) {
	jsonStr, err := v.client.Get("cdn/cdn/sites/"+id, map[string]string{})
	var cdn CDNWrapper
	var nilres CDN
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &cdn)
	if err != nil {
		return nilres, err
	}
	return cdn.Data, err
}
func (v *cdn) List(cdn_id string, params map[string]string) ([]CDN, error) {
	jsonStr, err := v.client.Get("cdn/cdn/sites/"+cdn_id, map[string]string{})
	var cdn CDNListWrapper
	var nilres []CDN
	if err != nil {
		return nilres, err
	}
	err = json.Unmarshal([]byte(jsonStr), &cdn)
	if err != nil {
		return nilres, err
	}
	return cdn.Data, err
}
func (v *cdn) Delete(id string) (ActionResponse, error) {
	return v.client.PerformDelete("cdn/cdn/sites/" + id)
}
func (v *cdn) Create(params map[string]interface{}) (string, error) {
	jsonStr, err := v.client.Post("cdn/cdn/sites", params)
	var response CDNCreateResponse
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", err
	}
	return response.Data.ID, nil
}

func (s *cdn) Update(id string, params map[string]interface{}) (ActionResponse, error) {
	return s.client.PerformUpdate("cdn/cdn/sites/"+id, params)
}
