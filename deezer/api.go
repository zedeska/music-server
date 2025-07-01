package deezer

import (
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://api.deezer.com/"

func Search(query string) *deezer_search_result {
	request := goaxios.GoAxios{
		Url: API_URL + "search",
		Query: map[string]string{
			"q":     strings.Replace(query, " ", "%20", -1),
			"limit": "10",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
		ResponseStruct: &deezer_search_result{},
	}

	res := request.RunRest()

	data, _ := res.Body.(*deezer_search_result)

	return data
}
