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

func MarshalHTTPHostNotification(host model.HTTPHostNotification) ([]byte, error) {
	type httpHostNotification struct {
		XMLName                  xml.Name `xml:"HttpHostNotification"`
		Version                  string   `xml:"version,attr"`
		XMLNS                    string   `xml:"xmlns,attr"`
		ID                       int      `xml:"id"`
		URL                      string   `xml:"url"`
		ProtocolType             string   `xml:"protocolType"`
		ParameterFormatType      string   `xml:"parameterFormatType"`
		AddressingFormatType     string   `xml:"addressingFormatType"`
		IPAddress                string   `xml:"ipAddress"`
		PortNo                   int      `xml:"portNo"`
		UserName                 string   `xml:"userName"`
		HTTPAuthenticationMethod string   `xml:"httpAuthenticationMethod"`
	}

	data := httpHostNotification{
		Version:                  "2.0",
		XMLNS:                    "http://www.isapi.org/ver20/XMLSchema",
		ID:                       host.ID,
		URL:                      host.URL,
		ProtocolType:             host.ProtocolType,
		ParameterFormatType:      host.ParameterFormatType,
		AddressingFormatType:     host.AddressingFormatType,
		IPAddress:                host.IPAddress,
		PortNo:                   host.PortNo,
		UserName:                 host.UserName,
		HTTPAuthenticationMethod: host.HTTPAuthenticationMethod,
	}

	body, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("marshal HTTP host notification: %w", err)
	}

	return append([]byte(xml.Header), body...), nil
}
