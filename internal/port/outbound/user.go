package outbound_port

import "prabogo/internal/model"

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
	LinkToTenant(userID string, tenantID string) error
	UpdatePassword(id string, hashedPassword string) error

	// Reset Token operations
	CreateResetToken(token *model.ResetToken) error
	GetResetToken(token string) (*model.ResetToken, error)
	MarkResetTokenUsed(id string) error
	DeleteExpiredResetTokens() error
}
