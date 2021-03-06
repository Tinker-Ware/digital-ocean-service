package usecases

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
)

// DOInteractor embeds an has the necessary fields to interact with the
// Digital Ocean API
type DOInteractor struct {
}

const authURL = "https://cloud.digitalocean.com/v1/oauth/authorize"

// GetOauthURL returns the oauth URL
func (interactor DOInteractor) GetOauthURL(id, redirectURI string, scope []string) string {

	scp := strings.Join(scope, " ")

	u, _ := url.Parse(authURL)
	q := u.Query()
	q.Set("client_id", id)
	q.Set("redirect_uri", redirectURI)
	q.Set("scope", scp)
	q.Set("response_type", "code")

	u.RawQuery = q.Encode()

	return u.String()
}

// GetToken returns a token from Digital Ocean
func (interactor DOInteractor) GetToken(code, id, secret, redirectURL string) (*domain.DOToken, error) {
	u, _ := url.Parse("https://cloud.digitalocean.com/v1/oauth/token")
	q := u.Query()

	q.Set("grant_type", "authorization_code")
	q.Set("code", code)
	q.Set("client_id", id)
	q.Set("client_secret", secret)
	q.Set("redirect_uri", redirectURL)

	u.RawQuery = q.Encode()

	res, err := http.Post(u.String(), "", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		errMap := map[string]string{}
		decoder.Decode(&errMap)
		log.Printf("OAUTH error %s: %s", errMap["error"], errMap["error_description"])
		return nil, errors.New("Cannot create token")
	}

	accessToken := domain.DOToken{}

	err = decoder.Decode(&accessToken)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
