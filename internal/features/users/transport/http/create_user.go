package users_transport_http

import "net/http"

type CreateUserRequest struct {
	FullName    string
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {

}
