package models

type ArcDataResponse struct {
	Data []ArcEventResponse `json:"data"`
	Meta ArcMetadata        `json:"meta"`
}

type ArcEventResponse struct {
	MapId       int   `json:"mapId"`
	ConditionId int   `json:"conditionId"`
	StartTime   int64 `json:"startTime"`
	EndTime     int64 `json:"endTime"`
}

type ArcEvent struct {
	MapName       string `json:"mapDisplayName"`
	ConditionName string `json:"conditionName"`
	StartTime     int64  `json:"startTimestamp"`
	EndTime       int64  `json:"endTimestamp"`
}

type ArcMetadata struct {
	Maps       map[int]string           `json:"maps"`
	Conditions map[int]ConditionDetails `json:"conditions"`
	ServerTime int64                    `json:"serverTime"`
	IconsUrl   string                   `json:"iconsUrl"`
}

type ConditionDetails struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Icon string `json:"icon"`
}
