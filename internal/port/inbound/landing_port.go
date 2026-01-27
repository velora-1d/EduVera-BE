package inbound_port

type LandingHttpPort interface {
	Home(any) error
	Register(any) error
	RegisterProcess(any) error
	Step2(any) error
	Step2Process(any) error
	Step3(any) error
	Step3Process(any) error
	Step4(any) error
	Step4Process(any) error
	Step5(any) error
	Step5Process(any) error
	Step6(any) error
	Step6Process(any) error
	Step7(any) error
	Step7Process(any) error
}
