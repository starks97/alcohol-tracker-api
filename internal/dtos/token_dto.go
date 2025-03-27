package dtos

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenDetailsDto represents the details of a JWT token.
type TokenDetailsDto struct {
	// Token is the JWT token string. It's a pointer to allow for optional values.
	Token *string `json:"token,omitempty" db:"token"`

	// TokenUUID is a unique identifier for the token.
	TokenUUID uuid.UUID `json:"token_uuid" db:"token_uuid"`

	// UserID is the identifier of the user associated with the token.
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	// ExpiresIn is the Unix timestamp representing the token's expiration time. It's a pointer to allow for optional values.
	ExpiresIn *int64 `json:"expires_in,omitempty" db:"expires_in"`
}

// TokenClaimsDto represents the claims within a JWT token.
type TokenClaimsDto struct {
	// Sub is the subject of the token (typically the user ID).
	Sub string `json:"sub" db:"sub"`

	// TokenUUID is the unique identifier of the token.
	TokenUUID string `json:"token_uuid" db:"token_uuid"`

	// Exp is the expiration time of the token as a Unix timestamp.
	Exp int64 `json:"exp" db:"exp"`

	// Iat is the issued-at time of the token as a Unix timestamp.
	Iat int64 `json:"iat" db:"iat"`

	// Nbf is the not-before time of the token as a Unix timestamp.
	Nbf int64 `json:"nbf" db:"nbf"`
}

// GetExpirationTime returns the expiration time claim (exp) as a jwt.NumericDate.
// Returns nil if no expiration time is set.
func (tc TokenClaimsDto) GetExpirationTime() (*jwt.NumericDate, error) {
	if tc.Exp == 0 {
		return nil, nil // No expiration time set
	}
	expirationTime := time.Unix(tc.Exp, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(expirationTime), nil
}

// GetIssuedAt returns the issued-at time claim (iat) as a jwt.NumericDate.
// Returns nil if no issued-at time is set.
func (tc TokenClaimsDto) GetIssuedAt() (*jwt.NumericDate, error) {
	if tc.Iat == 0 {
		return nil, nil // No issued time set
	}
	// Convert the Iat int64 to time.Time
	issuedAtTime := time.Unix(tc.Iat, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(issuedAtTime), nil
}

// GetNotBefore returns the not-before time claim (nbf) as a jwt.NumericDate.
// Returns nil if no not-before time is set.
func (tc TokenClaimsDto) GetNotBefore() (*jwt.NumericDate, error) {
	if tc.Nbf == 0 {
		return nil, nil // No not before time set
	}
	// Convert the Nbf int64 to time.Time
	notBeforeTime := time.Unix(tc.Nbf, 0) // Convert Unix timestamp to time.Time
	return jwt.NewNumericDate(notBeforeTime), nil
}

// GetIssuer returns the issuer claim (iss).
// In this implementation, the issuer is not part of the claims, so it returns an empty string and nil error.
func (tc TokenClaimsDto) GetIssuer() (string, error) {
	// The issuer might not be part of your claims, so we return an empty string or default value
	return "", nil
}

// GetSubject returns the subject claim (sub).
func (tc TokenClaimsDto) GetSubject() (string, error) {
	return tc.Sub, nil
}

// GetAudience returns the audience claim (aud) as jwt.ClaimStrings.
// If the audience is not present, it returns an empty slice and nil error.
func (tc TokenClaimsDto) GetAudience() (jwt.ClaimStrings, error) {
	// If you don't have audience in your token, just return an empty slice
	return jwt.ClaimStrings{}, nil
}
