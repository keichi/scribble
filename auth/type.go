package auth

import (
	"github.com/keichi/scribble/model"
)

// Context is a struct that holds an authentication context
type Context struct {
	IsLoggedIn bool
	User       *model.User
	Session    *model.Session
}

// Authorizer is an interface for models that authorizes requested action
type Authorizer interface {
	Authorize(user *model.User, action model.AuthorizedAction) bool
}
