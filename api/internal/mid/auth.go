package mid

import (
	"encoding/json"
	"errors"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
	"net/http"
	"strings"
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
func  (a0 *Auth0) CheckScope(scope, tokenString string) (bool, error) {
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

func (a0 *Auth0) GetManagementToken() {
	url :=  a0.Domain + "/oauth/token"
	fmt.Print("===============>>>>>>>>>>>>>>>")
	fmt.Print(url)
	//payload := strings.NewReader("grant_type=client_credentials&client_id=%24%7Baccount.clientId%7D&client_secret=YOUR_CLIENT_SECRET&audience=https%3A%2F%2F%24%7Baccount.namespace%7D%2Fapi%2Fv2%2F")
	//
	//req, _ := http.NewRequest("POST", url, payload)
	//
	//req.Header.Add("content-type", "application/x-www-form-urlencoded")
	//
	//res, _ := http.DefaultClient.Do(req)
	//
	//defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)
	//
	//fmt.Println(res)
	//fmt.Println(string(body))
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
