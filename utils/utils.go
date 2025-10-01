package utils

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-flac/go-flac/v2"
	"github.com/opensaucerer/goaxios"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomChar returns a random alphanumeric character
func RandomChar() byte {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return chars[rand.Intn(len(chars))]
}

// RandomString returns a random string of specified length
func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = RandomChar()
	}
	return string(result)
}

// RandomCharSet returns a random character from a custom character set
func RandomCharSet(charset string) byte {
	return charset[rand.Intn(len(charset))]
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

func getFlacFileDuration(path string) (int64, error) {
	f, err := flac.ParseFile(path)
	if err != nil {
		os.Remove(path)
		return 0, fmt.Errorf("error parsing FLAC file: %w", err)
	}
	defer f.Close()

	data, err := f.GetStreamInfo()
	if err != nil {
		os.Remove(path)
		return 0, fmt.Errorf("error getting stream info: %w", err)
	}

	return data.SampleCount / int64(data.SampleRate), nil
}

func downloadFile(path string, url string) error {
	request := goaxios.GoAxios{
		Url:    url,
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		IsDownload: true,
		DownloadDestination: goaxios.Destination{
			Location: path,
		},
	}

	res := request.RunRest()
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func DownloadAndCheckTime(path string, url string) error {
	err := downloadFile(path, url)
	if err != nil {
		return fmt.Errorf("error downloading file: %w", err)
	}

	time, err := getFlacFileDuration(path)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("error getting file duration: %w", err)
	}

	if time == 30 {
		os.Remove(path)
		return fmt.Errorf("file duration is too short: %d seconds", time)
	}

	return nil
}

func GetPlatformName(platform int) (string, error) {
	switch platform {
	case 0:
		return "", nil
	case 1:
		return "qobuz", nil
	case 2:
		return "deezer", nil
	default:
		return "Unknown", fmt.Errorf("unknown platform: %d", platform)
	}
}
