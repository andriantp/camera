package xml

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/model"
)

type httpHostNotificationList struct {
	Items []httpHostNotification `xml:"HttpHostNotification"`
}

type httpHostNotification struct {
	ID                       int    `xml:"id"`
	URL                      string `xml:"url"`
	ProtocolType             string `xml:"protocolType"`
	ParameterFormatType      string `xml:"parameterFormatType"`
	AddressingFormatType     string `xml:"addressingFormatType"`
	IPAddress                string `xml:"ipAddress"`
	PortNo                   int    `xml:"portNo"`
	UserName                 string `xml:"userName"`
	HTTPAuthenticationMethod string `xml:"httpAuthenticationMethod"`
}

func ParseHTTPHosts(resp *http.Response) ([]model.HTTPHostNotification, error) {
	var data httpHostNotificationList

	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode HTTP hosts: %w", err)
	}

	result := make([]model.HTTPHostNotification, 0, len(data.Items))

	for _, item := range data.Items {
		result = append(result, model.HTTPHostNotification{
			ID:                       item.ID,
			URL:                      item.URL,
			ProtocolType:             item.ProtocolType,
			ParameterFormatType:      item.ParameterFormatType,
			AddressingFormatType:     item.AddressingFormatType,
			IPAddress:                item.IPAddress,
			PortNo:                   item.PortNo,
			UserName:                 item.UserName,
			HTTPAuthenticationMethod: item.HTTPAuthenticationMethod,
		})
	}

	return result, nil
}
