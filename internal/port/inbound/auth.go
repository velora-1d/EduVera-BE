package inbound_port

type AuthHttpPort interface {
	Login(a any) error
	Me(a any) error
	Refresh(a any) error
	Logout(a any) error
}
