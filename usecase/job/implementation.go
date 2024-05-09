package job

import (
	"context"
	"errors"
	"time"

	"github.com/michaelwongycn/job-portal/domain/enum"
	"github.com/michaelwongycn/job-portal/domain/model"
	"github.com/michaelwongycn/job-portal/domain/request"
	"github.com/michaelwongycn/job-portal/repository/appDB"
)

type jobImpl struct {
	appDB                appDB.AppDBInterface
	refreshTokenDuration time.Duration
}

func NewJobImpl(appDB appDB.AppDBInterface, refreshTokenDuration time.Duration) JobUsecase {
	return &jobImpl{
		appDB:                appDB,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (u *jobImpl) GetAllJob(ctx context.Context) (*[]model.Job, error) {
	return u.appDB.GetAllJob(ctx)
}

func (u *jobImpl) GetJobById(ctx context.Context, req request.SearchJobByIdRequest) (*model.Job, error) {
	return u.appDB.GetJobById(ctx, req.JobId)
}

func (u *jobImpl) InsertJob(ctx context.Context, req request.InsertJobRequest) error {
	return u.appDB.InsertJob(ctx, req.EmployerId, req.Title, req.Description, req.Requirement)
}

func (u *jobImpl) GetApplicationsByJobId(ctx context.Context, req request.SearchApplicationByJobRequest) (*[]model.Application, error) {
	return u.appDB.GetApplicationsByJobId(ctx, req.JobId, req.EmployerId)
}

func (u *jobImpl) GetApplicationById(ctx context.Context, req request.SearchApplicationByIdRequest) (*model.Application, error) {
	if req.Role == 1 {
		return u.appDB.GetApplicationByIdAndTalentId(ctx, req.ApplicationId, req.UserId)
	}
	return u.appDB.GetApplicationByIdAndEmployeerId(ctx, req.ApplicationId, req.UserId)
}

func (u *jobImpl) InsertApplication(ctx context.Context, req request.InsertApplicationRequest) error {
	return u.appDB.InsertApplication(ctx, req.JobId, req.TalentId)
}

func (u *jobImpl) UpdateApplicationStatus(ctx context.Context, req request.UpdateApplicationStatusRequest) error {
	application, err := u.appDB.GetApplicationByIdAndEmployeerId(ctx, req.ApplicationId, req.EmployerId)
	if err != nil {
		return err
	}

	job, err := u.appDB.GetJobById(ctx, application.JobId)
	if err != nil {
		return err
	}

	if job.EmployerId != req.EmployerId {
		return errors.New("Unauthorized")
	}

	if req.Status != enum.InterviewStatus && req.Status != enum.AcceptedStatus && req.Status != enum.DeclinedStatus {
		return errors.New("Invalid Status")
	}
	return u.appDB.UpdateApplicationStatus(ctx, req.ApplicationId, req.Status)
}
