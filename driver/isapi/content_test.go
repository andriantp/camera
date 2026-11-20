package isapi

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"testing"
)

func TestEscapeXML(t *testing.T) {
	tests := []string{
		"plain",
		`p&ss<w>"'`,
		"rtsp://192.168.1.1/Streaming/tracks/101?starttime=20260101T000000Z&size=1024",
	}

	for _, in := range tests {
		var v struct {
			Value string `xml:"v"`
		}

		doc := "<r><v>" + escapeXML(in) + "</v></r>"
		if err := xml.Unmarshal([]byte(doc), &v); err != nil {
			t.Fatalf("escapeXML(%q) produced invalid XML: %v", in, err)
		}
		if v.Value != in {
			t.Errorf("round trip = %q, want %q", v.Value, in)
		}
	}
}

func TestContentDownloadEscapesRequest(t *testing.T) {
	const (
		password    = `p&ss<w>rd`
		playbackURI = "rtsp://cam/Streaming/tracks/101?starttime=20260101T000000Z&size=2048"
	)

	var got struct {
		PlaybackURI string `xml:"playbackURI"`
		UserName    string `xml:"userName"`
		Password    string `xml:"password"`
	}

	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := xml.Unmarshal(body, &got); err != nil {
			t.Errorf("camera received invalid XML: %v\n%s", err, body)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = w.Write([]byte("video-bytes"))
	})
	cam.Password = password

	_, body, err := client.ContentDownload(context.Background(), cam, playbackURI)
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()

	data, _ := io.ReadAll(body)
	if string(data) != "video-bytes" {
		t.Errorf("body = %q", data)
	}

	if got.Password != password || got.PlaybackURI != playbackURI || got.UserName != cam.Username {
		t.Errorf("camera received %+v", got)
	}
}
