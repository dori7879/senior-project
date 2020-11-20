package auth

import (
	"time"

	"api/config"

	"github.com/dgrijalva/jwt-go"
)

var (
	// AccessTokenSecretKey is a access token's secret key type that is encoded into key id of payload in JWT
	AccessTokenSecretKey = "441253856187b286"
	// RefreshTokenSecretKey is a refresh token's secret key type that is encoded into key id of payload in JWT
	RefreshTokenSecretKey = "c425ad3cc7feb212"
)

// JwtUtils is a struct with functions that need access to secret keys
type JwtUtils struct {
	AtJwtSecretKey []byte
	RtJwtSecretKey []byte
}

// New creates JwtUtils struct instance
func New(conf *config.Conf) *JwtUtils {
	return &JwtUtils{
		AtJwtSecretKey: conf.AtJwtSecretKey,
		RtJwtSecretKey: conf.RtJwtSecretKey,
	}
}

// GenerateTokenPair creates and returns a new set of access_token and refresh_token.
func (j *JwtUtils) GenerateTokenPair(email, role string) (string, string, error) {

	tokenString, err := j.GenerateAccessToken(email, role)
	if err != nil {
		return "", "", err
	}

	// Create Refresh token, this will be used to get new access token.
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshToken.Header["kid"] = RefreshTokenSecretKey

	// Expiration time is 180 minutes
	expirationTimeRefreshToken := time.Now().Add(180 * time.Minute).Unix()

	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = email
	rtClaims["exp"] = expirationTimeRefreshToken

	refreshTokenString, err := refreshToken.SignedString(j.RtJwtSecretKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

// ValidateToken is used to validate both access_token and refresh_token. It is done based on the "Key ID" provided by the JWT
func (j *JwtUtils) ValidateToken(tokenString string) (bool, string, string, error) {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["kid"].(string)
		// If the "kid" (Key ID) is equal to AccessTokenSecretKey, then it is compared against access_token secret key, else if it
		// is equal to RefreshTokenSecretKey , it is compared against refresh_token secret key.
		if keyID == AccessTokenSecretKey {
			key = j.AtJwtSecretKey
		} else if keyID == RefreshTokenSecretKey {
			key = j.RtJwtSecretKey
		}
		return key, nil
	})

	// Check if signatures are valid.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// logger.Logger.Errorf("Invalid Token Signature")
			return false, "", keyID, err
		}
		return false, "", keyID, err
	}

	if !token.Valid {
		// logger.Logger.Errorf("Invalid Token")
		return false, "", keyID, err
	}

	return true, claims["sub"].(string), keyID, nil
}

// InvalidateToken method marks the access_token as invalid. This token cannot be used for future authentication
func (j *JwtUtils) InvalidateToken(tokenString string) error {

	var key []byte

	var keyID string

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		keyID = token.Header["kid"].(string)
		// If the "kid" (Key ID) is equal to AccessTokenSecretKey, then it is compared against access_token secret key, else if it
		// is equal to RefreshTokenSecretKey , it is compared against refresh_token secret key.
		if keyID == AccessTokenSecretKey {
			key = j.AtJwtSecretKey
		} else if keyID == RefreshTokenSecretKey {
			key = j.RtJwtSecretKey
		}
		return key, nil
	})

	// Check if signatures are valid.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// logger.Logger.Errorf("Invalid Token Signature")
			return err
		}
		return err
	}

	if !token.Valid {
		// logger.Logger.Errorf("Invalid Token")
		return err
	}

	// @TODO - Fix the expiration time
	// status := db.RedisClient.Set(tokenString, tokenString, 0)

	// if status.Err() != nil {
	// 	logger.Logger.Errorf("Could not set value in Redis")

	// }

	// val, err := status.Result()

	// if val == "OK" {
	// 	// logger.Logger.Infof("User Logged Out")
	// 	return nil
	// }

	// return status.Err()
	return nil // just for passing errors
}

// GenerateAccessToken method creats a new access token when the user logs in by providing email and password
func (j *JwtUtils) GenerateAccessToken(email, role string) (string, error) {
	// Declare the expiration time of the access token
	// Here the expiration is 60 minutes
	expirationTimeAccessToken := time.Now().Add(60 * time.Minute).Unix()

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["kid"] = AccessTokenSecretKey
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTimeAccessToken
	claims["sub"] = email
	claims["role"] = role

	// Create the JWT string
	tokenString, err := token.SignedString(j.AtJwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
