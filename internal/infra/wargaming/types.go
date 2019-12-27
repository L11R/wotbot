package wargaming

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Status string          `json:"status"`
	Error  *Error          `json:"error"`
	Data   json.RawMessage `json:"data"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"code: %s; message: %s; field: %s; value: %s",
		e.Code,
		e.Message,
		e.Field,
		e.Value,
	)
}

type PlayerData struct {
	Nickname  string `json:"nickname"`
	AccountID int    `json:"account_id"`
}
