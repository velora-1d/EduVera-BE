package audit_log

import (
	"context"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

// AuditHelper provides convenient methods for logging common events
type AuditHelper struct {
	repo outbound_port.AuditLogDatabasePort
}

func NewAuditHelper(repo outbound_port.AuditLogDatabasePort) *AuditHelper {
	return &AuditHelper{repo: repo}
}

// LogSubscriptionEvent logs subscription-related events
func (h *AuditHelper) LogSubscriptionEvent(ctx context.Context, action, tenantID, description string) error {
	log := &model.AuditLog{
		Action:      action,
		TargetType:  "subscription",
		TargetID:    tenantID,
		Description: description,
	}
	return h.repo.Create(log)
}

// LogPaymentEvent logs payment-related events
func (h *AuditHelper) LogPaymentEvent(ctx context.Context, action, tenantID, orderId, amount string) error {
	log := &model.AuditLog{
		Action:      action,
		TargetType:  "payment",
		TargetID:    orderId,
		Description: "Tenant: " + tenantID + ", Amount: " + amount,
	}
	return h.repo.Create(log)
}

// LogLoginEvent logs authentication events
func (h *AuditHelper) LogLoginEvent(ctx context.Context, action, userID, email, ipAddress string) error {
	log := &model.AuditLog{
		Action:      action,
		TargetType:  "user",
		TargetID:    userID,
		AdminEmail:  email,
		IPAddress:   ipAddress,
		Description: "Login attempt for: " + email,
	}
	return h.repo.Create(log)
}

// LogWhatsAppEvent logs WhatsApp connection events
func (h *AuditHelper) LogWhatsAppEvent(ctx context.Context, action, tenantID, instanceName string) error {
	log := &model.AuditLog{
		Action:      action,
		TargetType:  "whatsapp",
		TargetID:    tenantID,
		Description: "Instance: " + instanceName,
	}
	return h.repo.Create(log)
}

// Quick static helpers for common events

// LogSubscriptionCreated logs when a new subscription is created
func LogSubscriptionCreated(repo outbound_port.AuditLogDatabasePort, tenantID, planType string) {
	h := NewAuditHelper(repo)
	_ = h.LogSubscriptionEvent(context.Background(), model.AuditActionSubscriptionCreated, tenantID, "Plan: "+planType)
}

// LogSubscriptionUpgraded logs when a subscription is upgraded
func LogSubscriptionUpgraded(repo outbound_port.AuditLogDatabasePort, tenantID, oldPlan, newPlan string) {
	h := NewAuditHelper(repo)
	_ = h.LogSubscriptionEvent(context.Background(), model.AuditActionSubscriptionUpgraded, tenantID, "From "+oldPlan+" to "+newPlan)
}

// LogPaymentSuccess logs successful payment
func LogPaymentSuccess(repo outbound_port.AuditLogDatabasePort, tenantID, orderId, amount string) {
	h := NewAuditHelper(repo)
	_ = h.LogPaymentEvent(context.Background(), model.AuditActionPaymentSuccess, tenantID, orderId, amount)
}

// LogPaymentFailed logs failed payment
func LogPaymentFailed(repo outbound_port.AuditLogDatabasePort, tenantID, orderId, reason string) {
	h := NewAuditHelper(repo)
	_ = h.LogPaymentEvent(context.Background(), model.AuditActionPaymentFailed, tenantID, orderId, reason)
}
