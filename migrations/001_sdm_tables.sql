-- SDM Module Tables Migration
-- Created: 2026-01-28
-- Description: Create tables for Employee, Payroll, Attendance, and PayrollConfig

-- =============================================
-- EMPLOYEES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    nip VARCHAR(50),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    role VARCHAR(50) NOT NULL, -- Guru, Staf, Kepala Sekolah, Wakil Kepala Sekolah, Tata Usaha
    employee_type VARCHAR(50) NOT NULL, -- PNS, Honorer, Kontrak, Tetap Yayasan
    department VARCHAR(100),
    join_date DATE,
    base_salary BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for employees
CREATE INDEX IF NOT EXISTS idx_employees_tenant_id ON employees(tenant_id);
CREATE INDEX IF NOT EXISTS idx_employees_role ON employees(role);
CREATE INDEX IF NOT EXISTS idx_employees_is_active ON employees(is_active);
CREATE UNIQUE INDEX IF NOT EXISTS idx_employees_nip_unique ON employees(tenant_id, nip) WHERE nip IS NOT NULL AND nip != '';

-- =============================================
-- PAYROLL_CONFIGS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS payroll_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    components JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint: one config per tenant
CREATE UNIQUE INDEX IF NOT EXISTS idx_payroll_configs_tenant ON payroll_configs(tenant_id);

-- =============================================
-- PAYROLLS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS payrolls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    period VARCHAR(7) NOT NULL, -- Format: YYYY-MM (e.g., 2026-01)
    base_salary BIGINT NOT NULL DEFAULT 0,
    allowances BIGINT NOT NULL DEFAULT 0, -- Total tunjangan
    deductions BIGINT NOT NULL DEFAULT 0, -- Total potongan
    net_salary BIGINT NOT NULL DEFAULT 0, -- Total terima (base + allowances - deductions)
    details JSONB DEFAULT '[]', -- Breakdown of each component
    status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, pending, paid
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for payrolls
CREATE INDEX IF NOT EXISTS idx_payrolls_tenant_id ON payrolls(tenant_id);
CREATE INDEX IF NOT EXISTS idx_payrolls_employee_id ON payrolls(employee_id);
CREATE INDEX IF NOT EXISTS idx_payrolls_period ON payrolls(period);
CREATE INDEX IF NOT EXISTS idx_payrolls_status ON payrolls(status);

-- Unique constraint: one payroll per employee per period
CREATE UNIQUE INDEX IF NOT EXISTS idx_payrolls_employee_period ON payrolls(employee_id, period);

-- =============================================
-- ATTENDANCES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS attendances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    check_in TIMESTAMP WITH TIME ZONE,
    check_out TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) NOT NULL DEFAULT 'hadir', -- hadir, sakit, izin, alpha
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for attendances
CREATE INDEX IF NOT EXISTS idx_attendances_tenant_id ON attendances(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances(employee_id);
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances(date);
CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances(status);

-- Unique constraint: one attendance per employee per date
CREATE UNIQUE INDEX IF NOT EXISTS idx_attendances_employee_date ON attendances(employee_id, date);

-- =============================================
-- TRIGGERS FOR updated_at
-- =============================================

-- Function to auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for employees
DROP TRIGGER IF EXISTS update_employees_updated_at ON employees;
CREATE TRIGGER update_employees_updated_at
    BEFORE UPDATE ON employees
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for payroll_configs
DROP TRIGGER IF EXISTS update_payroll_configs_updated_at ON payroll_configs;
CREATE TRIGGER update_payroll_configs_updated_at
    BEFORE UPDATE ON payroll_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for payrolls
DROP TRIGGER IF EXISTS update_payrolls_updated_at ON payrolls;
CREATE TRIGGER update_payrolls_updated_at
    BEFORE UPDATE ON payrolls
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- SAMPLE DATA (Optional - for testing)
-- =============================================
-- Uncomment below to insert sample employees for testing

-- INSERT INTO employees (tenant_id, nip, name, email, role, employee_type, department, join_date, base_salary)
-- VALUES 
--     ('your-tenant-uuid-here', '19650101001', 'Drs. H. Suherman, M.Pd', 'suherman@school.id', 'Kepala Sekolah', 'PNS', 'Kepala Sekolah', '2000-01-01', 5000000),
--     ('your-tenant-uuid-here', '19850101002', 'Budi Santoso, S.Kom', 'budi@school.id', 'Guru', 'PNS', 'Produktif TKJ', '2010-07-01', 3500000),
--     ('your-tenant-uuid-here', NULL, 'Ahmad Dahlan', 'ahmad@school.id', 'Tata Usaha', 'Honorer', 'Tata Usaha', '2020-01-15', 2000000);
