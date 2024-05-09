package model

import "time"

type Application struct {
	ID                int       `json:"id"`
	JobId             int       `json:"job_id"`
	TalentId          int       `json:"talent_id"`
	ApplicationStatus int       `json:"application_status"`
	ApplyDate         time.Time `json:"apply_date"`
}
