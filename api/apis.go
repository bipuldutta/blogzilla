package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/domain"
	"github.com/bipuldutta/blogzilla/usecases"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

/*
WebService is the main entry to the APIs. In real worl application we will collect several matrix such as
each endpoints request/response latency, counter etc.
*/
type WebService struct {
	conf           *config.Config
	authMiddleware *AuthMiddleware
	userManager    *usecases.UserManager
}

func NewWebService(conf *config.Config, userManager *usecases.UserManager) *WebService {
	logger = utils.Logger()
	return &WebService{
		conf:           conf,
		authMiddleware: NewAuthMiddleware(conf),
		userManager:    userManager,
	}
}

func (ws *WebService) Start() error {
	// Initialize HTTP router
	r := mux.NewRouter()
	// Define routes
	// Register a new user
	r.HandleFunc("/v1/register", ws.registerHandler).Methods("POST")
	// User login
	r.HandleFunc("/v1/login", ws.loginHandler).Methods("POST")
	// Get a user details
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.ReadUserPermission, http.HandlerFunc(ws.getUserHandler))).Methods("GET")
	// Update a user details
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.UpdateUserPermission, http.HandlerFunc(ws.updateUserHandler))).Methods("PUT")
	// Delete a user
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.DeleteUserPermission, http.HandlerFunc(ws.deleteUserHandler))).Methods("DELETE")

	// Create a blog, creator id will be extracted from the jwt token
	r.Handle("/v1/blogs", ws.authMiddleware.authorize(utils.CreateBlogPermission, http.HandlerFunc(ws.createBlogHandler))).Methods("POST")
	// Search all blogs, in real world application there will be filter mechanism and pagination
	r.Handle("/v1/blogs", ws.authMiddleware.authorize(utils.ReadBlogPermission, http.HandlerFunc(ws.getBlogsHandler))).Methods("GET")
	// Get the details about a blog, mainly for reading purpose
	r.Handle("/v1/blogs/{id}", ws.authMiddleware.authorize(utils.ReadBlogPermission, http.HandlerFunc(ws.getBlogHandler))).Methods("GET")
	// Update a blog, creator id will extracted from the jwt token
	r.Handle("/v1/blogs/{id}", ws.authMiddleware.authorize(utils.UpdateBlogPermission, http.HandlerFunc(ws.updateBlogHandler))).Methods("PUT")
	// Delete a blog, creator id will extracted from the jwt token
	r.Handle("/v1/blogs/{id}", ws.authMiddleware.authorize(utils.DeleteBlogPermission, http.HandlerFunc(ws.deleteBlogHandler))).Methods("GET")

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
	ctx := utils.CreateContext()
	id, err := ws.userManager.Create(ctx, &newUser)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	logger.Infof("successfully created user with id: %d", id)
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
	// Read the request body
	var blog domain.Blog
	err := json.NewDecoder(r.Body).Decode(&blog)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	// continue saving data
	userID := ws.getUserID(r)
	// set the user on the way in
	blog.UserID = userID
	logger.Infof("Successfully authorized user: %d to create the blog: %+v", userID, blog)

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

func (ws *WebService) getUserID(r *http.Request) int64 {
	userID := r.Context().Value("userId")
	if userID != nil {
		return userID.(int64)
	}
	return 0
}
