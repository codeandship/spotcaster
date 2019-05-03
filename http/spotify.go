package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/codeandship/spotcaster"
)

const (
	bearerTokenURL = "https://generic.wg.spotify.com/creator-auth-proxy/v1/web/token"
	cookieDomain   = ".spotify.com"
	metaDataURL    = "https://generic.wg.spotify.com/podcasters/v0/shows/"
	ClientID       = "05a1371ee5194c27860b3ff3ff3979d2"
)

type SpotifyAPI struct {
	clientID  string
	token     *spotcaster.Token
	cookieVal string
	client    *http.Client
}

func NewSpotifyAPI(clientID, cookieVal string) (*SpotifyAPI, error) {
	if clientID == "" {
		return nil, errors.New("spotify-api: empty client-id")
	}
	if cookieVal == "" {
		return nil, errors.New("spotify-api: empty cookie value")
	}
	client := http.DefaultClient
	client.Timeout = time.Second * 10
	return &SpotifyAPI{
		clientID:  clientID,
		cookieVal: cookieVal,
		client:    client,
	}, nil
}

func (s *SpotifyAPI) FetchBearerToken() (spotcaster.Token, error) {
	token := spotcaster.Token{}
	spotifyAuth, err := url.Parse(bearerTokenURL)
	if err != nil {
		return token, err
	}
	q := spotifyAuth.Query()
	q.Set("client_id", s.clientID)
	spotifyAuth.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, spotifyAuth.String(), nil)
	if err != nil {
		return token, err
	}

	req.AddCookie(&http.Cookie{Name: "sp_dc", Value: s.cookieVal, Domain: cookieDomain})

	res, err := s.client.Do(req)
	if err != nil {
		return token, err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	err = json.Unmarshal(result, &token)
	if err != nil {
		return token, err
	}
	token.ExpiresAt = time.Now().Add(time.Second * time.Duration(token.ExpiresIn)).UnixNano()
	return token, nil
}

func (s *SpotifyAPI) FetchMetaData(showID string) (spotcaster.MetaData, error) {
	meta := spotcaster.MetaData{}
	if time.Now().After(time.Unix(0, s.token.ExpiresAt)) {
		log.Println("api token expired, fetching new token")
		t, err := s.FetchBearerToken()
		if err != nil {
			return meta, err
		}
		s.token = &t
	}

	spotifyMeta, err := url.Parse(metaDataURL)
	if err != nil {
		return meta, err
	}

	spotifyMeta.Path = path.Join(spotifyMeta.Path, showID)
	spotifyMeta.Path = path.Join(spotifyMeta.Path, "metadata")

	req, err := http.NewRequest(http.MethodGet, spotifyMeta.String(), nil)
	if err != nil {
		return meta, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token.AccessToken))

	res, err := s.client.Do(req)
	if err != nil {
		return meta, err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	err = json.Unmarshal(result, &meta)
	if err != nil {
		return meta, err
	}
	return meta, nil
}

func (s *SpotifyAPI) BearerToken() *spotcaster.Token {
	return s.token
}

func (s *SpotifyAPI) SetBearerToken(t *spotcaster.Token) {
	s.token = t
}
