package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"

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
	dataset := dialect.Insert(tableUser).Rows(goqu.Record{
		"tenant_id":     user.TenantID,
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
