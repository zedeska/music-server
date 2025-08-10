package qobuz

import (
	"fmt"
	"music-server/utils"
	"strconv"

	"github.com/opensaucerer/goaxios"
)

var DAB_API_URL string = "https://dab.yeet.su/api/"
var SQUID_API_EU_URL string = "https://eu.qobuz.squid.wtf/api/"
var SQUID_API_US_URL string = "https://us.qobuz.squid.wtf/api/"

type DabDLResult struct {
	Url string `json:"url"`
}

type SquidDLResult struct {
	Success bool `json:"success"`
	Data    struct {
		Url string `json:"url"`
	} `json:"data"`
}

func dabDownload(id int, quality string, path string) error {
	request := goaxios.GoAxios{
		Url: DAB_API_URL + "stream",
		Query: map[string]string{
			"quality": quality,
			"trackId": strconv.Itoa(id),
		},
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		Method:         "GET",
		ResponseStruct: &DabDLResult{},
	}

	res := request.RunRest()

	if res.Error != nil || res.Response.StatusCode != 200 {
		return fmt.Errorf("error fetching download URL: %w", res.Error)
	}

	result, _ := res.Body.(*DabDLResult)

	err := utils.DownloadAndCheckTime(path, result.Url)
	if err != nil {
		return fmt.Errorf("error downloading and checking file: %w", err)
	}

	return nil
}

func squidDownload(id int, quality string, path string, api string) error {
	request := goaxios.GoAxios{
		Url: api + "download-music",
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

	if res.Error != nil || res.Response.StatusCode != 200 {
		return fmt.Errorf("error fetching download URL: %w", res.Error)
	}

	result, _ := res.Body.(*SquidDLResult)

	if !result.Success {
		return fmt.Errorf("failed to fetch download URL: %s", string(res.Bytes))
	}

	err := utils.DownloadAndCheckTime(path, result.Data.Url)
	if err != nil {
		return fmt.Errorf("error downloading and checking file: %w", err)
	}

	return nil
}

func Download(id int, quality string, path string) error {

	err := squidDownload(id, quality, path, SQUID_API_EU_URL)
	if err != nil {
		err := squidDownload(id, quality, path, SQUID_API_US_URL)
		if err != nil {
			err = dabDownload(id, quality, path)
			if err != nil {
				return fmt.Errorf("error downloading track: %w", err)
			}
		}
	}

	return nil
}
