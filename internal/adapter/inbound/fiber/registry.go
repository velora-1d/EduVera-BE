package fiber_inbound_adapter

import (
	"prabogo/internal/adapter/inbound/fiber/pesantren"
	"prabogo/internal/adapter/inbound/fiber/sekolah"
	"prabogo/internal/domain"
	inbound_port "prabogo/internal/port/inbound"
	outbound_port "prabogo/internal/port/outbound"
)

type adapter struct {
	domain  domain.Domain
	message outbound_port.MessagePort
}

func NewAdapter(
	domain domain.Domain,
	message outbound_port.MessagePort,
) inbound_port.HttpPort {
	return &adapter{
		domain:  domain,
		message: message,
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
	var whatsapp outbound_port.WhatsAppMessagePort
	if s.message != nil && s.message.WhatsApp() != nil {
		whatsapp = s.message.WhatsApp()
	}
	return NewOnboardingAdapter(s.domain, whatsapp)
}

func (s *adapter) Auth() inbound_port.AuthHttpPort {
	return NewAuthAdapter(s.domain)
}

func (s *adapter) Payment() inbound_port.PaymentHttpPort {
	var whatsapp outbound_port.WhatsAppMessagePort
	if s.message != nil && s.message.WhatsApp() != nil {
		whatsapp = s.message.WhatsApp()
	}
	return NewPaymentAdapter(s.domain, whatsapp)
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

func (a *adapter) SDM() inbound_port.SDMHttpPort {
	return NewSDMAdapter(a.domain)
}

func (a *adapter) Subscription() inbound_port.SubscriptionHttpPort {
	return NewSubscriptionAdapter(a.domain)
}

func (a *adapter) Analytics() inbound_port.AnalyticsHttpPort {
	return NewAnalyticsAdapter(a.domain)
}
