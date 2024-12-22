package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"

	account "DoAn"
	"DoAn/auth"
	"DoAn/entity"
	"DoAn/middleware"

	"github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"

	"golang.org/x/oauth2"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func init() {
	err := godotenv.Load("D:\\barbershop-system\\backend\\pkg\\account\\cmd\\account.env")
	if err != nil {
		log.Fatalln("Error loading account.env file")
	}
	fmt.Println("GOOGLE_CLIENT_ID: ", os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("GOOGLE_CLIENT_SECRET: ", os.Getenv("GOOGLE_CLIENT_SECRET"))

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8008/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type (
	RegisterUserRequest struct {
		user entity.Account
	}
	ChangePassFirstTimeRequest struct {
		Username        string `json:"username,omitempty"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	RegisterUserResponse struct {
		Msg string `json:"msg"`
		Err error  `json:"error,omitempty"`
	}

	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResponse struct {
		Username     string `json:"username"`
		RefreshToken string `json:"refreshToken"`
		AccessToken  string `json:"accessToken"`
	}
	RefreshResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	LogoutRequest struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}

	GoogleCallbackRequest struct {
		State string `json:"state"`
		Code  string `json:"code"`
	}

	UserProfileResponse struct {
		Data interface{} `json:"data"`
		Err  error       `json:"error,omitempty"`
	}

	Response struct {
		Message string      `json:"message"`
		Status  int         `json:"status"`
		Data    interface{} `json:"data"`
	}
)

func MakeLoginEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginRequest)
		accessToken, refreshToken, err := u.Login(ctx, entity.Account{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return nil, err
		}
		resp := LoginResponse{
			Username:     req.Username,
			RefreshToken: *refreshToken,
			AccessToken:  *accessToken,
		}
		return resp, err
	}
}
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err, ok := response.(error); ok {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func MakeRegisterUserEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterUserRequest)
		msg, err := u.Register(ctx, req.user)
		return RegisterUserResponse{
			Msg: msg,
			Err: err,
		}, err
	}
}

func MakeGetListBarberEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		data, err := u.GetAllBarber(ctx)
		return Response{
			Message: "success",
			Status:  200,
			Data:    data,
		}, err
	}
}

func MakeChangePassFirstTimeEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ChangePassFirstTimeRequest)
		msg, err := u.ChangePassFirstTime(ctx, req.Username, req.NewPassword)
		if err != nil {
			return Response{
				Message: err.Error(),
				Status:  500,
				Data:    nil,
			}, err
		}
		return Response{
			Message: msg.(string),
			Status:  200,
			Data:    nil,
		}, err
	}
}
func MakeRefreshEndpoints(u auth.AuthenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := request.(LogoutRequest)
		refreshToken, _, err := u.CreateRefreshToken(ctx, token.Username, token.Token)
		if err != nil {
			fmt.Printf("error while creating refresh token: %v", err)
			return nil, err
		}
		accessToken, err := u.CreateAccessToken(ctx, token.Username)
		if err != nil {
			fmt.Printf("error while creating access token: %v", err)
			return nil, err
		}
		if err != nil {
			return Response{
				Message: err.Error(),
				Status:  500,
				Data:    nil,
			}, err
		}
		return RefreshResponse{AccessToken: accessToken, RefreshToken: refreshToken}, err
	}
}

func MakeGetProfileEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// get request param in a string
		email := request.(string)
		data, err := u.GetProfile(ctx, email)
		return UserProfileResponse{
			Data: data,
			Err:  err,
		}, err
	}
}

func MakeLogoutEndpoints(u auth.AuthenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// get request param in a string
		token := request.(LogoutRequest)
		err = u.LogoutToken(ctx, token.Username, token.Token)
		if err != nil {
			return Response{
				Message: err.Error(),
				Status:  500,
				Data:    nil,
			}, err
		}
		return Response{
			Message: "success",
			Status:  200,
			Data:    nil,
		}, err
	}
}

func DecodeRegisterUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request.user); err != nil {
		return nil, err
	}
	if govalidator.IsNull(request.user.Username) || govalidator.IsNull(request.user.Email) || govalidator.IsNull(request.user.Password) {
		return nil, errors.New("all fields: username, email, password are required")
	}
	if !govalidator.IsEmail(request.user.Email) {
		return nil, errors.New("email is invalid")
	}
	return request, nil
}

func DecodeEmpty(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeGetProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	// userName, err := auth.GetSubjectFromToken(tokenString)
	// if err != nil {
	//	return nil, err
	// }
	tokenString = tokenString[len("Bearer "):]
	emailRP := r.URL.Query().Get("email")
	return emailRP, nil
}

func DecodeGoogleCallbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	return GoogleCallbackRequest{
		State: state,
		Code:  code,
	}, nil
}

func DecodeEmptyRequest(ctx context.Context, r *http.Request) (_ interface{}, err error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	// userName, err := auth.GetSubjectFromToken(tokenString)
	// if err != nil {
	//	return nil, err
	// }
	tokenString = tokenString[len("Bearer "):]
	username := middleware.GetSessionHandler(r)
	// err = auth.VerifyRefreshToken(tokenString)
	// if err != nil {
	//	return nil, err
	// }
	// collectionPostgres, err := database.ConnectPostgres()
	// if err != nil {
	//	return nil, err
	// }
	// var result entity.Account
	// err = collectionPostgres.QueryRow("SELECT id,username,email,role FROM account WHERE username=$1", userName).Scan(&result.ID, &result.Username, &result.Email, &result.Role)
	// if err != nil {
	//	return nil, err
	// }

	return LogoutRequest{
		Token:    tokenString,
		Username: username,
	}, nil
}

func DecodeChangePassFirstTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ChangePassFirstTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// func DecodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
//	tokenString := r.Header.Get("Authorization")
//	if tokenString == "" {
//		return nil, errors.New("token is required")
//	}
//	userName, err := auth.GetSubjectFromToken(tokenString)
//	if err != nil {
//		return nil, err
//	}
//	tokenString = tokenString[len("Bearer "):]
//	err = auth.VerifyToken(tokenString)
//	if err != nil {
//		return nil, err
//	}
//	collectionPostgres, err := database.ConnectPostgres()
//	if err != nil {
//		return nil, err
//	}
//	var result entity.Account
//	err = collectionPostgres.QueryRow("SELECT id,username,email,role FROM account WHERE username=$1", userName).Scan(&result.ID, &result.Username, &result.Email, &result.Role)
//	if err != nil {
//		return nil, err
//	}
//
//	return tokenString, nil
// }

func HandleMain(w http.ResponseWriter, request *http.Request) {
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`

	fmt.Fprintf(w, htmlIndex)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//
//	token, email, username, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
//	if err != nil {
//		fmt.Println(err.Error())
//		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
//		return
//	}
//	if email == nil || username == nil {
//		fmt.Errorf("email or username is nil")
//		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
//		return
//	}
//	result := map[string]interface{}{"token": token, "email": email, "username": username}
//	fmt.Println("email: ", *email)
//
//	//collectionPostgres, err := database.ConnectPostgres()
//	//if err != nil {
//	//	return
//	//}
//	//var user entity.Account
//	//errFindEmail := collectionPostgres.QueryRow("SELECT id,username,email,password,role FROM account WHERE email=$1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
//
//	data, err := u.GetProfile(ctx, email)
//	if err != nil {
//		fmt.Println("Account does not exist, create new account")
//		newUser := entity.Account{
//			Username:  *username,
//			Email:     *email,
//			Password:  "",
//			Role:      entity.RoleUser,
//			CreatedAt: time.Now().Unix(),
//			UpdatedAt: time.Now().Unix(),
//		}
//		req := RegisterUserRequest{user: newUser}
//		msg, err := u.Register(ctx, req.user)
//		if err != nil {
//			return
//		}
//		fmt.Println("msg:", msg)
//	}
//
//	//if errFindEmail != nil {
//	//	fmt.Println("Account does not exist, create new account")
//	//	password := ""
//	//	newUser := entity.Account{
//	//		Username:  *username,
//	//		Email:     *email,
//	//		Password:  password,
//	//		Role:      entity.RoleUser,
//	//		CreatedAt: time.Now().Unix(),
//	//		UpdatedAt: time.Now().Unix(),
//	//	}
//	//	_, errs := collectionPostgres.Exec("INSERT INTO account (username, email, password,role,phone_number,full_name,created_at,updated_at) VALUES ($1, $2, $3,$4,$5,$6,$7,$8)",
//	//		newUser.Username, newUser.Email, newUser.Password, newUser.Role, newUser.PhoneNumber, newUser.FullName, newUser.CreatedAt, newUser.UpdatedAt)
//	//	if errs != nil {
//	//		return
//	//	}
//	//}
//
//	accessToken, refreshToken, err := u.Login(ctx, entity.Account{
//		Username: *username,
//		Password: "",
//	})
//	if err != nil {
//		return nil, err
//	}
//	resp := LoginResponse{
//		Username:     *username,
//		RefreshToken: *refreshToken,
//		AccessToken:  accessToken,
//	}
//
//	fmt.Print("\nToken:")
//	fmt.Println(accessToken)
//
//	fmt.Fprintf(w, "Token: "+accessToken)
// }

// func GenerateToken(username string) (string, error) {
//	token, err := auth.CreateAccessToken(username)
//	if err != nil {
//		fmt.Println("errCreate")
//		fmt.Println(err)
//		return "", err
//	}
//
//	req := LoginRequest{
//		Username: username,
//		Password: username,
//	}
//	accessToken, refreshToken, err := u.Login(ctx, entity.Account{
//		Username: req.Username,
//		Password: req.Password,
//	})
//	if err != nil {
//		return nil, err
//	}
//	resp := LoginResponse{
//		Username:     req.Username,
//		RefreshToken: *refreshToken,
//		AccessToken:  *accessToken,
//	}
//	return resp, err
//	return token, nil
// }

func MakeGoogleCallbackEndpoints(u account.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		googleCallBack := request.(GoogleCallbackRequest)
		token, email, username, err := getUserInfo(googleCallBack.State, googleCallBack.Code)
		fmt.Println("token: ", token)
		if err != nil {
			fmt.Println(err.Error())
			// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		if email == nil || username == nil {
			fmt.Errorf("email or username is nil")
			// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		_, err = u.GetProfile(ctx, *email)
		if err != nil {
			fmt.Println("Account does not exist, create new account")
			newUser := entity.Account{
				Username:  *username,
				Email:     *email,
				Password:  "",
				Role:      entity.RoleUser,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}
			msg, errRegister := u.Register(ctx, newUser)
			if errRegister != nil {
				return
			}
			fmt.Println("msg:", msg)
		}

		accessToken, refreshToken, err := u.Login(ctx, entity.Account{
			Username: *username,
			Password: "",
		})
		if err != nil {
			return nil, err
		}
		resp := LoginResponse{
			Username:     *username,
			RefreshToken: *refreshToken,
			AccessToken:  *accessToken,
		}

		fmt.Print("\nToken resp:")
		fmt.Println(resp)
		return resp, err
	}
}

func getUserInfo(state string, code string) (*string, *string, *string, error) {
	if state != oauthStateString {
		return nil, nil, nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	contents, err := auth.CheckTokenOauth(token.AccessToken)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Invalid token oauth2: %s", err.Error())
	}

	// get email in contents
	var result map[string]interface{}
	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse account info: %s", err.Error())
	}
	email, ok := result["email"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("email is not a string")
	}
	parts := strings.Split(email, "@")
	if len(parts) < 2 {
		return nil, nil, nil, fmt.Errorf("invalid email format")
	}
	username := parts[0]
	fmt.Println("username: ", username)

	return &token.AccessToken, &email, &username, nil
}
