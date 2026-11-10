package model

type ShortURL struct {
	URL string `json:"url"`
}

type ShortURLResponse struct {
	URL string `json:"result"`
}
