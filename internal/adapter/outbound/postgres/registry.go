package postgres_outbound_adapter

import (
	"database/sql"

	"github.com/pkg/errors"

	outbound_port "eduvera/internal/port/outbound"
)

type adapter struct {
	db         *sql.DB
	dbexecutor outbound_port.DatabaseExecutor
}

func NewAdapter(db *sql.DB) outbound_port.DatabasePort {
	return &adapter{
		db: db,
	}
}

func (s *adapter) DoInTransaction(txFunc outbound_port.InTransaction) (out interface{}, err error) {
	var tx *sql.Tx
	reg := s
	if s.dbexecutor == nil {
		tx, err = s.db.Begin()
		if err != nil {
			return
		}
		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				switch x := p.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					// Fallback err (per specs, error strings should be lowercase w/o punctuation
					err = errors.New("unknown panic")
				}
			} else if err != nil {
				xerr := tx.Rollback() // err is non-nil; don't change it
				if xerr != nil {
					err = errors.Wrap(err, xerr.Error())
				}
			} else {
				err = tx.Commit() // err is nil; if Commit returns error update err
			}
		}()
		reg = &adapter{
			db:         s.db,
			dbexecutor: tx,
		}
	}
	out, err = txFunc(reg)
	if err != nil {
		if out != nil {
			return out, err
		}

		return nil, err
	}
	return
}

func (s *adapter) Client() outbound_port.ClientDatabasePort {
	if s.dbexecutor != nil {
		return NewClientAdapter(s.dbexecutor)
	}
	return NewClientAdapter(s.db)
}

func (s *adapter) Tenant() outbound_port.TenantDatabasePort {
	if s.dbexecutor != nil {
		return NewTenantAdapter(s.dbexecutor)
	}
	return NewTenantAdapter(s.db)
}

func (s *adapter) User() outbound_port.UserDatabasePort {
	if s.dbexecutor != nil {
		return NewUserAdapter(s.dbexecutor)
	}
	return NewUserAdapter(s.db)
}

func (s *adapter) Content() outbound_port.ContentDatabasePort {
	if s.dbexecutor != nil {
		return NewContentAdapter(s.dbexecutor)
	}
	return NewContentAdapter(s.db)
}

func (s *adapter) Payment() outbound_port.PaymentDatabasePort {
	if s.dbexecutor != nil {
		return NewPaymentAdapter(s.dbexecutor)
	}
	return NewPaymentAdapter(s.db)
}

func (s *adapter) Disbursement() outbound_port.DisbursementDatabasePort {
	// Currently disbursement adapter only supports non-transactional db
	// TODO: Update disbursement adapter to support DatabaseExecutor for transactions
	return NewDisbursementAdapter(s.db)
}

func (s *adapter) SPP() outbound_port.SPPDatabasePort {
	return NewSPPAdapter(s.db)
}

func (s *adapter) Notification() outbound_port.NotificationDatabasePort {
	return NewNotificationAdapter(s.db)
}

func (s *adapter) Sekolah() outbound_port.SekolahPort {
	if s.dbexecutor != nil {
		return NewSekolahAdapter(s.dbexecutor)
	}
	return NewSekolahAdapter(s.db)
}

func (s *adapter) AuditLog() outbound_port.AuditLogDatabasePort {
	return NewAuditLogAdapter(s.db)
}

func (s *adapter) ERapor() outbound_port.ERaporDatabasePort {
	return NewERaporAdapter(s.db)
}
