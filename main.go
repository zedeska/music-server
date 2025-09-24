package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	db "music-server/database"
	"music-server/qobuz"
	"music-server/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var SONG_FOLDER string
var MAX_QUALITY string = "6"
var dbConn *sql.DB

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
		TrackIDs   []int  `json:"track_ids"`
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

	for _, ID := range data.TrackIDs {
		e, err := db.IsTrackInPlaylist(dbConn, data.PlaylistID, ID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if e {
			continue
		}

		trackExist, _ := db.CheckIfTrackExists(dbConn, ID)

		if !trackExist {
			track, err := qobuz.GetTrack(ID)
			if err != nil {
				http.Error(w, "Failed to get track", http.StatusInternalServerError)
				return
			}

			err = db.AddPartialTrack(dbConn, track)
			if err != nil {
				http.Error(w, "Failed to add partial track", http.StatusInternalServerError)
				return
			}
		}

		err = db.AddTrackToPlaylist(dbConn, data.PlaylistID, ID)
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
		log.Println(err.Error())
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
	token := r.URL.Query().Get("token")

	if idStr == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}

	if token == "" {
		http.Error(w, "Missing authentication token", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
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

	filePath, err := play(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userId, err := db.GetUserID(dbConn, token)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = db.AddToListen(dbConn, userId, id)

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

	results, err := qobuz.Search(query)
	if err != nil {
		http.Error(w, "Invalid search query", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results.ToJSON())
}

func getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Missing album ID", http.StatusBadRequest)
		return
	}

	album, err := qobuz.GetAlbum(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(album.ToJSON())
}

func play(id int) (string, error) {
	err := checkAndAddTrack(id)
	if err != nil {
		return "", fmt.Errorf("failed to check and add track: %w", err)
	}

	track, err := db.GetTrack(dbConn, id)
	if err != nil {
		return "", fmt.Errorf("failed to get track: %w", err)
	}

	return track.Path, nil
}

func constitutePlaylist(playlist *db.Playlist, tracks []int) error {
	for _, trackID := range tracks {
		track, err := db.GetTrack(dbConn, trackID)
		if err != nil {
			return fmt.Errorf("failed to get track with ID %d: %w", trackID, err)
		}
		playlist.Tracks = append(playlist.Tracks, *track)
	}
	return nil
}

func downloadQobuzTrack(id int) (string, string, error) {
	file_name := utils.RandomString(50)
	file_path := filepath.Join(SONG_FOLDER, file_name)

	err := qobuz.Download(id, MAX_QUALITY, file_path)
	if err != nil {
		return "", "", errors.New("failed to cache track")
	}

	return file_path, file_name, nil
}

func checkAndAddTrack(trackID int) error {
	trackExists, needDownload := db.CheckIfTrackExists(dbConn, trackID)
	if !trackExists {
		qobuzTrack, err := qobuz.GetTrack(trackID)
		if err != nil {
			return fmt.Errorf("failed to get track with ID %d: %w", trackID, err)
		}

		if MAX_QUALITY == "6" {
			qobuzTrack.SampleRate = 44.1
			qobuzTrack.Bitrate = 16
		}

		file_name, file_path, err := downloadQobuzTrack(trackID)
		if err != nil {
			return fmt.Errorf("failed to download track: %w", err)
		}

		qobuzTrack.Path = file_path
		qobuzTrack.Filename = file_name

		err = db.AddTrack(dbConn, qobuzTrack)
		if err != nil {
			return fmt.Errorf("failed to add track to database: %w", err)
		}

	} else if trackExists && needDownload {
		file_path, file_name, err := downloadQobuzTrack(trackID)
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
