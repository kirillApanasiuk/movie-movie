package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kirillApanasiuk/kit/pkg/discovery"
	"github.com/kirillApanasiuk/movie-movie/internal/gateway"
	"github.com/kirillApanasiuk/movie-rating/model"
	rating "github.com/kirillApanasiuk/movie-rating/model"
	"log"
	"math/rand"
	"net/http"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordId model.RecordID,
	recordType model.RecordType) (float64,
	error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "rating")
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/rating"

	log.Printf("Calling metadata service. Request: GET: %s", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("+%v", recordType))
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		return 0, gateway.ErrNotFound
	}

	if resp.StatusCode/100 != 2 {
		return 0, gateway.NotSuccesfull
	}
	var v float64
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func (g *Gateway) PutRating(ctx context.Context, recordId model.RecordID, recordType model.RecordType,
	rating *rating.Rating) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, g.addr+"/rating/", nil)
	if err != nil {
		return err
	}

	values := req.URL.Query()
	values.Add("id", string(recordId))
	values.Add("type", fmt.Sprintf("+%v", recordType))
	values.Add("userId", fmt.Sprintf("+%v", rating.UserID))
	values.Add("value", fmt.Sprintf("+%v", rating.Value))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return gateway.NotSuccesfull
	}
	return nil
}
