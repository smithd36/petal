package utils

import (
    "time"
    "github.com/dgrijalva/jwt-go"
    "os"
    "log"
)

var JwtKey []byte

func InitializeJWTKey() {
    JwtKey = []byte(os.Getenv("JWT_KEY"))
    if len(JwtKey) == 0 {
        log.Fatal("JWT_KEY environment variable not set")
    }
}

type Claims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    jwt.StandardClaims
}

func GenerateJWT(username string, userID int) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: username,
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
