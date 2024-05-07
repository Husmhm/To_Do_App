package contract

import "go.dev/entity"

type UserWriteStore interface {
	Save(u entity.User)
}
type UserReadStore interface {
	Load() []entity.User
}
