package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/starks97/alcohol-tracker-api/internal/models"
)

func GenerateJwtToken(UserID string, ttl int64, privateKey string) (models.TokenDetails, error) {
	bytesPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("Error decoding base64: %v", err)
	}
	timeStamp := time.Now().Unix()

	expirationTime := time.Now().Add(time.Duration(ttl) * time.Minute).Unix()

	tokenDetails := models.TokenDetails{
		UserID:    UserID,
		TokenUUID: uuid.New().String(),
		Token:     nil,
	}

	tokenDetails.ExpiresIn = &expirationTime

	//token claims
	claims := models.TokenClaims{
		Sub:       tokenDetails.UserID,
		TokenUUID: tokenDetails.TokenUUID,
		Exp:       *tokenDetails.ExpiresIn,
		Iat:       timeStamp,
		Nbf:       timeStamp,
	}

	// Parse the private key to use in signing
	privateKeyObj, err := jwt.ParseRSAPrivateKeyFromPEM(bytesPrivateKey)
	if err != nil {
		return tokenDetails, err
	}

	// Create the token with the claims and sign with the private key using RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := token.SignedString(privateKeyObj) // Use the decoded private key for signing
	if err != nil {
		return tokenDetails, fmt.Errorf("Error signing token: %v", err)
	}

	tokenDetails.Token = &tokenString
	return tokenDetails, nil
}

func VerifyJwtToken(publicKey string, token string) (models.TokenDetails, error) {
	bytesPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("Error decoding base64: %v", err)
	}

	publicKeyObj, err := jwt.ParseRSAPublicKeyFromPEM(bytesPublicKey)
	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("error parsing public key: %v", err)
	}

	// Verify and decode the token
	claims := models.TokenClaims{}

	tokenObj, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKeyObj, nil
	})

	if err != nil {
		return models.TokenDetails{}, fmt.Errorf("error verifying token: %v", err)
	}

	if !tokenObj.Valid {
		return models.TokenDetails{}, fmt.Errorf("invalid token")
	}

	// Return token details
	tokenDetails := models.TokenDetails{
		TokenUUID: claims.TokenUUID,
		UserID:    claims.Sub,
	}
	return tokenDetails, nil

}
