package inbound_port

type HttpPort interface {
	Middleware() MiddlewareHttpPort
	Ping() PingHttpPort
	Client() ClientHttpPort
	Landing() LandingHttpPort
	Onboarding() OnboardingHttpPort
	Auth() AuthHttpPort
	Payment() PaymentHttpPort
	Owner() OwnerHttpPort
	Content() ContentHttpPort
}
