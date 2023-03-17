package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenResponse struct {
	Token string `json:"token"`
}

var githubAPIURL string = "https://api.github.com"

// Creates a new GitHub App with JWT authentication
func GetJWTContext(appId int, key []byte) (context.Context, error) {
	jwt, err := buildJWTToken(appId, key)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), "auth_token", jwt)
	return ctx, nil
}

// Builds the JWT token from the provided appId and key
func buildJWTToken(appId int, privateKey []byte) (string, error) {
	// Generate JWT token
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Add(time.Duration(-1) * time.Minute).Unix()
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["iss"] = appId

	// Parse RSA private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	// Sign the token
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Gets the access token to authenticate with github
func GetAccessToken(ctx context.Context, instId int) (string, error) {
	// Parse url
	url := fmt.Sprintf("%s/app/installations/%d/access_tokens", githubAPIURL, instId)

	// Create the request to github's api
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ctx.Value("auth_token")))

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the the response from the server
	body, err := ioutil.ReadAll(resp.Body)

	// Parse the response to the struct
	t := TokenResponse{}
	json.Unmarshal(body, &t)

	return t.Token, nil
}
