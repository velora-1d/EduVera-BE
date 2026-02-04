package whatsapp

import (
	"context"
	"errors"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

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
	evolutionPort outbound_port.EvolutionApiPort
}

func NewWhatsAppDomain(dbPort outbound_port.DatabasePort, evolutionPort outbound_port.EvolutionApiPort) WhatsAppDomain {
	return &whatsAppDomain{
		dbPort:        dbPort,
		evolutionPort: evolutionPort,
	}
}

// ConnectTenant creates Evolution API instance and returns QR code for scanning
func (d *whatsAppDomain) ConnectTenant(ctx context.Context, tenantID string) (*model.WhatsAppSession, error) {
	// Check if session already exists
	existing, err := d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
	if err == nil && existing != nil && existing.Status == model.WhatsAppStatusConnected {
		return existing, errors.New("already connected")
	}

	// Generate unique instance name and token
	instanceName := "tenant_" + tenantID[:8]
	token := uuid.New().String()

	// Create instance in Evolution API
	session, err := d.evolutionPort.CreateInstance(ctx, instanceName, token)
	if err != nil {
		return nil, err
	}

	// Get QR code
	qrCode, err := d.evolutionPort.ConnectInstance(ctx, instanceName, token)
	if err != nil {
		return nil, err
	}

	session.TenantID = tenantID
	session.QRCode = qrCode

	// Save to database
	if err := d.dbPort.WhatsApp().Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetStatus checks current connection status from Evolution API
func (d *whatsAppDomain) GetStatus(ctx context.Context, tenantID string) (*model.WhatsAppSession, error) {
	session, err := d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
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

	// Update status in database
	_ = d.dbPort.WhatsApp().UpdateStatus(ctx, session.ID, liveStatus.Status)

	return session, nil
}

// DisconnectTenant logs out and removes WhatsApp session
func (d *whatsAppDomain) DisconnectTenant(ctx context.Context, tenantID string) error {
	session, err := d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
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

// SendMessage sends WhatsApp message using tenant's connected number
func (d *whatsAppDomain) SendMessage(ctx context.Context, tenantID, phone, message string) error {
	session, err := d.dbPort.WhatsApp().GetByTenantID(ctx, tenantID)
	if err != nil || session == nil {
		return errors.New("tenant WhatsApp not connected")
	}

	if session.Status != model.WhatsAppStatusConnected {
		return errors.New("WhatsApp not connected, please reconnect")
	}

	return d.evolutionPort.SendMessage(ctx, session.InstanceName, session.APIKey, phone, message)
}
