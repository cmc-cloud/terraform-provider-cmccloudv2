package gocmcapiv2

import (
	"encoding/json"
	"strconv"
)

// BoolFromString is a custom type that implements json.Unmarshaler
type BoolFromString bool
type BoolFromInt bool
type IntFromString int

// UnmarshalJSON method to convert "true"/"false" string to bool
func (b *BoolFromInt) UnmarshalJSON(data []byte) error {
	var val int
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	switch val {
	case 1:
		*b = BoolFromInt(true)
	case 0:
		*b = BoolFromInt(false)
	default:
		return nil
	}
	return nil
}

// UnmarshalJSON method to convert "true"/"false" string to bool
func (b *BoolFromString) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string value
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// Parse the string value to a boolean
	switch str {
	case "true":
		*b = BoolFromString(true)
	case "false":
		*b = BoolFromString(false)
	default:
		return nil //fmt.Errorf("invalid value for bool: %s", str)
	}
	return nil
}

// UnmarshalJSON method to convert string to int
func (i *IntFromString) UnmarshalJSON(data []byte) error {
	// Remove quotes from the string value
	str := string(data)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	// Parse the string value to an integer
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	// Set the integer value to the custom type
	*i = IntFromString(intValue)
	return nil
}
