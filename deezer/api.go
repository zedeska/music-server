package deezer

import (
	"errors"
	db "music-server/database"
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://api.deezer.com/"

func Search(query string) (*db.Custom_search_result, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "search",
		Query: map[string]string{
			"q":     strings.Replace(query, " ", "%20", -1),
			"limit": "10",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_search_result{},
	}

	res := request.RunRest()
	if res.Error != nil {
		return nil, errors.New("Error: " + res.Error.Error())
	}

	temp_results, _ := res.Body.(*Deezer_search_result)

	var results db.Custom_search_result

	for _, track := range temp_results.Data {
		if track.Type == "track" {
			results.Tracks = append(results.Tracks, db.Track{
				ID:       track.ID,
				Title:    track.Title,
				Artist:   track.Artist.Name,
				ArtistID: track.Artist.ID,
				Album:    track.Album.Title,
				Duration: track.Duration,
				Cover:    track.Album.CoverMedium,
			})
		}
	}

	return &results, nil
}

func GetTrack(id int) (db.Track, error) {
	request := goaxios.GoAxios{
		Url:    API_URL + "track/" + string(rune(id)),
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_track{},
	}

	res := request.RunRest()
	if res.Error != nil {
		return db.Track{}, errors.New("Error: " + res.Error.Error())
	}

	temp_result, _ := res.Body.(*Deezer_track)

	track := db.Track{
		ID:          int(temp_result.ID),
		Title:       temp_result.Title,
		Artist:      temp_result.Artist.Name,
		ArtistID:    temp_result.Artist.ID,
		Album:       temp_result.Album.Title,
		Duration:    temp_result.Duration,
		Cover:       temp_result.Album.CoverMedium,
		Platform:    "deezer",
		Bitrate:     16,
		SampleRate:  44.1,
		TrackNumber: temp_result.TrackPosition,
	}

	return track, nil
}
