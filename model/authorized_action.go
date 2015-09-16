package model

type AuthorizedAction string

const (
	ACTION_CREATE AuthorizedAction = "common.create"
	ACTION_READ   AuthorizedAction = "common.read"
	ACTION_UPDATE AuthorizedAction = "common.update"
	ACTION_DELETE AuthorizedAction = "common.delete"
)
