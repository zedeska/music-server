package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"music-server/utils"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(db *sql.DB) {
	_, err := db.Exec(`
		BEGIN;

		CREATE TABLE IF NOT EXISTS user (
			id_user INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(20),
			password TEXT,
			token VARCHAR(50)
		);

		CREATE TABLE IF NOT EXISTS track (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			idqobuz INTEGER,
			iddeezer INTEGER,
			title VARCHAR(100),
			path TEXT,
			filename VARCHAR(50),
			artist VARCHAR(100),
			artist_id INTEGER,
			album VARCHAR(100),
			year INTEGER,
			duration INTEGER,
			cover TEXT,
			sample_rate INTEGER,
			bitrate INTEGER
		);

		CREATE TABLE IF NOT EXISTS playlist (
			id_playlist INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user INTEGER,
			name VARCHAR(50),
			FOREIGN KEY (id_user) REFERENCES user(id_user)
		);

		CREATE TABLE IF NOT EXISTS playlist (
			id_playlist INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user INTEGER,
			name VARCHAR(50),
			FOREIGN KEY (id_user) REFERENCES user(id)
		);

		CREATE TABLE IF NOT EXISTS in_playlist (
			id_playlist INTEGER,
			id_track INTEGER,
			FOREIGN KEY (id_playlist) REFERENCES playlist(id_playlist) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS listened (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user INTEGER,
			id_track INTEGER,
			timestamp INTEGER,
			FOREIGN KEY (id_user) REFERENCES user(id_user),
			FOREIGN KEY (id_track) REFERENCES track(id)
		);

		COMMIT;
	`)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	fmt.Println("Database initialized successfully")
}

func CheckIfTrackExists(db *sql.DB, id int, platform string) (bool, bool) {
	var track Track

	err := db.QueryRow(fmt.Sprintf("SELECT id, IFNULL(path, '') FROM track WHERE %s = ?", "id"+platform), id).Scan(&track.ID, &track.Path)
	if track.ID == 0 || err == sql.ErrNoRows {
		return false, true
	} else {
		if track.Path == "" {
			return true, true
		}
		_, err := os.Stat(track.Path)
		if errors.Is(err, os.ErrNotExist) {
			db.Exec("DELETE FROM track WHERE id = ?", track.ID)
			return false, true
		}
	}
	return true, false
}

func CheckIfTrackExistsByArtistAndAlbum(db *sql.DB, id int, platform string, artist string, album string, track_title string) (bool, bool) {
	var track Track
	err := db.QueryRow("SELECT id, IFNULL(path, '') FROM track WHERE artist = ? AND album = ? AND title = ?", artist, album, track_title).Scan(&track.ID, &track.Path)
	if track.ID == 0 || err == sql.ErrNoRows {
		return false, true
	} else {
		db.Exec(fmt.Sprintf("UPDATE track SET %s = ? WHERE id = ?", "id"+platform), id, track.ID)
		if track.Path == "" {
			return true, true
		} else {
			_, err := os.Stat(track.Path)
			if errors.Is(err, os.ErrNotExist) {
				db.Exec("DELETE FROM track WHERE id = ?", track.ID)
				return false, true
			}
		}
	}
	return true, false
}

func GetTrack(db *sql.DB, id int, platformName string) (*Track, error) {
	var track Track
	err := db.QueryRow(fmt.Sprintf("SELECT id, title, IFNULL(path, ''), IFNULL(filename, ''), artist, artist_id, album, year, duration, cover, sample_rate, bitrate FROM track WHERE %s = ?", "id"+platformName), id).Scan(&track.ID, &track.Title, &track.Path, &track.Filename, &track.Artist, &track.ArtistID, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.SampleRate, &track.Bitrate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &track, nil
}

func AddTrack(db *sql.DB, track Track) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO track (%s, title, path, filename, artist, artist_id, album, year, duration, cover, sample_rate, bitrate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", "id"+track.Platform),
		track.ID, track.Title, track.Path, track.Filename, track.Artist, track.ArtistID, track.Album, track.Year, track.Duration, track.Cover, track.SampleRate, track.Bitrate)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	return nil
}

func AddPartialTrack(db *sql.DB, track Track) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO track (%s, title, artist, artist_id, album, year, duration, cover, sample_rate, bitrate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", "id"+track.Platform),
		track.ID, track.Title, track.Artist, track.ArtistID, track.Album, track.Year, track.Duration, track.Cover, track.SampleRate, track.Bitrate)
	if err != nil {
		return fmt.Errorf("failed to insert partial track: %w", err)
	}

	return nil
}

func UpdateTrackPathAndFilename(db *sql.DB, id int, file_path string, file_name string) error {
	_, err := db.Exec("UPDATE track SET filename = ?, path = ? WHERE id = ?", file_name, file_path, id)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}

	return nil
}

func Login(db *sql.DB, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	var storedPassword string
	var token string
	err := db.QueryRow("SELECT password, token FROM user WHERE username = ?", username).Scan(&storedPassword, &token)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	if !utils.VerifyPassword(password, storedPassword) {
		return "", errors.New("invalid username or password")
	}

	return token, nil
}

