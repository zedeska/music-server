package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
			username VARCHAR(50),
			password TEXT,
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
			bitrate INTEGER,
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

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM track WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	return exists
}

func AddTrack(track Track) error {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO track (id, title, path, artist, album, year, duration, cover, sample_rate, bitrate) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		track.ID, track.Title, track.Path, track.Artist, track.Album, track.Year, track.Duration, track.Cover, track.SampleRate, track.Bitrate)
	if err != nil {
		return fmt.Errorf("failed to insert track: %w", err)
	}

	return nil
}
