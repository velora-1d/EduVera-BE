package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upSdmTables, downSdmTables)
}

func upSdmTables(ctx context.Context, tx *sql.Tx) error {
	// Employees Table
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			nip VARCHAR(50),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			phone VARCHAR(20),
			role VARCHAR(50) NOT NULL,
			employee_type VARCHAR(50) NOT NULL,
			department VARCHAR(100),
			join_date DATE,
			base_salary BIGINT DEFAULT 0,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_employees_tenant_id ON employees(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_employees_role ON employees(role);
		CREATE INDEX IF NOT EXISTS idx_employees_is_active ON employees(is_active);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_nip_unique ON employees(tenant_id, nip) WHERE nip IS NOT NULL AND nip != '';

		CREATE TABLE IF NOT EXISTS payroll_configs (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			components JSONB NOT NULL DEFAULT '[]',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_payroll_configs_tenant ON payroll_configs(tenant_id);

		CREATE TABLE IF NOT EXISTS payrolls (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
			period VARCHAR(7) NOT NULL,
			base_salary BIGINT NOT NULL DEFAULT 0,
			allowances BIGINT NOT NULL DEFAULT 0,
			deductions BIGINT NOT NULL DEFAULT 0,
			net_salary BIGINT NOT NULL DEFAULT 0,
			details JSONB DEFAULT '[]',
			status VARCHAR(20) NOT NULL DEFAULT 'draft',
			paid_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_payrolls_tenant_id ON payrolls(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_payrolls_employee_id ON payrolls(employee_id);
		CREATE INDEX IF NOT EXISTS idx_payrolls_period ON payrolls(period);
		CREATE INDEX IF NOT EXISTS idx_payrolls_status ON payrolls(status);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_payrolls_employee_period ON payrolls(employee_id, period);

		CREATE TABLE IF NOT EXISTS attendances (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
			date DATE NOT NULL,
			check_in TIMESTAMP WITH TIME ZONE,
			check_out TIMESTAMP WITH TIME ZONE,
			status VARCHAR(20) NOT NULL DEFAULT 'hadir',
			notes TEXT,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_attendances_tenant_id ON attendances(tenant_id);
		CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances(employee_id);
		CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);
		CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances(status);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_attendances_employee_date ON attendances(employee_id, date);
	`)
	if err != nil {
		return err
	}

	// Triggers
	// Note: update_updated_at_column should already exist from previous migrations, but we can verify or assume it exists.
	// We'll proceed assuming it exists.
	_, err = tx.Exec(`
		-- Function to auto-update updated_at timestamp
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
		CREATE TRIGGER update_employees_updated_at BEFORE UPDATE ON employees FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

		DROP TRIGGER IF EXISTS update_payroll_configs_updated_at ON payroll_configs;
		CREATE TRIGGER update_payroll_configs_updated_at BEFORE UPDATE ON payroll_configs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

		DROP TRIGGER IF EXISTS update_payrolls_updated_at ON payrolls;
		CREATE TRIGGER update_payrolls_updated_at BEFORE UPDATE ON payrolls FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
	`)
	return err
}

func downSdmTables(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE IF EXISTS attendances;
		DROP TABLE IF EXISTS payrolls;
		DROP TABLE IF EXISTS payroll_configs;
		DROP TABLE IF EXISTS employees;
	`)
	return err
}
