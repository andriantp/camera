package model

type HTTPHostNotification struct {
	ID                       int    `json:"id"`
	URL                      string `json:"url"`
	ProtocolType             string `json:"protocol_type"`
	ParameterFormatType      string `json:"parameter_format_type"`
	AddressingFormatType     string `json:"addressing_format_type"`
	IPAddress                string `json:"ip_address"`
	PortNo                   int    `json:"port_no"`
	UserName                 string `json:"user_name"`
	HTTPAuthenticationMethod string `json:"http_authentication_method"`
}
