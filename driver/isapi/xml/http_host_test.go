package xml

import (
	"encoding/xml"
	"net/http"
	"strings"
	"testing"

	"github.com/andriantp/camera/model"
)

func TestMarshalHTTPHostNotificationRoundTrip(t *testing.T) {
	host := model.HTTPHostNotification{
		ID:                       1,
		URL:                      "/events?a=1&b=<2>",
		ProtocolType:             "HTTP",
		ParameterFormatType:      "XML",
		AddressingFormatType:     "ipaddress",
		IPAddress:                "192.168.1.64",
		PortNo:                   8080,
		UserName:                 "user",
		HTTPAuthenticationMethod: "none",
	}

	body, err := MarshalHTTPHostNotification(host)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(string(body), xml.Header) {
		t.Error("missing XML header")
	}
	if !strings.Contains(string(body), `xmlns="http://www.isapi.org/ver20/XMLSchema"`) {
		t.Error("missing ISAPI namespace")
	}

	list := `<HttpHostNotificationList>` + strings.TrimPrefix(string(body), xml.Header) + `</HttpHostNotificationList>`

	got, err := ParseHTTPHosts(response(http.StatusOK, list))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != host {
		t.Errorf("round trip = %+v, want %+v", got, host)
	}
}
