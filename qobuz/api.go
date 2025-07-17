package qobuz

import (
	"errors"
	"fmt"
	db "music-server/database"
	"strconv"
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://www.qobuz.com/api.json/0.2/"

func Search(query string) (custom_search_result, error) {

	request := goaxios.GoAxios{
		Url: API_URL + "catalog/search",
		Query: map[string]string{
			"limit":  "10",
			"offset": "0",
			"query":  strings.Replace(query, " ", "%20", -1),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":   "798273057",
			"x-zone":     "US",
			"x-store":    "US-en",
		},
		ResponseStruct: &qobuz_search_result{},
	}

	res := request.RunRest()

	if res.Error != nil {
		fmt.Println("Error:", res.Error)
		return custom_search_result{}, res.Error
	}

	temp_results, _ := res.Body.(*qobuz_search_result)

	var results custom_search_result

	for _, track := range temp_results.Tracks.Items {
		results.Tracks = append(results.Tracks, db.Track{
			ID:         track.ID,
			Title:      track.Title,
			Artist:     track.Performer.Name,
			Album:      track.Album.Title,
			Duration:   track.Duration,
			Cover:      track.Album.Image.Large,
			Bitrate:    track.MaximumBitDepth,
			SampleRate: float32(track.MaximumSamplingRate),
		})
	}

	for _, album := range temp_results.Albums.Items {
		results.Albums = append(results.Albums, db.Album{
			ID:         album.ID,
			Title:      album.Title,
			Artist:     album.Artist.Name,
			Year:       album.ReleasedAt,
			Cover:      album.Image.Large,
			Bitrate:    album.MaximumBitDepth,
			SampleRate: float32(album.MaximumSamplingRate),
		})
	}

	return results, nil

}

func GetTrack(id int) (db.Track, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "track/get",
		Query: map[string]string{
			"track_id": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":   "798273057",
		},
		ResponseStruct: &QobuzTrack{},
	}

	res := request.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Track{}, res.Error
	}

	temp_results, _ := res.Body.(*QobuzTrack)

	year, _ := strconv.Atoi(strings.Split(temp_results.ReleaseDateOriginal, "-")[0])

	var track db.Track = db.Track{
		ID:         temp_results.ID,
		Title:      temp_results.Title,
		Artist:     temp_results.Performer.Name,
		Album:      temp_results.Album.Title,
		Duration:   temp_results.Duration,
		Year:       year,
		Cover:      temp_results.Album.Image.Large,
		Bitrate:    temp_results.MaximumBitDepth,
		SampleRate: float32(temp_results.MaximumSamplingRate),
	}

	return track, nil
}

func GetAlbum(id string) (db.Album, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "album/get",
		Query: map[string]string{
			"album_id": id,
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":   "798273057",
		},
		ResponseStruct: &QobuzAlbum{},
	}

	res := request.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Album{}, errors.New("Error fetching album")
	}

	temp_results, _ := res.Body.(*QobuzAlbum)

	var tracks []db.Track
	var album db.Album

	for _, track := range temp_results.Tracks.Items {
		tracks = append(tracks, db.Track{
			ID:         track.ID,
			Title:      track.Title,
			Artist:     track.Performer.Name,
			Album:      temp_results.Title,
			Duration:   track.Duration,
			Cover:      temp_results.Image.Small,
			Bitrate:    track.MaximumBitDepth,
			SampleRate: float32(track.MaximumSamplingRate),
		})
	}

	album = db.Album{
		ID:         temp_results.ID,
		Title:      temp_results.Title,
		Artist:     temp_results.Artist.Name,
		Year:       temp_results.ReleasedAt,
		Cover:      temp_results.Image.Large,
		Bitrate:    temp_results.MaximumBitDepth,
		SampleRate: float32(temp_results.MaximumSamplingRate),
		Tracks:     tracks,
		MediaCount: temp_results.MediaCount,
	}

	return album, nil
}
