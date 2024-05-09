package response

type ReadResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Time    string      `json:"time"`
}

type WriteResponse struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
