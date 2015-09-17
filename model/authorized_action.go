package model

// AuthorizedAction is a type representing the type of actions requested by the user
type AuthorizedAction string

const (
	// ActionCreate represents an action to create entities
	ActionCreate AuthorizedAction = "common.create"
	// ActionRead represents an action to read entities
	ActionRead AuthorizedAction = "common.read"
	// ActionUpdate represents an action to update entities
	ActionUpdate AuthorizedAction = "common.update"
	// ActionDelete represents an action to delete entities
	ActionDelete AuthorizedAction = "common.delete"
)
