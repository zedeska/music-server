package musicbrainz

import "time"

type musicbrainz_search_result struct {
	Created    time.Time `json:"created"`
	Count      int       `json:"count"`
	Offset     int       `json:"offset"`
	Recordings []struct {
		ID             string `json:"id"`
		Score          int    `json:"score"`
		ArtistCreditID string `json:"artist-credit-id"`
		Title          string `json:"title"`
		Length         int    `json:"length"`
		Video          any    `json:"video"`
		ArtistCredit   []struct {
			Joinphrase string `json:"joinphrase,omitempty"`
			Name       string `json:"name"`
			Artist     struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				SortName       string `json:"sort-name"`
				Disambiguation string `json:"disambiguation"`
				Aliases        []struct {
					SortName  string `json:"sort-name"`
					Name      string `json:"name"`
					Locale    any    `json:"locale"`
					Type      any    `json:"type"`
					Primary   any    `json:"primary"`
					BeginDate any    `json:"begin-date"`
					EndDate   any    `json:"end-date"`
					TypeID    string `json:"type-id,omitempty"`
				} `json:"aliases"`
			} `json:"artist,omitempty"`
		} `json:"artist-credit"`
		FirstReleaseDate string `json:"first-release-date"`
		Releases         []struct {
			ID             string `json:"id"`
			StatusID       string `json:"status-id"`
			ArtistCreditID string `json:"artist-credit-id"`
			Count          int    `json:"count"`
			Title          string `json:"title"`
			Status         string `json:"status"`
			ArtistCredit   []struct {
				Name   string `json:"name"`
				Artist struct {
					ID             string `json:"id"`
					Name           string `json:"name"`
					SortName       string `json:"sort-name"`
					Disambiguation string `json:"disambiguation"`
				} `json:"artist"`
			} `json:"artist-credit"`
			ReleaseGroup struct {
				ID               string   `json:"id"`
				TypeID           string   `json:"type-id"`
				PrimaryTypeID    string   `json:"primary-type-id"`
				Title            string   `json:"title"`
				PrimaryType      string   `json:"primary-type"`
				SecondaryTypes   []string `json:"secondary-types"`
				SecondaryTypeIds []string `json:"secondary-type-ids"`
			} `json:"release-group"`
			Date          string `json:"date"`
			Country       string `json:"country"`
			ReleaseEvents []struct {
				Date string `json:"date"`
				Area struct {
					ID            string   `json:"id"`
					Name          string   `json:"name"`
					SortName      string   `json:"sort-name"`
					Iso31661Codes []string `json:"iso-3166-1-codes"`
				} `json:"area"`
			} `json:"release-events"`
			TrackCount int `json:"track-count"`
			Media      []struct {
				ID       string `json:"id"`
				Position int    `json:"position"`
				Format   string `json:"format"`
				Track    []struct {
					ID     string `json:"id"`
					Number string `json:"number"`
					Title  string `json:"title"`
					Length int    `json:"length"`
				} `json:"track"`
				TrackCount  int `json:"track-count"`
				TrackOffset int `json:"track-offset"`
			} `json:"media"`
		} `json:"releases"`
	} `json:"recordings"`
}
