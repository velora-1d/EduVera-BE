package outbound_port

import (
	"context"

	"prabogo/internal/model"
)

// SubscriptionDatabasePort defines subscription database operations
type SubscriptionDatabasePort interface {
	// Subscription CRUD
	CreateSubscription(ctx context.Context, sub *model.Subscription) error
	GetSubscription(ctx context.Context, filter model.SubscriptionFilter) (*model.Subscription, error)
	GetSubscriptions(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *model.Subscription) error

	// Pricing Plans
	GetPricingPlan(ctx context.Context, planType, billingCycle string) (*model.PricingPlan, error)
	GetPricingPlans(ctx context.Context, filter model.PricingFilter) ([]model.PricingPlan, error)

	// Subscription History
	RecordSubscriptionHistory(ctx context.Context, subID, action, oldPlan, newPlan, oldStatus, newStatus string, amount int64, notes string) error

	// Batch operations for scheduler
	GetExpiredSubscriptions(ctx context.Context) ([]model.Subscription, error)
	GetGracePeriodExpiredSubscriptions(ctx context.Context) ([]model.Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, id, status string) error
}
