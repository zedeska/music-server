package db

type Track struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Path       string  `json:"path"`
	Filename   string  `json:"filename"`
	Artist     string  `json:"artist"`
	Album      string  `json:"album"`
	Year       int     `json:"year"`
	Duration   int     `json:"duration"`
	Cover      string  `json:"cover"`
	SampleRate float32 `json:"sample_rate"`
	Bitrate    int     `json:"bitrate"`
}
