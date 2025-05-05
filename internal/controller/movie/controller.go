package movie

import (
	"context"
	"errors"
	metadata "github.com/kirillApanasiuk/movie-metadata/model"
	"github.com/kirillApanasiuk/movie-movie/pkg/model"
	rating "github.com/kirillApanasiuk/movie-rating/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID rating.RecordID, recordType rating.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID rating.RecordID, recordType rating.RecordType, rating *rating.Rating) error
}
type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadata.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{
		ratingGateway:   ratingGateway,
		metadataGateway: metadataGateway,
	}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, ErrNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{
		Metadata: metadata,
	}

	rating, err := c.ratingGateway.GetAggregatedRating(ctx, rating.RecordID(id), rating.RecordTypeMovie)
	if err != nil && !errors.Is(err, ErrNotFound) {
		//TODO write logic here
	}
	if err != nil {
		return nil, err
	}
	details.Rating = &rating

	return details, nil
}
