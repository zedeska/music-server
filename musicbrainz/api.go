package musicbrainz

import (
	"fmt"
	"strings"

	"github.com/opensaucerer/goaxios"
)

var API_URL string = "https://musicbrainz.org/ws/2/"

func Search(query string) {
	request := goaxios.GoAxios{
		Url: API_URL + "recording",
		Query: map[string]string{
			"query": strings.Replace(query, " ", "%20", -1),
			"fmt":   "json",
			"limit": "10",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
	}

	res := request.RunRest()
	fmt.Println(res.Body)
}

func GetTrackInfo(id string) {
	request := goaxios.GoAxios{
		Url: API_URL + "recording/" + id,
		Query: map[string]string{
			"fmt": "json",
		},
		Method: "GET",
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:139.0) Gecko/20100101 Firefox/139.0",
		},
	}

	res := request.RunRest()
	fmt.Println(res.Body)
}
