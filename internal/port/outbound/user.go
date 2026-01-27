package outbound_port

import "eduvera/internal/model"

//go:generate mockgen -source=user.go -destination=./../../../tests/mocks/port/mock_user.go
type UserDatabasePort interface {
	Create(user *model.User) error
	Update(user *model.User) error
	FindByFilter(filter model.UserFilter) ([]model.User, error)
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	EmailExists(email string) (bool, error)
	UpdateLastLogin(id string) error
	Activate(id string) error
}
