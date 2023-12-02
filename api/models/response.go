package models

type Response struct {
	Data  map[string]interface{} `json:"data"`
	Match string                 `json:"match"`
}
