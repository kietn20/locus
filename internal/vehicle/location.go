package vehicle

import "encoding/json"

type LocationData struct {
	VehicleID string `json:"vehicle_id"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (Id *LocationData) FromJSON(data []byte) error {
	return json.Unmarshal(data, Id)
}

func (Id *LocationData) ToJSON() ([]byte, error) {
	return json.Marshal(Id)
}