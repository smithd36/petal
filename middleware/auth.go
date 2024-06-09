package middleware

import (
    "context"
    "net/http"
    "github.com/dgrijalva/jwt-go"
    "github.com/smithd36/petal/utils"
)

func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Missing token", http.StatusUnauthorized)
                return
            }
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        tokenString := cookie.Value
        claims := &utils.Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return utils.JwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
