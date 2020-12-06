package entity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//AccessToken is an representation of jwt access token.
type AccessToken struct {
	Token     string
	ExpiresAt int64
}

//RefreshToken is an representation of jwt refresh token that will be stored in mongoDB.
type RefreshToken struct {
	UUID      string `bson:"_id"`
	UserID    string `bson:"user_id"`
	Token     string `bson:"token"`
	ExpiresAt int64  `bson:"expires_at"`
	Used      bool   `bson:"used"`
}

//TokenPair is an representation of access and refresh token pair.
type TokenPair struct {
	AccessToken  AccessToken
	RefreshToken RefreshToken
}

//CustomClaimsAcessToken is a set of additional claims for jwt access token.
type CustomClaimsAcessToken struct {
	User_id string
	//This field helps to bind access token to refresh token.
	Refresh_uuid string
	jwt.StandardClaims
}

//CustomClaimsRefreshToken is a Set of additional claims for jwt refresh token.
type CustomClaimsRefreshToken struct {
	User_id string
	UUID    string
	jwt.StandardClaims
}

//CreateTokenPair creates a new pair of access and refresh tokens.
func CreateTokenPair(userID string) (*TokenPair, error) {

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 7).Unix()
	refreshTokenUUID := uuid.New().String()
	refreshToken, err := createRefreshToken(userID, refreshTokenUUID, refreshTokenExp)
	if err != nil {
		return nil, err
	}

	accessTokenExp := time.Now().Add(time.Hour * 24 * 7).Unix()
	accessToken, err := createAccessToken(userID, refreshTokenUUID, accessTokenExp)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	tokens := &TokenPair{
		AccessToken: AccessToken{
			Token:     accessToken,
			ExpiresAt: accessTokenExp,
		},
		RefreshToken: RefreshToken{
			UserID:    userID,
			UUID:      refreshTokenUUID,
			Token:     refreshToken,
			ExpiresAt: refreshTokenExp,
			Used:      false,
		},
	}
	return tokens, nil
}

//createAccessToken creates a new jwt refresh token.
func createAccessToken(userID string, refreshUUID string, expires int64) (string, error) {
	claims := CustomClaimsAcessToken{
		User_id:      userID,
		Refresh_uuid: refreshUUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

//createRefreshToken creates a new jwt refresh token.
func createRefreshToken(userID, UUID string, expires int64) (string, error) {
	claims := CustomClaimsRefreshToken{
		User_id: userID,
		UUID:    UUID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

//ParseRefreshToken checks validity of refresh token and returns it`s claims.
func ParseRefreshToken(tokenString string) (*CustomClaimsRefreshToken, error) {
	claims := &CustomClaimsRefreshToken{}
	token, err := ParseJWTToken(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("Refresh token is not valid. %s", err.Error())
	}

	claims, ok := token.Claims.(*CustomClaimsRefreshToken)
	if !ok || !token.Valid {
		return nil, errors.New("Refresh token is not valid")
	}
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("Refresh token is expired")
	}
	return claims, nil
}

//ParseAccessToken checks validity of access token and returns it`s claims.
func ParseAccessToken(tokenString string) (*CustomClaimsAcessToken, error) {
	claims := &CustomClaimsAcessToken{}
	token, err := ParseJWTToken(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("Access token is not valid. %s", err.Error())
	}

	claims, ok := token.Claims.(*CustomClaimsAcessToken)
	if err != nil {
		return nil, fmt.Errorf("Refresh token is not valid: %s", err.Error())
	}

	if !ok || !token.Valid {
		return nil, errors.New("Access token is not valid")
	}
	return claims, nil
}

//ParseJWTToken parses token string into jwt token.
func ParseJWTToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Token is not valid: %s", err.Error())
	}
	return token, nil
}

//GenerateHash generates bcrypt hash.
func GenerateHash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), 12)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	res := string(hash)
	return res, nil
}

//EncodeToken64 encodes into base64 encoding.
func EncodeToken64(token string) string {
	refreshToken := base64.StdEncoding.EncodeToString([]byte(token))
	return refreshToken
}

//DecodeToken64 decodes from base64 encoding.
func DecodeToken64(token string) (string, error) {
	refreshToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("Error decoding access token")
	}

	return string(refreshToken), nil
}
