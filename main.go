package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	db "music-server/database"
	"music-server/deezer"
	"music-server/qobuz"
	"music-server/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var SONG_FOLDER string
var MAX_QUALITY string = "6"
var dbConn *sql.DB

var (
	downloadMutex   = &sync.Mutex{}
	activeDownloads = make(map[string]*sync.WaitGroup)
)

func main() {

	var err error
	if _, err := os.Stat("./db.db"); err != nil {
		os.Create("db.db")
	}

	dbConn, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		return
	}
	defer dbConn.Close()

	db.InitDB(dbConn)
	homedir, _ := os.UserHomeDir()
	SONG_FOLDER = homedir + "/Music"
	if _, err := os.Stat(SONG_FOLDER); os.IsNotExist(err) {
		fmt.Println("Music folder does not exist, creating it...")
		os.Mkdir(SONG_FOLDER, os.ModePerm)
	}

	http.HandleFunc("/play", playHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/album", getAlbumHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/playlist", playlistHandler)
	http.HandleFunc("/user-playlists", userPlaylistsHandler)
	http.HandleFunc("/create-playlist", createPlaylistHandler)
	http.HandleFunc("/add-to-playlist", addToPlaylistHandler)
	http.HandleFunc("/delete-playlist", deletePlaylistHandler)
	http.HandleFunc("/delete-track-from-playlist", deleteTrackFromPlaylistHandler)
	http.HandleFunc("/listened", listenedHandler)
	http.ListenAndServe(":8488", nil)
}

func deleteTrackFromPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Token      string `json:"token"`
		PlaylistID int    `json:"playlist_id"`
		TrackID    int    `json:"track_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := db.GetUserID(dbConn, data.Token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	e, err := db.IsPlaylistOwnedByUser(dbConn, userID, data.PlaylistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !e {
		http.Error(w, "You do not own this playlist", http.StatusForbidden)
		return
	}

	err = db.DeleteTrackFromPlaylist(dbConn, data.PlaylistID, data.TrackID)
	if err != nil {
		http.Error(w, "Failed to delete track from playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Track deleted from playlist successfully"))
}

func deletePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Token      string `json:"token"`
		PlaylistID int    `json:"playlist_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := db.GetUserID(dbConn, data.Token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	e, err := db.IsPlaylistOwnedByUser(dbConn, userID, data.PlaylistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !e {
		http.Error(w, "You do not own this playlist", http.StatusForbidden)
		return
	}

	err = db.DeletePlaylist(dbConn, data.PlaylistID)
	if err != nil {
		http.Error(w, "Failed to delete playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Playlist deleted successfully"))
}

func addToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Token      string `json:"token"`
		PlaylistID int    `json:"playlist_id"`
		TrackIDs   []struct {
			ID       int `json:"id"`
			Platform int `json:"platform"`
		} `json:"track_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := db.GetUserID(dbConn, data.Token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	e, err := db.IsPlaylistOwnedByUser(dbConn, userID, data.PlaylistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !e {
		http.Error(w, "You do not own this playlist", http.StatusForbidden)
		return
	}

	for _, elt := range data.TrackIDs {
		e, err := db.IsTrackInPlaylist(dbConn, data.PlaylistID, elt.ID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if e {
			continue
		}

		platformeName, err := utils.GetPlatformName(elt.Platform)
		if err != nil {
			http.Error(w, "Invalid platform parameter", http.StatusBadRequest)
			return
		}

		trackExist, _ := db.CheckIfTrackExists(dbConn, elt.ID, platformeName)

		if !trackExist {
			track, err := searchTrackFromID(elt.ID, platformeName)
			if err != nil {
				http.Error(w, "Failed to find track", http.StatusInternalServerError)
				return
			}

			trackExist, _ = db.CheckIfTrackExistsByArtistAndAlbum(dbConn, elt.ID, platformeName, track.Artist, track.Album, track.Title)

			if !trackExist {
				err = db.AddPartialTrack(dbConn, track)
				if err != nil {
					http.Error(w, "Failed to add partial track", http.StatusInternalServerError)
					return
				}
			}
		}

		track, err := db.GetTrack(dbConn, elt.ID, platformeName)
		if err != nil {
			http.Error(w, "Failed to get track", http.StatusInternalServerError)
			return
		}

		err = db.AddTrackToPlaylist(dbConn, data.PlaylistID, track.ID)
		if err != nil {
			http.Error(w, "Failed to add track to playlist", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tracks added to playlist successfully"))
}

func createPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := db.GetUserID(dbConn, data.Token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	playlistID, err := db.CreatePlaylist(dbConn, userID, data.Name)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(playlistID)))
}

func userPlaylistsHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	userID, err := db.GetUserID(dbConn, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	playlists, err := db.GetPlaylistsByUserID(dbConn, userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(playlists.ToJSON())
}

func playlistHandler(w http.ResponseWriter, r *http.Request) {
	playlistIDStr := r.URL.Query().Get("id")
	if playlistIDStr == "" {
		http.Error(w, "Missing playlist ID", http.StatusBadRequest)
		return
	}

	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	playlist, err := db.GetPlaylistByID(dbConn, playlistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tracks, err := db.GetPlaylistTracks(dbConn, playlistID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = constitutePlaylist(playlist, tracks)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(playlist.ToJSON())
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing artist ID", http.StatusBadRequest)
		return
	}

	artist, err := qobuz.GetArtist(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(artist.ToJSON())
}

func listenedHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10" // Default limit if not specified
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	userID, err := db.GetUserID(dbConn, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch listened tracks for the user
	listenedTracks, err := db.GetListenedTracks(dbConn, userID, limit)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(listenedTracks)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	platformStr := r.URL.Query().Get("p")
	token := r.URL.Query().Get("token")

	log.Println("Play request for track ID:", idStr, "on platform:", platformStr)

	if idStr == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}

	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	if platformStr == "" {
		http.Error(w, "Missing platform parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
		return
	}

	platform, err := strconv.Atoi(platformStr)
	if err != nil {
		http.Error(w, "Invalid platform parameter", http.StatusBadRequest)
		return
	}
	platformName, err := utils.GetPlatformName(platform)
	if err != nil {
		http.Error(w, "Invalid platform parameter", http.StatusBadRequest)
		return
	}

	tokenValid, err := db.CheckToken(dbConn, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !tokenValid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	filePath, err := play(id, platformName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := db.GetUserID(dbConn, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.AddToListen(dbConn, userId, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "audio/flac")

	http.ServeFile(w, r, filePath)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := db.Login(dbConn, creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := db.Register(dbConn, creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(token))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	results_qobuz, err := qobuz.Search(query)
	if err != nil {
		http.Error(w, "Invalid search query", http.StatusBadRequest)
		return
	}
	results_deezer, err := deezer.Search(query)
	if err != nil {
		http.Error(w, "Invalid search query", http.StatusBadRequest)
		return
	}

	filteredQobuzTracks := make([]db.Track, 0, len(results_qobuz.Tracks))
	for _, q := range results_qobuz.Tracks {
		dup := false
		for _, d := range results_deezer.Tracks {
			if utils.Normalize(d.Album) == utils.Normalize(q.Album) && utils.Normalize(d.Artist) == utils.Normalize(q.Artist) && utils.Normalize(d.Title) == utils.Normalize(q.Title) {
				dup = true
				break
			}
		}
		if !dup {
			filteredQobuzTracks = append(filteredQobuzTracks, q)
		}
	}

	filteredQobuzAlbums := make([]db.Album, 0, len(results_qobuz.Albums))
	for _, q := range results_qobuz.Albums {
		dup := false
		for _, d := range results_deezer.Albums {
			if utils.Normalize(d.Title) == utils.Normalize(q.Title) && utils.Normalize(d.Artist) == utils.Normalize(q.Artist) {
				dup = true
				break
			}
		}
		if !dup {
			filteredQobuzAlbums = append(filteredQobuzAlbums, q)
		}
	}

	results_deezer.Tracks = append(results_deezer.Tracks, filteredQobuzTracks...)
	results_deezer.Albums = append(results_deezer.Albums, filteredQobuzAlbums...)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results_deezer.ToJSON())
}

func getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	idStr := r.URL.Query().Get("id")
	platformStr := r.URL.Query().Get("p")

	if platformStr == "" {
		http.Error(w, "Missing platform parameter", http.StatusBadRequest)
		return
	}

	if idStr == "" {
		http.Error(w, "Missing album ID", http.StatusBadRequest)
		return
	}

	platform, err := strconv.Atoi(platformStr)
	if err != nil {
		http.Error(w, "Invalid platform parameter", http.StatusBadRequest)
		return
	}

	var album db.Album
	platformName, err := utils.GetPlatformName(platform)

	if platformName == "deezer" {
		var id int
		id, err = strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid album ID", http.StatusBadRequest)
			return
		}
		album, err = deezer.GetAlbum(id)
	} else if platformName == "qobuz" {
		album, err = qobuz.GetAlbum(idStr)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(album.ToJSON())
}

func play(id int, platform string) (string, error) {
	err := checkAndAddTrack(id, platform)
	if err != nil {
		return "", fmt.Errorf("failed to check and add track: %w", err)
	}

	track, err := db.GetTrack(dbConn, id, platform)
	if err != nil {
		return "", fmt.Errorf("failed to get track: %w", err)
	}

	return track.Path, nil
}

func constitutePlaylist(playlist *db.Playlist, tracks []int) error {
	for _, trackID := range tracks {
		track, err := db.GetTrack(dbConn, trackID, "")
		if err != nil {
			return fmt.Errorf("failed to get track with ID %d: %w", trackID, err)
		}
		playlist.Tracks = append(playlist.Tracks, *track)
	}
	return nil
}

func downloadTrack(id int, platform string) (string, string, error) {
	file_name := utils.RandomString(50)
	temp_path := filepath.Join(SONG_FOLDER, file_name+".tmp") // Temporary file
	final_path := filepath.Join(SONG_FOLDER, file_name)       // Final file

	var err error

	if platform == "qobuz" {
		err = qobuz.Download(id, MAX_QUALITY, temp_path) // Download to temp file
		if err != nil {
			os.Remove(temp_path) // Clean up on failure
			return "", "", errors.New("failed to cache track")
		}
	} else if platform == "deezer" {
		err = deezer.Download(id, temp_path) // Download to temp file
		if err != nil {
			os.Remove(temp_path) // Clean up on failure
			return "", "", errors.New("failed to cache track")
		}
	} else if platform == "" {
		ids, err := db.GetTrackIds(dbConn, id)
		if err != nil {
			return "", "", fmt.Errorf("failed to get track IDs: %w", err)
		}
		if ids[0] != 0 {
			err = qobuz.Download(ids[0], MAX_QUALITY, temp_path)
		} else if ids[1] != 0 {
			err = deezer.Download(ids[1], temp_path)
		} else {
			return "", "", errors.New("no platform available for this track")
		}
		if err != nil {
			os.Remove(temp_path) // Clean up on failure
			return "", "", errors.New("failed to cache track")
		}
	}

	// Verify the downloaded file exists and has content
	if info, err := os.Stat(temp_path); err != nil || info.Size() == 0 {
		os.Remove(temp_path)
		return "", "", errors.New("downloaded file is empty or corrupted")
	}

	// Atomically move temp file to final location
	err = os.Rename(temp_path, final_path)
	if err != nil {
		os.Remove(temp_path)
		return "", "", fmt.Errorf("failed to finalize download: %w", err)
	}

	log.Printf("Successfully downloaded track to: %s", final_path)
	return final_path, file_name, nil
}

func searchTrackFromID(id int, platform string) (db.Track, error) {
	if platform == "qobuz" {
		return qobuz.GetTrack(id)
	} else if platform == "deezer" {
		return deezer.GetTrack(id)
	}

	return db.Track{}, fmt.Errorf("track not found")
}

func checkAndAddTrack(trackID int, platform string) error {
	trackKey := fmt.Sprintf("%s_%d", platform, trackID)

	downloadMutex.Lock()

	// Check if this track is currently being downloaded
	if wg, exists := activeDownloads[trackKey]; exists {
		downloadMutex.Unlock()
		// Wait for the download to complete
		wg.Wait()
		// Track should now exist, just return
		return nil
	}

	// Check if track exists now (before starting download)
	trackExists, needDownload := db.CheckIfTrackExists(dbConn, trackID, platform)

	if !needDownload {
		downloadMutex.Unlock()
		return nil
	}

	// This track needs downloading and no one else is downloading it
	wg := &sync.WaitGroup{}
	wg.Add(1)
	activeDownloads[trackKey] = wg
	downloadMutex.Unlock()

	// Ensure cleanup happens
	defer func() {
		downloadMutex.Lock()
		delete(activeDownloads, trackKey)
		downloadMutex.Unlock()
		wg.Done()
	}()

	// Now do the actual download work
	var file_path string
	var file_name string
	var err error

	if !trackExists {
		track, err := searchTrackFromID(trackID, platform)
		if err != nil {
			return fmt.Errorf("failed to search track: %w", err)
		}

		trackExists, needDownload = db.CheckIfTrackExistsByArtistAndAlbum(dbConn, trackID, platform, track.Artist, track.Album, track.Title)
		if needDownload {
			file_path, file_name, err = downloadTrack(trackID, platform)
			if err != nil {
				return fmt.Errorf("failed to download track: %w", err)
			}
		}
		if !trackExists {
			if MAX_QUALITY == "6" {
				track.SampleRate = 44.1
				track.Bitrate = 16
			}

			track.Path = file_path
			track.Filename = file_name

			err = db.AddTrack(dbConn, track)
			if err != nil {
				return fmt.Errorf("failed to add track to database: %w", err)
			}
		} else if trackExists && needDownload {
			err = db.UpdateTrackPathAndFilename(dbConn, trackID, file_path, file_name)
			if err != nil {
				return fmt.Errorf("failed to update track in database: %w", err)
			}
		}
	} else if trackExists && needDownload {
		file_path, file_name, err = downloadTrack(trackID, platform)
		if err != nil {
			return fmt.Errorf("failed to download track: %w", err)
		}
		err = db.UpdateTrackPathAndFilename(dbConn, trackID, file_path, file_name)
		if err != nil {
			return fmt.Errorf("failed to update track in database: %w", err)
		}
	}

	return nil
}
