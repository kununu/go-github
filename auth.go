package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Auth struct {
	JWTToken string
	Token    string
}

// Builds the JWT token from the provided appId and key
func (ghApp *GitHubApp) buildJWTToken() error {
	// Generate JWT token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Add(time.Duration(-1) * time.Minute).Unix()
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["iss"] = ghApp.Config.ApplicationID

	// Parse RSA private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(ghApp.Config.PrivateKey)
	if err != nil {
		return err
	}

	// Sign the token
	jwtToken, err := token.SignedString(key)
	if err != nil {
		return err
	}

	ghApp.Auth.JWTToken = jwtToken
	return nil
}

// TokenResponse is a struct that represents the response from the API
type TokenResponse struct {
	Token string `json:"token"`
}

// Gets the access token to authenticate with github
func (ghApp *GitHubApp) GetAccessToken() (string, error) {
	// Parse url
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", githubAPIURL, ghApp.Config.InstallationID)

	// Create the request to github's api
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ghApp.Auth.JWTToken))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", errors.New("couldn't get the token")
	}

	// Read the the response from the server
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the response to the struct
	t := TokenResponse{}
	json.Unmarshal(body, &t)

	return t.Token, nil
}
