package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
)

type GiphyUC struct {
	apiKey string
}

type GiphyResponse struct {
	Data []struct {
		Type   string `json:"type"`
		Title  string `json:"title"`
		Images map[string]struct {
			URL string `json:"url"`
		}
	} `json:"data"`
	Meta struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
}

func NewGiphyUC(c *config.Config) *GiphyUC {
	return &GiphyUC{apiKey: c.GiphyApiKey}
}

func (g *GiphyUC) GetTrendingGifs(ctx context.Context, offset string) ([]*domain.Gif, error) {
	return g.baseRequest(ctx, "GET", "https://api.giphy.com/v1/gifs/trending", map[string]string{
		"offset": offset,
	})
}

func (g *GiphyUC) GetGifsBySearch(ctx context.Context, search string, offset string) ([]*domain.Gif, error) {
	return g.baseRequest(ctx, "GET", "https://api.giphy.com/v1/gifs/search", map[string]string{
		"q":      search,
		"offset": offset,
	})
}

func (g *GiphyUC) baseRequest(ctx context.Context, method string, url string, params map[string]string) ([]*domain.Gif, error) {
	r, err := http.NewRequestWithContext(ctx, method, url, nil)

	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	q.Add("api_key", g.apiKey)
	q.Add("limit", "40")

	for k, v := range params {
		q.Add(k, v)
	}

	r.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respObj := new(GiphyResponse)

	if err := json.NewDecoder(resp.Body).Decode(&respObj); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(respObj.Meta.Msg)
	}

	gifs := make([]*domain.Gif, len(respObj.Data))

	for i, g := range respObj.Data {
		gifs[i] = &domain.Gif{
			Url:   g.Images["original"].URL,
			Title: g.Title,
			Type:  g.Type,
		}
	}

	return gifs, nil
}
