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

const pathAccessSecret = "/vol/secrets/access_secret.env"
const pathRefreshSecret = "/vol/secrets/refresh_secret.env"
const accessSecretLength = 32
const refreshSecretLength = 64

const accessDuration = time.Minute * 15
const refreshDuration = time.Hour * 24 * 7

var accessSecret string
var refreshSecret string

var activeTokens map[string]EntityInformation

type Token struct {
	AccessToken string
	RefreshToken string
	AccessUUID string
	RefreshUUID string
	AccessExpires int64
	RefreshExpires int64
}

type AccessToken struct {
	AccessUUID string
	Username string
}

type EntityInformation struct {
	Username string
	ExpiresAt time.Time
}

func InitTokenManager() {
	readRefreshSecret()
	readAccessSecret()
	activeTokens = make(map[string]EntityInformation)
	go ttlCheck()
}



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

func SaveToken(username string, token *Token) {
	acExp := time.Unix(token.AccessExpires, 0)
	refExp := time.Unix(token.RefreshExpires, 0)

	activeTokens[token.AccessUUID] = EntityInformation{username, acExp}
	activeTokens[token.RefreshUUID] = EntityInformation{username, refExp}
}

func ExtractToken(r *http.Request) string {
	bear := r.Header.Get("Authorization")
	split := strings.Split(bear, " ")
	if len(split) == 2 {
		return split[1]
	}
	return ""
}

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
			Username: username,
		}, nil
	}
	return nil, err
}

func FetchAuth(auth *AccessToken) (username string, ok bool) {
	entity, ok := activeTokens[auth.AccessUUID]
	username = entity.Username
	if !ok {
		return "", ok
	}
	return
}

func DeleteToken(uuid string)  {
	delete(activeTokens, uuid)
}

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
		accessSecret = string(secret)
	}
}

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
		refreshSecret = string(secret)
	}
}

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