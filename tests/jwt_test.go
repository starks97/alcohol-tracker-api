package tests

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/starks97/alcohol-tracker-api/config"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJwtToken(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	privateKey := cfg.AccessTokenPrivateKey
	bytesPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	assert.NoError(t, err)

	userID := "testUser"
	ttl := int64(1) // 1 minute

	tokenDetails, err := services.GenerateJwtToken(userID, ttl, string(bytesPrivateKey)) // Pass the byte slice
	assert.NoError(t, err)

	assert.NotNil(t, tokenDetails.Token)
	assert.NotNil(t, tokenDetails.ExpiresIn)
	assert.Equal(t, userID, tokenDetails.UserID)

	// Verify the expiry time is within a reasonable range.
	expectedExpiry := time.Now().Add(time.Minute).Unix()
	actualExpiry := *tokenDetails.ExpiresIn
	assert.GreaterOrEqual(t, expectedExpiry+10, actualExpiry) // allow 10 seconds of clock drift
	assert.LessOrEqual(t, expectedExpiry-10, actualExpiry)
}

func TestVerifyJwtToken(t *testing.T) {
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	privateKey := cfg.AccessTokenPrivateKey
	publicKey := cfg.AccessTokenPublicKey

	bytesPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	assert.NoError(t, err)
	bytesPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	assert.NoError(t, err)

	userID := "testUser"
	ttl := int64(1)

	tokenDetails, err := services.GenerateJwtToken(userID, ttl, string(bytesPrivateKey)) // Pass the byte slice
	assert.NoError(t, err)

	verifiedTokenDetails, err := services.VerifyJwtToken(string(bytesPublicKey), *tokenDetails.Token) // Pass the byte slice
	assert.NoError(t, err)
	assert.Equal(t, userID, verifiedTokenDetails.UserID)
	assert.Equal(t, tokenDetails.TokenUUID, verifiedTokenDetails.TokenUUID)

	// Test with an invalid token
	_, err = services.VerifyJwtToken(string(bytesPublicKey), "invalid_token")
	assert.Error(t, err)

	// Test with an invalid public key
	invalidPublicKey := "invalid_public_key"
	invalidBytesPublicKey, err := base64.StdEncoding.DecodeString(invalidPublicKey)
	assert.NoError(t, err)

	_, err = services.VerifyJwtToken(string(invalidBytesPublicKey), *tokenDetails.Token)
	assert.Error(t, err)

	// Test with expired token
	expiredToken, err := services.GenerateJwtToken(userID, 0, string(bytesPrivateKey)) // Expired immediately
	assert.NoError(t, err)
	time.Sleep(1 * time.Second) // Ensure the token expires
	_, err = services.VerifyJwtToken(string(bytesPublicKey), *expiredToken.Token)
	assert.Error(t, err)
}

func TestGenerateJwtToken_InvalidPrivateKey(t *testing.T) {
	invalidPrivateKey := "invalid_key"
	invalidBytesPrivateKey, err := base64.StdEncoding.DecodeString(invalidPrivateKey)
	assert.NoError(t, err)
	_, err = services.GenerateJwtToken("user", 60, string(invalidBytesPrivateKey))
	assert.Error(t, err)
}

func TestVerifyJwtToken_InvalidPublicKey(t *testing.T) {
	invalidPublicKey := "invalid_key"
	invalidBytesPublicKey, err := base64.StdEncoding.DecodeString(invalidPublicKey)
	assert.NoError(t, err)

	_, err = services.VerifyJwtToken(string(invalidBytesPublicKey), "token")
	assert.Error(t, err)
}
