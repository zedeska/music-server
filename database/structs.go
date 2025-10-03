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
	Platform    string  `json:"platform"`
}

type Album struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	ArtistID   int     `json:"artist_id"`
	Year       int     `json:"year"`
	Cover      string  `json:"cover"`
	Tracks     []Track `json:"tracks"`
	SampleRate float32 `json:"sample_rate"`
	Bitrate    int     `json:"bitrate"`
	TrackCount int     `json:"track_count"`
	Platform   string  `json:"platform"`
}

type Artist struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Image       string  `json:"image"`
	TopTracks   []Track `json:"top_tracks"`
	LastRelease []Album `json:"last_release"`
}

type Playlist struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Tracks []Track `json:"tracks"`
}

type Playlists struct {
	Playlists []Playlist `json:"playlists"`
}

type Custom_search_result struct {
	Tracks []Track `json:"tracks"`
	Albums []Album `json:"albums"`
}

func (p Custom_search_result) ToJSON() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return data
}

func (p *Playlists) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p *Album) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p *Playlist) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p *Track) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}

func (p *Artist) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}
