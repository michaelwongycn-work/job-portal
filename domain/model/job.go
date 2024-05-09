package model

import "time"

type Job struct {
	ID          int       `json:"id"`
	EmployerId  int       `json:"employer_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Requirement string    `json:"requirement"`
	CreateDate  time.Time `json:"create_date"`
}
