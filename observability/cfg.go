package observability

import (
	"encoding/json"
	"strconv"
)

type Cfg struct {
	LogPretty    bool         `json:"log_pretty"`
	LogLevel     string       `json:"log_level"`
	Trace        bool         `json:"trace_enable"`
	OltpPass     SecretString `json:"oltp_pass"`
	OltpEndpoint string       `json:"oltp_endpoint"`
}

type SecretString string

func (b SecretString) MarshalJSON() ([]byte, error) {
	return json.Marshal("********" + "_" + strconv.Itoa(len(b)))
}

func (b SecretString) GetSecret() string {
	return string(b)
}
