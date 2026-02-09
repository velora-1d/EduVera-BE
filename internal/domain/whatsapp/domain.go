package whatsapp

import (
	"context"
	"errors"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
	"prabogo/utils/redis"

	"github.com/google/uuid"
)

type WhatsAppDomain interface {
	ConnectTenant(ctx context.Context, tenantID string) (*model.WhatsAppSession, error)
	GetStatus(ctx context.Context, tenantID string) (*model.WhatsAppSession, error)
	DisconnectTenant(ctx context.Context, tenantID string) error
	SendMessage(ctx context.Context, tenantID, phone, message string) error
}

type whatsAppDomain struct {
	dbPort        outbound_port.DatabasePort
	evolutionPort outbound_port.WhatsAppClientPort
}

func NewWhatsAppDomain(dbPort outbound_port.DatabasePort, evolutionPort outbound_port.WhatsAppClientPort) WhatsAppDomain {
	return &whatsAppDomain{
		dbPort:        dbPort,
		evolutionPort: evolutionPort,
	}
}

// ConnectTenant creates Evolution API instance and returns QR code for scanning
// For owner sessions, pass empty tenantID
func (d *whatsAppDomain) ConnectTenant(ctx context.Context, tenantID string) (*model.WhatsAppSession, error) {
	// Determine instance name based on tenantID
	var instanceName string
	if tenantID == "" {
		// Owner session
		instanceName = "eduvera_owner"
		// Check by instance name for owner
		existing, err := d.dbPort.WhatsApp().GetByInstanceName(ctx, instanceName)
		if err == nil && existing != nil && existing.Status == model.WhatsAppStatusConnected {
			return existing, errors.New("already connected")
		}
	} else {
		// Tenant session - check by tenant ID
		existing, err := d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
		if err == nil && existing != nil && existing.Status == model.WhatsAppStatusConnected {
			return existing, errors.New("already connected")
		}
		instanceName = "tenant_" + tenantID[:8]
	}

	// Generate unique token
	token := uuid.New().String()

	// Create instance in Evolution API
	session, err := d.evolutionPort.CreateInstance(ctx, instanceName, token)
	if err != nil {
		// Instance might already exist, try to get QR directly
		qrCode, qrErr := d.evolutionPort.ConnectInstance(ctx, instanceName, token)
		if qrErr != nil {
			return nil, err
		}
		// Return minimal session with QR
		return &model.WhatsAppSession{
			InstanceName: instanceName,
			TenantID:     tenantID,
			QRCode:       qrCode,
			Status:       model.WhatsAppStatusConnecting,
		}, nil
	}

	// Get QR code
	qrCode, err := d.evolutionPort.ConnectInstance(ctx, instanceName, token)
	if err != nil {
		return nil, err
	}

	session.TenantID = tenantID
	session.InstanceName = instanceName
	session.QRCode = qrCode
	session.Status = model.WhatsAppStatusConnecting

	// Save to database
	if err := d.dbPort.WhatsApp().Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetStatus checks current connection status from Evolution API
func (d *whatsAppDomain) GetStatus(ctx context.Context, tenantID string) (*model.WhatsAppSession, error) {
	var session *model.WhatsAppSession
	var err error

	if tenantID == "" {
		// Owner session - lookup by instance name
		session, err = d.dbPort.WhatsApp().GetByInstanceName(ctx, "eduvera_owner")
	} else {
		session, err = d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
	}

	if err != nil {
		return nil, err
	}
	if session == nil {
		return &model.WhatsAppSession{
			TenantID: tenantID,
			Status:   model.WhatsAppStatusDisconnected,
		}, nil
	}

	// Fetch live status from Evolution API
	liveStatus, err := d.evolutionPort.FetchInstance(ctx, session.InstanceName, session.APIKey)
	if err != nil {
		session.Status = model.WhatsAppStatusDisconnected
		return session, nil
	}

	session.Status = liveStatus.Status
	if liveStatus.PhoneNumber != "" {
		session.PhoneNumber = liveStatus.PhoneNumber
	}

	// Update status in database
	_ = d.dbPort.WhatsApp().UpdateStatus(ctx, session.ID, liveStatus.Status)

	return session, nil
}

// DisconnectTenant logs out and removes WhatsApp session
func (d *whatsAppDomain) DisconnectTenant(ctx context.Context, tenantID string) error {
	var session *model.WhatsAppSession
	var err error

	if tenantID == "" {
		// Owner session
		session, err = d.dbPort.WhatsApp().GetByInstanceName(ctx, "eduvera_owner")
	} else {
		session, err = d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
	}

	if err != nil || session == nil {
		return errors.New("no active session found")
	}

	// Logout from Evolution API
	if err := d.evolutionPort.LogoutInstance(ctx, session.InstanceName, session.APIKey); err != nil {
		// Log but continue to clean up database
	}

	// Delete instance
	_ = d.evolutionPort.DeleteInstance(ctx, session.InstanceName, session.APIKey)

	// Delete from database
	return d.dbPort.WhatsApp().Delete(ctx, session.ID)
}

// SendMessage sends WhatsApp message using tenant's or owner's connected number
func (d *whatsAppDomain) SendMessage(ctx context.Context, tenantID, phone, message string) error {
	var session *model.WhatsAppSession
	var err error

	if tenantID == "" {
		// Owner session - no rate limit for owner
		session, err = d.dbPort.WhatsApp().GetByInstanceName(ctx, "eduvera_owner")
	} else {
		// SECURITY: Check message rate limit for tenant (100 messages/day)
		if err := d.checkMessageRateLimit(ctx, tenantID); err != nil {
			return err
		}
		session, err = d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
	}

	if err != nil || session == nil {
		return errors.New("WhatsApp not connected")
	}

	if session.Status != model.WhatsAppStatusConnected {
		return errors.New("WhatsApp not connected, please reconnect")
	}

	// Send message
	sendErr := d.evolutionPort.SendMessage(ctx, session.InstanceName, session.APIKey, phone, message)
	if sendErr != nil {
		return sendErr
	}

	// Increment message counter on successful send (only for tenant)
	if tenantID != "" {
		d.incrementMessageCount(ctx, tenantID)
	}

	return nil
}

// checkMessageRateLimit checks if tenant has exceeded daily message limit
func (d *whatsAppDomain) checkMessageRateLimit(ctx context.Context, tenantID string) error {
	key := "wa_msg_count:" + tenantID
	maxMessages := int64(100) // 100 messages per day per tenant

	count, err := redis.GetInt(ctx, key)
	if err != nil {
		// Key doesn't exist yet, allow
		return nil
	}

	if count >= maxMessages {
		return errors.New("Batas pengiriman WA harian tercapai (100 pesan/hari). Upgrade ke Premium untuk batas lebih tinggi.")
	}

	return nil
}

// incrementMessageCount increments the daily message counter
func (d *whatsAppDomain) incrementMessageCount(ctx context.Context, tenantID string) {
	key := "wa_msg_count:" + tenantID

	count, _ := redis.Incr(ctx, key)
	if count == 1 {
		// Set TTL to end of day (24 hours)
		_ = redis.Expire(ctx, key, 24*time.Hour)
	}
}
