package outbound_port

import "database/sql"

//go:generate mockgen -source=registry_database.go -destination=./../../../tests/mocks/port/mock_registry_database.go
type InTransaction func(repoRegistry DatabasePort) (interface{}, error)

type DatabasePort interface {
	Client() ClientDatabasePort
	Tenant() TenantDatabasePort
	User() UserDatabasePort
	Content() ContentDatabasePort
	Payment() PaymentDatabasePort
	Disbursement() DisbursementDatabasePort
	SPP() SPPDatabasePort
	Notification() NotificationDatabasePort
	Sekolah() SekolahPort
	AuditLog() AuditLogDatabasePort
	ERapor() ERaporDatabasePort
	SDM() SDMDatabasePort
	Subscription() SubscriptionDatabasePort
	PesantrenDashboard() PesantrenDashboardPort
	WhatsApp() WhatsAppDatabasePort
	Student() StudentDatabasePort
	NotificationTemplate() NotificationTemplateDatabasePort
	DoInTransaction(txFunc InTransaction) (out interface{}, err error)
}

type DatabaseExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}
