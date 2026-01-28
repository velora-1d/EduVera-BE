-- Subscription & Billing Tables Migration
-- Created: 2026-01-28
-- Description: Create tables for subscriptions and pricing plans

-- =============================================
-- PRICING_PLANS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS pricing_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_type VARCHAR(20) NOT NULL, -- sekolah, pesantren, hybrid
    billing_cycle VARCHAR(20) NOT NULL, -- monthly, annual
    price BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Unique constraint: one pricing per plan + cycle
CREATE UNIQUE INDEX IF NOT EXISTS idx_pricing_plans_unique ON pricing_plans(plan_type, billing_cycle) WHERE is_active = TRUE;

-- Seed initial pricing (from hardcoded PlanPricing)
INSERT INTO pricing_plans (plan_type, billing_cycle, price, description) VALUES
    ('sekolah', 'monthly', 499000, 'Paket Sekolah Bulanan'),
    ('sekolah', 'annual', 4990000, 'Paket Sekolah Tahunan (Hemat 2 Bulan)'),
    ('pesantren', 'monthly', 499000, 'Paket Pesantren Bulanan'),
    ('pesantren', 'annual', 4990000, 'Paket Pesantren Tahunan (Hemat 2 Bulan)'),
    ('hybrid', 'monthly', 799000, 'Paket Hybrid Bulanan (Sekolah + Pesantren)'),
    ('hybrid', 'annual', 7990000, 'Paket Hybrid Tahunan (Hemat 2 Bulan)')
ON CONFLICT DO NOTHING;

-- =============================================
-- SUBSCRIPTIONS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plan_type VARCHAR(20) NOT NULL,
    billing_cycle VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, grace_period, suspended, cancelled, terminated
    current_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    current_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    grace_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    scheduled_plan_type VARCHAR(20), -- For scheduled downgrade
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for subscriptions
CREATE INDEX IF NOT EXISTS idx_subscriptions_tenant_id ON subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_period_end ON subscriptions(current_period_end);

-- Unique constraint: one active subscription per tenant
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_tenant_active 
    ON subscriptions(tenant_id) 
    WHERE status IN ('active', 'grace_period');

-- =============================================
-- SUBSCRIPTION_HISTORY TABLE (Audit Trail)
-- =============================================
CREATE TABLE IF NOT EXISTS subscription_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL, -- created, upgraded, downgraded, renewed, suspended, cancelled
    old_plan_type VARCHAR(20),
    new_plan_type VARCHAR(20),
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    amount BIGINT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_subscription_history_subscription ON subscription_history(subscription_id);

-- =============================================
-- TRIGGERS
-- =============================================

-- Auto-update updated_at
DROP TRIGGER IF EXISTS update_pricing_plans_updated_at ON pricing_plans;
CREATE TRIGGER update_pricing_plans_updated_at
    BEFORE UPDATE ON pricing_plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
