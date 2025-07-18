package server

type echoObject struct {
	Method        string              `json:"method,omitempty"`
	Host          string              `json:"host,omitempty"`
	Header        map[string][]string `json:"headers,omitempty"`
	BodyString    string              `json:"body_string,omitempty"`
	BodyParsed    any                 `json:"body_parsed,omitempty"`
	RemoteAddress string              `json:"remote_address,omitempty"`
	Path          string              `json:"path,omitempty"`
	Query         map[string][]string `json:"query,omitempty"`
}
