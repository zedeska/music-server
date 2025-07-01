package qobuz

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://www.qobuz.com/api.json/0.2/"

func Search(query string) (*qobuz_search_result, error) {

	request := goaxios.GoAxios{
		Url: API_URL + "catalog/search",
		Query: map[string]string{
			"limit":  "10",
			"offset": "0",
			"query":  strings.Replace(query, " ", "%20", -1),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":   "798273057",
			"x-zone":     "US",
			"x-store":    "US-en",
		},
		ResponseStruct: &qobuz_search_result{},
	}

	res := request.RunRest()

	if res.Error != nil {
		fmt.Println("Error:", res.Error)
		return nil, res.Error
	}

	result, _ := res.Body.(*qobuz_search_result)

	return result, nil

}

func GetTrack(id int) *QobuzTrack {
	request := goaxios.GoAxios{
		Url: API_URL + "track/get",
		Query: map[string]string{
			"track_id": strconv.Itoa(id),
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
			"X-App-Id":   "798273057",
		},
		ResponseStruct: &QobuzTrack{},
	}

	res := request.RunRest()

	result, _ := res.Body.(*QobuzTrack)

	return result
}
