package inbound_port

type OnboardingHttpPort interface {
	Register(a any) error
	Institution(a any) error
	CheckSubdomain(a any) error // New: realtime subdomain validation
	Subdomain(a any) error
	BankAccount(a any) error
	Confirm(a any) error
	Status(a any) error
}
