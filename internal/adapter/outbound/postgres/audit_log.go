package postgres_outbound_adapter

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

const tableAuditLog = "admin_audit_logs"

type auditLogAdapter struct {
	db *sql.DB
}

func NewAuditLogAdapter(db *sql.DB) outbound_port.AuditLogDatabasePort {
	return &auditLogAdapter{db: db}
}

func (a *auditLogAdapter) Create(log *model.AuditLog) error {
	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()

	dialect := goqu.Dialect("postgres")
	dataset := dialect.Insert(tableAuditLog).Rows(goqu.Record{
		"id":          log.ID,
		"admin_id":    log.AdminID,
		"admin_email": log.AdminEmail,
		"action":      log.Action,
		"target_type": log.TargetType,
		"target_id":   log.TargetID,
		"old_value":   log.OldValue,
		"new_value":   log.NewValue,
		"ip_address":  log.IPAddress,
		"user_agent":  log.UserAgent,
		"description": log.Description,
		"created_at":  log.CreatedAt,
	})

	query, _, err := dataset.ToSQL()
	if err != nil {
		return err
	}

	_, err = a.db.Exec(query)
	return err
}

func (a *auditLogAdapter) FindByFilter(filter model.AuditLogFilter) ([]model.AuditLog, error) {
	dialect := goqu.Dialect("postgres")
	dataset := dialect.From(tableAuditLog).Order(goqu.I("created_at").Desc())

	if filter.AdminID != "" {
		dataset = dataset.Where(goqu.Ex{"admin_id": filter.AdminID})
	}
	if filter.Action != "" {
		dataset = dataset.Where(goqu.Ex{"action": filter.Action})
	}
	if filter.TargetType != "" {
		dataset = dataset.Where(goqu.Ex{"target_type": filter.TargetType})
	}
	if filter.TargetID != "" {
		dataset = dataset.Where(goqu.Ex{"target_id": filter.TargetID})
	}
	if filter.StartDate != nil {
		dataset = dataset.Where(goqu.I("created_at").Gte(*filter.StartDate))
	}
	if filter.EndDate != nil {
		dataset = dataset.Where(goqu.I("created_at").Lte(*filter.EndDate))
	}

	limit := filter.Limit
	if limit == 0 {
		limit = 100
	}
	dataset = dataset.Limit(uint(limit)).Offset(uint(filter.Offset))

	query, _, err := dataset.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var log model.AuditLog
		err := rows.Scan(
			&log.ID, &log.AdminID, &log.AdminEmail, &log.Action,
			&log.TargetType, &log.TargetID, &log.OldValue, &log.NewValue,
			&log.IPAddress, &log.UserAgent, &log.Description, &log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (a *auditLogAdapter) GetStats() (map[string]interface{}, error) {
	dialect := goqu.Dialect("postgres")

	// Count total logs
	countQuery, _, _ := dialect.From(tableAuditLog).Select(goqu.COUNT("*").As("total")).ToSQL()
	var total int
	_ = a.db.QueryRow(countQuery).Scan(&total)

	// Count by action type
	stats := map[string]interface{}{
		"total_logs": total,
	}

	return stats, nil
}
