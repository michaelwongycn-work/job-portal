package appDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/michaelwongycn/job-portal/domain/model"
	"github.com/michaelwongycn/job-portal/lib/log"
)

const (
	noRowsFoundErrorMsg      = "no rows found for the query"
	errorScanningRowErrorMsg = "error when scanning row"
	errorQueryingSQLErrorMsg = "error when querying SQL"
)

type appDBImpl struct {
	db      *sql.DB
	timeout time.Duration
}

func NewAppDBImpl(timeout time.Duration, db *sql.DB) AppDBInterface {
	return &appDBImpl{
		db:      db,
		timeout: timeout * time.Second,
	}
}

func (d *appDBImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*model.User, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	var data model.User
	row := d.db.QueryRowContext(ctx, getUserByEmailAndPasswordQuery, email, password)

	err := row.Scan(&data.ID, &data.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &data, nil
}

func (d *appDBImpl) InsertUser(ctx context.Context, email, password string, role int) (int, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	var id int

	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, insertUserQuery, email, password, role).Scan(&id)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *appDBImpl) GetUserToken(ctx context.Context, userId int) (*model.UserToken, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	UserToken := model.UserToken{
		UserId: userId,
	}
	row := d.db.QueryRowContext(ctx, getUserTokenQuery, userId)

	err := row.Scan(&UserToken.AccessToken, &UserToken.RefreshToken, &UserToken.ExpirationTime)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &UserToken, nil
}

func (d *appDBImpl) InsertUserToken(ctx context.Context, userId int, accessToken, refreshToken string, expirationTime int64) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, insertUserTokenQuery, userId, accessToken, refreshToken, expirationTime)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *appDBImpl) DeleteUserToken(ctx context.Context, userId int) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, deleteUserTokenQuery, userId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *appDBImpl) GetAllJob(ctx context.Context) (*[]model.Job, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	rows, err := d.db.QueryContext(ctx, getAllJobQuery)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return nil, err
	}

	var data []model.Job
	for rows.Next() {
		var job model.Job
		err := rows.Scan(&job.ID, &job.EmployerId, &job.Title, &job.Description, &job.Requirement, &job.CreateDate)
		if err != nil {
			if err == sql.ErrNoRows {
				log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
				return nil, err
			} else {
				log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)

				return nil, err
			}
		}
		data = append(data, job)
	}
	return &data, nil
}

func (d *appDBImpl) GetJobById(ctx context.Context, jobId int) (*model.Job, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	row := d.db.QueryRowContext(ctx, getJobByIdQuery, jobId)

	var job model.Job
	err := row.Scan(&job.ID, &job.EmployerId, &job.Title, &job.Description, &job.Requirement, &job.CreateDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &job, nil
}

func (d *appDBImpl) InsertJob(ctx context.Context, employerId int, title, description, requirement string) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, insertJobQuery, employerId, title, description, requirement)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *appDBImpl) GetApplicationsByJobId(ctx context.Context, jobId, employerId int) (*[]model.Application, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	rows, err := d.db.QueryContext(ctx, getApplicationsByJobIdQuery, jobId, employerId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return nil, err
	}

	var data []model.Application
	for rows.Next() {
		var application model.Application
		err := rows.Scan(&application.ID, &application.JobId, &application.TalentId, &application.ApplicationStatus, &application.ApplyDate)
		if err != nil {
			if err == sql.ErrNoRows {
				log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
				return nil, err
			} else {
				log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)

				return nil, err
			}
		}
		data = append(data, application)
	}
	return &data, nil
}

func (d *appDBImpl) GetApplicationByIdAndEmployeerId(ctx context.Context, applicationId, employerId int) (*model.Application, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	var data model.Application
	row := d.db.QueryRowContext(ctx, getApplicationByIdAndEmployerIdQuery, applicationId, employerId)

	err := row.Scan(&data.ID, &data.JobId, &data.TalentId, &data.ApplicationStatus, &data.ApplyDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &data, nil
}

func (d *appDBImpl) GetApplicationByIdAndTalentId(ctx context.Context, applicationId, talentId int) (*model.Application, error) {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	var data model.Application
	row := d.db.QueryRowContext(ctx, getApplicationByIdAndTalentIdQuery, applicationId, talentId)

	err := row.Scan(&data.ID, &data.JobId, &data.TalentId, &data.ApplicationStatus, &data.ApplyDate)
	if err != nil {
		if err == sql.ErrNoRows {
			log.PrintLogErr(ctx, noRowsFoundErrorMsg, err)
			return nil, err
		} else {
			log.PrintLogErr(ctx, errorScanningRowErrorMsg, err)
			return nil, err
		}
	}
	return &data, nil
}

func (d *appDBImpl) InsertApplication(ctx context.Context, jobId, talentId int) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, insertApplicationQuery, jobId, talentId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *appDBImpl) UpdateApplicationStatus(ctx context.Context, applicationId, status int) error {
	ctx, cancelfunc := context.WithTimeout(ctx, d.timeout)
	defer cancelfunc()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, updateApplicationStatusQuery, status, applicationId)
	if err != nil {
		log.PrintLogErr(ctx, errorQueryingSQLErrorMsg, err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
