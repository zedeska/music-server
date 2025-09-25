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
			ArtistID:   track.Performer.ID,
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
			ArtistID:   album.Artist.ID,
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
		ID:          temp_results.ID,
		Title:       temp_results.Title,
		Artist:      temp_results.Performer.Name,
		ArtistID:    temp_results.Performer.ID,
		Album:       temp_results.Album.Title,
		AlbumID:     temp_results.Album.ID,
		Duration:    temp_results.Duration,
		Year:        year,
		Cover:       temp_results.Album.Image.Large,
		Bitrate:     temp_results.MaximumBitDepth,
		SampleRate:  float32(temp_results.MaximumSamplingRate),
		TrackNumber: temp_results.TrackNumber,
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
			ID:          track.ID,
			Title:       track.Title,
			Artist:      track.Performer.Name,
			ArtistID:    track.Performer.ID,
			Album:       temp_results.Title,
			Duration:    track.Duration,
			Cover:       temp_results.Image.Small,
			Bitrate:     track.MaximumBitDepth,
			SampleRate:  float32(track.MaximumSamplingRate),
			TrackNumber: track.TrackNumber,
		})
	}

	album = db.Album{
		ID:         temp_results.ID,
		Title:      temp_results.Title,
		Artist:     temp_results.Artist.Name,
		ArtistID:   temp_results.Artist.ID,
		Year:       temp_results.ReleasedAt,
		Cover:      temp_results.Image.Large,
		Bitrate:    temp_results.MaximumBitDepth,
		SampleRate: float32(temp_results.MaximumSamplingRate),
		Tracks:     tracks,
		TrackCount: temp_results.TracksCount,
	}

	return album, nil
}

func GetArtist(id string) (db.Artist, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "artist/page",
		Query: map[string]string{
			"sort":      "release_date",
			"artist_id": id,
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent":        "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":          "798273057",
			"X-User-Auth-Token": "mPhMKfJnkp1M2GcWnznMdzlE9rT6jBA4O24c_KzGZ5uiMA7M3HCxZ4tGXMZljV-1kzA-hc86oMWyFdN1CD4vYw",
		},
		ResponseStruct: &QobuzArtist{},
	}

	res := request.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Artist{}, errors.New("Error fetching artist")
	}

	temp_results, _ := res.Body.(*QobuzArtist)

	var tracks []db.Track
	var albums []db.Album

	for e, track := range temp_results.TopTracks {
		tracks = append(tracks, db.Track{
			ID:          track.ID,
			Title:       track.Title,
			Artist:      track.Artist.Name.Display,
			ArtistID:    track.Artist.ID,
			Album:       track.Album.Title,
			AlbumID:     track.Album.ID,
			Duration:    track.Duration,
			Cover:       track.Album.Image.Large,
			TrackNumber: e + 1,
		})
	}

	for _, album := range temp_results.Releases[0].Items {
		albums = append(albums, db.Album{
			ID:       album.ID,
			Title:    album.Title,
			Artist:   album.Artist.Name.Display,
			ArtistID: album.Artist.ID,
			Cover:    album.Image.Large,
		})
	}

	for _, album := range temp_results.Releases[3].Items {
		albums = append(albums, db.Album{
			ID:       album.ID,
			Title:    album.Title,
			Artist:   album.Artist.Name.Display,
			ArtistID: album.Artist.ID,
			Cover:    album.Image.Large,
		})
	}

	var image string = "https://static.qobuz.com/images/artists/covers/medium/" + temp_results.Images.Portrait.Hash + "." + temp_results.Images.Portrait.Format

	var artist db.Artist = db.Artist{
		ID:          temp_results.ID,
		Name:        temp_results.Name.Display,
		Image:       image,
		TopTracks:   tracks,
		LastRelease: albums,
	}

	return artist, nil
}
