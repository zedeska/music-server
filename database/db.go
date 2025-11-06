package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"music-server/utils"
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
			artist VARCHAR(100),
			album VARCHAR(100),
			year INTEGER,
			duration INTEGER,
			cover TEXT,
			artistqobuz INTEGER,
			artistdeezer INTEGER
		);

		CREATE TABLE IF NOT EXISTS quality (
			id INTEGER,
			path TEXT,
			bitrate INTEGER,
			sample_rate FLOAT,
			FOREIGN KEY (id) REFERENCES track(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS playlist (
			id_playlist INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user INTEGER,
			name VARCHAR(50),
			FOREIGN KEY (id_user) REFERENCES user(id_user)
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

func CheckIfTrackExists(db *sql.DB, id int, platform string, quality ...utils.QualityLevel) (bool, bool, bool) {
	var trackID int
	var artistID int
	var trackExists, needDownload, artistExists bool

	err := db.QueryRow(fmt.Sprintf("SELECT id, IFNULL(%s, 0) FROM track WHERE %s = ?", "artist"+platform, "id"+platform), id).Scan(&trackID, &artistID)
	if trackID == 0 || err == sql.ErrNoRows {
		trackExists = false
	} else {
		trackExists = true
	}

	if artistID != 0 {
		artistExists = true
	} else {
		artistExists = false
	}

	if len(quality) > 0 {
		var path string
		err := db.QueryRow("SELECT IFNULL(path, '') FROM quality WHERE id = ? AND bitrate = ?", trackID, quality[0].Bitrate).Scan(&path)
		if err == sql.ErrNoRows || path == "" {
			needDownload = true
		} else {
			needDownload = false
		}
	} else {
		needDownload = true
	}

	return trackExists, needDownload, artistExists
}

func CheckIfTrackExistsByArtistAndAlbum(db *sql.DB, id int, platform string, artist string, album string, track_title string, quality ...utils.QualityLevel) (bool, bool) {
	var trackId int
	err := db.QueryRow("SELECT id FROM track WHERE artist = ? AND album = ? AND title = ?", artist, album, track_title).Scan(&trackId)
	if trackId == 0 || err == sql.ErrNoRows {
		return false, true
	} else if len(quality) > 0 {
		db.Exec(fmt.Sprintf("UPDATE track SET %s = ? WHERE id = ?", "id"+platform), id, trackId)
		var path string
		err := db.QueryRow("SELECT IFNULL(path, '') FROM quality WHERE id = ? AND bitrate = ?", trackId, quality[0].Bitrate).Scan(&path)
		if err == sql.ErrNoRows || path == "" {
			return true, true
		}
	}
	return true, false
}

func GetTrack(db *sql.DB, id int, platformName string, quality ...utils.QualityLevel) (*Track, error) {
	var track Track
	var idQobuz, idDeezer int
	var artistQobuz, artistDeezer int

	if platformName == "" {
		err := db.QueryRow("SELECT id, IFNULL(idqobuz, 0), IFNULL(iddeezer, 0), title, artist, album, year, duration, cover, IFNULL(artistqobuz, 0), IFNULL(artistdeezer, 0) FROM track WHERE id = ?", id).Scan(&track.ID, &idQobuz, &idDeezer, &track.Title, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &artistQobuz, &artistDeezer)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		if idQobuz != 0 {
			track.ID = idQobuz
			track.Platform = "qobuz"
			track.ArtistID = artistQobuz
		} else if idDeezer != 0 {
			track.ID = idDeezer
			track.Platform = "deezer"
			track.ArtistID = artistDeezer
		}
	} else {
		err := db.QueryRow(fmt.Sprintf("SELECT id, title, artist, album, year, duration, cover, %s FROM track WHERE %s = ?", "artist"+platformName, "id"+platformName), id).Scan(&track.ID, &track.Title, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.ArtistID)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	}
	if len(quality) > 0 {
		err := db.QueryRow("SELECT path, bitrate, sample_rate FROM quality WHERE id = ? AND bitrate = ?", track.ID, quality[0].Bitrate).Scan(&track.Path, &track.Bitrate, &track.SampleRate)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	}
	return &track, nil
}

func AddTrack(db *sql.DB, track Track) error {
	res, err := db.Exec(fmt.Sprintf("INSERT INTO track (%s, title, artist, album, year, duration, cover, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", "id"+track.Platform, "artist"+track.Platform),
		track.ID, track.Title, track.Artist, track.Album, track.Year, track.Duration, track.Cover, track.ArtistID)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	_, err = db.Exec("INSERT INTO quality (id, path, bitrate, sample_rate) VALUES (?, ?, ?, ?)",
		id, track.Path, track.Bitrate, track.SampleRate)
	if err != nil {
		return fmt.Errorf("failed to insert quality: %w", err)
	}

	return nil
}

func AddPartialTrack(db *sql.DB, track Track) error {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO track (%s, title, artist, album, year, duration, cover, %s) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", "id"+track.Platform, "artist"+track.Platform),
		track.ID, track.Title, track.Artist, track.Album, track.Year, track.Duration, track.Cover, track.ArtistID)
	if err != nil {
		return fmt.Errorf("failed to insert partial track: %w", err)
	}

	return nil
}

func UpdateTrackPathAndFilename(db *sql.DB, id int, platform string, quality utils.QualityLevel, file_path string) error {
	var TrackId int
	err := db.QueryRow(fmt.Sprintf("SELECT id FROM track WHERE %s = ?", "id"+platform), id).Scan(&TrackId)
	if err != nil {
		return fmt.Errorf("failed to update track: %w", err)
	}
	_, err = db.Exec("INSERT INTO quality (id, path, bitrate, sample_rate) VALUES (?, ?, ?, ?)", TrackId, file_path, quality.Bitrate, quality.SampleRate)
	if err != nil {
		return fmt.Errorf("failed to insert quality: %w", err)
	}
	return nil
}

func UpdateTrackArtist(db *sql.DB, id int, platform string, artistID int) error {
	_, err := db.Exec(fmt.Sprintf("UPDATE track SET %s = ? WHERE %s = ?", "artist"+platform, "id"+platform), artistID, id)
	if err != nil {
		return fmt.Errorf("failed to update track artist: %w", err)
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
	rows, err := db.Query("SELECT DISTINCT(t.id), IFNULL(t.idqobuz, 0), IFNULL(t.iddeezer, 0), t.title, t.artist, t.album, t.year, t.duration, t.cover, IFNULL(t.artistqobuz, 0), IFNULL(t.artistdeezer, 0) FROM listened l JOIN track t ON l.id_track = t.id WHERE l.id_user = ? ORDER BY l.timestamp DESC LIMIT ?", userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tracks struct {
		Data []Track `json:"data"`
	}

	for rows.Next() {
		var track Track
		var idQobuz, idDeezer int
		var artistQobuz, artistDeezer int
		if err := rows.Scan(&track.ID, &idQobuz, &idDeezer, &track.Title, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &artistQobuz, &artistDeezer); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if idQobuz != 0 {
			track.ID = idQobuz
			track.Platform = "qobuz"
			track.ArtistID = artistQobuz
		} else if idDeezer != 0 {
			track.ID = idDeezer
			track.Platform = "deezer"
			track.ArtistID = artistDeezer
		} else {
			continue
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
