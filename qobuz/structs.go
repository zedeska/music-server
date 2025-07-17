package qobuz

import (
	"encoding/json"

	db "music-server/database"
)

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

type custom_search_result struct {
	Tracks []db.Track `json:"tracks"`
	Albums []db.Album `json:"albums"`
}

func (p custom_search_result) ToJSON() []byte {
	data, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return data
}

type QobuzAlbum struct {
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
	MaximumChannelCount int    `json:"maximum_channel_count"`
	ID                  string `json:"id"`
	MaximumSamplingRate int    `json:"maximum_sampling_rate"`
	Articles            []any  `json:"articles"`
	ReleaseDateOriginal string `json:"release_date_original"`
	ReleaseDateDownload string `json:"release_date_download"`
	ReleaseDateStream   string `json:"release_date_stream"`
	Purchasable         bool   `json:"purchasable"`
	Streamable          bool   `json:"streamable"`
	Previewable         bool   `json:"previewable"`
	Sampleable          bool   `json:"sampleable"`
	Downloadable        bool   `json:"downloadable"`
	Displayable         bool   `json:"displayable"`
	PurchasableAt       int    `json:"purchasable_at"`
	StreamableAt        int    `json:"streamable_at"`
	Hires               bool   `json:"hires"`
	HiresStreamable     bool   `json:"hires_streamable"`
	Awards              []any  `json:"awards"`
	Description         string `json:"description"`
	DescriptionLanguage string `json:"description_language"`
	Goodies             []any  `json:"goodies"`
	Area                any    `json:"area"`
	Catchline           string `json:"catchline"`
	Composer            struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		AlbumsCount int    `json:"albums_count"`
		Picture     any    `json:"picture"`
		Image       any    `json:"image"`
	} `json:"composer"`
	CreatedAt                      int      `json:"created_at"`
	GenresList                     []string `json:"genres_list"`
	Period                         any      `json:"period"`
	Copyright                      string   `json:"copyright"`
	IsOfficial                     bool     `json:"is_official"`
	MaximumTechnicalSpecifications string   `json:"maximum_technical_specifications"`
	ProductSalesFactorsMonthly     int      `json:"product_sales_factors_monthly"`
	ProductSalesFactorsWeekly      int      `json:"product_sales_factors_weekly"`
	ProductSalesFactorsYearly      int      `json:"product_sales_factors_yearly"`
	ProductType                    string   `json:"product_type"`
	ProductURL                     string   `json:"product_url"`
	RecordingInformation           string   `json:"recording_information"`
	RelativeURL                    string   `json:"relative_url"`
	ReleaseTags                    []any    `json:"release_tags"`
	ReleaseType                    string   `json:"release_type"`
	Slug                           string   `json:"slug"`
	Subtitle                       string   `json:"subtitle"`
	TrackIds                       []int    `json:"track_ids"`
	Tracks                         struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
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
			Work     any `json:"work"`
			Composer struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			} `json:"composer"`
			Isrc                string `json:"isrc"`
			Title               string `json:"title"`
			Version             any    `json:"version"`
			Duration            int    `json:"duration"`
			ParentalWarning     bool   `json:"parental_warning"`
			TrackNumber         int    `json:"track_number"`
			MaximumChannelCount int    `json:"maximum_channel_count"`
			ID                  int    `json:"id"`
			MediaNumber         int    `json:"media_number"`
			MaximumSamplingRate int    `json:"maximum_sampling_rate"`
			ReleaseDateOriginal string `json:"release_date_original"`
			ReleaseDateDownload string `json:"release_date_download"`
			ReleaseDateStream   string `json:"release_date_stream"`
			ReleaseDatePurchase string `json:"release_date_purchase"`
			Purchasable         bool   `json:"purchasable"`
			Streamable          bool   `json:"streamable"`
			Previewable         bool   `json:"previewable"`
			Sampleable          bool   `json:"sampleable"`
			Downloadable        bool   `json:"downloadable"`
			Displayable         bool   `json:"displayable"`
			PurchasableAt       int    `json:"purchasable_at"`
			StreamableAt        int    `json:"streamable_at"`
			Hires               bool   `json:"hires"`
			HiresStreamable     bool   `json:"hires_streamable"`
		} `json:"items"`
	} `json:"tracks"`
}

