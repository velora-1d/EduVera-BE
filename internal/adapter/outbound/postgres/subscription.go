package postgres_outbound_adapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
)

type subscriptionAdapter struct {
	db *sql.DB
}

func NewSubscriptionAdapter(db *sql.DB) outbound_port.SubscriptionDatabasePort {
	return &subscriptionAdapter{
		db: db,
	}
}

// CreateSubscription creates a new subscription record
func (a *subscriptionAdapter) CreateSubscription(ctx context.Context, sub *model.Subscription) error {
	query, args, err := goqu.Insert("subscriptions").Rows(
		goqu.Record{
			"id":                   sub.ID,
			"tenant_id":            sub.TenantID,
			"plan_type":            sub.PlanType,
			"billing_cycle":        sub.BillingCycle,
			"status":               sub.Status,
			"current_period_start": sub.CurrentPeriodStart,
			"current_period_end":   sub.CurrentPeriodEnd,
			"grace_period_end":     sub.GracePeriodEnd,
			"created_at":           sub.CreatedAt,
			"updated_at":           sub.UpdatedAt,
		},
	).ToSQL()

	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}

// GetSubscription retrieves a single subscription based on filter
func (a *subscriptionAdapter) GetSubscription(ctx context.Context, filter model.SubscriptionFilter) (*model.Subscription, error) {
	subs, err := a.GetSubscriptions(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return nil, nil // Not found
	}
	return &subs[0], nil
}

