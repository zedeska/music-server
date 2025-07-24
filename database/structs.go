package db

import "encoding/json"

type Track struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Path        string  `json:"path"`
	Filename    string  `json:"filename"`
	Artist      string  `json:"artist"`
	ArtistID    int     `json:"artist_id"`
	Album       string  `json:"album"`
	AlbumID     string  `json:"album_id"`
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

type Artist struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	TopTracks   []Track `json:"top_tracks"`
	LastRelease []Album `json:"last_release"`
}

func (p *Album) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p *Track) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}
