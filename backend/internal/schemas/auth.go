package schemas

type AuthResponse struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message,omitempty"`
}
