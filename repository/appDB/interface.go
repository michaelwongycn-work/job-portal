package appDB

import (
	"context"

	"github.com/michaelwongycn/job-portal/domain/model"
)

type AppDBInterface interface {
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error)
	InsertUser(ctx context.Context, email, password string, role int) (int, error)
	GetUserToken(ctx context.Context, userId int) (*model.UserToken, error)
	InsertUserToken(ctx context.Context, userId int, accessToken, refreshToken string, expirationTime int64) error
	DeleteUserToken(ctx context.Context, userId int) error

	GetAllJob(ctx context.Context) (*[]model.Job, error)
	GetJobById(ctx context.Context, jobId int) (*model.Job, error)
	InsertJob(ctx context.Context, employerId int, title, description, requirement string) error
	GetApplicationsByJobId(ctx context.Context, jobId, employerId int) (*[]model.Application, error)
	GetApplicationByIdAndEmployeerId(ctx context.Context, applicationId, employerId int) (*model.Application, error)
	GetApplicationByIdAndTalentId(ctx context.Context, applicationId, talentId int) (*model.Application, error)
	InsertApplication(ctx context.Context, jobId, talentId int) error
	UpdateApplicationStatus(ctx context.Context, applicationId, status int) error
}
