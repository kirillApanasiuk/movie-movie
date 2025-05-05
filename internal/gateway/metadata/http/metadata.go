package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kirillApanasiuk/kit/pkg/discovery"
	metadata "github.com/kirillApanasiuk/movie-metadata/model"
	"github.com/kirillApanasiuk/movie-movie/internal/gateway"
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

func (g *Gateway) Get(ctx context.Context, id string) (*metadata.Metadata, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadata service. Request: GET: %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	//TODO  strange behavior
	var v *metadata.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

//func (r *Request) WithContext(ctx context.Context) *Request {
//	if ctx == nil {
//		panic("nil context")
//	}
//	r2 := new(Request)
//	*r2 = *r
//	r2.ctx = ctx
//	return r2
//}
