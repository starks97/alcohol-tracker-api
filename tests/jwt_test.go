package tests

import (
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
	"github.com/starks97/alcohol-tracker-api/internal/services"
	"github.com/stretchr/testify/assert"
)

var testTokenPrivateBase64 string = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2QUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktZd2dnU2lBZ0VBQW9JQkFRQ3pNR0w5b3FsSjAyM0QKYWdjbzZpUjJDc1Y2eDBodzVNRTI2UUxkMUhtUWJnN1pZcGh5bVdKTUJYUXpxZ1ZiRnZhb01ybW1GZzZiWENFQgorbnZVODlZK2NxYjBXZHJKL2ZZb0VwN2YxWXcwUkhidGtpSmdBR2hlV1RUNHFDUkdYUWs4Y1F6Z2tIbDQ4SW9TCmxJdTViOFJKZGxyRk81YjRHdkpsWVhpeXNsaWFqcHRXY1A2MFRiaEFSSndhL0pjQTJVSkhqTzd2NnFCd2pjdlIKU2hheUwwTTJCckFqczYzTG4yb3U4Vm9KZTVMY3R3Uys2b0JON2RrR2Fzdm1NNGp1RmhVeXZYdm53SkhidU9yYQpMaFczWko3T09GWHk4Vk11UUYwM25pMFZVSnE2ak5YZTFkSzJnRzROcCtiM3VIdzFNdTZxMG5mdHpLMzlVakZvCkx3RkIvRC9qQWdNQkFBRUNnZ0VBSEIzbUVvMXhDZHdLdDZTTi9nNExiWmhRRjIxb3dRb3NCVDAzelc0WEt5SVIKNDJ0MHAxckpFVXV6eVoyT25KWDBXejBtWTFqSHJ2b2NWYzZqbXEwdU8zdExGa0Y1TXNQT1djaGVOSm95RjB0OAo2OWRIM0krRDBQWW5lVE1OQ2h0MEpROUtLWHlTQ3ZlWGVzWGpUTlFzVlNpa29wa3duYnJBdVVhN3BUS1Y4NTVVCnFRWm0xaDZrYkdnR3JkYmJzaFZ5K1B1L3h0a3I1Y28zU3N2Z1pHQXpBalYzOGNDTHdGYUpVWEpTYjNxZkFwQSsKREphdTk1cWhRN2tCLzROWGRISkhHOHF6ZnFPOU0zRlB5T0FmVXc5THZaUk1iZ3lNcW90ek5JMnBlTzRFU1ptRQpCcmQxYVZlV1htR245ZThpeTRhU0d1WEZ2eHNzaE1zN2RRczBwN2o1cVFLQmdRRFoyZEpzblZnK3kwRE9JYlVVCmZpSFpzOTRUQlJWa1hvaVIyYTVwV1Nqd2dxaUpMTzFzZFZFc3Y1bzVSMDlWRWFiMXBBNGpCU2ZkOS90Y1dJcUsKQzQwMmRQZ2pQcXB0djA5Z3MvUUNiZzkzVTVDUm1GTEFHNVJFcS9NWk9mVnNTYmZFZzVyenNYeEY2VVJnYmZhZwpZd1FJRzA4ZzM1YnBIbjVpWnhXS3BpRmxlUUtCZ1FEU2tWNEN5UkM0V2VnQ0RrMDFzRHRpU3ZMd04xV1NJWG9YCkp2MTRNVGR1Q1cyQ0ZYZ2pFWkRBNXVFSVY1bGd4d2s1Qk1Mb25EZHFJMkNFdzBMdUZpalcwcnlRWUxhTkZlUEEKbHFCdmZYN2JaOGpXdXMreXJJRGNYUmd0REN5UHY5MHVJTnB6MnRVOUJFbDZSb1AzUTA1V3ZDQm9INkFqbUFXYQpjaGNrcUtLRk93S0JnQUpDcUZSSUxhbzVJYXNCM29jZjUrb0NXOE9Hd2ZvUW9Rb2lZQlRudit6KzdoQytUcGhaCmYwWWZsdElSVTFsbW5YemUvdWFPSHlQR2R1MDJYZm1ndFE1am1FK0ZUdTlrbE1aRUY3d091RXBjcTV0WElVU3QKQkpUUjArdm5GZ3pSbHY2Sy93aVlSdG5TMmNyR1dWREF0a0gvUm9yb3h3QVVPT3Q4ZGxUQjlJYkJBb0dBVlVNcApod1UxY1FCdXNvNXA4eWhtRTFuMzN3NzQ1bEFKNk9BUDJLQk5LcEJFdUZ6TEphQVNOaG9HMnVMbHAwdFF6N2ErCjJZT1A2TGxrZHIyK0Z6di8wMlRIbDhxaGdLVnhjR1ZObDNlQWE0VXR6TTBlRnVKRTEzWVd1UDdwK0ZjZlQzTmMKSVhkbHl1dzJlSDJmSi9zbitIVDZ4azZ3QUZtcFF5MlpjMjJaU1VzQ2dZQlU4NHorTGNJckVENndTNlNMZGp0Wgpka3RZRDlrb0ZmaS8zd2pYVVQwcnloV00vMWY4a0hYc21vMmc4V1c1L1lMTCsyZGdRWnhEQ2JueEZiek0xSnhvCitiVFNmT1F6WEdwbjhBVDBZa1VQWnhScnRMUkhwbUtGN1RDaUpENXNQREdLbFRGVWd5MEYraDd4RStrZDFDczYKWFYzVHRoV3l0Z3RRV2EvOFZPM2JaQT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"

