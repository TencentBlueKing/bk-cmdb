package v3

// v3Resp v3 api response data struct
type v3Resp struct {
	Result  bool        `json:"result"`
	Code    int         `json:"bk_error_code"`
	Message string      `json:"bk_error_msg"`
	Data    interface{} `json:"data"`
}
