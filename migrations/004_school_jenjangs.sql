-- School Jenjangs (Multi-jenjang) Migration
-- Created: 2026-02-01
-- Description: Add school_jenjangs column for multi-select jenjang support (TK, SD, MI, SMP, MTs, SMA, MA, SMK)

-- =============================================
-- ADD SCHOOL_JENJANGS COLUMN TO TENANTS TABLE
-- =============================================
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS school_jenjangs TEXT[] DEFAULT '{}';

COMMENT ON COLUMN tenants.school_jenjangs IS 'Array of school levels: TK, SD, MI, SMP, MTs, SMA, MA, SMK';

-- Create index for efficient filtering by jenjang
CREATE INDEX IF NOT EXISTS idx_tenants_school_jenjangs ON tenants USING GIN(school_jenjangs);
