package main

import (
	"encoding/json"
	"fmt"
	db "music-server/database"
	"music-server/qobuz"
	"music-server/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var SONG_FOLDER string

func main() {

	db.InitDB()
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
	http.ListenAndServe(":8488", nil)

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

	tokenValid, err := db.CheckToken(token)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	token, err := db.Login(creds.Username, creds.Password)
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

	token, err := db.Register(creds.Username, creds.Password)
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

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	album, err := qobuz.GetAlbum(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(album.ToJSON())
}

func play(id int) (string, error) {
	if !db.CheckIfTrackExists(id) {
		qobuzTrack, err := qobuz.GetTrack(id)
		if err != nil {
			return "", err
		}

		file_name := utils.RandomString(50)
		file_path := filepath.Join(SONG_FOLDER, file_name)

		err = qobuz.Download(qobuzTrack.ID, "27", file_path)
		if err != nil {
			return "", err
		}

		qobuzTrack.Path = file_path
		qobuzTrack.Filename = file_name

		err = db.AddTrack(*qobuzTrack)
		if err != nil {
			return "", err
		}

		return file_path, nil
	} else {
		track, err := db.GetTrack(id)
		if err != nil {
			return "", fmt.Errorf("failed to get track: %w", err)
		}

		return track.Path, nil
	}
}
