package relay

import "encoding/json"

const (
	MethodPing           = "ping"
	MethodListTables     = "list_tables"
	MethodGetColumns     = "get_columns"
	MethodGetPrimaryKeys = "get_primary_keys"
	MethodExecuteQuery   = "execute_query"
)


type Request struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	SQL    string `json:"sql,omitempty"`
	Schema string `json:"schema,omitempty"`
	Table  string `json:"table,omitempty"`
}


type Response struct {
	ID     string          `json:"id"`
	Error  string          `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}