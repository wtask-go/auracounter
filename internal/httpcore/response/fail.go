package response

type ErrorDescription struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

type Fail struct {
	Error ErrorDescription `json:"error"`
}
