package model

type MovieDetails struct {
	Rating   *float64 `json:"rating,omitempty"`
	Metadata model.Metadata
}
