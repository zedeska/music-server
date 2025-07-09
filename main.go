package main

import (
	"encoding/json"
	"errors"
	"fmt"
	db "music-server/database"
	"music-server/qobuz"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var SONG_FOLDER string

type custom_search_result struct {
	Tracks []db.Track `json:"tracks"`
}

func (p custom_search_result) ToJSON() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return data
}

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
	http.ListenAndServe(":8488", nil)

}

func playHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid track ID", http.StatusBadRequest)
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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	temp_results, err := qobuz.Search(query)
	if err != nil {
		http.Error(w, "Invalid search query", http.StatusBadRequest)
		return
	}

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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results.ToJSON())
}

func play(id int) (string, error) {
	if !db.CheckIfTrackExists(id) {
		qobuzTrack := qobuz.GetTrack(id)
		if qobuzTrack == nil {
			return "", errors.New("Track not found")
		}

		file_name := RandomString(50)
		file_path := filepath.Join(SONG_FOLDER, file_name)

		err := qobuz.Download(qobuzTrack.ID, "27", file_path)
		if err != nil {
			return "", err
		}

		year, _ := strconv.Atoi(strings.Split(qobuzTrack.ReleaseDateOriginal, "-")[0])

		track := &db.Track{
			ID:         qobuzTrack.ID,
			Title:      qobuzTrack.Title,
			Path:       file_path,
			Filename:   file_name,
			Year:       year,
			Artist:     qobuzTrack.Performer.Name,
			Album:      qobuzTrack.Album.Title,
			Duration:   qobuzTrack.Duration,
			Cover:      qobuzTrack.Album.Image.Large,
			Bitrate:    qobuzTrack.MaximumBitDepth,
			SampleRate: qobuzTrack.MaximumSamplingRate,
		}

		err = db.AddTrack(*track)
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
