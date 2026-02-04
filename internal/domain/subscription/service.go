package subscription

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"prabogo/internal/domain/audit_log"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/google/uuid"
)

type SubscriptionDomain interface {
	CreateSubscription(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error)
	GetSubscription(ctx context.Context, tenantID string) (*model.Subscription, error)
	UpgradePlan(ctx context.Context, input model.UpgradeInput) (*model.UpgradeCalculation, error)
	DowngradePlan(ctx context.Context, input model.DowngradeInput) (*model.Subscription, error)
	CalculateUpgrade(ctx context.Context, input model.UpgradeInput) (*model.UpgradeCalculation, error)
	RenewSubscription(ctx context.Context, tenantID, orderID string) error
	CheckExpiredSubscriptions(ctx context.Context) error
	GetPricingPlans(ctx context.Context, filter model.PricingFilter) ([]model.PricingPlan, error)
}

type subscriptionDomain struct {
	repo       outbound_port.SubscriptionDatabasePort
	tenantRepo outbound_port.DatabasePort // for tenant updates
}

func NewSubscriptionDomain(repo outbound_port.SubscriptionDatabasePort, tenantRepo outbound_port.DatabasePort) SubscriptionDomain {
	return &subscriptionDomain{
		repo:       repo,
		tenantRepo: tenantRepo,
	}
}

