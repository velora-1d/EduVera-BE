package postgres_outbound_adapter

import (
	"context"
	"database/sql"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type sppAdapter struct {
	db *sql.DB
}

func NewSPPAdapter(db *sql.DB) outbound_port.SPPDatabasePort {
	return &sppAdapter{
		db: db,
	}
}

func (a *sppAdapter) Create(ctx context.Context, spp *model.SPPTransaction) error {
	query := `
		INSERT INTO spp_transactions (tenant_id, student_id, student_name, amount, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return a.db.QueryRowContext(ctx, query,
		spp.TenantID, spp.StudentID, spp.StudentName, spp.Amount, spp.Status, spp.Description,
	).Scan(&spp.ID, &spp.CreatedAt, &spp.UpdatedAt)
}

func (a *sppAdapter) ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, status, gateway_ref, description, created_at, updated_at
		FROM spp_transactions
		WHERE tenant_id = $1
		ORDER BY created_at DESC
	`
	rows, err := a.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.SPPTransaction
	for rows.Next() {
		var t model.SPPTransaction
		var studentID, paymentMethod, gatewayRef, description sql.NullString
		if err := rows.Scan(
			&t.ID, &t.TenantID, &studentID, &t.StudentName, &t.Amount,
			&paymentMethod, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.StudentID = studentID.String
		t.PaymentMethod = paymentMethod.String
		t.GatewayRef = gatewayRef.String
		t.Description = description.String
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (a *sppAdapter) ListAll(ctx context.Context) ([]model.SPPTransaction, error) {
	query := `
		SELECT s.id, s.tenant_id, t.name as tenant_name, s.student_id, s.student_name, s.amount, s.payment_method, s.status, s.gateway_ref, s.description, s.created_at, s.updated_at
		FROM spp_transactions s
		JOIN tenants t ON s.tenant_id = t.id
		ORDER BY s.created_at DESC
	`
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.SPPTransaction
	for rows.Next() {
		var t model.SPPTransaction
		var tenantName string
		var studentID, paymentMethod, gatewayRef, description sql.NullString
		if err := rows.Scan(
			&t.ID, &t.TenantID, &tenantName, &studentID, &t.StudentName, &t.Amount,
			&paymentMethod, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.StudentID = studentID.String
		t.PaymentMethod = paymentMethod.String
		t.GatewayRef = gatewayRef.String
		t.Description = description.String
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (a *sppAdapter) FindByID(ctx context.Context, id string) (*model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, status, gateway_ref, description, created_at, updated_at
		FROM spp_transactions
		WHERE id = $1
	`
	var t model.SPPTransaction
	var studentID, paymentMethod, gatewayRef, description sql.NullString
	err := a.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.TenantID, &studentID, &t.StudentName, &t.Amount,
		&paymentMethod, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.StudentID = studentID.String
	t.PaymentMethod = paymentMethod.String
	t.GatewayRef = gatewayRef.String
	t.Description = description.String
	return &t, nil
}

func (a *sppAdapter) UpdateStatus(ctx context.Context, id string, status model.SPPStatus, paymentMethod string) error {
	query := `
		UPDATE spp_transactions 
		SET status = $1, payment_method = $2, updated_at = NOW() 
		WHERE id = $3
	`
	_, err := a.db.ExecContext(ctx, query, status, paymentMethod, id)
	return err
}

func (a *sppAdapter) GetStatsByTenant(ctx context.Context, tenantID string) (*model.SPPStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(amount), 0) as total_billed,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN amount ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN status = 'pending' THEN amount ELSE 0 END), 0) as total_pending,
			COUNT(*) as transaction_count,
			COUNT(CASE WHEN status = 'paid' THEN 1 END) as paid_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count
		FROM spp_transactions
		WHERE tenant_id = $1
	`
	var stats model.SPPStats
	err := a.db.QueryRowContext(ctx, query, tenantID).Scan(
		&stats.TotalBilled, &stats.TotalPaid, &stats.TotalPending,
		&stats.TransactionCount, &stats.PaidCount, &stats.PendingCount,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
