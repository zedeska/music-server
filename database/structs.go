package db

import "encoding/json"

type Track struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Path        string  `json:"path"`
	Filename    string  `json:"filename"`
	Artist      string  `json:"artist"`
	Album       string  `json:"album"`
	Year        int     `json:"year"`
	Duration    int     `json:"duration"`
	Cover       string  `json:"cover"`
	SampleRate  float32 `json:"sample_rate"`
	Bitrate     int     `json:"bitrate"`
	TrackNumber int     `json:"media_count"`
}

type Album struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Year       int     `json:"year"`
	Cover      string  `json:"cover"`
	Tracks     []Track `json:"tracks"`
	SampleRate float32 `json:"sample_rate"`
	Bitrate    int     `json:"bitrate"`
	TrackCount int     `json:"track_count"`
}

func (p *Album) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}
