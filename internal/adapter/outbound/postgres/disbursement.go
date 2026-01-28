package postgres_outbound_adapter

import (
	"context"
	"database/sql"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
)

type disbursementAdapter struct {
	db *sql.DB
}

func NewDisbursementAdapter(db *sql.DB) outbound_port.DisbursementDatabasePort {
	return &disbursementAdapter{
		db: db,
	}
}

func (r *disbursementAdapter) GetAll(ctx context.Context) ([]model.Disbursement, error) {
	query := `
		SELECT 
			d.id, d.tenant_id, t.name as tenant_name, d.amount, d.bank_name, 
			d.account_number, d.account_holder, d.status, d.requested_at, d.processed_at
		FROM disbursements d
		JOIN tenants t ON d.tenant_id = t.id
		ORDER BY d.requested_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disbursements []model.Disbursement
	for rows.Next() {
		var d model.Disbursement
		var processedAt sql.NullTime
		if err := rows.Scan(
			&d.ID, &d.TenantID, &d.TenantName, &d.Amount, &d.BankName,
			&d.AccountNumber, &d.AccountHolder, &d.Status, &d.RequestedAt, &processedAt,
		); err != nil {
			return nil, err
		}
		if processedAt.Valid {
			processedAtTime := processedAt.Time
			d.ProcessedAt = &processedAtTime
		}
		disbursements = append(disbursements, d)
	}

	return disbursements, nil
}

func (r *disbursementAdapter) Approve(ctx context.Context, id string) error {
	query := `
		UPDATE disbursements 
		SET status = 'completed', processed_at = NOW() 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *disbursementAdapter) Reject(ctx context.Context, id string, reason string) error {
	query := `
		UPDATE disbursements 
		SET status = 'rejected', admin_notes = $1, processed_at = NOW() 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, reason, id)
	return err
}
