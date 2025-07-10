package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"music-server/utils"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
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
			db.QueryRow("DELETE FROM track WHERE id = ?", track.ID)
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

	if !VerifyPassword(password, storedPassword) {
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

	hashedPassword, err := HashPassword(password)
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

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
