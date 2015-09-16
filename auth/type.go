package auth

import (
	"github.com/keichi/scribble/model"
)

type AuthContext struct {
	IsLoggedIn bool
	User       *model.User
	Session    *model.Session
}

type Authorizer interface {
	Authorize(user *model.User, action model.AuthorizedAction) bool
}
