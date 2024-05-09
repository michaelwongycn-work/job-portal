package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/michaelwongycn/job-portal/domain/request"
	"github.com/michaelwongycn/job-portal/domain/response"
	"github.com/michaelwongycn/job-portal/lib/auth"
	"github.com/michaelwongycn/job-portal/usecase/job"
	"github.com/michaelwongycn/job-portal/usecase/user"
)

const (
	invalidCredentialsErrorMsg = "Invalid Credentials"
	passwordNotMatchErrorMsg   = "Password doesn't match"
	unableToParseTokenErrorMsg = "Unable to parse token"
	notFoundErrorMsg           = "Not Found"
	internalServerErrorMsg     = "Internal Server Error"
)

type controllerImpl struct {
	userUsecase user.UserUsecase
	jobUsecase  job.JobUsecase
}

func NewControllerImpl(userUsecase user.UserUsecase, jobUsecase job.JobUsecase) Controller {
	return &controllerImpl{
		userUsecase: userUsecase,
		jobUsecase:  jobUsecase,
	}
}

func setResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func parseHeader(authorizationHeader string) (accessToken string, claims jwt.MapClaims, err error) {
	accessToken = strings.Split(authorizationHeader, " ")[1]
	claims, err = auth.ParseToken(accessToken)
	if err != nil {
		return "", nil, err
	}
	return accessToken, claims, nil
}

// @Summary Ping endpoint
// @Description Used for Health Check"
// @Tags Ping
// @Accept json
// @Produce json
// @Success 200 {string} string "Pong!"
// @Router /ping [get]
func (c *controllerImpl) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

// @Summary User Login
// @Description Authenticates a user and returns access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserLoginRequest true "User Login Request"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} response.ReadResponse
// @Router /login [post]
func (c *controllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.UserLoginRequest{}
	authResponse := response.AuthResponse{}
	response := response.ReadResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	accessToken, refreshToken, err := c.userUsecase.Login(ctx, req)
	if err != nil {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	authResponse.AccessToken = *accessToken
	authResponse.RefreshToken = *refreshToken

	response.Message = ""
	response.Data = authResponse
	setResponse(w, http.StatusOK, response)
}

// @Summary User Registration
// @Description Registers a new user and returns access and refresh tokens
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.UserRegisterRequest true "User Registration Request"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} response.ReadResponse
// @Failure 409 {object} response.ReadResponse
// @Failure 500 {object} response.ReadResponse
// @Router /register [post]
func (c *controllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.UserRegisterRequest{}
	authResponse := response.AuthResponse{}
	response := response.ReadResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	accessToken, refreshToken, err := c.userUsecase.Register(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			setResponse(w, http.StatusConflict, response)
			return
		}
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	authResponse.AccessToken = *accessToken
	authResponse.RefreshToken = *refreshToken

	response.Message = ""
	response.Data = authResponse
	setResponse(w, http.StatusOK, response)
}

