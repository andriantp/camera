package model

type AnalyticsCapability struct {
	ROI                  bool `json:"roi"`
	AudioDetection       bool `json:"audio_detection"`
	FaceDetection        bool `json:"face_detection"`
	LineDetection        bool `json:"line_detection"`
	FieldDetection       bool `json:"field_detection"`
	RegionEntrance       bool `json:"region_entrance"`
	RegionExiting        bool `json:"region_exiting"`
	Loitering            bool `json:"loitering"`
	Group                bool `json:"group"`
	RapidMove            bool `json:"rapid_move"`
	Parking              bool `json:"parking"`
	UnattendedBaggage    bool `json:"unattended_baggage"`
	AttendedBaggage      bool `json:"attended_baggage"`
	SmartCalibration     bool `json:"smart_calibration"`
	IntelliTrace         bool `json:"intelli_trace"`
	PeopleDetection      bool `json:"people_detection"`
	DefocusDetection     bool `json:"defocus_detection"`
	SceneChangeDetection bool `json:"scene_change_detection"`
	StorageDetection     bool `json:"storage_detection"`
	ChannelResource      bool `json:"channel_resource"`
	SmartOverlapParams   bool `json:"smart_overlap_params"`
}
