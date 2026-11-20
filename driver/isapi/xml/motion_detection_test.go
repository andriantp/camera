package xml

import (
	"net/http"
	"strings"
	"testing"
)

const motionDetectionXML = `<MotionDetection version="2.0">
	<enabled>true</enabled>
	<regionType>grid</regionType>
	<Grid><rowGranularity>18</rowGranularity><columnGranularity>22</columnGranularity></Grid>
	<MotionDetectionLayout version="2.0">
		<sensitivityLevel>20</sensitivityLevel>
		<layout><gridMap>fffffc</gridMap></layout>
	</MotionDetectionLayout>
</MotionDetection>`

func TestSetMotionSensitivity(t *testing.T) {
	got, err := SetMotionSensitivity([]byte(motionDetectionXML), 75)
	if err != nil {
		t.Fatal(err)
	}

	s := string(got)
	if !strings.Contains(s, "<sensitivityLevel>75</sensitivityLevel>") {
		t.Errorf("sensitivity not updated:\n%s", s)
	}
	if strings.Replace(s, "75", "20", 1) != motionDetectionXML {
		t.Errorf("other content changed:\n%s", s)
	}
}

func TestSetMotionSensitivityMissing(t *testing.T) {
	if _, err := SetMotionSensitivity([]byte(`<MotionDetection></MotionDetection>`), 50); err == nil {
		t.Fatal("expected error when sensitivityLevel is missing")
	}
}

func TestParseMotionDetection(t *testing.T) {
	got, err := ParseMotionDetection(response(http.StatusOK, motionDetectionXML))
	if err != nil {
		t.Fatal(err)
	}

	if !got.Enabled || got.RegionType != "grid" || got.Grid.RowGranularity != 18 ||
		got.Layout.SensitivityLevel != 20 || got.Layout.GridMap != "fffffc" {
		t.Errorf("got %+v", got)
	}
}

func TestParseResponseStatus(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"ok", `<ResponseStatus><statusCode>1</statusCode><statusString>OK</statusString><subStatusCode>ok</subStatusCode></ResponseStatus>`, false},
		{"rejected", `<ResponseStatus><statusCode>4</statusCode><statusString>Invalid Operation</statusString><subStatusCode>notSupport</subStatusCode></ResponseStatus>`, true},
		{"invalid xml", `not xml`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParseResponseStatus(response(http.StatusOK, tt.body))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
