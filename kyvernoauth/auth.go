package kyvernoauth

// CheckRequest mirrors the structure of the Envoy CheckRequest for HTTP requests
type CheckRequest struct {
	Attributes struct {
		Request struct {
			HTTP struct {
				Method  string            `json:"method"`
				Headers map[string]string `json:"headers"`
				Path    string            `json:"path"`
				Host    string            `json:"host"`
				Scheme  string            `json:"scheme"`
				Query   string            `json:"query"`
			} `json:"http"`
		} `json:"request"`
	} `json:"attributes"`
}

// CheckResponse mirrors the structure of the Envoy CheckResponse
type CheckResponse struct {
	Status struct {
		Code int32 `json:"code"`
	} `json:"status"`
	HTTPResponse struct {
		Status int32 `json:"status"`
	} `json:"httpResponse"`
}
