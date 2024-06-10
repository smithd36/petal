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
                http.Redirect(w, r, "/login", http.StatusSeeOther)
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

        // Store user ID and username in context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        ctx = context.WithValue(ctx, "username", claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
