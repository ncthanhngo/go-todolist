package common

const (
	CurrentUser = "current_user"
)

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdminOrMode(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
