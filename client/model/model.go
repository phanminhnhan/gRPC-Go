package model

type ResponseData struct {
	Message  string `json:"message"`
	StattusCode int `json:"stattus_code"`
	Data interface{} `json:"data"`
}