func Register(db *sql.DB, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}
	var exists int
	// Check if username already exists
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	if exists > 0 {
		return "", errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	token := utils.RandomString(50)

	_, err = db.Exec("INSERT INTO user (username, password, token) VALUES (?, ?, ?)", username, hashedPassword, token)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	return token, nil
}

func CheckToken(db *sql.DB, token string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE token = ?", token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return exists > 0, nil
}

func AddToListen(db *sql.DB, userID, trackID int) error {
	_, err := db.Exec("INSERT INTO listened (id_user, id_track, timestamp) VALUES (?, ?, ?)", userID, trackID, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to insert listened record: %w", err)
	}

	return nil
}

func GetUserID(db *sql.DB, token string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id_user FROM user WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return userID, nil
}

func GetListenedTracks(db *sql.DB, userID, limit int) ([]byte, error) {
	rows, err := db.Query("SELECT DISTINCT(t.id), t.title, t.artist, t.artist_id, t.album, t.year, t.duration, t.cover, t.sample_rate, t.bitrate FROM listened l JOIN track t ON l.id_track = t.id WHERE l.id_user = ? ORDER BY l.timestamp DESC LIMIT ?", userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tracks struct {
		Data []Track `json:"data"`
	}

	for rows.Next() {
		var track Track
		if err := rows.Scan(&track.ID, &track.Title, &track.Artist, &track.ArtistID, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.SampleRate, &track.Bitrate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tracks.Data = append(tracks.Data, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	data, err := json.Marshal(tracks)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tracks: %w", err)
	}

	return data, nil
}

func GetPlaylistsByUserID(db *sql.DB, userID int) (*Playlists, error) {
	rows, err := db.Query("SELECT id_playlist, name FROM playlist WHERE id_user = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var playlists Playlists

	for rows.Next() {
		var playlist Playlist
		if err := rows.Scan(&playlist.ID, &playlist.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		playlists.Playlists = append(playlists.Playlists, playlist)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return &playlists, nil
}

func GetPlaylistTracks(db *sql.DB, playlistID int) ([]int, error) {
	rows, err := db.Query("SELECT id_track FROM in_playlist WHERE id_playlist = ?", playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var trackIDs []int
	for rows.Next() {
		var trackID int
		if err := rows.Scan(&trackID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		trackIDs = append(trackIDs, trackID)
	}

	return trackIDs, nil
}

func GetPlaylistByID(db *sql.DB, playlistID int) (*Playlist, error) {
	var playlist Playlist
	err := db.QueryRow("SELECT id_playlist, name FROM playlist WHERE id_playlist = ?", playlistID).Scan(&playlist.ID, &playlist.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &playlist, nil
}

func AddTrackToPlaylist(db *sql.DB, playlistID, trackID int) error {
	_, err := db.Exec("INSERT INTO in_playlist (id_playlist, id_track) VALUES (?, ?)", playlistID, trackID)
	if err != nil {
		return fmt.Errorf("failed to insert track into playlist: %w", err)
	}

	return nil
}

func CreatePlaylist(db *sql.DB, userID int, name string) (int, error) {
	if name == "" {
		return 0, errors.New("playlist name cannot be empty")
	}

	result, err := db.Exec("INSERT INTO playlist (id_user, name) VALUES (?, ?)", userID, name)
	if err != nil {
		return 0, fmt.Errorf("failed to create playlist: %w", err)
	}

	playlistID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return int(playlistID), nil
}

func IsTrackInPlaylist(db *sql.DB, playlistID, trackID int) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM in_playlist WHERE id_playlist = ? AND id_track = ?", playlistID, trackID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return exists > 0, nil
}

func IsPlaylistOwnedByUser(db *sql.DB, userID int, playlistID int) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM playlist WHERE id_user = ? AND id_playlist = ?", userID, playlistID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return exists > 0, nil
}

func DeletePlaylist(db *sql.DB, playlistID int) error {
	_, err := db.Exec("DELETE FROM playlist WHERE id_playlist = ?", playlistID)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	_, err = db.Exec("DELETE FROM in_playlist WHERE id_playlist = ?", playlistID)
	if err != nil {
		return fmt.Errorf("failed to delete tracks from playlist: %w", err)
	}

	return nil
}

func DeleteTrackFromPlaylist(db *sql.DB, playlistID, trackID int) error {
	_, err := db.Exec("DELETE FROM in_playlist WHERE id_playlist = ? AND id_track = ?", playlistID, trackID)
	if err != nil {
		return fmt.Errorf("failed to delete track from playlist: %w", err)
	}

	return nil
}

func GetTrackIds(db *sql.DB, trackID int) ([]int, error) {
	var qobuzID, deezerID int
	err := db.QueryRow("SELECT IFNULL(idqobuz, 0), IFNULL(iddeezer, 0) FROM track WHERE id = ?", trackID).Scan(&qobuzID, &deezerID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	trackIDs := []int{qobuzID, deezerID}
	return trackIDs, nil
}
