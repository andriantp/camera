package model

type ImageCapabilities struct {
	WDR          WDRCapabilities          `json:"wdr"`
	BLC          BLCCapabilities          `json:"blc"`
	IRCutFilter  IRCutFilterCapabilities  `json:"ir_cut_filter"`
	WhiteBalance WhiteBalanceCapabilities `json:"white_balance"`
	Exposure     ExposureCapabilities     `json:"exposure"`
	Focus        FocusCapabilities        `json:"focus"`
	EIS          EISCapabilities          `json:"eis"`
	Dehaze       DehazeCapabilities       `json:"dehaze"`
}

type WDRCapabilities struct {
	Modes []string `json:"modes"`
}

type BLCCapabilities struct {
	Modes []string `json:"modes"`
}

type IRCutFilterCapabilities struct {
	Types []string `json:"types"`
}

type WhiteBalanceCapabilities struct {
	Styles []string `json:"styles"`
}

type ExposureCapabilities struct {
	Types []string `json:"types"`
}

type FocusCapabilities struct {
	Styles []string `json:"styles"`
}

type EISCapabilities struct {
	Supported bool `json:"supported"`
}

type DehazeCapabilities struct {
	Modes []string `json:"modes"`
}
