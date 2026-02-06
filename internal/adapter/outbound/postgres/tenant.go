package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

const tableTenant = "tenants"

type tenantAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewTenantAdapter(
	db outbound_port.DatabaseExecutor,
) outbound_port.TenantDatabasePort {
	return &tenantAdapter{
		db: db,
	}
}

func (a *tenantAdapter) Create(tenant *model.Tenant) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableTenant).Rows(goqu.Record{
		"name":              tenant.Name,
		"subdomain":         tenant.Subdomain,
		"plan_type":         tenant.PlanType,
		"subscription_tier": tenant.SubscriptionTier,
		"institution_type":  tenant.InstitutionType,
		"address":           tenant.Address,
		"bank_name":         tenant.BankName,
		"account_number":    tenant.AccountNumber,
		"account_holder":    tenant.AccountHolder,
		"status":            tenant.Status,
		"trial_ends_at":     tenant.TrialEndsAt,
		"created_at":        tenant.CreatedAt,
		"updated_at":        tenant.UpdatedAt,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&tenant.ID)
}

func (a *tenantAdapter) Update(tenant *model.Tenant) error {
	tenant.UpdatedAt = time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableTenant).
		Set(goqu.Record{
			"name":              tenant.Name,
			"subdomain":         tenant.Subdomain,
			"plan_type":         tenant.PlanType,
			"subscription_tier": tenant.SubscriptionTier,
			"institution_type":  tenant.InstitutionType,
			"address":           tenant.Address,
			"bank_name":         tenant.BankName,
			"account_number":    tenant.AccountNumber,
			"account_holder":    tenant.AccountHolder,
			"status":            tenant.Status,
			"trial_ends_at":     tenant.TrialEndsAt,
			"updated_at":        tenant.UpdatedAt,
		}).
		Where(goqu.Ex{"id": tenant.ID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *tenantAdapter) FindByFilter(filter model.TenantFilter) ([]model.Tenant, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableTenant).Select(
		"id", "name", "subdomain", "plan_type", "subscription_tier",
		"institution_type", "address", "bank_name",
		"account_number", "account_holder", "status",
		"trial_ends_at", "created_at", "updated_at",
	)
	dataset = addTenantFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []model.Tenant
	for rows.Next() {
		var t model.Tenant
		err := rows.Scan(
			&t.ID, &t.Name, &t.Subdomain, &t.PlanType, &t.SubscriptionTier,
			&t.InstitutionType, &t.Address, &t.BankName,
			&t.AccountNumber, &t.AccountHolder, &t.Status,
			&t.TrialEndsAt, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}

	return tenants, nil
}

func (a *tenantAdapter) FindByID(id string) (*model.Tenant, error) {
	tenants, err := a.FindByFilter(model.TenantFilter{IDs: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, sql.ErrNoRows
	}
	return &tenants[0], nil
}

func (a *tenantAdapter) FindBySubdomain(subdomain string) (*model.Tenant, error) {
	tenants, err := a.FindByFilter(model.TenantFilter{Subdomains: []string{subdomain}})
	if err != nil {
		return nil, err
	}
	if len(tenants) == 0 {
		return nil, sql.ErrNoRows
	}
	return &tenants[0], nil
}

func (a *tenantAdapter) SubdomainExists(subdomain string) (bool, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableTenant).
		Select(goqu.L("1")).
		Where(goqu.Ex{"subdomain": subdomain}).
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

func (a *tenantAdapter) UpdateStatus(id string, status string) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tableTenant).
		Set(goqu.Record{
			"status":     status,
			"updated_at": time.Now(),
		}).
		Where(goqu.Ex{"id": id})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func addTenantFilter(dataset *goqu.SelectDataset, filter model.TenantFilter) *goqu.SelectDataset {
	if len(filter.IDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}
	if len(filter.Subdomains) > 0 {
		dataset = dataset.Where(goqu.Ex{"subdomain": filter.Subdomains})
	}
	if len(filter.PlanTypes) > 0 {
		dataset = dataset.Where(goqu.Ex{"plan_type": filter.PlanTypes})
	}
	if len(filter.Statuses) > 0 {
		dataset = dataset.Where(goqu.Ex{"status": filter.Statuses})
	}
	return dataset
}

// CountTableRecords counts records in a specific table for a tenant
// Supported tables: santri, siswa, guru, ustadz, staf, kelas
func (a *tenantAdapter) CountTableRecords(tenantID string, tableName string) (int, error) {
	// Whitelist of allowed tables to prevent SQL injection
	allowedTables := map[string]bool{
		"santri":   true,
		"siswa":    true,
		"guru":     true,
		"ustadz":   true,
		"staf":     true,
		"kelas":    true,
		"rombel":   true,
		"asrama":   true,
		"mapel":    true,
		"students": true, // Unified students table
	}

	if !allowedTables[tableName] {
		return 0, nil // Unknown table, return 0
	}

	// Use raw SQL for dynamic table name (safe because whitelist validated)
	query := `SELECT COUNT(*) FROM ` + tableName + ` WHERE tenant_id = $1`

	var count int
	err := a.db.QueryRow(query, tenantID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
