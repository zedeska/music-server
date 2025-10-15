package deezer

import (
	"errors"
	db "music-server/database"
	"strconv"
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://api.deezer.com/"

func searchTrack(query string) (*Deezer_track_search, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "search/track",
		Query: map[string]string{
			"q":     strings.Replace(query, " ", "%20", -1),
			"limit": "10",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_track_search{},
	}

	res := request.RunRest()
	if res.Error != nil {
		return nil, errors.New("Error: " + res.Error.Error())
	}

	result, _ := res.Body.(*Deezer_track_search)

	return result, nil

}

func searchAlbum(query string) (*Deezer_album_search, error) {
	request := goaxios.GoAxios{
		Url: API_URL + "search/album",
		Query: map[string]string{
			"q":     strings.Replace(query, " ", "%20", -1),
			"limit": "10",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_album_search{},
	}

	res := request.RunRest()
	if res.Error != nil {
		return nil, errors.New("Error: " + res.Error.Error())
	}

	result, _ := res.Body.(*Deezer_album_search)

	return result, nil

}

func Search(query string) (*db.Custom_search_result, error) {
	temp_tracks, err := searchTrack(query)
	if err != nil {
		return nil, err
	}
	temp_albums, err := searchAlbum(query)
	if err != nil {
		return nil, err
	}

	var results db.Custom_search_result

	for _, track := range temp_tracks.Data {
		results.Tracks = append(results.Tracks, db.Track{
			ID:       int(track.ID),
			Title:    track.Title,
			Artist:   track.Artist.Name,
			ArtistID: track.Artist.ID,
			Album:    track.Album.Title,
			Duration: track.Duration,
			Cover:    track.Album.CoverXl,
			Platform: "deezer",
		})
	}

	for _, album := range temp_albums.Data {
		results.Albums = append(results.Albums, db.Album{
			ID:       strconv.Itoa(album.ID),
			Title:    album.Title,
			Artist:   album.Artist.Name,
			ArtistID: album.Artist.ID,
			Cover:    album.CoverXl,
			Platform: "deezer",
		})
	}

	return &results, nil
}

func GetTrack(id int) (db.Track, error) {
	request := goaxios.GoAxios{
		Url:    API_URL + "track/" + strconv.Itoa(id),
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

	year, _ := strconv.Atoi(strings.Split(temp_result.ReleaseDate, "-")[0])

	track := db.Track{
		ID:          int(temp_result.ID),
		Title:       temp_result.Title,
		Artist:      temp_result.Artist.Name,
		ArtistID:    temp_result.Artist.ID,
		Album:       temp_result.Album.Title,
		Duration:    temp_result.Duration,
		Cover:       temp_result.Album.CoverXl,
		Platform:    "deezer",
		Bitrate:     16,
		SampleRate:  44.1,
		TrackNumber: temp_result.TrackPosition,
		Year:        year,
	}

	return track, nil
}

func GetAlbum(id int) (db.Album, error) {
	var tracks []db.Track
	var album db.Album

	request := goaxios.GoAxios{
		Url:    API_URL + "album/" + strconv.Itoa(id),
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_album{},
	}

	res := request.RunRest()
	if res.Error != nil {
		return db.Album{}, errors.New("Error: " + res.Error.Error())
	}

	temp_result_album, _ := res.Body.(*Deezer_album)

	request_track := goaxios.GoAxios{
		Url:    API_URL + "album/" + strconv.Itoa(id) + "/tracks",
		Method: "GET",
		Query: map[string]string{
			"limit": strconv.Itoa(temp_result_album.NbTracks),
		},
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_album_track{},
	}
	res = request_track.RunRest()
	if res.Error != nil {
		return db.Album{}, errors.New("Error: " + res.Error.Error())
	}

	temp_results_track, _ := res.Body.(*Deezer_album_track)

	for _, track := range temp_results_track.Data {
		tracks = append(tracks, db.Track{
			ID:          int(track.ID),
			Title:       track.Title,
			Artist:      track.Artist.Name,
			ArtistID:    track.Artist.ID,
			Album:       temp_result_album.Title,
			Duration:    track.Duration,
			Cover:       temp_result_album.CoverXl,
			TrackNumber: track.TrackPosition,
			Platform:    "deezer",
		})
	}

	album = db.Album{
		ID:         strconv.Itoa(temp_result_album.ID),
		Title:      temp_result_album.Title,
		Artist:     temp_result_album.Artist.Name,
		ArtistID:   temp_result_album.Artist.ID,
		Cover:      temp_result_album.CoverXl,
		Tracks:     tracks,
		TrackCount: temp_result_album.NbTracks,
		Platform:   "deezer",
	}

	return album, nil
}

func GetArtist(id string) (db.Artist, error) {
	request_artist := goaxios.GoAxios{
		Url:    API_URL + "artist/" + id,
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_artist{},
	}

	res := request_artist.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Artist{}, errors.New("Error fetching artist")
	}

	temp_results_artist, _ := res.Body.(*Deezer_artist)

	request_top := goaxios.GoAxios{
		Url:    API_URL + "artist/" + id + "/top",
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_artist_top{},
	}

	res = request_top.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Artist{}, errors.New("Error fetching artist")
	}

	temp_results_artist_top, _ := res.Body.(*Deezer_artist_top)

	request_album := goaxios.GoAxios{
		Url:    API_URL + "artist/" + id + "/albums",
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &Deezer_artist_albums{},
	}

	res = request_album.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return db.Artist{}, errors.New("Error fetching artist")
	}

	temp_results_artist_albums, _ := res.Body.(*Deezer_artist_albums)

	var tracks []db.Track
	var albums []db.Album

	for e, track := range temp_results_artist_top.Data {
		tracks = append(tracks, db.Track{
			ID:          int(track.ID),
			Title:       track.Title,
			Artist:      track.Artist.Name,
			ArtistID:    track.Artist.ID,
			Album:       track.Album.Title,
			Duration:    track.Duration,
			Cover:       track.Album.CoverXl,
			TrackNumber: e + 1,
			Platform:    "deezer",
		})
	}

	for _, album := range temp_results_artist_albums.Data {
		albums = append(albums, db.Album{
			ID:       strconv.Itoa(album.ID),
			Title:    album.Title,
			Artist:   temp_results_artist.Name,
			ArtistID: temp_results_artist.ID,
			Cover:    album.CoverXl,
			Platform: "deezer",
		})
	}

	var artist db.Artist = db.Artist{
		ID:          temp_results_artist.ID,
		Name:        temp_results_artist.Name,
		Image:       temp_results_artist.PictureXl,
		TopTracks:   tracks,
		LastRelease: albums,
	}

	return artist, nil
}
