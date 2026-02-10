package whatsmeow_adapter

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// Manager handles multi-tenant WhatsMeow sessions.
// One whatsmeow.Client per tenant, keyed by instanceName.
type Manager struct {
	mu      sync.RWMutex
	clients map[string]*clientEntry // instanceName -> clientEntry
	baseDir string                  // base directory for session SQLite files
	qrCodes map[string]string       // instanceName -> latest base64 QR
	qrMu    sync.RWMutex
}

type clientEntry struct {
	client       *whatsmeow.Client
	store        *sqlstore.Container
	instanceName string
	tenantID     string
	connected    bool
}

// Compile-time check: Manager implements WhatsAppClientPort
var _ outbound_port.WhatsAppClientPort = (*Manager)(nil)

// NewManager creates a new WhatsMeow multi-session manager.
// baseDir is the root directory for session storage (e.g., "./wa_sessions").
func NewManager() *Manager {
	baseDir := os.Getenv("WA_SESSION_DIR")
	if baseDir == "" {
		baseDir = "./wa_sessions"
	}
	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Printf("[WHATSMEOW] Failed to create session dir %s: %v", baseDir, err)
	}

	return &Manager{
		clients: make(map[string]*clientEntry),
		baseDir: baseDir,
		qrCodes: make(map[string]string),
	}
}

// CreateInstance creates a new WhatsMeow client for the given instance.
// The token parameter is stored but not used by WhatsMeow directly (kept for interface compat).
func (m *Manager) CreateInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If client already exists, return existing
	if entry, ok := m.clients[instanceName]; ok {
		status := model.WhatsAppStatusDisconnected
		if entry.client.IsConnected() {
			status = model.WhatsAppStatusConnected
		}
		return &model.WhatsAppSession{
			InstanceName: instanceName,
			APIKey:       token,
			Status:       status,
		}, nil
	}

	// Create session directory for this instance
	sessionDir := filepath.Join(m.baseDir, instanceName)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat session dir: %w", err)
	}

	dbPath := filepath.Join(sessionDir, "session.db")
	container, err := sqlstore.New(context.Background(), "sqlite3",
		fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbPath),
		waLog.Noop,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat session store: %w", err)
	}

	// Get or create device store
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan device store: %w", err)
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)

	entry := &clientEntry{
		client:       client,
		store:        container,
		instanceName: instanceName,
		connected:    false,
	}
	m.clients[instanceName] = entry

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		APIKey:       token,
		Status:       model.WhatsAppStatusConnecting,
	}, nil
}

// FetchInstance returns the current connection state of an instance.
func (m *Manager) FetchInstance(ctx context.Context, instanceName string, token string) (*model.WhatsAppSession, error) {
	m.mu.RLock()
	entry, ok := m.clients[instanceName]
	m.mu.RUnlock()

	if !ok {
		// Try to load from disk
		loaded, err := m.loadFromDisk(instanceName)
		if err != nil || loaded == nil {
			return nil, fmt.Errorf("instance %s tidak ditemukan", instanceName)
		}
		entry = loaded
	}

	status := model.WhatsAppStatusDisconnected
	phoneNumber := ""

	if entry.client.IsConnected() {
		status = model.WhatsAppStatusConnected
		if entry.client.Store.ID != nil {
			phoneNumber = entry.client.Store.ID.User
		}
	} else if entry.client.IsLoggedIn() {
		// Has session but not connected - try reconnect
		status = model.WhatsAppStatusConnecting
	}

	return &model.WhatsAppSession{
		InstanceName: instanceName,
		Status:       status,
		PhoneNumber:  phoneNumber,
	}, nil
}

// ConnectInstance connects the instance and returns a Base64 QR code for pairing.
// If already logged in, reconnects without QR.
func (m *Manager) ConnectInstance(ctx context.Context, instanceName string, token string) (string, error) {
	m.mu.RLock()
	entry, ok := m.clients[instanceName]
	m.mu.RUnlock()

	if !ok {
		// Auto-create if not exists
		_, err := m.CreateInstance(ctx, instanceName, token)
		if err != nil {
			return "", err
		}
		m.mu.RLock()
		entry = m.clients[instanceName]
		m.mu.RUnlock()
	}

	// If already connected, no QR needed
	if entry.client.IsConnected() {
		return "", nil
	}

	// If already logged in (has session), just reconnect
	if entry.client.Store.ID != nil {
		err := entry.client.Connect()
		if err != nil {
			return "", fmt.Errorf("gagal reconnect: %w", err)
		}
		m.mu.Lock()
		entry.connected = true
		m.mu.Unlock()
		return "", nil
	}

	// Need QR code for new pairing
	qrChan, _ := entry.client.GetQRChannel(ctx)
	err := entry.client.Connect()
	if err != nil {
		return "", fmt.Errorf("gagal connect untuk QR: %w", err)
	}

	// Register event handler for this instance
	entry.client.AddEventHandler(m.makeEventHandler(instanceName))

	// Wait for QR code with timeout
	select {
	case evt := <-qrChan:
		if evt.Event == "code" {
			// Generate QR code as Base64 PNG
			qrPNG, err := qrcode.Encode(evt.Code, qrcode.Medium, 512)
			if err != nil {
				return "", fmt.Errorf("gagal generate QR image: %w", err)
			}
			base64QR := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrPNG)

			// Cache QR code
			m.qrMu.Lock()
			m.qrCodes[instanceName] = base64QR
			m.qrMu.Unlock()

			return base64QR, nil
		} else if evt.Event == "success" {
			m.mu.Lock()
			entry.connected = true
			m.mu.Unlock()
			return "", nil // Already paired
		} else if evt.Event == "timeout" {
			return "", fmt.Errorf("QR code timeout, silakan coba lagi")
		}
	case <-time.After(30 * time.Second):
		return "", fmt.Errorf("timeout menunggu QR code")
	}

	return "", fmt.Errorf("gagal mendapatkan QR code")
}

