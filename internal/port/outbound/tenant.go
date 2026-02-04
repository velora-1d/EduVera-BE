package outbound_port

import "prabogo/internal/model"

//go:generate mockgen -source=tenant.go -destination=./../../../tests/mocks/port/mock_tenant.go
type TenantDatabasePort interface {
	Create(tenant *model.Tenant) error
	Update(tenant *model.Tenant) error
	FindByFilter(filter model.TenantFilter) ([]model.Tenant, error)
	FindByID(id string) (*model.Tenant, error)
	FindBySubdomain(subdomain string) (*model.Tenant, error)
	SubdomainExists(subdomain string) (bool, error)
	UpdateStatus(id string, status string) error
	CountTableRecords(tenantID string, tableName string) (int, error) // For data limit enforcement
}
