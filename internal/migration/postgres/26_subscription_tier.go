package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationNoTxContext(upSubscriptionTier, downSubscriptionTier)
}

// upSubscriptionTier adds tier field for feature gating (basic=manual, premium=PG)
func upSubscriptionTier(ctx context.Context, db *sql.DB) error {
	queries := []struct {
		name  string
		query string
	}{
		{
			name:  "add subscription_tier to tenants",
			query: `ALTER TABLE tenants ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(20) DEFAULT 'basic'`,
		},
		{
			name:  "add subscription_tier to subscriptions",
			query: `ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(20) DEFAULT 'basic'`,
		},
		{
			name:  "add subscription_tier to pricing_plans",
			query: `ALTER TABLE pricing_plans ADD COLUMN IF NOT EXISTS subscription_tier VARCHAR(20) DEFAULT 'basic'`,
		},
		{
			name:  "create index on subscription_tier",
			query: `CREATE INDEX IF NOT EXISTS idx_tenants_tier ON tenants(subscription_tier)`,
		},
		// Insert premium pricing (2x basic price)
		{
			name: "insert premium sekolah monthly",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('sekolah', 'monthly', 999000, 'Sekolah Premium - Bulanan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
		{
			name: "insert premium sekolah annual",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('sekolah', 'annual', 9990000, 'Sekolah Premium - Tahunan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
		{
			name: "insert premium pesantren monthly",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('pesantren', 'monthly', 999000, 'Pesantren Premium - Bulanan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
		{
			name: "insert premium pesantren annual",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('pesantren', 'annual', 9990000, 'Pesantren Premium - Tahunan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
		{
			name: "insert premium hybrid monthly",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('hybrid', 'monthly', 1599000, 'Hybrid Premium - Bulanan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
		{
			name: "insert premium hybrid annual",
			query: `INSERT INTO pricing_plans (plan_type, billing_cycle, price, description, subscription_tier)
			        VALUES ('hybrid', 'annual', 15990000, 'Hybrid Premium - Tahunan (Auto PG)', 'premium')
			        ON CONFLICT DO NOTHING`,
		},
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q.query); err != nil {
			fmt.Printf("Warning: %s: %v\n", q.name, err)
		}
	}

	return nil
}

func downSubscriptionTier(ctx context.Context, db *sql.DB) error {
	queries := []string{
		`DELETE FROM pricing_plans WHERE subscription_tier = 'premium'`,
		`DROP INDEX IF EXISTS idx_tenants_tier`,
		`ALTER TABLE tenants DROP COLUMN IF EXISTS subscription_tier`,
		`ALTER TABLE subscriptions DROP COLUMN IF EXISTS subscription_tier`,
		`ALTER TABLE pricing_plans DROP COLUMN IF EXISTS subscription_tier`,
	}

	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			fmt.Printf("Warning during rollback: %v\n", err)
		}
	}

	return nil
}
