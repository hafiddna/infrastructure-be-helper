package helper

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	uuid2 "github.com/google/uuid"
	"github.com/hafiddna/infrastructure-be-helper/config"
	"strconv"
	"time"
)

type jwtGeneralClaim struct {
	jwt.RegisteredClaims
	Data interface{} `json:"data"`
	//Data EncryptedData `json:"data"`
}

type JwtAuthClaim struct {
	Roles []string `json:"roles"`
}

type JwtRememberClaim struct {
	RememberToken string `json:"remember_token"`
}

func GenerateRS512Token(privateKey, key string, userID uint64, data interface{}, duration time.Time) string {
	uuid := uuid2.New()

	var claims jwt.Claims
	var err error

	//marshalledData := JSONMarshal(data)
	//encryptedData, err := EncryptAES256CBC([]byte(marshalledData), []byte(key))
	//if err != nil {
	//	log.Fatalf("Error encrypting data: %v", err)
	//}

	jwtAuthClaimData, ok := data.(JwtAuthClaim)
	if ok {
		claims = &jwtGeneralClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.Config.App.ServerName,
				Subject:   strconv.Itoa(int(userID)),
				ExpiresAt: jwt.NewNumericDate(duration),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.String(),
				Audience: jwt.ClaimStrings{
					config.Config.App.Server.URL,
				},
			},
			Data: jwtAuthClaimData,
		}
	}

	jwtRememberClaimData, ok := data.(JwtRememberClaim)
	if ok {
		claims = &jwtGeneralClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    config.Config.App.ServerName,
				Subject:   strconv.Itoa(int(userID)),
				ExpiresAt: jwt.NewNumericDate(duration),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        uuid.String(),
				Audience: jwt.ClaimStrings{
					config.Config.App.Server.URL,
				},
			},
			Data: jwtRememberClaimData,
		}
	}

	var rsaPrivateKey *rsa.PrivateKey

	bytePrivateKey := []byte(privateKey)

	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bytePrivateKey)

	if err != nil {
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	t, err := token.SignedString(rsaPrivateKey)

	if err != nil {
		return ""
	}

	return t
}

func ValidateRS512Token(publicKey, token string) (*jwt.Token, error) {
	bytePublicKey := []byte(publicKey)

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(bytePublicKey)
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return rsaPublicKey, nil
	})
}
