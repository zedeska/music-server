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
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func dabDownload(id int) (string, error) {
	request := goaxios.GoAxios{
		Url: DAB_API_URL + "stream",
		Query: map[string]string{
			"trackId": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &DabDLResult{},
	}

	res := request.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return "", fmt.Errorf("error fetching download URL: %w", res.Error)
	}

	result, _ := res.Body.(*DabDLResult)
	return result.Url, nil

}

func squidDownload(id int) (string, error) {
	request := goaxios.GoAxios{
		Url: SQUID_API_URL + "download-music",
		Query: map[string]string{
			"quality":  "27",
			"track_id": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &SquidDLResult{},
	}

	res := request.RunRest()
	if res.Error != nil || res.Response.StatusCode != 200 {
		return "", fmt.Errorf("error fetching squid download URL: %w", res.Error)
	}

	result, _ := res.Body.(*SquidDLResult)
	return result.Data[0].Url, nil
}

func Download(id int, path string) error {

	url, err := squidDownload(id)
	if err != nil {
		url, err = dabDownload(id)
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
