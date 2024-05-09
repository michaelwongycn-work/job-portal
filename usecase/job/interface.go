package job

import (
	"context"

	"github.com/michaelwongycn/job-portal/domain/model"
	"github.com/michaelwongycn/job-portal/domain/request"
)

type JobUsecase interface {
	GetAllJob(ctx context.Context) (*[]model.Job, error)
	GetJobById(ctx context.Context, req request.SearchJobByIdRequest) (*model.Job, error)
	InsertJob(ctx context.Context, req request.InsertJobRequest) error

	GetApplicationsByJobId(ctx context.Context, req request.SearchApplicationByJobRequest) (*[]model.Application, error)
	GetApplicationById(ctx context.Context, req request.SearchApplicationByIdRequest) (*model.Application, error)
	InsertApplication(ctx context.Context, req request.InsertApplicationRequest) error
	UpdateApplicationStatus(ctx context.Context, req request.UpdateApplicationStatusRequest) error
}
