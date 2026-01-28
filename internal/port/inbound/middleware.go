package inbound_port

type MiddlewareHttpPort interface {
	InternalAuth(a any) error
	ClientAuth(a any) error
	OwnerAuth(a any) error
}
