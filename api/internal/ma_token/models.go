// This package contains model definitions for the management api token for auth0
// This token is distinct from the regular access token used by the client
// to receive api resources. The management api token allows the backend to
// update the user's auth0 account, setting the app_metadata with the inter.
package ma_token

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	ID          string    `db:"ma_token_id" json:"id"`
	AccessToken string    `db:"token" json:"access_token"`
	Expiration  string    `db:"expiration" json:"expiration"`
	Created     time.Time `db:"created" json:"created"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

