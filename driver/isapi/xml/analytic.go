package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type analyticsCapability struct {
	ROI                  bool `xml:"isSupportROI"`
	AudioDetection       bool `xml:"isSupportAudioDetection"`
	FaceDetection        bool `xml:"isSupportFaceDetect"`
	LineDetection        bool `xml:"isSupportLineDetection"`
	FieldDetection       bool `xml:"isSupportFieldDetection"`
	RegionEntrance       bool `xml:"isSupportRegionEntrance"`
	RegionExiting        bool `xml:"isSupportRegionExiting"`
	Loitering            bool `xml:"isSupportLoitering"`
	Group                bool `xml:"isSupportGroup"`
	RapidMove            bool `xml:"isSupportRapidMove"`
	Parking              bool `xml:"isSupportParking"`
	UnattendedBaggage    bool `xml:"isSupportUnattendedBaggage"`
	AttendedBaggage      bool `xml:"isSupportAttendedBaggage"`
	SmartCalibration     bool `xml:"isSupportSmartCalibration"`
	IntelliTrace         bool `xml:"isSupportIntelliTrace"`
	PeopleDetection      bool `xml:"isSupportPeopleDetection"`
	DefocusDetection     bool `xml:"isSupportDefocusDetection"`
	SceneChangeDetection bool `xml:"isSupportSceneChangeDetection"`
	StorageDetection     bool `xml:"isSupportStorageDetection"`
	ChannelResource      bool `xml:"isSupportChannelResource"`
	SmartOverlapParams   bool `xml:"isSupportSmartOverlapParams"`
}

func ParseAnalyticsCapabilities(resp *http.Response) (*model.AnalyticsCapability, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}

	var data analyticsCapability

	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal analytics capability: %w", err)
	}

	return &model.AnalyticsCapability{
		ROI:                  data.ROI,
		AudioDetection:       data.AudioDetection,
		FaceDetection:        data.FaceDetection,
		LineDetection:        data.LineDetection,
		FieldDetection:       data.FieldDetection,
		RegionEntrance:       data.RegionEntrance,
		RegionExiting:        data.RegionExiting,
		Loitering:            data.Loitering,
		Group:                data.Group,
		RapidMove:            data.RapidMove,
		Parking:              data.Parking,
		UnattendedBaggage:    data.UnattendedBaggage,
		AttendedBaggage:      data.AttendedBaggage,
		SmartCalibration:     data.SmartCalibration,
		IntelliTrace:         data.IntelliTrace,
		PeopleDetection:      data.PeopleDetection,
		DefocusDetection:     data.DefocusDetection,
		SceneChangeDetection: data.SceneChangeDetection,
		StorageDetection:     data.StorageDetection,
		ChannelResource:      data.ChannelResource,
		SmartOverlapParams:   data.SmartOverlapParams,
	}, nil
}
