package geofence

import "time"

// GeoJSONFeature represents the structure for creating a new geofence.
type GeoJSONFeature struct {
	Type       string `json:"type"`
	Properties struct {
		Name string `json:"name"`
	} `json:"properties"`
	Geometry struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
}

// Geofence represents the structure of a geofence in the database.
type Geofence struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Area      string    `json:"area_geojson"`
	CreatedAt time.Time `json:"created_at"`
}