type QobuzTrack struct {
	MaximumBitDepth int    `json:"maximum_bit_depth"`
	Copyright       string `json:"copyright"`
	Performers      string `json:"performers"`
	AudioInfo       struct {
		ReplaygainTrackGain float64 `json:"replaygain_track_gain"`
		ReplaygainTrackPeak float64 `json:"replaygain_track_peak"`
	} `json:"audio_info"`
	Performer struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"performer"`
	Album struct {
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
		MaximumChannelCount int    `json:"maximum_channel_count"`
		ID                  string `json:"id"`
		MaximumSamplingRate int    `json:"maximum_sampling_rate"`
		Articles            []any  `json:"articles"`
		ReleaseDateOriginal string `json:"release_date_original"`
		ReleaseDateDownload string `json:"release_date_download"`
		ReleaseDateStream   string `json:"release_date_stream"`
		Purchasable         bool   `json:"purchasable"`
		Streamable          bool   `json:"streamable"`
		Previewable         bool   `json:"previewable"`
		Sampleable          bool   `json:"sampleable"`
		Downloadable        bool   `json:"downloadable"`
		Displayable         bool   `json:"displayable"`
		PurchasableAt       int    `json:"purchasable_at"`
		StreamableAt        int    `json:"streamable_at"`
		Hires               bool   `json:"hires"`
		HiresStreamable     bool   `json:"hires_streamable"`
		Awards              []any  `json:"awards"`
		Description         string `json:"description"`
		DescriptionLanguage string `json:"description_language"`
		Goodies             []any  `json:"goodies"`
		Area                any    `json:"area"`
		Catchline           string `json:"catchline"`
		Composer            struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Slug        string `json:"slug"`
			AlbumsCount int    `json:"albums_count"`
			Picture     any    `json:"picture"`
			Image       any    `json:"image"`
		} `json:"composer"`
		CreatedAt                      int      `json:"created_at"`
		GenresList                     []string `json:"genres_list"`
		Period                         any      `json:"period"`
		Copyright                      string   `json:"copyright"`
		IsOfficial                     bool     `json:"is_official"`
		MaximumTechnicalSpecifications string   `json:"maximum_technical_specifications"`
		ProductSalesFactorsMonthly     int      `json:"product_sales_factors_monthly"`
		ProductSalesFactorsWeekly      int      `json:"product_sales_factors_weekly"`
		ProductSalesFactorsYearly      int      `json:"product_sales_factors_yearly"`
		ProductType                    string   `json:"product_type"`
		ProductURL                     string   `json:"product_url"`
		RecordingInformation           string   `json:"recording_information"`
		RelativeURL                    string   `json:"relative_url"`
		ReleaseTags                    []any    `json:"release_tags"`
		ReleaseType                    string   `json:"release_type"`
		Slug                           string   `json:"slug"`
		Subtitle                       string   `json:"subtitle"`
	} `json:"album"`
	Work     any `json:"work"`
	Composer struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"composer"`
	Isrc                string `json:"isrc"`
	Title               string `json:"title"`
	Version             any    `json:"version"`
	Duration            int    `json:"duration"`
	ParentalWarning     bool   `json:"parental_warning"`
	TrackNumber         int    `json:"track_number"`
	MaximumChannelCount int    `json:"maximum_channel_count"`
	ID                  int    `json:"id"`
	MediaNumber         int    `json:"media_number"`
	MaximumSamplingRate int    `json:"maximum_sampling_rate"`
	Articles            []any  `json:"articles"`
	ReleaseDateOriginal string `json:"release_date_original"`
	ReleaseDateDownload string `json:"release_date_download"`
	ReleaseDateStream   string `json:"release_date_stream"`
	ReleaseDatePurchase string `json:"release_date_purchase"`
	Purchasable         bool   `json:"purchasable"`
	Streamable          bool   `json:"streamable"`
	Previewable         bool   `json:"previewable"`
	Sampleable          bool   `json:"sampleable"`
	Downloadable        bool   `json:"downloadable"`
	Displayable         bool   `json:"displayable"`
	PurchasableAt       int    `json:"purchasable_at"`
	StreamableAt        int    `json:"streamable_at"`
	Hires               bool   `json:"hires"`
	HiresStreamable     bool   `json:"hires_streamable"`
}
