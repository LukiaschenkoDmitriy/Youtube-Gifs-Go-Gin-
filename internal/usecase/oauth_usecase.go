package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmytrii/youtube-gifs-chat/config"
	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/dmytrii/youtube-gifs-chat/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
)

type OAuthUC struct {
	Config *oauth2.Config
}

type BaseUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type YoutubeChannelInfo struct {
	Items []struct {
		Snippet struct {
			CustomUrl string `json:"customUrl"`
		} `json:"snippet"`
	} `json:"items"`
}

func NewOAuthUC(c config.Config) *OAuthUC {
	return &OAuthUC{
		Config: getConfig(c),
	}
}

func (oa *OAuthUC) getToken(ctx context.Context, code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("OAuth2: Authorisation failed")
	}

	token, err := oa.Config.Exchange(ctx, code)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (oa *OAuthUC) GetLoginRedirectUrl() string {
	state, _ := utils.RandomHex(32)
	return oa.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

func (oa *OAuthUC) GetUser(ctx context.Context, authCode string) (*domain.User, error) {
	token, err := oa.getToken(ctx, authCode)

	if err != nil {
		return nil, err
	}

	g, ctx := errgroup.WithContext(ctx)

	var client = oa.Config.Client(ctx, token)

	var baseInfo *BaseUserInfo
	var customUrl *string

	g.Go(func() error {
		var err error
		baseInfo, err = oa.getUserInfo(ctx, client)
		return err
	})

	g.Go(func() error {
		var err error
		customUrl, err = oa.getUserCustomURL(ctx, client)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &domain.User{
		OAuthId:   baseInfo.ID,
		Name:      baseInfo.Name,
		Email:     baseInfo.Email,
		Picture:   baseInfo.Picture,
		CustomUrl: *customUrl,
	}, nil
}

func (oa *OAuthUC) getUserInfo(ctx context.Context, client *http.Client) (*BaseUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("OAuth2: Failed to get user info: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("OAuth2: Token is invalid")
	}

	baseInfo := new(BaseUserInfo)

	if err := json.NewDecoder(resp.Body).Decode(&baseInfo); err != nil {
		return nil, err
	}

	return baseInfo, nil
}

func (oa *OAuthUC) getUserCustomURL(ctx context.Context, client *http.Client) (*string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&mine=true", nil)

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, errors.New("OAuth2: Failed to get Youtube channel info: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("OAuth2: Token is invalid")
	}

	channelInfo := new(YoutubeChannelInfo)

	if err := json.NewDecoder(resp.Body).Decode(&channelInfo); err != nil {
		return nil, err
	}

	return &channelInfo.Items[0].Snippet.CustomUrl, nil
}

func getConfig(c config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.GoogleClientId,
		ClientSecret: c.GoogleClientSecret,
		RedirectURL:  "http://localhost:8080/auth/2l8s118z69mkq91m3y6r161bq8yp4hmsgaveoqzivvzfs45kb1/callback",
		Scopes: []string{
			"email",
			"profile",
			"https://www.googleapis.com/auth/youtube.readonly",
		},
		Endpoint: google.Endpoint,
	}
}