// @Summary User Logout
// @Description Logs out a user by invalidating the access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Access Token"
// @Success 200 {object} response.WriteResponse
// @Failure 500 {object} response.WriteResponse
// @Router /logout [post]
func (c *controllerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.UserLogoutRequest{}
	response := response.WriteResponse{}
	response.Time = requestTime

	accessToken, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.UserId = int(claims["sub"].(float64))
	req.AccessToken = accessToken
	err = c.userUsecase.Logout(ctx, req)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

// @Summary Refresh Access Token
// @Description Refreshes the access token using a refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Refresh Token"
// @Param request body request.UserRefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} response.AuthResponse
// @Failure 400 {object} response.ReadResponse
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 500 {object} response.ReadResponse
// @Router /refresh [post]
func (c *controllerImpl) RefreshToken(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.UserRefreshTokenRequest{}
	authResponse := response.AuthResponse{}
	response := response.ReadResponse{}
	response.Time = requestTime

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	accessTokenUserId := int(claims["sub"].(float64))
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}
	req.Role = int(claims["rle"].(float64))
	req.UserId = accessTokenUserId

	claims, err = auth.ParseToken(req.RefreshToken)
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	refreshTokenUserId := int(claims["sub"].(float64))
	if refreshTokenUserId != accessTokenUserId {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	newAccessToken, newRefreshToken, err := c.userUsecase.RefreshToken(ctx, req)
	if err != nil {
		response.Message = invalidCredentialsErrorMsg
		setResponse(w, http.StatusOK, response)
		return
	}

	authResponse.AccessToken = *newAccessToken
	authResponse.RefreshToken = *newRefreshToken

	response.Message = ""
	response.Data = authResponse
	setResponse(w, http.StatusOK, response)
}

// @Summary Get all jobs
// @Description Retrieves all jobs
// @Tags Job
// @Produce json
// @Success 200 {array} model.Job
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 404 {object} response.ReadResponse "Not Found"
// @Failure 500 {object} response.ReadResponse
// @Router /jobs [get]
func (c *controllerImpl) GetAllJob(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	response := response.ReadResponse{}
	response.Time = requestTime

	jobs, err := c.jobUsecase.GetAllJob(ctx)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	if len(*jobs) == 0 {
		setResponse(w, http.StatusNotFound, response)
		return
	}

	response.Message = ""
	response.Data = jobs
	setResponse(w, http.StatusOK, response)
}

// @Summary Get a job by ID
// @Description Retrieves a job by ID
// @Tags Job
// @Produce json
// @Param jobId path int true "Job ID"
// @Success 200 {object} response.ReadResponse
// @Failure 400 {object} response.ReadResponse "Bad Request"
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 404 {object} response.ReadResponse "Not Found"
// @Failure 500 {object} response.ReadResponse "Internal Server Error"
// @Router /job/{jobId} [get]
func (c *controllerImpl) GetJobById(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.SearchJobByIdRequest{}
	response := response.ReadResponse{}
	response.Time = requestTime

	jobIdStr := chi.URLParam(r, "jobId")

	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	req.JobId = jobId
	jobs, err := c.jobUsecase.GetJobById(ctx, req)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Message = notFoundErrorMsg
			setResponse(w, http.StatusNotFound, response)
			return
		}
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	response.Data = jobs
	setResponse(w, http.StatusOK, response)
}

// @Summary Insert a new job
// @Description Inserts a new job into the database
// @Tags Job
// @Produce json
// @Param job body request.InsertJobRequest true "Job data"
// @Success 200 {object} response.WriteResponse
// @Failure 400 {object} response.WriteResponse "Bad Request"
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 500 {object} response.WriteResponse "Internal Server Error"
// @Router /job [post]
func (c *controllerImpl) InsertJob(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.InsertJobRequest{}
	response := response.WriteResponse{}
	response.Time = requestTime

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.EmployerId = int(claims["sub"].(float64))
	err = c.jobUsecase.InsertJob(ctx, req)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}
	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

