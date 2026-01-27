package tenant

import (
	"context"

	"github.com/palantir/stacktrace"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type TenantDomain interface {
	Create(ctx context.Context, input *model.TenantInput) (*model.Tenant, error)
	FindByID(ctx context.Context, id string) (*model.Tenant, error)
	FindBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error)
	SubdomainExists(ctx context.Context, subdomain string) (bool, error)
	UpdateInstitution(ctx context.Context, id string, input *model.TenantInput) error
	UpdateBankAccount(ctx context.Context, id string, bankName, accountNumber, accountHolder string) error
	Activate(ctx context.Context, id string) error
}

type tenantDomain struct {
	databasePort outbound_port.DatabasePort
}

func NewTenantDomain(databasePort outbound_port.DatabasePort) TenantDomain {
	return &tenantDomain{
		databasePort: databasePort,
	}
}

func (d *tenantDomain) Create(ctx context.Context, input *model.TenantInput) (*model.Tenant, error) {
	// Check if subdomain already exists
	exists, err := d.databasePort.Tenant().SubdomainExists(input.Subdomain)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to check subdomain")
	}
	if exists {
		return nil, stacktrace.NewError("subdomain already exists")
	}

	tenant := model.TenantPrepare(input)
	err = d.databasePort.Tenant().Create(tenant)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to create tenant")
	}

	return tenant, nil
}

func (d *tenantDomain) FindByID(ctx context.Context, id string) (*model.Tenant, error) {
	tenant, err := d.databasePort.Tenant().FindByID(id)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find tenant")
	}
	return tenant, nil
}

func (d *tenantDomain) FindBySubdomain(ctx context.Context, subdomain string) (*model.Tenant, error) {
	tenant, err := d.databasePort.Tenant().FindBySubdomain(subdomain)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to find tenant by subdomain")
	}
	return tenant, nil
}

func (d *tenantDomain) SubdomainExists(ctx context.Context, subdomain string) (bool, error) {
	return d.databasePort.Tenant().SubdomainExists(subdomain)
}

func (d *tenantDomain) UpdateInstitution(ctx context.Context, id string, input *model.TenantInput) error {
	tenant, err := d.databasePort.Tenant().FindByID(id)
	if err != nil {
		return stacktrace.Propagate(err, "failed to find tenant")
	}

	tenant.Name = input.Name
	tenant.InstitutionType = input.InstitutionType
	tenant.Address = input.Address

	return d.databasePort.Tenant().Update(tenant)
}

func (d *tenantDomain) UpdateBankAccount(ctx context.Context, id string, bankName, accountNumber, accountHolder string) error {
	tenant, err := d.databasePort.Tenant().FindByID(id)
	if err != nil {
		return stacktrace.Propagate(err, "failed to find tenant")
	}

	tenant.BankName = bankName
	tenant.AccountNumber = accountNumber
	tenant.AccountHolder = accountHolder

	return d.databasePort.Tenant().Update(tenant)
}

func (d *tenantDomain) Activate(ctx context.Context, id string) error {
	return d.databasePort.Tenant().UpdateStatus(id, model.TenantStatusActive)
}
