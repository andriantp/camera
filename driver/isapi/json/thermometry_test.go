package json

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestParseThermometryRulesTemperatureInfo(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{
			"ThermometryRulesTemperatureInfoList": {
				"ThermometryRulesTemperatureInfo": [{
					"maxTemperature": 36.8,
					"minTemperature": 22.1,
					"averageTemperature": 30.5,
					"MaxTemperaturePoint": {"positionX": 0.25, "positionY": 0.5},
					"MinTemperaturePoint": {"positionX": 0.75, "positionY": 0.1},
					"isFreezedata": true
				}]
			}
		}`)),
	}

	got, err := ParseThermometryRulesTemperatureInfo(resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("items = %d, want 1", len(got))
	}
	item := got[0]
	if item.MaxTemperature != 36.8 || item.MinTemperature != 22.1 || item.AverageTemperature != 30.5 ||
		item.MaxTemperaturePoint.PositionX != 0.25 || item.MinTemperaturePoint.PositionY != 0.1 || !item.IsFreezeData {
		t.Errorf("got %+v", item)
	}
}

func TestParseThermometryRulesTemperatureInfoInvalid(t *testing.T) {
	resp := &http.Response{Body: io.NopCloser(strings.NewReader("not json"))}

	if _, err := ParseThermometryRulesTemperatureInfo(resp); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
