package qobuz

import "encoding/json"

type qobuz_search_result struct {
	Query  string `json:"query"`
	Albums struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
		Items  []struct {
			MaximumBitDepth int `json:"maximum_bit_depth"`
			Image           struct {
				Small     string `json:"small"`
				Thumbnail string `json:"thumbnail"`
				Large     string `json:"large"`
				Back      any    `json:"back"`
			} `json:"image"`
			MediaCount int `json:"media_count"`
			Artist     struct {
				Image       any    `json:"image"`
				Name        string `json:"name"`
				ID          int    `json:"id"`
				AlbumsCount int    `json:"albums_count"`
				Slug        string `json:"slug"`
				Picture     any    `json:"picture"`
			} `json:"artist"`
			Artists []struct {
				ID    int      `json:"id"`
				Name  string   `json:"name"`
				Roles []string `json:"roles"`
			} `json:"artists"`
			Upc        string `json:"upc"`
			ReleasedAt int    `json:"released_at"`
			Label      struct {
				Name        string `json:"name"`
				ID          int    `json:"id"`
				AlbumsCount int    `json:"albums_count"`
				SupplierID  int    `json:"supplier_id"`
				Slug        string `json:"slug"`
			} `json:"label"`
			Title           string `json:"title"`
			QobuzID         int    `json:"qobuz_id"`
			Version         any    `json:"version"`
			URL             string `json:"url"`
			Duration        int    `json:"duration"`
			ParentalWarning bool   `json:"parental_warning"`
			Popularity      int    `json:"popularity"`
			TracksCount     int    `json:"tracks_count"`
			Genre           struct {
				Path  []int  `json:"path"`
				Color string `json:"color"`
				Name  string `json:"name"`
				ID    int    `json:"id"`
				Slug  string `json:"slug"`
			} `json:"genre"`
			MaximumChannelCount int     `json:"maximum_channel_count"`
			ID                  string  `json:"id"`
			MaximumSamplingRate float64 `json:"maximum_sampling_rate"`
			Articles            []any   `json:"articles"`
			ReleaseDateOriginal string  `json:"release_date_original"`
			ReleaseDateDownload string  `json:"release_date_download"`
			ReleaseDateStream   string  `json:"release_date_stream"`
			Purchasable         bool    `json:"purchasable"`
			Streamable          bool    `json:"streamable"`
			Previewable         bool    `json:"previewable"`
			Sampleable          bool    `json:"sampleable"`
			Downloadable        bool    `json:"downloadable"`
			Displayable         bool    `json:"displayable"`
			PurchasableAt       int     `json:"purchasable_at"`
			StreamableAt        int     `json:"streamable_at"`
			Hires               bool    `json:"hires"`
			HiresStreamable     bool    `json:"hires_streamable"`
		} `json:"items"`
	} `json:"albums"`
	Tracks struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
		Items  []struct {
			MaximumBitDepth int    `json:"maximum_bit_depth"`
			Copyright       string `json:"copyright"`
			Performers      string `json:"performers"`
			AudioInfo       struct {
				ReplaygainTrackPeak float64 `json:"replaygain_track_peak"`
				ReplaygainTrackGain float64 `json:"replaygain_track_gain"`
			} `json:"audio_info"`
			Performer struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"performer"`
			Album struct {
				Image struct {
					Small     string `json:"small"`
					Thumbnail string `json:"thumbnail"`
					Large     string `json:"large"`
				} `json:"image"`
				MaximumBitDepth int `json:"maximum_bit_depth"`
				MediaCount      int `json:"media_count"`
				Artist          struct {
					Image       any    `json:"image"`
					Name        string `json:"name"`
					ID          int    `json:"id"`
					AlbumsCount int    `json:"albums_count"`
					Slug        string `json:"slug"`
					Picture     any    `json:"picture"`
				} `json:"artist"`
				Upc        string `json:"upc"`
				ReleasedAt int    `json:"released_at"`
				Label      struct {
					Name        string `json:"name"`
					ID          int    `json:"id"`
					AlbumsCount int    `json:"albums_count"`
					SupplierID  int    `json:"supplier_id"`
					Slug        string `json:"slug"`
				} `json:"label"`
				Title           string `json:"title"`
				QobuzID         int    `json:"qobuz_id"`
				Version         any    `json:"version"`
				Duration        int    `json:"duration"`
				ParentalWarning bool   `json:"parental_warning"`
				TracksCount     int    `json:"tracks_count"`
				Popularity      int    `json:"popularity"`
				Genre           struct {
					Path  []int  `json:"path"`
					Color string `json:"color"`
					Name  string `json:"name"`
					ID    int    `json:"id"`
					Slug  string `json:"slug"`
				} `json:"genre"`
				MaximumChannelCount int     `json:"maximum_channel_count"`
				ID                  string  `json:"id"`
				MaximumSamplingRate float64 `json:"maximum_sampling_rate"`
				Previewable         bool    `json:"previewable"`
				Sampleable          bool    `json:"sampleable"`
				Displayable         bool    `json:"displayable"`
				Streamable          bool    `json:"streamable"`
				StreamableAt        int     `json:"streamable_at"`
				Downloadable        bool    `json:"downloadable"`
				PurchasableAt       any     `json:"purchasable_at"`
				Purchasable         bool    `json:"purchasable"`
				ReleaseDateOriginal string  `json:"release_date_original"`
				ReleaseDateDownload string  `json:"release_date_download"`
				ReleaseDateStream   string  `json:"release_date_stream"`
				ReleaseDatePurchase string  `json:"release_date_purchase"`
				Hires               bool    `json:"hires"`
				HiresStreamable     bool    `json:"hires_streamable"`
			} `json:"album"`
			Work     any `json:"work"`
			Composer struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"composer"`
			Isrc                string  `json:"isrc"`
			Title               string  `json:"title"`
			Version             any     `json:"version"`
			Duration            int     `json:"duration"`
			ParentalWarning     bool    `json:"parental_warning"`
			TrackNumber         int     `json:"track_number"`
			MaximumChannelCount int     `json:"maximum_channel_count"`
			ID                  int     `json:"id"`
			MediaNumber         int     `json:"media_number"`
			MaximumSamplingRate float64 `json:"maximum_sampling_rate"`
			ReleaseDateOriginal string  `json:"release_date_original"`
			ReleaseDateDownload string  `json:"release_date_download"`
			ReleaseDateStream   string  `json:"release_date_stream"`
			ReleaseDatePurchase string  `json:"release_date_purchase"`
			Purchasable         bool    `json:"purchasable"`
			Streamable          bool    `json:"streamable"`
			Previewable         bool    `json:"previewable"`
			Sampleable          bool    `json:"sampleable"`
			Downloadable        bool    `json:"downloadable"`
			Displayable         bool    `json:"displayable"`
			PurchasableAt       int     `json:"purchasable_at"`
			StreamableAt        int     `json:"streamable_at"`
			Hires               bool    `json:"hires"`
			HiresStreamable     bool    `json:"hires_streamable"`
		} `json:"items"`
	} `json:"tracks"`
	Artists struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
		Items  []struct {
			Picture string `json:"picture"`
			Image   struct {
				Small      string `json:"small"`
				Medium     string `json:"medium"`
				Large      string `json:"large"`
				Extralarge string `json:"extralarge"`
				Mega       string `json:"mega"`
			} `json:"image"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			AlbumsCount int    `json:"albums_count"`
			ID          int    `json:"id"`
		} `json:"items"`
	} `json:"artists"`
	Playlists struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
		Items  []struct {
			ImageRectangleMini []string `json:"image_rectangle_mini"`
			IsPublished        bool     `json:"is_published"`
			FeaturedArtists    []any    `json:"featured_artists"`
			Description        string   `json:"description"`
			CreatedAt          int      `json:"created_at"`
			TimestampPosition  int      `json:"timestamp_position"`
			Images300          []string `json:"images300"`
			Duration           int      `json:"duration"`
			UpdatedAt          int      `json:"updated_at"`
			PublishedTo        int      `json:"published_to"`
			Genres             []struct {
				ID      int     `json:"id"`
				Color   string  `json:"color"`
				Name    string  `json:"name"`
				Path    []int   `json:"path"`
				Slug    string  `json:"slug"`
				Percent float64 `json:"percent"`
			} `json:"genres"`
			ImageRectangle []string `json:"image_rectangle"`
			ID             int      `json:"id"`
			Slug           string   `json:"slug"`
			Owner          struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"owner"`
			UsersCount      int      `json:"users_count"`
			Images150       []string `json:"images150"`
			Images          []string `json:"images"`
			IsCollaborative bool     `json:"is_collaborative"`
			Stores          []string `json:"stores"`
			Tags            []struct {
				FeaturedTagID string `json:"featured_tag_id"`
				NameJSON      string `json:"name_json"`
				Slug          string `json:"slug"`
				Color         string `json:"color"`
				GenreTag      any    `json:"genre_tag"`
				IsDiscover    bool   `json:"is_discover"`
			} `json:"tags"`
			TracksCount   int    `json:"tracks_count"`
			PublicAt      int    `json:"public_at"`
			Name          string `json:"name"`
			IsPublic      bool   `json:"is_public"`
			PublishedFrom int    `json:"published_from"`
			IsFeatured    bool   `json:"is_featured"`
		} `json:"items"`
	} `json:"playlists"`
	Stories struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
		Items  []struct {
			ID               string   `json:"id"`
			SectionSlugs     []string `json:"section_slugs"`
			Title            string   `json:"title"`
			DescriptionShort string   `json:"description_short"`
			Authors          []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"authors"`
			Image  string `json:"image"`
			Images []struct {
				Format string `json:"format"`
				URL    string `json:"url"`
			} `json:"images"`
			DisplayDate int `json:"display_date"`
		} `json:"items"`
	} `json:"stories"`
}

type QobuzTrack struct {
	Album               Album     `json:"album"`
	Articles            []any     `json:"articles"`
	AudioInfo           AudioInfo `json:"audio_info"`
	Composer            Person    `json:"composer"`
	Copyright           string    `json:"copyright"`
	Displayable         bool      `json:"displayable"`
	Downloadable        bool      `json:"downloadable"`
	Duration            int       `json:"duration"`
	Hires               bool      `json:"hires"`
	HiresStreamable     bool      `json:"hires_streamable"`
	ID                  int       `json:"id"`
	ISRC                string    `json:"isrc"`
	MaximumBitDepth     int       `json:"maximum_bit_depth"`
	MaximumChannelCount int       `json:"maximum_channel_count"`
	MaximumSamplingRate float32   `json:"maximum_sampling_rate"`
	MediaNumber         int       `json:"media_number"`
	ParentalWarning     bool      `json:"parental_warning"`
	Performer           Person    `json:"performer"`
	Performers          []string  `json:"performers"`
	Previewable         bool      `json:"previewable"`
	Purchasable         bool      `json:"purchasable"`
	PurchasableAt       int64     `json:"purchasable_at"`
	ReleaseDateDownload string    `json:"release_date_download"`
	ReleaseDateOriginal string    `json:"release_date_original"`
	ReleaseDatePurchase string    `json:"release_date_purchase"`
	ReleaseDateStream   string    `json:"release_date_stream"`
	Sampleable          bool      `json:"sampleable"`
	Streamable          bool      `json:"streamable"`
	StreamableAt        int64     `json:"streamable_at"`
	Title               string    `json:"title"`
	TrackNumber         int       `json:"track_number"`
	Version             string    `json:"version"`
}

type Album struct {
	Area                *string      `json:"area"`
	Articles            []any        `json:"articles"`
	Artist              Person       `json:"artist"`
	Artists             []ArtistRole `json:"artists"`
	Awards              []any        `json:"awards"`
	Catchline           string       `json:"catchline"`
	Composer            Person       `json:"composer"`
	Copyright           string       `json:"copyright"`
	CreatedAt           int64        `json:"created_at"`
	Description         string       `json:"description"`
	DescriptionLanguage string       `json:"description_language"`
	Displayable         bool         `json:"displayable"`
	Downloadable        bool         `json:"downloadable"`
	Duration            int          `json:"duration"`
	Genre               Genre        `json:"genre"`
	GenresList          []string     `json:"genres_list"`
	Goodies             []any        `json:"goodies"`
	Hires               bool         `json:"hires"`
	HiresStreamable     bool         `json:"hires_streamable"`
	ID                  string       `json:"id"`
	Image               Image        `json:"image"`
	IsOfficial          bool         `json:"is_official"`
	Label               Label        `json:"label"`
	MaximumBitDepth     int          `json:"maximum_bit_depth"`
	MaximumChannelCount int          `json:"maximum_channel_count"`
	MaximumSamplingRate float32      `json:"maximum_sampling_rate"`
	MediaCount          int          `json:"media_count"`
	ParentalWarning     bool         `json:"parental_warning"`
	Popularity          int          `json:"popularity"`
	Previewable         bool         `json:"previewable"`
	ProductSalesFactors struct {
		Monthly int `json:"product_sales_factors_monthly"`
		Weekly  int `json:"product_sales_factors_weekly"`
		Yearly  int `json:"product_sales_factors_yearly"`
	} `json:"product_sales_factors"`
	ProductType          string        `json:"product_type"`
	ProductURL           string        `json:"product_url"`
	Purchasable          bool          `json:"purchasable"`
	PurchasableAt        int64         `json:"purchasable_at"`
	QobuzID              float64       `json:"qobuz_id"`
	RecordingInformation RecordingInfo `json:"recording_information"`
	RelativeURL          string        `json:"relative_url"`
	ReleaseDateDownload  string        `json:"release_date_download"`
	ReleaseDateOriginal  string        `json:"release_date_original"`
	ReleaseDateStream    string        `json:"release_date_stream"`
	ReleaseTags          []string      `json:"release_tags"`
	ReleaseType          string        `json:"release_type"`
	ReleasedAt           int64         `json:"released_at"`
	Sampleable           bool          `json:"sampleable"`
	Slug                 string        `json:"slug"`
	Streamable           bool          `json:"streamable"`
	StreamableAt         int64         `json:"streamable_at"`
	Subtitle             string        `json:"subtitle"`
	Title                string        `json:"title"`
	TracksCount          int           `json:"tracks_count"`
	UPC                  string        `json:"upc"`
	URL                  string        `json:"url"`
	Version              string        `json:"version"`
}

type AudioInfo struct {
	ReplayGainTrackGain float64 `json:"replaygain_track_gain"`
	ReplayGainTrackPeak float64 `json:"replaygain_track_peak"`
}

type Person struct {
	ID          float64 `json:"id"`
	Name        string  `json:"name"`
	AlbumsCount float64 `json:"albums_count,omitempty"`
	Image       *string `json:"image"`
	Picture     *string `json:"picture"`
	Slug        string  `json:"slug,omitempty"`
}

type ArtistRole struct {
	ID    float64  `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type Genre struct {
	Color string  `json:"color"`
	ID    float64 `json:"id"`
	Name  string  `json:"name"`
	Path  []int   `json:"path"`
	Slug  string  `json:"slug"`
}

type Image struct {
	Back      *string `json:"back"`
	Large     string  `json:"large"`
	Small     string  `json:"small"`
	Thumbnail string  `json:"thumbnail"`
}

type Label struct {
	AlbumsCount int     `json:"albums_count"`
	ID          float64 `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	SupplierID  int     `json:"supplier_id"`
}

type RecordingInfo struct {
	RelativeURL string `json:"relative_url"`
}

func (p qobuz_search_result) ToJSON() []byte {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return []byte("{}")
	}
	return jsonData
}
