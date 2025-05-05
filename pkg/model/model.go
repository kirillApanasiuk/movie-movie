package model

import "github.com/kirillApanasiuk/movie-metadata/model"

type MovieDetails struct {
	Rating   *float64        `json:"rating,omitempty"`
	Metadata *model.Metadata `json:"metadata,omitempty"`
}