var testTokenPublicBase64 string = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUFzekJpL2FLcFNkTnR3Mm9IS09vawpkZ3JGZXNkSWNPVEJOdWtDM2RSNWtHNE8yV0tZY3BsaVRBVjBNNm9GV3hiMnFESzVwaFlPbTF3aEFmcDcxUFBXClBuS205Rm5heWYzMktCS2UzOVdNTkVSMjdaSWlZQUJvWGxrMCtLZ2tSbDBKUEhFTTRKQjVlUENLRXBTTHVXL0UKU1haYXhUdVcrQnJ5WldGNHNySlltbzZiVm5EK3RFMjRRRVNjR3Z5WEFObENSNHp1NytxZ2NJM0wwVW9Xc2k5RApOZ2F3STdPdHk1OXFMdkZhQ1h1UzNMY0V2dXFBVGUzWkJtckw1ak9JN2hZVk1yMTc1OENSMjdqcTJpNFZ0MlNlCnpqaFY4dkZUTGtCZE41NHRGVkNhdW96VjN0WFN0b0J1RGFmbTk3aDhOVEx1cXRKMzdjeXQvVkl4YUM4QlFmdy8KNHdJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg=="

func decodeBase64Key(encodedKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedKey)
}

func TestGenerateAndVerifyJwtToken(t *testing.T) {

	privateKey, err := decodeBase64Key(testTokenPrivateBase64)
	assert.NoError(t, err, "Failed to decode private key")

	publicKey, err := decodeBase64Key(testTokenPublicBase64)
	assert.NoError(t, err, "Failed to decode public key")

	userID := uuid.New()
	ttl := int64(30) // 30 minutes

	t.Run("Generate JWT Token", func(t *testing.T) {
		tokenDetails, err := services.GenerateJwtToken(userID, ttl, string(privateKey))
		assert.NoError(t, err, "Error generating JWT token")
		assert.NotNil(t, tokenDetails, "Token details should not be nil")
		assert.NotNil(t, tokenDetails.Token, "Generated token should not be nil")
		assert.NotEmpty(t, tokenDetails.TokenUUID, "TokenUUID should not be empty")
	})

	t.Run("Verify JWT Token", func(t *testing.T) {
		tokenDetails, _ := services.GenerateJwtToken(userID, ttl, string(privateKey))
		verifiedTokenDetails, err := services.VerifyJwtToken(string(publicKey), *tokenDetails.Token)

		assert.NoError(t, err, "Error verifying JWT token")
		assert.NotNil(t, verifiedTokenDetails, "Verified token details should not be nil")
		assert.Equal(t, userID, verifiedTokenDetails.UserID, "UserID should match")
		assert.Equal(t, tokenDetails.TokenUUID, verifiedTokenDetails.TokenUUID, "TokenUUID should match")
	})
}

func TestGenerateJwtToken_InvalidPrivateKey(t *testing.T) {
	userID := uuid.New()
	ttl := int64(30)

	_, err := services.GenerateJwtToken(userID, ttl, "invalid-private-key")
	assert.Error(t, err)
}

func TestVerifyJwtToken_InvalidPublicKey(t *testing.T) {
	userID := uuid.New()
	ttl := int64(30)
	tokenDetails, _ := services.GenerateJwtToken(userID, ttl, "some-private-key")

	_, err := services.VerifyJwtToken("invalid-public-key", *tokenDetails.Token)
	assert.Error(t, err)
}

func TestVerifyJwtToken_InvalidToken(t *testing.T) {
	_, err := services.VerifyJwtToken("some-public-key", "invalid-token")
	assert.Error(t, err)
}

func TestVerifyJwtToken_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	ttl := int64(-30) // Expired 30 minutes ago
	tokenDetails, _ := services.GenerateJwtToken(userID, ttl, "some-private-key")

	_, err := services.VerifyJwtToken("some-public-key", *tokenDetails.Token)
	assert.Error(t, err)
}

func TestVerifyJwtToken_ModifiedClaims(t *testing.T) {
	userID := uuid.New()
	ttl := int64(30)
	tokenDetails, _ := services.GenerateJwtToken(userID, ttl, "some-private-key")
	modifiedToken := *tokenDetails.Token + "modified"

	_, err := services.VerifyJwtToken("some-public-key", modifiedToken)
	assert.Error(t, err)
}
