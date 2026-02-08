package scheduler

import (
	"context"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"

	"prabogo/internal/adapter/outbound/notification"
	"prabogo/internal/domain"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/utils/log"
)

// Scheduler manages scheduled tasks
type Scheduler struct {
	ctx      context.Context
	cron     *cron.Cron
	domain   domain.Domain
	db       outbound_port.DatabasePort
	telegram *notification.TelegramNotifier
}

// NewScheduler creates a new scheduler instance
func NewScheduler(ctx context.Context, d domain.Domain, db outbound_port.DatabasePort) *Scheduler {
	return &Scheduler{
		ctx:      ctx,
		cron:     cron.New(),
		domain:   d,
		db:       db,
		telegram: notification.NewTelegramNotifier(),
	}
}

// Start begins all scheduled jobs
func (s *Scheduler) Start() {
	log.WithContext(s.ctx).Info("Starting scheduler...")

	// Run subscription reminder every day at 08:00 WIB (01:00 UTC)
	s.cron.AddFunc("0 0 1 * * *", func() {
		s.checkSubscriptionReminders()
	})

	// Run monthly invoice generation every 1st of month at 00:00
	s.cron.AddFunc("0 0 1 * * *", func() {
		s.generateMonthlyInvoices()
	})

	// Also run at startup for testing (delayed by 10 seconds)
	go func() {
		time.Sleep(10 * time.Second)
		log.WithContext(s.ctx).Info("Running initial subscription check...")
		s.checkSubscriptionReminders()
	}()

	s.cron.Start()
	log.WithContext(s.ctx).Info("Scheduler started successfully")
}

// Stop gracefully stops the scheduler
func (s *Scheduler) Stop() {
	log.WithContext(s.ctx).Info("Stopping scheduler...")
	s.cron.Stop()
}

// checkSubscriptionReminders checks for expiring subscriptions and sends notifications
func (s *Scheduler) checkSubscriptionReminders() {
	ctx := s.ctx
	log.WithContext(ctx).Info("Checking subscription reminders...")

	// Check for subscriptions expiring in 7, 3, and 1 days
	reminderDays := []int{7, 3, 1}

	for _, days := range reminderDays {
		s.sendRemindersForDays(days)
	}

	log.WithContext(ctx).Info("Subscription reminder check completed")
}

// sendRemindersForDays sends reminders for subscriptions expiring in N days
func (s *Scheduler) sendRemindersForDays(days int) {
	ctx := s.ctx

	// Get subscriptions expiring in N days
	subscriptions, err := s.getExpiringSubscriptions(days)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to get subscriptions expiring in %d days", days)
		return
	}

	if len(subscriptions) == 0 {
		log.WithContext(ctx).Infof("No subscriptions expiring in %d days", days)
		return
	}

	log.WithContext(ctx).WithField("count", len(subscriptions)).Infof("Found %d subscriptions expiring in %d days", len(subscriptions), days)

	for _, sub := range subscriptions {
		s.sendReminderNotification(sub, days)
	}
}

// getExpiringSubscriptions queries subscriptions expiring in N days
func (s *Scheduler) getExpiringSubscriptions(days int) ([]model.Subscription, error) {
	ctx := s.ctx

	// Get all active subscriptions via database port
	filter := model.SubscriptionFilter{
		Statuses: []string{model.SubscriptionStatusActive},
	}

	subscriptions, err := s.db.Subscription().GetSubscriptions(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Filter by expiry date
	var expiring []model.Subscription
	targetDate := time.Now().AddDate(0, 0, days).Truncate(24 * time.Hour)

	for _, sub := range subscriptions {
		expiryDate := sub.CurrentPeriodEnd.Truncate(24 * time.Hour)
		if expiryDate.Equal(targetDate) {
			expiring = append(expiring, sub)
		}
	}

	return expiring, nil
}

// sendReminderNotification sends Telegram and WhatsApp notifications
func (s *Scheduler) sendReminderNotification(sub model.Subscription, daysLeft int) {
	ctx := s.ctx

	// Get tenant info via domain
	tenant, err := s.domain.Tenant().FindByID(ctx, sub.TenantID)
	if err != nil || tenant == nil {
		log.WithContext(ctx).WithError(err).Errorf("Failed to get tenant %s for reminder", sub.TenantID)
		return
	}

	institutionName := tenant.Name
	subdomain := tenant.Subdomain
	expiryDate := sub.CurrentPeriodEnd.Format("02 January 2006")

	// Send Telegram notification to owner
	err = s.telegram.SendSubscriptionReminder(institutionName, subdomain, daysLeft, expiryDate)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to send Telegram subscription reminder")
	} else {
		log.WithContext(ctx).WithFields(logrus.Fields{
			"tenant":    institutionName,
			"days_left": daysLeft,
		}).Info("Sent Telegram subscription reminder")
	}

	// TODO: Send WhatsApp to tenant admin
}

// generateMonthlyInvoices generates SPP invoices for all active students
func (s *Scheduler) generateMonthlyInvoices() {
	ctx := s.ctx
	log.WithContext(ctx).Info("Generating monthly invoices...")

	if err := s.domain.SPP().GenerateInvoices(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("Failed to generate monthly invoices")
	} else {
		log.WithContext(ctx).Info("Monthly invoices generated successfully")
	}
}
