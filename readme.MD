# Job Portal App

The Job Portal is an API-based platform that facilitates the connection between job seekers (talents) and employers. It features two main user flows: one for talents to search and apply for jobs, and another for employers to post jobs, review applications, and manage the hiring process. The project is fully integrated with JWT (JSON Web Token) for user authentication and authorization.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)

## Installation

```
git clone https://github.com/michaelwongycn-work/job-portal.git
cd job-portal
mv application_config_example.json application_config.json
go mod download
```

## Usage

```
go run main.go
```

The server will start running on http://localhost:2000.

## Endpoint

The following endpoints are available:

GET /ping
Ping endpoint, used for Health Check

POST /register
Registers a new user and returns access and refresh tokens.

POST /login
Authenticates a user and returns access and refresh tokens.

POST /logout
Logs out a user by invalidating the access token.

POST /refresh
Refreshes the access token using a refresh token.

POST /job
Inserts a new job into the database.

GET /jobs
Retrieves all jobs.

GET /job/{jobId}
Retrieves a job by ID.

GET /job/{jobId}/applications
Retrieves applications for a specific job.

POST /job/{jobId}
Inserts a new application into the database.

GET /application/{applicationId}
Retrieves an application by ID.

PUT /application/{applicationId}
Updates the status of an application in the database.

All endpoints require authentication except for /login, /register, and /refresh-token.
