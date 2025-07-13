package qobuz

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/opensaucerer/goaxios"
)

var DAB_API_URL string = "https://dab.yeet.su/api/"
var SQUID_API_URL string = "https://eu.qobuz.squid.wtf/api/"

type DabDLResult struct {
	Url string `json:"url"`
}

type SquidDLResult struct {
	Success bool `json:"success"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func dabDownload(id int, quality string) (string, error) {
	request := goaxios.GoAxios{
		Url: DAB_API_URL + "stream",
		Query: map[string]string{
			"quality": quality,
			"trackId": strconv.Itoa(id),
		},
		Method:         "GET",
		ResponseStruct: &DabDLResult{},
	}

	res := request.RunRest()

	fmt.Println("Response Status Code:", res.Response.StatusCode)

	if res.Error != nil || res.Response.StatusCode != 200 {
		return "", fmt.Errorf("error fetching download URL: %w", res.Error)
	}

	result, _ := res.Body.(*DabDLResult)
	return result.Url, nil

}

func squidDownload(id int, quality string) (string, error) {
	request := goaxios.GoAxios{
		Url: SQUID_API_URL + "download-music",
		Query: map[string]string{
			"quality":  quality,
			"track_id": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &SquidDLResult{},
	}

	res := request.RunRest()

	fmt.Println("Response Status Code:", res.Response.StatusCode)

	if res.Error != nil || res.Response.StatusCode != 200 {
		return "", fmt.Errorf("error fetching download URL: %w", res.Error)
	}

	result, _ := res.Body.(*SquidDLResult)

	if len(result.Data) == 0 {
		return "", fmt.Errorf("no download URL found")
	}

	return result.Data[0].Url, nil
}

func Download(id int, quality string, path string) error {

	var url string
	var err error

	url, err = squidDownload(id, quality)
	if err != nil {
		url, err = dabDownload(id, quality)
		if err != nil {
			return fmt.Errorf("error downloading track: %w", err)
		}
	}

	request := goaxios.GoAxios{
		Url:    url,
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
	}

	res := request.RunRest()
	if res.Error != nil {
		return res.Error
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bytes.NewReader(res.Bytes)

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
