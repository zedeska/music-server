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

func InitDB() {
	if _, err := os.Stat("./db.db"); err != nil {
		os.Create("db.db")
	}

	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err) // Proper error logging
	}
	defer db.Close() // Ensure db is closed at the end

	_, err = db.Exec(`
		BEGIN;

		CREATE TABLE IF NOT EXISTS user (
			id_user INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(20),
			password TEXT,
			token VARCHAR(50)
		);

		CREATE TABLE IF NOT EXISTS track (
			id INTEGER PRIMARY KEY,
			title VARCHAR(100),
			path TEXT,
			filename VARCHAR(50),
			artist VARCHAR(100),
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
			FOREIGN KEY (id_user) REFERENCES user(id)
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

func CheckIfTrackExists(id int) bool {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	var track Track
	err = db.QueryRow("SELECT * FROM track WHERE id = ?", id).Scan(&track.ID, &track.Title, &track.Path, &track.Filename, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.SampleRate, &track.Bitrate)
	if track.ID == 0 || err != nil {
		return false
	} else {
		_, err := os.Stat(track.Path)
		if errors.Is(err, os.ErrNotExist) {
			db.Exec("DELETE FROM track WHERE id = ?", track.ID)
			return false
		}
	}

	return true
}

func GetTrack(id int) (*Track, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var track Track
	err = db.QueryRow("SELECT * FROM track WHERE id = ?", id).Scan(&track.ID, &track.Title, &track.Path, &track.Filename, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.SampleRate, &track.Bitrate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &track, nil
}

func AddTrack(track Track) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO track (id, title, path, filename, artist, album, year, duration, cover, sample_rate, bitrate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		track.ID, track.Title, track.Path, track.Filename, track.Artist, track.Album, track.Year, track.Duration, track.Cover, track.SampleRate, track.Bitrate)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	return nil
}

func Login(username, password string) (string, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var storedPassword string
	var token string
	err = db.QueryRow("SELECT password, token FROM user WHERE username = ?", username).Scan(&storedPassword, &token)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %w", err)
	}

	if !utils.VerifyPassword(password, storedPassword) {
		return "", errors.New("invalid username or password")
	}

	return token, nil
}

func Register(username, password string) (string, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return "", fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var exists int
	// Check if username already exists	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE username = ?", username).Scan(&exists)
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

func CheckToken(token string) (bool, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return false, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM user WHERE token = ?", token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	return exists > 0, nil
}

func AddToListen(userID, trackID int) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO listened (id_user, id_track, timestamp) VALUES (?, ?, ?)", userID, trackID, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to insert listened record: %w", err)
	}

	return nil
}

func GetUserID(token string) (int, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return 0, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	var userID int
	err = db.QueryRow("SELECT id_user FROM user WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return userID, nil
}

func GetListenedTracks(userID, limit int) ([]byte, error) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT t.id, t.title, t.path, t.artist, t.album, t.year, t.duration, t.cover, t.sample_rate, t.bitrate FROM listened l JOIN track t ON l.id_track = t.id WHERE l.id_user = ? ORDER BY l.timestamp DESC LIMIT ?", userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var track Track
		if err := rows.Scan(&track.ID, &track.Title, &track.Path, &track.Artist, &track.Album, &track.Year, &track.Duration, &track.Cover, &track.SampleRate, &track.Bitrate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		tracks = append(tracks, track)
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
