package transport

import (
	"DoAn/database"
	"DoAn/entity"
	"DoAn/pkg/user"
	"DoAn/pkg/user/auth"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"strings"
	"time"

	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	fmt.Println("GOOGLE_CLIENT_ID: ", os.Getenv("GOOGLE_CLIENT_ID"))
	fmt.Println("GOOGLE_CLIENT_SECRET: ", os.Getenv("GOOGLE_CLIENT_SECRET"))

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8000/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

type (
	RegisterUserRequest struct {
		user entity.User
	}
	ChangePassFirstTimeRequest struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password"`
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
		Token string `json:"token"`
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

func MakeRegisterUserEndpoints(u user.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterUserRequest)
		msg, err := u.Register(ctx, req.user)
		return RegisterUserResponse{
			Msg: msg,
			Err: err,
		}, err
	}
}
func MakeChangePassFirstTimeEndpoints(u user.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ChangePassFirstTimeRequest)
		msg, err := u.ChangePassFirstTime(ctx, req.Username, req.Password)
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
func MakeLoginEndpoints(u user.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginRequest)
		token, err := u.Login(ctx, entity.User{
			Username: req.Username,
			Password: req.Password,
		})
		return LoginResponse{Token: token}, err
	}
}
func MakeGetProfileEndpoints(u user.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//get request param in a string
		email := request.(string)
		data, err := u.GetProfile(ctx, email)
		return UserProfileResponse{
			Data: data,
			Err:  err,
		}, err
	}
}
func MakeLogoutEndpoints(u user.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//get request param in a string
		token := request.(string)
		data, err := u.Logout(ctx, token)
		return UserProfileResponse{
			Data: data,
			Err:  err,
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
func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func DecodeGetProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	userName, err := auth.GetSubjectFromToken(tokenString)
	if err != nil {
		return nil, err
	}
	tokenString = tokenString[len("Bearer "):]
	err = auth.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		return nil, err
	}
	var result entity.User
	err = collectionPostgres.QueryRow("SELECT id,username,email,roles FROM userinfo WHERE username=$1", userName).Scan(&result.ID, &result.Username, &result.Email, &result.Roles)
	if err != nil {
		return nil, err
	}
	emailRP := r.URL.Query().Get("email")
	if result.Email != emailRP {
		return nil, errors.New("email is invalid")
	}
	return emailRP, nil
}
func DecodeChangePassFirstTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	userName, err := auth.GetSubjectFromToken(tokenString)
	if err != nil {
		return nil, err
	}
	tokenString = tokenString[len("Bearer "):]
	err = auth.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		return nil, err
	}
	var result entity.User
	err = collectionPostgres.QueryRow("SELECT id,username,email,roles FROM userinfo WHERE username=$1", userName).Scan(&result.ID, &result.Username, &result.Email, &result.Roles)
	if err != nil {
		return nil, err
	}
	var request ChangePassFirstTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	if govalidator.IsNull(request.Password) {
		return nil, errors.New("password is required")
	}
	request.Username = userName
	return request, nil
}
func DecodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	userName, err := auth.GetSubjectFromToken(tokenString)
	if err != nil {
		return nil, err
	}
	tokenString = tokenString[len("Bearer "):]
	err = auth.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		return nil, err
	}
	var result entity.User
	err = collectionPostgres.QueryRow("SELECT id,username,email,roles FROM userinfo WHERE username=$1", userName).Scan(&result.ID, &result.Username, &result.Email, &result.Roles)
	if err != nil {
		return nil, err
	}

	return tokenString, nil
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

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, email, username, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if email == nil || username == nil {
		fmt.Errorf("email or username is nil")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//result := map[string]interface{}{"token": token, "email": email, "username": username}
	fmt.Println("email: ", *email)

	collectionPostgres, err := database.ConnectPostgres()
	if err != nil {
		return
	}
	var user entity.User
	errFindEmail := collectionPostgres.QueryRow("SELECT id,username,email,password,roles FROM userinfo WHERE email=$1", email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Roles)

	if errFindEmail != nil {
		fmt.Println("User does not exist, create new user")
		password := ""
		newUser := entity.User{Username: *username, Email: *email, Password: password, Roles: "ROLE_USER"}
		_, errs := collectionPostgres.Exec("INSERT INTO userinfo (username, email, password,roles) VALUES ($1, $2, $3,$4)", newUser.Username, newUser.Email, newUser.Password, newUser.Roles)
		if errs != nil {
			return
		}
	}

	fmt.Println("User already exists, create jwt token")
	accessToken, err := GenerateToken(*username)
	if err != nil {
		return
	}
	// Connect to MongoDB
	collectionMongo := database.ConnectMongo(os.Getenv("TokenCollectionMongo"))
	newToken := bson.M{"token": token, "user": username, "created_at": time.Now()}
	_, errs := collectionMongo.InsertOne(context.TODO(), newToken)
	if errs != nil {
		return
	}

	fmt.Print("\nToken:")
	fmt.Println(accessToken)

	fmt.Fprintf(w, "Token: "+accessToken)
}
func GenerateToken(username string) (string, error) {
	token, err := auth.Create(username)
	if err != nil {
		fmt.Println("errCreate")
		fmt.Println(err)
		return "", err
	}
	return token, nil
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

	//get email in contents
	var result map[string]interface{}
	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse user info: %s", err.Error())
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
