package model

type EventCapability struct {
	Storage StorageEventCapability `json:"storage"`
	Network NetworkEventCapability `json:"network"`
	Access  AccessEventCapability  `json:"access"`
	Video   VideoEventCapability   `json:"video"`
	System  SystemEventCapability  `json:"system"`
}

type StorageEventCapability struct {
	HDFull  bool `json:"hdd_full"`
	HDError bool `json:"hdd_error"`
}

type NetworkEventCapability struct {
	NICBroken  bool `json:"nic_broken"`
	IPConflict bool `json:"ip_conflict"`
}

type AccessEventCapability struct {
	IllegalAccess bool `json:"illegal_access"`
}

type VideoEventCapability struct {
	MotionDetection bool `json:"motion_detection"`
	TamperDetection bool `json:"tamper_detection"`
}

type SystemEventCapability struct {
	AbnormalReboot bool `json:"abnormal_reboot"`
}
