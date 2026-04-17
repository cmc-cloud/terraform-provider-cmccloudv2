package gocmcapiv2

type DatastoreWrapper struct {
	Data Datastore `json:"data"`
}
type DatastoreListWrapper struct {
	Data struct {
		Docs      []Datastore `json:"docs"`
		Page      int         `json:"page"`
		Size      int         `json:"size"`
		Total     int         `json:"total"`
		TotalPage int         `json:"totalPage"`
	} `json:"data"`
}
type Datastore struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	VersionInfos []struct {
		ID          string `json:"id"`
		VersionName string `json:"versionName"`
		ModeInfo    []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"modeInfo"`
	} `json:"versionInfos"`
}
type CreateResponse struct {
	ID string `json:"id"`
}
