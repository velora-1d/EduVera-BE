package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

const tablePayment = "payments"

type paymentAdapter struct {
	db outbound_port.DatabaseExecutor
}

func NewPaymentAdapter(
	db outbound_port.DatabaseExecutor,
) outbound_port.PaymentDatabasePort {
	return &paymentAdapter{
		db: db,
	}
}

func (a *paymentAdapter) Create(payment *model.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tablePayment).Rows(goqu.Record{
		"tenant_id":    payment.TenantID,
		"order_id":     payment.OrderID,
		"amount":       payment.Amount,
		"status":       payment.Status,
		"payment_type": payment.PaymentType,
		"snap_token":   payment.SnapToken,
		"midtrans_id":  payment.MidtransID,
		"created_at":   payment.CreatedAt,
		"updated_at":   payment.UpdatedAt,
	}).Returning("id")

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	return a.db.QueryRow(query).Scan(&payment.ID)
}

func (a *paymentAdapter) Update(payment *model.Payment) error {
	payment.UpdatedAt = time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tablePayment).
		Set(goqu.Record{
			"status":       payment.Status,
			"payment_type": payment.PaymentType,
			"snap_token":   payment.SnapToken,
			"midtrans_id":  payment.MidtransID,
			"paid_at":      payment.PaidAt,
			"updated_at":   payment.UpdatedAt,
		}).
		Where(goqu.Ex{"id": payment.ID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *paymentAdapter) FindByFilter(filter model.PaymentFilter) ([]model.Payment, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tablePayment)
	dataset = addPaymentFilter(dataset, filter)

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.OrderID, &p.Amount,
			&p.Status, &p.PaymentType, &p.SnapToken,
			&p.MidtransID, &p.PaidAt,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	return payments, nil
}

func (a *paymentAdapter) FindByID(id string) (*model.Payment, error) {
	payments, err := a.FindByFilter(model.PaymentFilter{IDs: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(payments) == 0 {
		return nil, sql.ErrNoRows
	}
	return &payments[0], nil
}

func (a *paymentAdapter) FindByOrderID(orderID string) (*model.Payment, error) {
	payments, err := a.FindByFilter(model.PaymentFilter{OrderIDs: []string{orderID}})
	if err != nil {
		return nil, err
	}
	if len(payments) == 0 {
		return nil, sql.ErrNoRows
	}
	return &payments[0], nil
}

func (a *paymentAdapter) UpdateStatus(orderID string, status string, paymentType string, midtransID string) error {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tablePayment).
		Set(goqu.Record{
			"status":       status,
			"payment_type": paymentType,
			"midtrans_id":  midtransID,
			"updated_at":   time.Now(),
		}).
		Where(goqu.Ex{"order_id": orderID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *paymentAdapter) MarkAsPaid(orderID string, paymentType string, midtransID string) error {
	now := time.Now()
	dialect := goqu.Dialect("postgres")
	dataset := dialect.Update(tablePayment).
		Set(goqu.Record{
			"status":       model.PaymentStatusPaid,
			"payment_type": paymentType,
			"midtrans_id":  midtransID,
			"paid_at":      now,
			"updated_at":   now,
		}).
		Where(goqu.Ex{"order_id": orderID})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func addPaymentFilter(dataset *goqu.SelectDataset, filter model.PaymentFilter) *goqu.SelectDataset {
	if len(filter.IDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"id": filter.IDs})
	}
	if len(filter.TenantIDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"tenant_id": filter.TenantIDs})
	}
	if len(filter.OrderIDs) > 0 {
		dataset = dataset.Where(goqu.Ex{"order_id": filter.OrderIDs})
	}
	if len(filter.Statuses) > 0 {
		dataset = dataset.Where(goqu.Ex{"status": filter.Statuses})
	}
	return dataset
}
