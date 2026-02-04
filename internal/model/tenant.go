package model

import (
	"time"

	"github.com/lib/pq"
)

// Plan types for tenant (horizontal - which modules)
const (
	PlanTypeSekolah   = "sekolah"
	PlanTypePesantren = "pesantren"
	PlanTypeHybrid    = "hybrid"
)

// Subscription tiers (vertical - feature level)
const (
	TierBasic   = "basic"   // Manual payment confirmation
	TierPremium = "premium" // Auto payment gateway
)

// Tenant status
const (
	TenantStatusPending   = "pending"
	TenantStatusActive    = "active"
	TenantStatusSuspended = "suspended"
)

// Institution types
const (
	InstitutionTypeNegeri  = "Negeri"
	InstitutionTypeSwasta  = "Swasta"
	InstitutionTypeYayasan = "Yayasan"
)

type Tenant struct {
	ID               string         `json:"id" db:"id"`
	Name             string         `json:"name" db:"name"`                                 // Primary name / Yayasan name
	SchoolName       string         `json:"school_name,omitempty" db:"school_name"`         // For hybrid: specific school name
	PesantrenName    string         `json:"pesantren_name,omitempty" db:"pesantren_name"`   // For hybrid: specific pesantren name
	SchoolJenjangs   pq.StringArray `json:"school_jenjangs,omitempty" db:"school_jenjangs"` // Multi-select: TK, SD, MI, SMP, MTs, SMA, MA, SMK
	Subdomain        string         `json:"subdomain" db:"subdomain"`
	PlanType         string         `json:"plan_type" db:"plan_type"`
	SubscriptionTier string         `json:"subscription_tier" db:"subscription_tier"` // basic or premium
	InstitutionType  string         `json:"institution_type,omitempty" db:"institution_type"`
	Address          string         `json:"address,omitempty" db:"address"`
	BankName         string         `json:"bank_name,omitempty" db:"bank_name"`
	AccountNumber    string         `json:"account_number,omitempty" db:"account_number"`
	AccountHolder    string         `json:"account_holder,omitempty" db:"account_holder"`
	Status           string         `json:"status" db:"status"`
	IsSandbox        bool           `json:"is_sandbox" db:"is_sandbox"`       // True if sandbox tenant for owner testing
	OwnerID          *string        `json:"owner_id,omitempty" db:"owner_id"` // Owner user ID if sandbox
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

type TenantInput struct {
	Name            string   `json:"name" validate:"required"`
	SchoolName      string   `json:"school_name,omitempty"`     // For hybrid packages
	PesantrenName   string   `json:"pesantren_name,omitempty"`  // For hybrid packages
	SchoolJenjangs  []string `json:"school_jenjangs,omitempty"` // Multi-select jenjang
	Subdomain       string   `json:"subdomain" validate:"required,min=3,max=50"`
	PlanType        string   `json:"plan_type" validate:"required,oneof=sekolah pesantren hybrid"`
	InstitutionType string   `json:"institution_type,omitempty"`
	Address         string   `json:"address,omitempty"`
	BankName        string   `json:"bank_name,omitempty"`
	AccountNumber   string   `json:"account_number,omitempty"`
	AccountHolder   string   `json:"account_holder,omitempty"`
}

type TenantFilter struct {
	IDs        []string `json:"ids"`
	Subdomains []string `json:"subdomains"`
	PlanTypes  []string `json:"plan_types"`
	Statuses   []string `json:"statuses"`
}

func TenantPrepare(input *TenantInput) *Tenant {
	return &Tenant{
		Name:            input.Name,
		SchoolName:      input.SchoolName,
		PesantrenName:   input.PesantrenName,
		SchoolJenjangs:  pq.StringArray(input.SchoolJenjangs),
		Subdomain:       input.Subdomain,
		PlanType:        input.PlanType,
		InstitutionType: input.InstitutionType,
		Address:         input.Address,
		BankName:        input.BankName,
		AccountNumber:   input.AccountNumber,
		AccountHolder:   input.AccountHolder,
		Status:          TenantStatusPending,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (f TenantFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.Subdomains) == 0 && len(f.PlanTypes) == 0 && len(f.Statuses) == 0
}
