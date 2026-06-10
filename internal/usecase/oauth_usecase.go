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
)

type OAuthUC struct {
	Config *oauth2.Config
	client *http.Client
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

func (oa *OAuthUC) getToken(code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("OAuth2: Authorisation failed")
	}

	token, err := oa.Config.Exchange(context.Background(), code)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (oa *OAuthUC) InitClient(authCode string) error {
	if oa.client != nil {
		return nil
	}

	token, err := oa.getToken(authCode)

	if err != nil {
		return err
	}

	oa.client = oa.Config.Client(context.Background(), token)

	return nil
}

func (oa *OAuthUC) GetLoginRedirectUrl() string {
	state, _ := utils.RandomHex(32)
	return oa.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
}

func (oa *OAuthUC) GetUser() (*domain.User, error) {
	baseInfo, err := oa.getUserInfo()

	if err != nil {
		return nil, err
	}

	customUrl, err := oa.getUserCustomURL()

	if err != nil {
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

func (oa *OAuthUC) getUserInfo() (*BaseUserInfo, error) {
	resp, err := oa.client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, err
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

func (oa *OAuthUC) getUserCustomURL() (*string, error) {
	resp, err := oa.client.Get("https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&mine=true")

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
