package controller

import "net/http"

type Controller interface {
	Ping(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)

	GetAllJob(w http.ResponseWriter, r *http.Request)
	GetJobById(w http.ResponseWriter, r *http.Request)
	InsertJob(w http.ResponseWriter, r *http.Request)
	GetApplicationsByJobId(w http.ResponseWriter, r *http.Request)
	GetApplicationById(w http.ResponseWriter, r *http.Request)
	InsertApplication(w http.ResponseWriter, r *http.Request)
	UpdateApplicationStatus(w http.ResponseWriter, r *http.Request)
}
