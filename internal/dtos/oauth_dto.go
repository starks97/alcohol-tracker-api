package dtos

// OAuthDto represents the structure of user information received from Google's OAuth2 userinfo endpoint.
type OAuthDto struct {
	ID            string `json:"id"`             // Unique identifier for the Google user.
	Email         string `json:"email"`          // User's email address.
	VerifiedEmail bool   `json:"verified_email"` // Indicates whether the user's email address is verified.
	Name          string `json:"name"`           // User's full name.
	GivenName     string `json:"given_name"`     // User's given name (first name).
	FamilyName    string `json:"family_name"`    // User's family name (last name).
	Picture       string `json:"picture"`        // URL of the user's profile picture.
}
