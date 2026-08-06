package dto

type RoadCuttingInfo struct {
	ID              string  `json:"id,omitempty"`
	RoadType        string  `json:"roadType"`
	RoadCuttingArea float64 `json:"roadCuttingArea"`
	Status          string  `json:"status,omitempty"`
}
