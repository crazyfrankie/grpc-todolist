package mws

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SecretKey = []byte("sD4pP0qA8sO6fZ2fF9iZ5lN9nM1rF3vL")
)

type AuthBuild struct {
	paths map[string]struct{}
}

func NewAuthBuilder() *AuthBuild {
	return &AuthBuild{paths: make(map[string]struct{})}
}

func (a *AuthBuild) IgnorePath(path string) *AuthBuild {
	a.paths[path] = struct{}{}
	return a
}

func (a *AuthBuild) Auth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := a.paths[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("todolist_auth")
		if !cookie.HttpOnly || !cookie.Secure {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var token string
		if err == nil {
			token = cookie.Value
		} else {
			tokenHeader := r.Header.Get("Authorization")
			token = extractToken(tokenHeader)
		}

		claims, err := parseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}

		// 严重的安全问题
		if claims["user_agent"] != r.UserAgent() {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}

		userId := claims["user_id"].(float64)
		uid := strconv.FormatFloat(userId, 'f', 0, 64)
		r = r.WithContext(context.WithValue(r.Context(), "user_id", uid))

		next.ServeHTTP(w, r)
	}
}

func extractToken(token string) string {
	if token == "" {
		return ""
	}

	strs := strings.Split(token, " ")
	if strs[0] == "Bearer" {
		return strs[1]
	}

	return ""
}

func parseToken(token string) (jwt.MapClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.MapClaims); ok && tokenClaims.Valid {
			return *claims, nil
		}
	}

	return nil, errors.New("token is invalid")
}

func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "todolist_auth",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24小时
	})
}
