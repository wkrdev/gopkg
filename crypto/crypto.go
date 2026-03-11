package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func JwtEncode(claims any, secret string) (string, error) {
	var mapJwt map[string]any
	bytes, _ := json.Marshal(claims)
	if err := json.Unmarshal(bytes, &mapJwt); err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims(mapJwt),
	)
	token, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func JwtDecode[T any](token string, secret string) (T, error) {
	var result T
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !jwtToken.Valid {
		return result, fmt.Errorf("invalid token: %w", err)
	}
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		bytes, _ := json.Marshal(claims)
		if err := json.Unmarshal(bytes, &result); err != nil {
			return result, err
		}
		return result, nil
	}
	return result, errors.New("failed to cast claims")
}

func JwtParse[T any](token string) T {
	var result T
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return result
	}
	decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return result
	}
	json.Unmarshal(decoded, &result)
	return result
}

func HashPassword(pass string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	return string(bytes), err
}

func ComparePassword(hash string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

const CHAR_HASH string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func StringHash(n uint8) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHAR_HASH[rand.IntN(len(CHAR_HASH))]
	}
	return string(b)
}

func SHA256Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func MD5Hash(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}
