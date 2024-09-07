package interceptor

import (
	"blockchain-transactions/internal/env"
	"blockchain-transactions/internal/logger"
	"blockchain-transactions/internal/models"
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
)

// Contenido de las llaves pública y privada
var (
	verifyKey *rsa.PublicKey
)

// init lee los archivos de firma y validación RSA
func init() {
	c := env.NewConfiguration()
	verifyBytes, err := ioutil.ReadFile(c.App.RSAPublicKey)
	if err != nil {
		logger.Error.Printf("leyendo el archivo público de confirmación: %s", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		logger.Error.Printf("realizando el parse en jwt RSA public: %s", err)
	}
}

// UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
	jwt.StandardClaims
	User models.User `json:"user"`
	Role int         `json:"role"`
}

// Verify verifies the access token string and return a user claim if the token is valid
func Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