// LogoutInstance logs out and disconnects the WhatsApp session.
func (m *Manager) LogoutInstance(ctx context.Context, instanceName string, token string) error {
	m.mu.Lock()
	entry, ok := m.clients[instanceName]
	m.mu.Unlock()

	if !ok {
		return fmt.Errorf("instance %s tidak ditemukan", instanceName)
	}

	// Logout from WhatsApp
	err := entry.client.Logout(ctx)
	if err != nil {
		log.Printf("[WHATSMEOW] Logout error for %s: %v", instanceName, err)
	}

	entry.client.Disconnect()

	m.mu.Lock()
	entry.connected = false
	m.mu.Unlock()

	return nil
}

// DeleteInstance removes the instance completely (logout + delete session files).
func (m *Manager) DeleteInstance(ctx context.Context, instanceName string, token string) error {
	// Logout first
	_ = m.LogoutInstance(ctx, instanceName, token)

	m.mu.Lock()
	if entry, ok := m.clients[instanceName]; ok {
		entry.client.Disconnect()
		delete(m.clients, instanceName)
	}
	m.mu.Unlock()

	// Remove QR cache
	m.qrMu.Lock()
	delete(m.qrCodes, instanceName)
	m.qrMu.Unlock()

	// Delete session directory
	sessionDir := filepath.Join(m.baseDir, instanceName)
	if err := os.RemoveAll(sessionDir); err != nil {
		return fmt.Errorf("gagal hapus session dir: %w", err)
	}

	return nil
}

// SendMessage sends a text message via the specified instance.
func (m *Manager) SendMessage(ctx context.Context, instanceName string, token string, phone string, message string) error {
	m.mu.RLock()
	entry, ok := m.clients[instanceName]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("instance %s tidak ditemukan", instanceName)
	}

	if !entry.client.IsConnected() {
		return fmt.Errorf("instance %s tidak terhubung", instanceName)
	}

	// Normalize phone number to JID
	jid := phoneToJID(phone)

	_, err := entry.client.SendMessage(ctx, jid, &waE2E.Message{
		Conversation: proto.String(message),
	})
	if err != nil {
		return fmt.Errorf("gagal kirim pesan: %w", err)
	}

	return nil
}

// Shutdown gracefully disconnects all clients.
func (m *Manager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, entry := range m.clients {
		log.Printf("[WHATSMEOW] Disconnecting %s...", name)
		entry.client.Disconnect()
	}
}

// ReconnectAll loads all existing sessions from disk and reconnects them.
// Call this on application startup.
func (m *Manager) ReconnectAll() {
	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		log.Printf("[WHATSMEOW] No sessions to reconnect: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		instanceName := entry.Name()
		dbPath := filepath.Join(m.baseDir, instanceName, "session.db")
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			continue
		}

		log.Printf("[WHATSMEOW] Loading session: %s", instanceName)
		loaded, err := m.loadFromDisk(instanceName)
		if err != nil {
			log.Printf("[WHATSMEOW] Failed to load %s: %v", instanceName, err)
			continue
		}

		if loaded.client.Store.ID != nil {
			log.Printf("[WHATSMEOW] Reconnecting %s...", instanceName)
			loaded.client.AddEventHandler(m.makeEventHandler(instanceName))
			err := loaded.client.Connect()
			if err != nil {
				log.Printf("[WHATSMEOW] Failed to reconnect %s: %v", instanceName, err)
			} else {
				m.mu.Lock()
				loaded.connected = true
				m.mu.Unlock()
				log.Printf("[WHATSMEOW] Reconnected %s", instanceName)
			}
		}
	}
}

// loadFromDisk loads a session from SQLite and creates a client.
func (m *Manager) loadFromDisk(instanceName string) (*clientEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already loaded
	if entry, ok := m.clients[instanceName]; ok {
		return entry, nil
	}

	sessionDir := filepath.Join(m.baseDir, instanceName)
	dbPath := filepath.Join(sessionDir, "session.db")

	container, err := sqlstore.New(context.Background(), "sqlite3",
		fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbPath),
		waLog.Noop,
	)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, waLog.Noop)

	entry := &clientEntry{
		client:       client,
		store:        container,
		instanceName: instanceName,
	}
	m.clients[instanceName] = entry

	return entry, nil
}

// makeEventHandler creates event handler for a specific instance.
func (m *Manager) makeEventHandler(instanceName string) whatsmeow.EventHandler {
	return func(evt interface{}) {
		// Handle connection events
		switch evt.(type) {
		case *events.Connected:
			log.Printf("[WHATSMEOW] %s: Connected", instanceName)
			m.mu.Lock()
			if entry, ok := m.clients[instanceName]; ok {
				entry.connected = true
			}
			m.mu.Unlock()
		case *events.Disconnected:
			log.Printf("[WHATSMEOW] %s: Disconnected", instanceName)
			m.mu.Lock()
			if entry, ok := m.clients[instanceName]; ok {
				entry.connected = false
			}
			m.mu.Unlock()
		}
	}
}

// phoneToJID converts a phone number string to a WhatsApp JID.
// Handles: +628xxx, 628xxx, 08xxx formats.
func phoneToJID(phone string) types.JID {
	phone = strings.TrimSpace(phone)
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.TrimSuffix(phone, "@s.whatsapp.net")

	// Convert 08xxx to 628xxx
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}

	return types.NewJID(phone, types.DefaultUserServer)
}
