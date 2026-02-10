package postgres_outbound_adapter

import (
	"context"
	"database/sql"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
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
		INSERT INTO spp_transactions (tenant_id, student_id, student_name, amount, status, description, payment_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return a.db.QueryRowContext(ctx, query,
		spp.TenantID, spp.StudentID, spp.StudentName, spp.Amount, spp.Status, spp.Description, spp.PaymentType,
	).Scan(&spp.ID, &spp.CreatedAt, &spp.UpdatedAt)
}

func (a *sppAdapter) ListByTenant(ctx context.Context, tenantID string) ([]model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, payment_type, status, gateway_ref, description, created_at, updated_at
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
		var studentID, paymentMethod, paymentType, gatewayRef, description sql.NullString
		if err := rows.Scan(
			&t.ID, &t.TenantID, &studentID, &t.StudentName, &t.Amount,
			&paymentMethod, &paymentType, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.StudentID = studentID.String
		t.PaymentMethod = paymentMethod.String
		t.PaymentType = model.PaymentType(paymentType.String)
		t.GatewayRef = gatewayRef.String
		t.Description = description.String
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (a *sppAdapter) ListAll(ctx context.Context) ([]model.SPPTransaction, error) {
	query := `
		SELECT s.id, s.tenant_id, t.name as tenant_name, s.student_id, s.student_name, s.amount, s.payment_method, s.payment_type, s.status, s.gateway_ref, s.description, s.created_at, s.updated_at
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
		var studentID, paymentMethod, paymentType, gatewayRef, description sql.NullString
		if err := rows.Scan(
			&t.ID, &t.TenantID, &tenantName, &studentID, &t.StudentName, &t.Amount,
			&paymentMethod, &paymentType, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.StudentID = studentID.String
		t.PaymentMethod = paymentMethod.String
		t.PaymentType = model.PaymentType(paymentType.String)
		t.GatewayRef = gatewayRef.String
		t.Description = description.String
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (a *sppAdapter) FindByID(ctx context.Context, tenantID, id string) (*model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, payment_type, status, gateway_ref, description, created_at, updated_at
		FROM spp_transactions
		WHERE id = $1 AND tenant_id = $2
	`
	var t model.SPPTransaction
	var studentID, paymentMethod, paymentType, gatewayRef, description sql.NullString
	err := a.db.QueryRowContext(ctx, query, id, tenantID).Scan(
		&t.ID, &t.TenantID, &studentID, &t.StudentName, &t.Amount,
		&paymentMethod, &paymentType, &t.Status, &gatewayRef, &description, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	t.StudentID = studentID.String
	t.PaymentMethod = paymentMethod.String
	t.PaymentType = model.PaymentType(paymentType.String)
	t.GatewayRef = gatewayRef.String
	t.Description = description.String
	return &t, nil
}

func (a *sppAdapter) UpdateStatus(ctx context.Context, tenantID, id string, status model.SPPStatus, paymentMethod string) error {
	query := `
		UPDATE spp_transactions 
		SET status = $1, payment_method = $2, updated_at = NOW() 
		WHERE id = $3 AND tenant_id = $4
	`
	_, err := a.db.ExecContext(ctx, query, status, paymentMethod, id, tenantID)
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

// Update updates an SPP transaction
func (a *sppAdapter) Update(ctx context.Context, spp *model.SPPTransaction) error {
	query := `
		UPDATE spp_transactions 
		SET student_name = $1, amount = $2, description = $3, due_date = $4, period = $5, payment_type = $6, updated_at = NOW()
		WHERE id = $7
	`
	_, err := a.db.ExecContext(ctx, query, spp.StudentName, spp.Amount, spp.Description, spp.DueDate, spp.Period, spp.PaymentType, spp.ID)
	return err
}

// Delete removes an SPP transaction
func (a *sppAdapter) Delete(ctx context.Context, tenantID, id string) error {
	query := `DELETE FROM spp_transactions WHERE id = $1 AND tenant_id = $2`
	_, err := a.db.ExecContext(ctx, query, id, tenantID)
	return err
}

// UploadProof saves the payment proof URL
func (a *sppAdapter) UploadProof(ctx context.Context, tenantID, id string, proofURL string) error {
	query := `
		UPDATE spp_transactions 
		SET payment_proof = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`
	_, err := a.db.ExecContext(ctx, query, proofURL, id, tenantID)
	return err
}

// ConfirmPayment marks payment as paid by admin
func (a *sppAdapter) ConfirmPayment(ctx context.Context, tenantID, id string, confirmedBy string) error {
	query := `
		UPDATE spp_transactions 
		SET status = 'paid', confirmed_by = $1, paid_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
	`
	_, err := a.db.ExecContext(ctx, query, confirmedBy, id, tenantID)
	return err
}

// ListByPeriod returns transactions for a specific period (e.g., "2024-01")
func (a *sppAdapter) ListByPeriod(ctx context.Context, tenantID string, period string) ([]model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, payment_type, status, 
		       gateway_ref, description, payment_proof, confirmed_by, paid_at, due_date, period, created_at, updated_at
		FROM spp_transactions
		WHERE tenant_id = $1 AND period = $2
		ORDER BY student_name ASC
	`
	rows, err := a.db.QueryContext(ctx, query, tenantID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSPPTransactions(rows)
}

// ListOverdue returns pending transactions past due date
func (a *sppAdapter) ListOverdue(ctx context.Context, tenantID string) ([]model.SPPTransaction, error) {
	query := `
		SELECT id, tenant_id, student_id, student_name, amount, payment_method, payment_type, status, 
		       gateway_ref, description, payment_proof, confirmed_by, paid_at, due_date, period, created_at, updated_at
		FROM spp_transactions
		WHERE tenant_id = $1 AND status = 'pending' AND due_date < NOW()
		ORDER BY due_date ASC
	`
	rows, err := a.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSPPTransactions(rows)
}

// Helper function to scan SPP transactions with all fields
func scanSPPTransactions(rows *sql.Rows) ([]model.SPPTransaction, error) {
	var transactions []model.SPPTransaction
	for rows.Next() {
		var t model.SPPTransaction
		var studentID, paymentMethod, paymentType, gatewayRef, description, paymentProof, confirmedBy, period sql.NullString
		var paidAt, dueDate sql.NullTime
		if err := rows.Scan(
			&t.ID, &t.TenantID, &studentID, &t.StudentName, &t.Amount,
			&paymentMethod, &paymentType, &t.Status, &gatewayRef, &description,
			&paymentProof, &confirmedBy, &paidAt, &dueDate, &period,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		t.StudentID = studentID.String
		t.PaymentMethod = paymentMethod.String
		t.PaymentType = model.PaymentType(paymentType.String)
		t.GatewayRef = gatewayRef.String
		t.Description = description.String
		t.PaymentProof = paymentProof.String
		t.ConfirmedBy = confirmedBy.String
		t.Period = period.String
		if paidAt.Valid {
			t.PaidAt = &paidAt.Time
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}
