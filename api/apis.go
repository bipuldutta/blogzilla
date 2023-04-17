package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bipuldutta/blogzilla/config"
	"github.com/bipuldutta/blogzilla/usecases"
	"github.com/bipuldutta/blogzilla/utils"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	logger *logrus.Logger

	// we create a new custom metric of type counter
	createBlogCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_create_blog_count", // metric name
			Help: "Count of create blog.",
		},
		[]string{}, // labels
	)
)

func init() {
	// we need to register the counter so prometheus can collect this metric
	prometheus.MustRegister(createBlogCount)
}

/*
WebService is the main entry to the APIs exposed by Blogzilla. In real world application we will collect several matrix such as
each endpoints request/response latency, counter etc.
*/
type WebService struct {
	conf           *config.Config
	authMiddleware *AuthMiddleware
	userManager    *usecases.UserManager
	blogManager    *usecases.BlogManager
}

func NewWebService(conf *config.Config, authManager *usecases.AuthManager, userManager *usecases.UserManager, blogManager *usecases.BlogManager) *WebService {
	logger = utils.Logger()
	return &WebService{
		conf:           conf,
		authMiddleware: NewAuthMiddleware(conf, authManager),
		userManager:    userManager,
		blogManager:    blogManager,
	}
}

func (ws *WebService) Start() error {
	// Initialize HTTP router
	r := mux.NewRouter()

	// Define routes

	// Register a new user
	r.Handle("/v1/register", http.HandlerFunc(ws.registerHandler)).Methods("POST")
	// User login
	r.Handle("/v1/login", http.HandlerFunc(ws.loginHandler)).Methods("POST")
	// Get a user details
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.ReadUserPermission, http.HandlerFunc(ws.getUserHandler))).Methods("GET")
	// Update a user details
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.UpdateUserPermission, http.HandlerFunc(ws.updateUserHandler))).Methods("PUT")
	// Delete a user
	r.Handle("/v1/users/{id}", ws.authMiddleware.authorize(utils.DeleteUserPermission, http.HandlerFunc(ws.deleteUserHandler))).Methods("DELETE")

	// Create a blog, creator id will be extracted from the jwt token
	r.Handle("/v1/blogs", ws.authMiddleware.authorize(utils.CreateBlogPermission, http.HandlerFunc(ws.createBlogHandler))).Methods("POST")
	// Update a blog, creator id will extracted from the jwt token
	r.Handle("/v1/blogs", ws.authMiddleware.authorize(utils.UpdateBlogPermission, http.HandlerFunc(ws.updateBlogHandler))).Methods("PUT")
	// Search all blogs, in real world application there will be filter mechanism and pagination
	r.Handle("/v1/blogs", ws.authMiddleware.authorize(utils.ReadBlogPermission, http.HandlerFunc(ws.searchBlogsHandler))).Methods("GET")
	// Get the details about a blog, mainly for reading purpose
	r.Handle("/v1/blogs/{id}", ws.authMiddleware.authorize(utils.ReadBlogPermission, http.HandlerFunc(ws.getBlogHandler))).Methods("GET")
	// Delete a blog, creator id will extracted from the jwt token
	r.Handle("/v1/blogs/{id}", ws.authMiddleware.authorize(utils.DeleteBlogPermission, http.HandlerFunc(ws.deleteBlogHandler))).Methods("DELETE")

	// Start the server
	logger.Printf("Server listening on port %d", ws.conf.Server.Port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ws.conf.Server.Port), r))

	return nil
}

func (ws *WebService) registerHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequestV1

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := utils.CreateContext()
	newUser := convertCreateUserRequestToDomain(&request)

	createdUser, err := ws.userManager.Create(ctx, newUser)
	if err != nil {
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}
	logger.Infof("successfully created user with id: %d", createdUser.ID)
	response := convertUserDomainObjToAPI(createdUser)
	ws.setResponse(w, http.StatusOK, response)
}

func (ws *WebService) setResponse(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (ws *WebService) loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the username and password
	var request LoginRequestV1
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := utils.CreateContext()
	token, err := ws.userManager.Login(ctx, request.Username, request.Password)
	if err != nil {
		// this could also be internal server error (DB outage, etc.),
		// but it will take extra time to have a proper error handling
		logger.WithError(err).Error("failed to authenticate user")
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	response := &LoginResponseV1{
		Token: token,
	}
	ws.setResponse(w, http.StatusOK, response)
}

func (ws *WebService) getUserHandler(w http.ResponseWriter, r *http.Request) {
}

func (ws *WebService) updateUserHandler(w http.ResponseWriter, r *http.Request) {
}

func (ws *WebService) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
}

func (ws *WebService) createBlogHandler(w http.ResponseWriter, r *http.Request) {

	// increment the counter
	createBlogCount.With(nil).Inc()

	// Read the request body
	var blogRequest CreateBlogRequestV1
	err := json.NewDecoder(r.Body).Decode(&blogRequest)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	// continue saving data
	userID := ws.getUserID(r)
	newBlog := convertCreateBlogRequestToDomain(userID, &blogRequest)

	// Insert the blog post into the database
	ctx := utils.CreateContext()
	blogID, err := ws.blogManager.Create(ctx, newBlog)
	if err != nil {
		http.Error(w, "failed to create blog post", http.StatusInternalServerError)
		return
	}

	// Return the blog id to the client
	response := &CreateBlogResponseV1{
		ID: blogID,
	}
	// Return a success response
	ws.setResponse(w, http.StatusCreated, response)
}

func (ws *WebService) searchBlogsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0 // Default offset value
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10 // Default limit value
	}
	ctx := utils.CreateContext()
	blogs, err := ws.blogManager.Search(ctx, offset, limit, query)
	if err != nil {
		http.Error(w, "failed to search blogs", http.StatusInternalServerError)
		return
	}
	ws.setResponse(w, http.StatusCreated, blogs)
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
