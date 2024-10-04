package middleware

import (
	"DoAn/auth"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func JWTMiddleware(next http.Handler, svc auth.AuthenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		username, err := svc.VerifyToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		fmt.Println("Username: ", username)
		setSessionHandler(w, r, username)

		next.ServeHTTP(w, r)
	})
}
func JWTMiddlewareRefreshToken(next http.Handler, svc auth.AuthenService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		username, err := svc.VerifyRefreshToken(context.Background(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		fmt.Println("Username: ", username)
		setSessionHandler(w, r, username)

		next.ServeHTTP(w, r)
	})
}

func setSessionHandler(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := store.Get(r, "my-session")

	// Save a variable (e.g., username) in the session
	session.Values["username"] = username

	// Save the session
	err := session.Save(r, w)
	if err != nil {
		return
	}

	//fmt.Fprintln(w, "Username saved in session!")

}

func GetSessionHandler(r *http.Request) string {
	session, _ := store.Get(r, "my-session")

	// Retrieve the variable from the session
	username, ok := session.Values["username"].(string)
	if !ok {
		return ""
	}
	return username
}
