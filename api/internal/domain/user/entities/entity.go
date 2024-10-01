package entities

type Entity interface {
	GetOwnerID() string
	IsPublic() bool
}
