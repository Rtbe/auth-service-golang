package model

//TokenPair is a type for api JSON representations of request with provided pair of access/refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

//RefreshToken is a type for api JSON representation of request with provided refresh token.
type RefreshToken struct {
	Token string `json:"refresh_token"`
}
