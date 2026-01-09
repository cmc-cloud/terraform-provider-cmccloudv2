package gocmcapiv2

import (
	"encoding/json"
)

// AccountService interface
type AccountService interface {
	Get() (Account, error)
}

// Account object
type Account struct {
	ID               string `json:"id"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
	Username         string `json:"username"`
	Enabled          bool   `json:"enabled"`
	Totp             bool   `json:"totp"`
	EmailVerified    bool   `json:"emailVerified"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
}
type account struct {
	client *Client
}

// Get backup detail
func (v *account) Get() (Account, error) {
	jsonStr, err := v.client.Get("account/info", map[string]string{})
	var backup Account
	if err == nil {
		err = json.Unmarshal([]byte(jsonStr), &backup)
	}
	return backup, err
}
