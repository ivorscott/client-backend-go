package mid

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

type Auth0 struct {
	Audience string
	Domain   string
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

// Authenticate middleware verifies the access token sent to the api
func (auth0 *Auth0) Authenticate(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		Debug: true,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(auth0.Audience, false)
			if !checkAud {
				return token, errors.New("invalid audience.")
			}
			iss := "https://" + auth0.Domain + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer.")
			}
			cert, err := auth0.getPemCert(token)
			if err != nil {
				return nil, err
			}
			return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware.Handler(next)
}

// CheckScope middleware verifies the Access Token has the correct scope before returning a successful response
func (auth0 *Auth0) CheckScope(scope string, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := auth0.getPemCert(token)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	hasScope := false

	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}
	return hasScope, nil
}

// getPemCert takes a token and returns the associated certificate in PEM format so it can be parsed
func (auth0 *Auth0) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + auth0.Domain + "/.well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}
	return cert, nil
}
