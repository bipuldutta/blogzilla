package api

import (
	"blogs/config"
	"blogs/domain"
	"blogs/usecases"
	"blogs/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

type WebService struct {
	conf        *config.Config
	userManager *usecases.UserManager
}

func NewWebService(conf *config.Config, userManager *usecases.UserManager) *WebService {
	logger = utils.Logger()
	return &WebService{
		conf:        conf,
		userManager: userManager,
	}
}

func (ws *WebService) Start() error {
	// Initialize HTTP router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/register", ws.registerHandler).Methods("POST")
	r.HandleFunc("/login", ws.loginHandler).Methods("POST")
	r.HandleFunc("/users/{id}", ws.getUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", ws.updateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", ws.deleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users/{id}/blogs", ws.createBlogHandler).Methods("POST")
	r.HandleFunc("/users/{id}/blogs", ws.getBlogsHandler).Methods("GET")
	r.HandleFunc("/users/{id}/blogs/{blogID}", ws.getBlogHandler).Methods("GET")
	r.HandleFunc("/users/{id}/blogs/{blogID}", ws.updateBlogHandler).Methods("PUT")
	r.HandleFunc("/users/{id}/blogs/{blogID}", ws.deleteBlogHandler).Methods("DELETE")

	// Start the server
	logger.Printf("Server listening on port %d", ws.conf.Server.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ws.conf.Server.Port), r))

	return nil
}

func (ws *WebService) registerHandler(w http.ResponseWriter, r *http.Request) {
	var newUser domain.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logger.Infof("user to be created: %+v", newUser)
	ctx := utils.CreateContext()
	id, err := ws.userManager.Create(ctx, &newUser)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	logger.Infof("successfully create user with id: %d", id)
	w.WriteHeader(http.StatusCreated)
}

func (ws *WebService) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the username and password
	var credentials domain.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	logger.Infof("user to be authenticated: %+v", credentials)

	ctx := utils.CreateContext()
	token, err := ws.userManager.Login(ctx, credentials.Username, credentials.Password)
	if err != nil {
		// this could also be internal server error (DB outage, etc.),
		// but it will take extra time to have a proper error handling
		logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	// Return the session token to the client
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (ws *WebService) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) updateUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) deleteUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) createBlogHandler(w http.ResponseWriter, r *http.Request) {
	code, err := ws.authorize(r, utils.CreateBlogPermission)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	// Read the request body
	var blog domain.Blog
	err = json.NewDecoder(r.Body).Decode(&blog)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// continue saving data
	logger.Info("Successfully authorized to create the blog: %+v", blog)

	// Insert the blog post into the database
	/*
		err = createBlogPost(userID, blog.Title, blog.Content)
		if err != nil {
			http.Error(w, "Failed to create blog post", http.StatusInternalServerError)
			return
		}
	*/

	// Return a success response
	w.WriteHeader(http.StatusCreated)
}

func (ws *WebService) getBlogsHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) getBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) updateBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) deleteBlogHandler(w http.ResponseWriter, r *http.Request) {

}

func (ws *WebService) authorize(r *http.Request, permission string) (int, error) {
	// Read the user ID from the URL path
	userID, code, err := ws.getUserID(r)
	if err != nil {
		return code, err
	}

	token, err := ws.extractTokenFromHeader(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	claims, err := ws.validateToken(userID, token)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	if !claims.HasPermission(permission) {
		return http.StatusUnauthorized, fmt.Errorf("user does not have permission")
	}
	return http.StatusOK, nil
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header in the format "Bearer {token}".
func (ws *WebService) extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return bearerToken[1], nil
}

func (ws *WebService) validateToken(userID int64, tokenString string) (*domain.CustomClaims, error) {
	// Parse the token without verifying the signature.
	token, err := jwt.ParseWithClaims(tokenString, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ws.conf.Login.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Verify the token signature and expiration.
	if !token.Valid {
		return nil, fmt.Errorf("invalid token signature")
	}

	claims, ok := token.Claims.(*domain.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.UserID != userID {
		return nil, fmt.Errorf("invalid token used")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	logger.Infof("##### permissions from JWT: %+v", claims.Permissions)

	return claims, nil
}

func (ws *WebService) getUserID(r *http.Request) (int64, int, error) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return -1, http.StatusBadRequest, fmt.Errorf("invalid user ID")
	}
	return int64(userID), http.StatusOK, nil
}