// @Summary Get applications by Job ID
// @Description Retrieves applications for a specific job
// @Tags Application
// @Produce json
// @Param jobId path int true "Job ID"
// @Success 200 {object} response.ReadResponse
// @Failure 400 {object} response.ReadResponse "Bad Request"
// @Failure 401 {object} response.ReadResponse "Unauthorized"
// @Failure 500 {object} response.ReadResponse "Internal Server Error"
// @Router /job/{jobId}/applications [get]
func (c *controllerImpl) GetApplicationsByJobId(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.SearchApplicationByJobRequest{}
	response := response.ReadResponse{}
	response.Time = requestTime

	jobIdStr := chi.URLParam(r, "jobId")
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		response.Message = ""
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.JobId = jobId
	req.EmployerId = int(claims["sub"].(float64))
	jobs, err := c.jobUsecase.GetApplicationsByJobId(ctx, req)
	if err != nil {
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	response.Data = jobs
	setResponse(w, http.StatusOK, response)
}

// @Summary Get an application by ID
// @Description Retrieves an application by ID
// @Tags Application
// @Produce json
// @Param applicationId path int true "Application ID"
// @Success 200 {object} response.ReadResponse
// @Failure 400 {object} response.ReadResponse "Bad Request"
// @Failure 401 {object} response.ReadResponse "Unauthorized"
// @Failure 404 {object} response.ReadResponse "Not Found"
// @Failure 500 {object} response.ReadResponse "Internal Server Error"
// @Router /application/{applicationId} [get]
func (c *controllerImpl) GetApplicationById(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.SearchApplicationByIdRequest{}
	response := response.ReadResponse{}
	response.Time = requestTime

	applicationIdStr := chi.URLParam(r, "applicationId")
	applicationId, err := strconv.Atoi(applicationIdStr)
	if err != nil {
		response.Message = ""
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.ApplicationId = applicationId
	req.UserId = int(claims["sub"].(float64))
	req.Role = int(claims["rle"].(float64))
	jobs, err := c.jobUsecase.GetApplicationById(ctx, req)
	if err != nil {
		if err == sql.ErrNoRows {
			response.Message = notFoundErrorMsg
			setResponse(w, http.StatusNotFound, response)
			return
		}
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	response.Data = jobs
	setResponse(w, http.StatusOK, response)
}

// @Summary Apply for the Job
// @Description Inserts a new application into the database
// @Tags Job
// @Produce json
// @Param jobId path int true "Job ID"
// @Success 200 {object} response.WriteResponse
// @Failure 400 {object} response.WriteResponse "Bad Request"
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 404 {object} response.WriteResponse "Not Found"
// @Failure 409 {object} response.WriteResponse "Conflict"
// @Failure 500 {object} response.WriteResponse "Internal Server Error"
// @Router /job/{jobId} [post]
func (c *controllerImpl) InsertApplication(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.InsertApplicationRequest{}
	response := response.WriteResponse{}
	response.Time = requestTime

	jobIdStr := chi.URLParam(r, "jobId")
	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		response.Message = ""
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.JobId = jobId
	req.TalentId = int(claims["sub"].(float64))
	err = c.jobUsecase.InsertApplication(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			setResponse(w, http.StatusConflict, response)
			return
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			setResponse(w, http.StatusNotFound, response)
			return
		}
		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}

// @Summary Update the status of an application
// @Description Updates the status of an application in the database
// @Tags Application
// @Produce json
// @Param applicationId path int true "Application ID"
// @Param request body request.UpdateApplicationStatusRequest true "Update Application Status Request"
// @Success 200 {object} response.WriteResponse
// @Failure 400 {object} response.WriteResponse "Bad Request"
// @Failure 401 {object} response.WriteResponse "Unauthorized"
// @Failure 404 {object} response.WriteResponse "Not Found"
// @Failure 500 {object} response.WriteResponse "Internal Server Error"
// @Router /application/{applicationId} [put]
func (c *controllerImpl) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	requestTime := time.Now().Format(time.RFC3339)
	ctx := r.Context()
	req := request.UpdateApplicationStatusRequest{}
	response := response.WriteResponse{}
	response.Time = requestTime

	applicationIdStr := chi.URLParam(r, "applicationId")
	applicationId, err := strconv.Atoi(applicationIdStr)
	if err != nil {
		response.Message = ""
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = err.Error()
		setResponse(w, http.StatusBadRequest, response)
		return
	}

	_, claims, err := parseHeader(r.Header.Get("Authorization"))
	if err != nil {
		response.Message = unableToParseTokenErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	req.ApplicationId = applicationId
	req.EmployerId = int(claims["sub"].(float64))
	err = c.jobUsecase.UpdateApplicationStatus(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") {
			setResponse(w, http.StatusNotFound, response)
			return
		}
		if strings.Contains(err.Error(), "Invalid Status") {
			setResponse(w, http.StatusBadRequest, response)
			return
		}
		if err == sql.ErrNoRows {
			response.Message = notFoundErrorMsg
			setResponse(w, http.StatusNotFound, response)
			return
		}

		response.Message = internalServerErrorMsg
		setResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = ""
	setResponse(w, http.StatusOK, response)
}
