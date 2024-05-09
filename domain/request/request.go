package request

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogoutRequest struct {
	UserId      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type UserRefreshTokenRequest struct {
	UserId       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	Role         int    `json:"role"`
}

type SearchJobByIdRequest struct {
	JobId int `json:"job_id"`
}

type InsertJobRequest struct {
	EmployerId  int    `json:"employer_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Requirement string `json:"requirement"`
}

type SearchApplicationByJobRequest struct {
	JobId      int `json:"job_id"`
	EmployerId int `json:"employer_id"`
}

type SearchApplicationByIdRequest struct {
	ApplicationId int `json:"application_id"`
	UserId        int `json:"user_id"`
	Role          int `json:"role"`
}

type SearchApplicationhByTalentRequest struct {
	TalentId int `json:"talent_id"`
}

type InsertApplicationRequest struct {
	JobId    int `json:"job_id"`
	TalentId int `json:"talent_id"`
}

type UpdateApplicationStatusRequest struct {
	ApplicationId int `json:"application_id"`
	EmployerId    int `json:"employer_id"`
	Status        int `json:"status"`
}
