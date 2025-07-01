package db

type Track struct {
	ID         int
	Title      string
	Path       string
	Filename   string
	Artist     string
	Album      string
	Year       int
	Duration   int
	Cover      string
	SampleRate float32
	Bitrate    int
}
