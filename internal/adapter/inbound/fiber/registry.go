package fiber_inbound_adapter

import (
	"eduvera/internal/adapter/inbound/fiber/pesantren"
	"eduvera/internal/adapter/inbound/fiber/sekolah"
	"eduvera/internal/domain"
	inbound_port "eduvera/internal/port/inbound"
)

type adapter struct {
	domain domain.Domain
}

func NewAdapter(
	domain domain.Domain,
) inbound_port.HttpPort {
	return &adapter{
		domain: domain,
	}
}

func (s *adapter) Ping() inbound_port.PingHttpPort {
	return NewPingAdapter()
}

func (s *adapter) Middleware() inbound_port.MiddlewareHttpPort {
	return NewMiddlewareAdapter(s.domain)
}

func (s *adapter) Client() inbound_port.ClientHttpPort {
	return NewClientAdapter(s.domain)
}

func (s *adapter) Landing() inbound_port.LandingHttpPort {
	return NewLandingAdapter(s.domain)
}

func (s *adapter) Onboarding() inbound_port.OnboardingHttpPort {
	return NewOnboardingAdapter(s.domain)
}

func (s *adapter) Auth() inbound_port.AuthHttpPort {
	return NewAuthAdapter(s.domain)
}

func (s *adapter) Payment() inbound_port.PaymentHttpPort {
	return NewPaymentAdapter(s.domain)
}

func (s *adapter) Owner() inbound_port.OwnerHttpPort {
	return NewOwnerAdapter(s.domain)
}

func (a *adapter) Content() inbound_port.ContentHttpPort {
	return NewContentAdapter(a.domain)
}

func (a *adapter) PesantrenDashboard() inbound_port.PesantrenDashboardHttpPort {
	return pesantren.NewDashboardAdapter(a.domain.PesantrenDashboard())
}

func (a *adapter) SPP() inbound_port.SPPHttpPort {
	return NewSPPAdapter(a.domain)
}

func (a *adapter) Sekolah() inbound_port.SekolahHttpPort {
	return sekolah.NewAkademikHandler(a.domain.Sekolah())
}

func (a *adapter) ERapor() inbound_port.ERaporHttpPort {
	return NewERaporAdapter(a.domain)
}
