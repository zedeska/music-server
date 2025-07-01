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

type DabSearchResult struct {
	Url string `json:"url"`
}

func dabDownload(id int) string {
	request := goaxios.GoAxios{
		Url: DAB_API_URL + "stream",
		Query: map[string]string{
			"trackId": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &DabSearchResult{},
	}

	res := request.RunRest()
	fmt.Println(res.Response.Request.URL)
	result, _ := res.Body.(*DabSearchResult)
	return result.Url

}

func Download(id int, path string) error {
	url := dabDownload(id)

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
