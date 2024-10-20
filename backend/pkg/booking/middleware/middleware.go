package middleware

import (
	"DoAn/common"
	"DoAn/pb"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func JWTMiddleware(next http.Handler, svc *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		client := pb.NewUserServiceClient(svc)

		userId, err := client.VerifyToken(context.Background(), &pb.VerifyTokenRequest{Token: tokenString})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		fmt.Println("userId: ", userId)

		accountId, ok := strconv.Atoi(userId.Value)
		if ok != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		checkBarber, err := client.CheckExistedBarber(context.Background(), &pb.CheckExistedBarberRequest{Id: int32(accountId)})
		if err != nil {
			fmt.Printf("error when checking account: %v", err)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		if checkBarber == nil ||
			checkBarber.Value == common.RoleUnknown {
			http.Error(w, "Unauthorized account", http.StatusUnauthorized)
			return
		}

		setSessionHandler(w, r, userId.Value)

		next.ServeHTTP(w, r)
	})
}
func JWTMiddlewareGetListBooking(next http.Handler, svc *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		client := pb.NewUserServiceClient(svc)

		userId, err := client.VerifyToken(context.Background(), &pb.VerifyTokenRequest{Token: tokenString})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		accountId, ok := strconv.Atoi(userId.Value)
		if ok != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		checkBarber, err := client.CheckExistedBarber(context.Background(), &pb.CheckExistedBarberRequest{Id: int32(accountId)})
		if err != nil {
			fmt.Printf("error when checking account: %v", err)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		if checkBarber == nil ||
			checkBarber.Value == common.RoleUnknown {
			http.Error(w, "Unauthorized account", http.StatusUnauthorized)
			return
		}

		setSessionHandler(w, r, userId.Value)
		r.Header.Set("role_name", checkBarber.GetValue())
		r.Header.Set("account_id", userId.Value)
		next.ServeHTTP(w, r)
	})
}
func JWTMiddlewareAdmin(next http.Handler, svc *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		client := pb.NewUserServiceClient(svc)

		userId, err := client.VerifyToken(context.Background(), &pb.VerifyTokenRequest{Token: tokenString})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		fmt.Println("userId: ", userId)

		accountId, ok := strconv.Atoi(userId.Value)
		if ok != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		checkBarber, err := client.CheckExistedBarber(context.Background(), &pb.CheckExistedBarberRequest{Id: int32(accountId)})
		if err != nil {
			fmt.Printf("error when checking account: %v", err)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		if checkBarber == nil ||
			checkBarber.Value == common.RoleUnknown ||
			checkBarber.Value == common.RoleBarber ||
			checkBarber.Value == common.RoleUser {
			http.Error(w, "Unauthorized account", http.StatusUnauthorized)
			return
		}

		setSessionHandler(w, r, userId.Value)

		next.ServeHTTP(w, r)
	})
}
func JWTMiddlewareBarber(next http.Handler, svc *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		client := pb.NewUserServiceClient(svc)

		userId, err := client.VerifyToken(context.Background(), &pb.VerifyTokenRequest{Token: tokenString})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		accountId, ok := strconv.Atoi(userId.Value)
		if ok != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		checkBarber, err := client.CheckExistedBarber(context.Background(), &pb.CheckExistedBarberRequest{Id: int32(accountId)})
		if err != nil {
			fmt.Printf("error when checking account: %v", err)
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		if checkBarber == nil ||
			checkBarber.Value == common.RoleUnknown ||
			checkBarber.Value == common.RoleUser {
			http.Error(w, "Unauthorized account", http.StatusUnauthorized)
			return
		}

		setSessionHandler(w, r, userId.Value)

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
