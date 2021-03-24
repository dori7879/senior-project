package jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dori7879/senior-project/api"
	"golang.org/x/crypto/bcrypt"
)

// Ensure service implements interface.
var _ api.AuthService = (*AuthService)(nil)

// AuthService represents a service for managing auths.
type AuthService struct {
	SignKey   []byte
	VerifyKey []byte
}

// NewGroupService returns a new instance of AuthService.
func NewGroupService(secretkey []byte) *AuthService {
	return &AuthService{SignKey: secretkey}
}

func (a *AuthService) Login(auth *api.Auth, user *api.User) (string, error) {
	// Perform basic field validation.
	if err := auth.Validate(); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(auth.Password)); err != nil {
		return "", err
	}

	token, err := a.generateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) generateToken(id int) (string, error) {

	// Declare the expiration time of the access token
	// Here the expiration is 60 minutes
	expirationTimeAccessToken := time.Now().Add(60 * time.Minute).Unix()

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expirationTimeAccessToken
	claims["sub"] = strconv.Itoa(id)

	// Create the JWT string
	tokenStr, err := token.SignedString(a.SignKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Validate is used to validate both access_token.
func (a *AuthService) Validate(tokenStr string) (int, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return a.VerifyKey, nil
	})

	// Check if signatures are valid.
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, nil
	}

	// Retrieve user ID
	id, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		return 0, err
	}

	return id, nil
}
