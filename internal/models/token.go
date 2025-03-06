package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenDetails struct {
	Token     *string `json:"token,omitempty" db:"token"`
	TokenUUID string  `json:"token_uuid,omitempty" db:"token_uuid"`
	UserID    string  `json:"user_id,omitempty" db:"user_id"`
	ExpiresIn *int64  `json:"expires_in,omitempty" db:"expires_in"`
}

type TokenClaims struct {
	Sub       string `json:"sub,omitempty" db:"sub"`
	TokenUUID string `json:"token_uuid,omitempty" db:"token_uuid"`
	Exp       int64  `json:"exp,omitempty" db:"exp"`
	Iat       int64  `json:"iat,omitempty" db:"iat"`
	Nbf       int64  `json:"nbf,omitempty" db:"nbf"`
}

func (tc TokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if tc.Exp == 0 {
		return nil, nil // No expiration time set
	}
	expirationTime := time.Unix(tc.Exp, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(expirationTime), nil

}

func (tc TokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	if tc.Iat == 0 {
		return nil, nil // No issued time set
	}
	// Convert the Iat int64 to time.Time
	issuedAtTime := time.Unix(tc.Iat, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(issuedAtTime), nil
}

// GetNotBefore returns the Not Before (nbf) claim as NumericDate
func (tc TokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	if tc.Nbf == 0 {
		return nil, nil // No not before time set
	}
	// Convert the Nbf int64 to time.Time
	notBeforeTime := time.Unix(tc.Nbf, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(notBeforeTime), nil
}

// GetIssuer returns the Issuer (iss) claim
func (tc TokenClaims) GetIssuer() (string, error) {
	// The issuer might not be part of your claims, so we return an empty string or default value
	return "", nil
}

// GetSubject returns the Subject (sub) claim
func (tc TokenClaims) GetSubject() (string, error) {
	return tc.Sub, nil
}

// GetAudience returns the Audience (aud) claim
func (tc TokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	// If you don't have audience in your token, just return an empty slice
	return jwt.ClaimStrings{}, nil
}
