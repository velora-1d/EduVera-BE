package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

const tableUser = "users"

type userAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewUserAdapter(
	db outbound_port.DatabaseExecutor,
) outbound_port.UserDatabasePort {
	return &userAdapter{
		db: db,
	}
}

func (a *userAdapter) Create(user *model.User) error {
	dialect := goqu.Dialect("postgres")

	// Handle empty TenantID as NULL (avoiding invalid UUID error)
	var tenantID interface{}
	if user.TenantID == "" {
		tenantID = nil
	} else {
		tenantID = user.TenantID
	}

	dataset := dialect.Insert(tableUser).Rows(goqu.Record{
		"tenant_id":     tenantID,
		"name":          user.Name,
		"email":         user.Email,
		"whatsapp":      user.WhatsApp,
		"password_hash": user.PasswordHash,
		"role":          user.Role,
		"is_active":     user.IsActive,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&user.ID)
}

func (a *userAdapter) Update(user *model.User) error {
	user.UpdatedAt = time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableUser).
		Set(goqu.Record{
			"name":       user.Name,
			"email":      user.Email,
			"whatsapp":   user.WhatsApp,
			"role":       user.Role,
			"is_active":  user.IsActive,
			"updated_at": user.UpdatedAt,
		}).
		Where(goqu.Ex{"id": user.ID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) FindByFilter(filter model.UserFilter) ([]model.User, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableUser)
	dataset = addUserFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(
			&u.ID, &u.TenantID, &u.Name, &u.Email, &u.WhatsApp,
			&u.PasswordHash, &u.Role, &u.IsActive,
			&u.EmailVerifiedAt, &u.LastLoginAt,
			&u.CreatedAt, &u.UpdatedAt,
			&u.IsOwner,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (a *userAdapter) FindByID(id string) (*model.User, error) {
	users, err := a.FindByFilter(model.UserFilter{IDs: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return &users[0], nil
}

func (a *userAdapter) FindByEmail(email string) (*model.User, error) {
	users, err := a.FindByFilter(model.UserFilter{Emails: []string{email}})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return &users[0], nil
}

func (a *userAdapter) EmailExists(email string) (bool, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableUser).
		Select(goqu.L("1")).
		Where(goqu.Ex{"email": email}).
		Limit(1)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return false, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (a *userAdapter) UpdateLastLogin(id string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableUser).
		Set(goqu.Record{
			"last_login_at": now,
			"updated_at":    now,
		}).
		Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) Activate(id string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableUser).
		Set(goqu.Record{
			"is_active":         true,
			"email_verified_at": now,
			"updated_at":        now,
		}).
		Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *userAdapter) LinkToTenant(userID string, tenantID string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableUser).
		Set(goqu.Record{
			"tenant_id":  tenantID,
			"updated_at": now,
		}).
		Where(goqu.Ex{"id": userID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func addUserFilter(dataset *goqu.SelectDataset, filter model.UserFilter) *goqu.SelectDataset {
	if len(filter.IDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}
	if len(filter.TenantIDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"tenant_id": filter.TenantIDs})
	}
	if len(filter.Emails) > 0 {
		dataset = dataset.Where(goqu.Ex{"email": filter.Emails})
	}
	if len(filter.Roles) > 0 {
		dataset = dataset.Where(goqu.Ex{"role": filter.Roles})
	}
	return dataset
}

// UpdatePassword updates user's password hash
func (a *userAdapter) UpdatePassword(id string, hashedPassword string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableUser).
		Set(goqu.Record{
			"password_hash": hashedPassword,
			"updated_at":    now,
		}).
		Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

const tableResetToken = "reset_tokens"

// CreateResetToken creates a new reset token
func (a *userAdapter) CreateResetToken(token *model.ResetToken) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableResetToken).Rows(goqu.Record{
		"user_id":    token.UserID,
		"token":      token.Token,
		"expires_at": token.ExpiresAt,
		"created_at": token.CreatedAt,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&token.ID)
}

// GetResetToken retrieves a reset token by token string
func (a *userAdapter) GetResetToken(tokenStr string) (*model.ResetToken, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableResetToken).
		Where(goqu.Ex{"token": tokenStr})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	var token model.ResetToken
	err = a.db.QueryRow(query).Scan(
		&token.ID, &token.UserID, &token.Token,
		&token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// MarkResetTokenUsed marks a reset token as used
func (a *userAdapter) MarkResetTokenUsed(id string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableResetToken).
		Set(goqu.Record{
			"used_at": now,
		}).
		Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

// DeleteExpiredResetTokens deletes all expired and used reset tokens
func (a *userAdapter) DeleteExpiredResetTokens() error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Delete(tableResetToken).
		Where(goqu.Or(
			goqu.C("expires_at").Lt(time.Now()),
			goqu.C("used_at").IsNotNull(),
		))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}
