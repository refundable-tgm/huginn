package rest

import (
	"bufio"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// pathAccessSecret is the file path to the secret string to encode access tokens
const pathAccessSecret = "/vol/secrets/access_secret.env"
// pathRefreshSecret is the file path to the secret string to encode access tokens
const pathRefreshSecret = "/vol/secrets/refresh_secret.env"
// accessSecretLength is the length of the access secret
const accessSecretLength = 32
// refreshSecretLength is the length of the refresh secret
const refreshSecretLength = 64

// accessDuration is the time for which an access token is valid (default 15 mins)
const accessDuration = time.Minute * 15
// refreshDuration is the time for which a refresh token is valid (default 7 days)
const refreshDuration = time.Hour * 24 * 7

// accessSecret is the secret used to encode access tokens
var accessSecret string
// refreshSecret is the secret used to encode refresh tokens
var refreshSecret string

// activeTokens stores all token information of active tokens
var activeTokens map[string]EntityInformation

// Token represents a token pair
type Token struct {
	// AccessToken is the access token itself
	AccessToken    string
	// RefreshToken is the refresh token itself
	RefreshToken   string
	// AccessUUID is the uuid the access token is referenced by
	AccessUUID     string
	// RefreshUUID is the uuid the refresh token is referenced by
	RefreshUUID    string
	// AccessExpires is the date the access tokens expires at
	AccessExpires  int64
	// RefreshExpires is the date the refresh tokens expires at
	RefreshExpires int64
}

// AccessToken represents the further information about an access token
type AccessToken struct {
	// AccessUUID is the uuid of the access token
	AccessUUID string
	// Username is the username of the user this token belongs to
	Username   string
}

// EntityInformation represents information about tokens
type EntityInformation struct {
	// Username identifies the user this token belongs to
	Username  string
	// ExpiresAt marks the time the token expires at
	ExpiresAt time.Time
}

// InitTokenManager initializes the token manager
// it reads or generates both secrets, creates the map of active tokens, and starts the thread to remove expired tokens
func InitTokenManager() {
	readRefreshSecret()
	readAccessSecret()
	activeTokens = make(map[string]EntityInformation)
	go ttlCheck()
}

// CreateToken creates a token pair based on a username
func CreateToken(username string) (*Token, error) {
	token := &Token{}
	token.AccessExpires = time.Now().Add(accessDuration).Unix()
	token.AccessUUID = uuid.New().String()

	token.RefreshExpires = time.Now().Add(refreshDuration).Unix()
	token.RefreshUUID = uuid.New().String()

	// Access Token
	acClaims := jwt.MapClaims{}
	acClaims["authorized"] = true
	acClaims["access_uuid"] = token.AccessUUID
	acClaims["username"] = username
	acClaims["exp"] = time.Now().Add(accessDuration).Unix()
	acBase := jwt.NewWithClaims(jwt.SigningMethodHS256, acClaims)
	var err error
	token.AccessToken, err = acBase.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refClaims := jwt.MapClaims{}
	refClaims["refresh_uuid"] = token.RefreshUUID
	refClaims["username"] = username
	refClaims["exp"] = token.RefreshExpires
	refBase := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	token.RefreshToken, err = refBase.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return token, nil
}

// SaveToken saves a token in the active token map with its corresponding username as key
func SaveToken(username string, token *Token) {
	acExp := time.Unix(token.AccessExpires, 0)
	refExp := time.Unix(token.RefreshExpires, 0)

	activeTokens[token.AccessUUID] = EntityInformation{username, acExp}
	activeTokens[token.RefreshUUID] = EntityInformation{username, refExp}
}

// ExtractToken parses the token string out of a request
func ExtractToken(r *http.Request) string {
	bear := r.Header.Get("Authorization")
	split := strings.Split(bear, " ")
	if len(split) == 2 {
		return split[1]
	}
	return ""
}

// VerifyToken verifies that the token provided in the request originates from this API
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	extr := ExtractToken(r)
	token, err := jwt.Parse(extr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid checks whether a token is still valid, therefore also verifies it
func TokenValid(r *http.Request) (bool, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false, err
	}
	return true, nil
}

// ExtractTokenMeta extracts the meta information encoded in the token and returns both uuid and username
func ExtractTokenMeta(r *http.Request) (*AccessToken, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		acccessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		username, ok := claims["username"].(string)
		if !ok {
			return nil, err
		}
		return &AccessToken{
			AccessUUID: acccessUUID,
			Username:   username,
		}, nil
	}
	return nil, err
}

// DeleteToken deletes a token
func DeleteToken(uuid string) {
	delete(activeTokens, uuid)
}

// readAccessSecret manages the refresh secret generation
func readAccessSecret() {
	if _, err := os.Stat(pathAccessSecret); os.IsNotExist(err) {
		const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
		secret := make([]byte, accessSecretLength)
		for i := range secret {
			secret[i] = char[rand.Intn(len(char))]
		}
		accessSecret = string(secret)
		file, err := os.Create(pathAccessSecret)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(string(secret))
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	} else {
		file, err := os.Open(pathAccessSecret)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		secret, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		accessSecret = string(secret)
	}
}

// readRefreshSecret manages the refresh secret generation
func readRefreshSecret() {
	if _, err := os.Stat(pathRefreshSecret); os.IsNotExist(err) {
		const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
		secret := make([]byte, refreshSecretLength)
		for i := range secret {
			secret[i] = char[rand.Intn(len(char))]
		}
		refreshSecret = string(secret)
		file, err := os.Create(pathRefreshSecret)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(string(secret))
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
	} else {
		file, err := os.Open(pathRefreshSecret)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		secret, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		refreshSecret = string(secret)
	}
}

// ttlCheck checks whether tokens expired and removes them
func ttlCheck() {
	for {
		now := time.Now()
		for key, value := range activeTokens {
			if value.ExpiresAt.Before(now) {
				delete(activeTokens, key)
			}
		}
		time.Sleep(time.Minute)
	}
}
