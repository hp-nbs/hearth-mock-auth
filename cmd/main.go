package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"github.com/MicahParks/jwkset"
	"github.com/golang-jwt/jwt/v4"
	echo2 "github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	logFmt = "%s\nError: %s"
	keyId  = "my-key-id"
)

func main() {
	ctx := context.Background()
	logger := log.New(os.Stdout, "", 0)

	jwkSet := jwkset.NewMemory[any]()

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.Fatalf(logFmt, "Failed to generate RSA key.", err)
	}

	err = jwkSet.Store.WriteKey(ctx, jwkset.NewKey[any](key, keyId))
	if err != nil {
		logger.Fatalf(logFmt, "Failed to store RSA key.", err)
	}

	server := &server{
		key,
		jwkSet,
	}

	echo := echo2.New()
	echo.POST("/generate-token", server.generateToken)
	echo.GET("/.well-known/jwks.json", server.getJwks)
	echo.GET("/healthz", server.health)

	if err := server.start(echo); err != nil {
		log.Panic(err)
	}
}

type server struct {
	key    *rsa.PrivateKey
	jwkSet jwkset.JWKSet[any]
}

func (s *server) start(e *echo2.Echo) error {
	return e.Start(":7000")
}

type generateTokenRequest struct {
	CustomClaims map[string]interface{} `json:"customClaims"`
}

// generateToken ...hard coded userServiceId currently but can be taken through input
func (s *server) generateToken(c echo2.Context) error {
	requestBody := generateTokenRequest{}
	err := extractRequestBody(c, &requestBody)

	claims := jwt.MapClaims{
		"sub": "",
	}

	// add the claims to the map
	for k, v := range requestBody.CustomClaims {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(s.key)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.String(http.StatusOK, tokenString)
}

// getJwks ...send back the jwks json which will include the jwt data for the generateToken api
func (s *server) getJwks(c echo2.Context) error {
	response, err := s.jwkSet.JSONPublic(c.Request().Context())
	if err != nil {
		log.Printf(logFmt, "Failed to get JWK Set JSON.", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("sending jwks")
	return c.JSON(http.StatusOK, response)

}

func (s *server) health(c echo2.Context) error {
	return c.String(http.StatusOK, "")

}

func extractRequestBody(ctx echo2.Context, destination interface{}) error {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &destination); err != nil {
		return err
	}

	return nil
}
