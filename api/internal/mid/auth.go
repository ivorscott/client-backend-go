package mid

import (
	"encoding/json"
	"errors"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Auth0 struct {
	Audience  string
	Domain    string
	M2MClient string
	M2MSecret string
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

type token struct {
	AccessToken string `json:"access_token"`
}

// Authenticate middleware verifies the access token sent to the api
func (a0 *Auth0) Authenticate() web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {

			jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
				Debug: false,
				ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
					checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(a0.Audience, false)
					if !checkAud {
						return token, errors.New("invalid audience.")
					}
					iss := "https://" + a0.Domain + "/"
					checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
					if !checkIss {
						return token, errors.New("invalid issuer.")
					}
					cert, err := a0.GetPemCert(token)
					if err != nil {
						return nil, err
					}
					return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
				},
				SigningMethod: jwt.SigningMethodRS256,
			})

			if err := jwtMiddleware.CheckJWT(w, r); err != nil {
				return web.NewRequestError(err, http.StatusForbidden)
			}

			return after(w, r)
		}

		return h
	}

	return f
}

// CheckScope middleware verifies the Access Token has the correct scope before returning a successful response
func (a0 *Auth0) CheckScope(scope, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a0.GetPemCert(token)
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
func (a0 *Auth0) GetPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + a0.Domain + "/.well-known/jwks.json")
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

func (a0 *Auth0) GetManagementToken() (*token, error) {
	// get token from database

	// if exists and not expired return it

	// else generate new token

	baseUrl := "https://" + a0.Domain
	resource := "/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a0.M2MClient)
	data.Set("client_secret", a0.M2MSecret)
	data.Set("audience", a0.Audience)

	u, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}

	u.Path = resource
	urlStr := u.String()

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	// Clean table. Maintain only 1 valid token at a time (DELETE FROM auth0_management;)

	// save generated token

	return &token, nil
}

func (a0 *Auth0) GetAuthTokenSubject(r *http.Request) string {
	claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["sub"])
}

//func CheckExpiration() {
//	var claims Helpers.JwtCustomClaims
//	_ , err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
//		return []byte(Helpers.GetConfig("JWT_SECRET","***")), nil
//	})
//
//	v, _ := err.(*jwt.ValidationError)
//
//	if v.Errors == jwt.ValidationErrorExpired && claims.ExpiresAt > time.Now().Unix() - (86400 * 14){
//		/** refresh logic **/
//	}
//}

//func ParseManagementToken(r *http.Request) {
//	user := r.Context().Value("user")
//	fmt.Printf("This is an authenticated request===================>>>>>>>>>>>>>>>>>>>>>>")
//	fmt.Printf("Claim content:\n")
//	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
//		fmt.Printf("%s :\t%#v\n", k, v)
//	}
//}
