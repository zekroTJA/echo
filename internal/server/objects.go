package server

type echoObject struct {
	Method        string              `json:"method,omitempty"`
	Host          string              `json:"host,omitempty"`
	Header        map[string][]string `json:"headers,omitempty"`
	BodyString    string              `json:"body_string,omitempty"`
	RemoteAddress string              `json:"remote_address,omitempty"`
	RawPath       string              `json:"raw_path,omitempty"`
	Query         map[string][]string `json:"query,omitempty"`
}
