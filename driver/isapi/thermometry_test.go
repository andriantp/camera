package isapi

import (
	"strings"
	"testing"
)

func TestSetThermometryEnabledXML(t *testing.T) {
	const doc = `<ThermometryBasicParam version="2.0">
	<id>1</id>
	<enabled>false</enabled>
	<streamOverlay>true</streamOverlay>
</ThermometryBasicParam>`

	got, err := setThermometryEnabledXML([]byte(doc), true)
	if err != nil {
		t.Fatal(err)
	}

	s := string(got)
	if !strings.Contains(s, "<enabled>true</enabled>") {
		t.Errorf("enabled not set:\n%s", s)
	}
	if !strings.Contains(s, "<streamOverlay>true</streamOverlay>") || !strings.Contains(s, "<id>1</id>") {
		t.Errorf("other fields changed:\n%s", s)
	}

	got, err = setThermometryEnabledXML(got, false)
	if err != nil || !strings.Contains(string(got), "<enabled>false</enabled>") {
		t.Errorf("disable failed: %v\n%s", err, got)
	}
}

func TestSetThermometryEnabledXMLMissing(t *testing.T) {
	if _, err := setThermometryEnabledXML([]byte(`<Other><enabled>true</enabled></Other>`), true); err == nil {
		t.Fatal("expected error when ThermometryBasicParam is missing")
	}
}
