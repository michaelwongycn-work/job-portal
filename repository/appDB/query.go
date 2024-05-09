package appDB

const (
	getUserByEmailAndPasswordQuery = "SELECT id, role FROM users WHERE email = $1 AND password = $2"
	insertUserQuery                = "INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id"
	getUserTokenQuery              = "SELECT access_token, refresh_token, expiration_time FROM user_tokens WHERE user_id = $1"
	insertUserTokenQuery           = "INSERT INTO user_tokens (user_id, access_token, refresh_token, expiration_time) VALUES ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET access_token = EXCLUDED.access_token, refresh_token = EXCLUDED.refresh_token, expiration_time = EXCLUDED.expiration_time"
	deleteUserTokenQuery           = "DELETE FROM user_tokens WHERE user_id = $1"

	getAllJobQuery                       = "SELECT * FROM jobs"
	getJobByIdQuery                      = "SELECT * FROM jobs WHERE id = $1"
	insertJobQuery                       = "INSERT INTO jobs (employer_id, title, description, requirement) VALUES ($1, $2, $3, $4)"
	getApplicationsByJobIdQuery          = "SELECT a.* FROM applications a JOIN jobs j ON a.job_id = j.id WHERE j.id = $1 AND employer_id = $2"
	getApplicationByIdAndEmployerIdQuery = "SELECT a.* FROM applications a JOIN jobs j ON a.job_id = j.id WHERE a.id = $1 AND employer_id = $2"
	getApplicationByIdAndTalentIdQuery   = "SELECT * FROM applications WHERE id = $1 AND talent_id = $2"
	insertApplicationQuery               = "INSERT INTO applications (job_id, talent_id) VALUES ($1, $2)"
	updateApplicationStatusQuery         = "UPDATE applications SET application_status = $1 WHERE id = $2"
)
