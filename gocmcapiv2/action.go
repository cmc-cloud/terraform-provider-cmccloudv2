package gocmcapiv2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// response when create an action
type DBActionResponse struct {
	Data struct {
		ActionID string `json:"actionId"`
	} `json:"data"`
}

// response when get status of an action
type DBListActionResponse struct {
	Status  int    `json:"status"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		ActionID string            `json:"actionId"`
		Status   string            `json:"status"`
		Docs     []json.RawMessage `json:"docs"`
	} `json:"data"`
}

func parseDocs[T any](raw []json.RawMessage) ([]T, error) {
	out := make([]T, 0, len(raw))
	for _, r := range raw {
		var item T
		if err := json.Unmarshal(r, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
func WaitForActionResult[T any](
	client *Client,
	url string,
	actionID string,
	timeoutMinutes int,
) ([]T, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMinutes)*time.Minute)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for action %s", actionID)

		case <-ticker.C:
			bodyStr, err := client.Get(url, map[string]string{})
			if err != nil {
				return nil, err
			}

			var resp DBListActionResponse
			if err := json.Unmarshal([]byte(bodyStr), &resp); err != nil {
				return nil, err
			}

			if !resp.Success || resp.Status != 1 {
				return nil, fmt.Errorf("api error: %s", resp.Msg)
			}

			switch resp.Data.Status {
			case "completed":
				return parseDocs[T](resp.Data.Docs)

			case "failed":
				return nil, fmt.Errorf("action failed")

			default:
				// pending / running → tiếp tục poll
			}
		}
	}
}