// GetSubscriptions retrieves list of subscriptions based on filter
func (a *subscriptionAdapter) GetSubscriptions(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error) {
	dataset := goqu.From("subscriptions")

	if len(filter.IDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}
	if len(filter.TenantIDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"tenant_id": filter.TenantIDs})
	}
	if len(filter.Statuses) > 0 {
		dataset = dataset.Where(goqu.Ex{"status": filter.Statuses})
	}
	if len(filter.PlanTypes) > 0 {
		dataset = dataset.Where(goqu.Ex{"plan_type": filter.PlanTypes})
	}

	query, args, err := dataset.Select(
		"id", "tenant_id", "plan_type", "billing_cycle", "status",
		"current_period_start", "current_period_end", "grace_period_end",
		"cancelled_at", "scheduled_plan_type", "created_at", "updated_at",
	).ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var s model.Subscription
		err := rows.Scan(
			&s.ID, &s.TenantID, &s.PlanType, &s.BillingCycle, &s.Status,
			&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.GracePeriodEnd,
			&s.CancelledAt, &s.ScheduledPlanType, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

// UpdateSubscription updates existing subscription
func (a *subscriptionAdapter) UpdateSubscription(ctx context.Context, sub *model.Subscription) error {
	query, args, err := goqu.Update("subscriptions").Set(
		goqu.Record{
			"plan_type":           sub.PlanType,
			"billing_cycle":       sub.BillingCycle,
			"status":              sub.Status,
			"current_period_end":  sub.CurrentPeriodEnd,
			"grace_period_end":    sub.GracePeriodEnd,
			"cancelled_at":        sub.CancelledAt,
			"scheduled_plan_type": sub.ScheduledPlanType,
			"updated_at":          time.Now(),
		},
	).Where(goqu.Ex{"id": sub.ID}).ToSQL()

	if err != nil {
		return err
	}

	result, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

// GetPricingPlan retrieves a specific pricing plan
func (a *subscriptionAdapter) GetPricingPlan(ctx context.Context, planType, billingCycle string) (*model.PricingPlan, error) {
	plans, err := a.GetPricingPlans(ctx, model.PricingFilter{
		PlanTypes:     []string{planType},
		BillingCycles: []string{billingCycle},
	})
	if err != nil {
		return nil, err
	}
	if len(plans) == 0 {
		return nil, nil
	}
	return &plans[0], nil
}

// GetPricingPlans retrieves list of pricing plans
func (a *subscriptionAdapter) GetPricingPlans(ctx context.Context, filter model.PricingFilter) ([]model.PricingPlan, error) {
	dataset := goqu.From("pricing_plans")

	if len(filter.PlanTypes) > 0 {
		dataset = dataset.Where(goqu.Ex{"plan_type": filter.PlanTypes})
	}
	if len(filter.BillingCycles) > 0 {
		dataset = dataset.Where(goqu.Ex{"billing_cycle": filter.BillingCycles})
	}
	if filter.IsActive != nil {
		dataset = dataset.Where(goqu.Ex{"is_active": *filter.IsActive})
	}

	query, args, err := dataset.Select(
		"id", "plan_type", "billing_cycle", "price", "currency",
		"description", "is_active", "created_at", "updated_at",
	).ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []model.PricingPlan
	for rows.Next() {
		var p model.PricingPlan
		err := rows.Scan(
			&p.ID, &p.PlanType, &p.BillingCycle, &p.Price, &p.Currency,
			&p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

// RecordSubscriptionHistory logs audit trail for changes
func (a *subscriptionAdapter) RecordSubscriptionHistory(ctx context.Context, subID, action, oldPlan, newPlan, oldStatus, newStatus string, amount int64, notes string) error {
	query, args, err := goqu.Insert("subscription_history").Rows(
		goqu.Record{
			"subscription_id": subID,
			"action":          action,
			"old_plan_type":   oldPlan,
			"new_plan_type":   newPlan,
			"old_status":      oldStatus,
			"new_status":      newStatus,
			"amount":          amount,
			"notes":           notes,
		},
	).ToSQL()

	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}

// GetExpiredSubscriptions returns subscriptions that have passed period end
func (a *subscriptionAdapter) GetExpiredSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	// Active status but Period End < NOW
	// Should transition to Grace Period
	query, args, err := goqu.From("subscriptions").
		Where(
			goqu.Ex{"status": model.SubscriptionStatusActive},
			goqu.C("current_period_end").Lt(time.Now()),
		).Select(
		"id", "tenant_id", "plan_type", "billing_cycle", "status",
		"current_period_start", "current_period_end", "grace_period_end",
		"cancelled_at", "scheduled_plan_type", "created_at", "updated_at",
	).ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var s model.Subscription
		err := rows.Scan(
			&s.ID, &s.TenantID, &s.PlanType, &s.BillingCycle, &s.Status,
			&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.GracePeriodEnd,
			&s.CancelledAt, &s.ScheduledPlanType, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

// GetGracePeriodExpiredSubscriptions returns subscriptions that have passed grace period
func (a *subscriptionAdapter) GetGracePeriodExpiredSubscriptions(ctx context.Context) ([]model.Subscription, error) {
	// Grace Period status but Grace Period End < NOW
	// Should transition to Suspended
	query, args, err := goqu.From("subscriptions").
		Where(
			goqu.Ex{"status": model.SubscriptionStatusGracePeriod},
			goqu.C("grace_period_end").Lt(time.Now()),
		).Select(
		"id", "tenant_id", "plan_type", "billing_cycle", "status",
		"current_period_start", "current_period_end", "grace_period_end",
		"cancelled_at", "scheduled_plan_type", "created_at", "updated_at",
	).ToSQL()

	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var s model.Subscription
		err := rows.Scan(
			&s.ID, &s.TenantID, &s.PlanType, &s.BillingCycle, &s.Status,
			&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.GracePeriodEnd,
			&s.CancelledAt, &s.ScheduledPlanType, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

// UpdateSubscriptionStatus batch update status helper
func (a *subscriptionAdapter) UpdateSubscriptionStatus(ctx context.Context, id, status string) error {
	query, args, err := goqu.Update("subscriptions").Set(
		goqu.Record{
			"status":     status,
			"updated_at": time.Now(),
		},
	).Where(goqu.Ex{"id": id}).ToSQL()

	if err != nil {
		return err
	}

	_, err = a.db.ExecContext(ctx, query, args...)
	return err
}