func (d *subscriptionDomain) CreateSubscription(ctx context.Context, input model.CreateSubscriptionInput) (*model.Subscription, error) {
	// Calculate period
	start := time.Now()
	end := calculatePeriodEnd(start, input.BillingCycle)

	sub := &model.Subscription{
		ID:                 uuid.New().String(),
		TenantID:           input.TenantID,
		PlanType:           input.PlanType,
		BillingCycle:       input.BillingCycle,
		Status:             model.SubscriptionStatusActive,
		CurrentPeriodStart: start,
		CurrentPeriodEnd:   end,
		GracePeriodEnd:     end.AddDate(0, 0, model.GracePeriodDays),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := d.repo.CreateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	// Audit Log
	audit_log.LogSubscriptionCreated(d.tenantRepo.AuditLog(), sub.TenantID, sub.PlanType)

	// Also record history
	d.repo.RecordSubscriptionHistory(ctx, sub.ID, "created", "", input.PlanType, "", model.SubscriptionStatusActive, 0, "Initial subscription")

	return sub, nil
}

func (d *subscriptionDomain) GetSubscription(ctx context.Context, tenantID string) (*model.Subscription, error) {
	return d.repo.GetSubscription(ctx, model.SubscriptionFilter{
		TenantIDs: []string{tenantID},
		Statuses:  []string{model.SubscriptionStatusActive, model.SubscriptionStatusGracePeriod},
	})
}

func (d *subscriptionDomain) CalculateUpgrade(ctx context.Context, input model.UpgradeInput) (*model.UpgradeCalculation, error) {
	sub, err := d.GetSubscription(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("active subscription not found")
	}

	// Check if upgrade is valid (e.g. not same plan)
	if sub.PlanType == input.NewPlanType {
		return nil, errors.New("cannot upgrade to same plan")
	}

	// Get Prices
	// Assume billing cycle remains same for upgrade calculation simplicity for now
	// Or force upgrade to Annual? Let's stick to same cycle
	oldPricePlan, err := d.repo.GetPricingPlan(ctx, sub.PlanType, sub.BillingCycle)
	if err != nil {
		return nil, err
	}
	newPricePlan, err := d.repo.GetPricingPlan(ctx, input.NewPlanType, sub.BillingCycle)
	if err != nil {
		return nil, err
	}

	if oldPricePlan == nil || newPricePlan == nil {
		return nil, errors.New("pricing plan not found")
	}

	// Calculate remaining days
	daysRemaining := sub.DaysRemaining()
	totalDays := sub.TotalDays()
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// Prorata Credit: (Remaining / Total) * OldPrice
	prorataCredit := int64(float64(daysRemaining) / float64(totalDays) * float64(oldPricePlan.Price))

	// New Price for full period?
	// Usually upgrade resets the cycle OR charges diff for remainder.
	// Scenario: Reset cycle starts NOW. Credit applied.
	// Amount Due = NewPrice - ProrataCredit
	amountDue := newPricePlan.Price - prorataCredit
	if amountDue < 0 {
		amountDue = 0 // No refund, just 0 charge (or carry over credit, but let's keep simple)
	}

	return &model.UpgradeCalculation{
		OldPlanType:   sub.PlanType,
		NewPlanType:   input.NewPlanType,
		RemainingDays: daysRemaining,
		OldPrice:      oldPricePlan.Price,
		NewPrice:      newPricePlan.Price,
		ProrataCredit: prorataCredit,
		AmountDue:     amountDue,
	}, nil
}

func (d *subscriptionDomain) UpgradePlan(ctx context.Context, input model.UpgradeInput) (*model.UpgradeCalculation, error) {
	// 1. Calculate
	calc, err := d.CalculateUpgrade(ctx, input)
	if err != nil {
		return nil, err
	}

	// Getting subscription to know billing cycle and current tier
	sub, _ := d.GetSubscription(ctx, input.TenantID)
	// We assume upgrade is to Premium if current is Basic
	targetTier := "premium" // Default upgrade target
	if sub != nil && sub.SubscriptionTier == "premium" {
		// Already premium? Keep premium
		targetTier = "premium"
	}

	// Calculate price using pricing model
	isAnnual := false
	if sub != nil {
		isAnnual = sub.BillingCycle == model.BillingCycleAnnual
	}
	if input.BillingCycle == model.BillingCycleAnnual {
		isAnnual = true
	}

	price := model.GetPlanPriceWithTier(input.NewPlanType, targetTier, isAnnual)

	calc.NewPrice = price
	calc.TargetTier = targetTier

	// Note: Payment creation is now handled by HTTP adapter calling PaymentDomain.CreateSnapTransaction
	// The PaymentURL will be set at the adapter level after calling Midtrans

	return calc, nil
}

func (d *subscriptionDomain) DowngradePlan(ctx context.Context, input model.DowngradeInput) (*model.Subscription, error) {
	sub, err := d.GetSubscription(ctx, input.TenantID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("active subscription not found")
	}

	// Downgrade is scheduled for END of period
	if sub.ScheduledPlanType != nil {
		return nil, errors.New("downgrade already scheduled")
	}

	// TODO: Validate data hiding if !input.Confirmed (Option 2 Logic)
	// For now, we trust confirmed=true from frontend

	sub.ScheduledPlanType = &input.NewPlanType
	if err := d.repo.UpdateSubscription(ctx, sub); err != nil {
		return nil, err
	}

	d.repo.RecordSubscriptionHistory(ctx, sub.ID, "downgrade_scheduled", sub.PlanType, input.NewPlanType, sub.Status, sub.Status, 0, "Downgrade scheduled for period end")

	return sub, nil
}

func (d *subscriptionDomain) RenewSubscription(ctx context.Context, tenantID, orderID string) error {
	sub, err := d.GetSubscription(ctx, tenantID)
	if err != nil {
		return err
	}
	if sub == nil {
		return errors.New("subscription not found")
	}

	// If there was a scheduled downgrade, apply it now
	if sub.ScheduledPlanType != nil {
		oldPlan := sub.PlanType
		sub.PlanType = *sub.ScheduledPlanType
		sub.ScheduledPlanType = nil
		d.repo.RecordSubscriptionHistory(ctx, sub.ID, "downgrade_executed", oldPlan, sub.PlanType, sub.Status, sub.Status, 0, "Scheduled downgrade executed")
	}

	// Extend period
	// Start from previous end if not expired long ago, else from now
	start := sub.CurrentPeriodEnd
	if time.Now().After(start.AddDate(0, 1, 0)) { // If expired > 1 month ago, reset
		start = time.Now()
	}

	sub.CurrentPeriodStart = start
	sub.CurrentPeriodEnd = calculatePeriodEnd(start, sub.BillingCycle)
	sub.GracePeriodEnd = sub.CurrentPeriodEnd.AddDate(0, 0, model.GracePeriodDays)
	sub.Status = model.SubscriptionStatusActive // Reset status from grace/suspended

	if err := d.repo.UpdateSubscription(ctx, sub); err != nil {
		return err
	}

	// Audit Log
	helper := audit_log.NewAuditHelper(d.tenantRepo.AuditLog())
	_ = helper.LogSubscriptionEvent(ctx, model.AuditActionSubscriptionRenewed, sub.TenantID, fmt.Sprintf("Renewed via Order %s", orderID))

	d.repo.RecordSubscriptionHistory(ctx, sub.ID, "renewed", sub.PlanType, sub.PlanType, sub.Status, model.SubscriptionStatusActive, 0, fmt.Sprintf("Renewed via Order %s", orderID))
	return nil
}

func (d *subscriptionDomain) CheckExpiredSubscriptions(ctx context.Context) error {
	// 1. Check Active -> Grace Period
	expiredSubs, err := d.repo.GetExpiredSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, sub := range expiredSubs {
		sub.Status = model.SubscriptionStatusGracePeriod
		if err := d.repo.UpdateSubscription(ctx, &sub); err != nil {
			log.Printf("Failed to update sub %s to grace period: %v", sub.ID, err)
			continue
		}
		d.repo.RecordSubscriptionHistory(ctx, sub.ID, "status_change", sub.PlanType, sub.PlanType, model.SubscriptionStatusActive, model.SubscriptionStatusGracePeriod, 0, "Auto-expired to grace period")
	}

	// 2. Check Grace Period -> Suspended
	graceExpiredSubs, err := d.repo.GetGracePeriodExpiredSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, sub := range graceExpiredSubs {
		sub.Status = model.SubscriptionStatusSuspended
		if err := d.repo.UpdateSubscription(ctx, &sub); err != nil {
			log.Printf("Failed to update sub %s to suspended: %v", sub.ID, err)
			continue
		}

		// Audit Log
		helper := audit_log.NewAuditHelper(d.tenantRepo.AuditLog())
		_ = helper.LogSubscriptionEvent(ctx, model.AuditActionSubscriptionSuspended, sub.TenantID, "Grace period ended, account suspended")

		d.repo.RecordSubscriptionHistory(ctx, sub.ID, "status_change", sub.PlanType, sub.PlanType, model.SubscriptionStatusGracePeriod, model.SubscriptionStatusSuspended, 0, "Grace period ended, suspended")
	}

	return nil
}

func calculatePeriodEnd(start time.Time, cycle string) time.Time {
	if cycle == model.BillingCycleAnnual {
		return start.AddDate(1, 0, 0)
	}
	return start.AddDate(0, 1, 0)
}

func (d *subscriptionDomain) GetPricingPlans(ctx context.Context, filter model.PricingFilter) ([]model.PricingPlan, error) {
	return d.repo.GetPricingPlans(ctx, filter)
}
