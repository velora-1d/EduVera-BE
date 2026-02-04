package model

import "time"

// Relationship types for owner institutions
const (
	OwnerRelationshipOwned   = "owned"
	OwnerRelationshipManaged = "managed"
)

// OwnerInstitution represents the relationship between owner and their institutions
type OwnerInstitution struct {
	ID               string    `json:"id" db:"id"`
	OwnerUserID      string    `json:"owner_user_id" db:"owner_user_id"`
	TenantID         string    `json:"tenant_id" db:"tenant_id"`
	RelationshipType string    `json:"relationship_type" db:"relationship_type"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	Tenant *Tenant `json:"tenant,omitempty" db:"-"`
}

type OwnerInstitutionInput struct {
	TenantID         string `json:"tenant_id" validate:"required"`
	RelationshipType string `json:"relationship_type,omitempty"`
}

type OwnerInstitutionFilter struct {
	IDs          []string `json:"ids"`
	OwnerUserIDs []string `json:"owner_user_ids"`
	TenantIDs    []string `json:"tenant_ids"`
}

func OwnerInstitutionPrepare(ownerUserID string, input *OwnerInstitutionInput) *OwnerInstitution {
	relType := input.RelationshipType
	if relType == "" {
		relType = OwnerRelationshipOwned
	}

	return &OwnerInstitution{
		OwnerUserID:      ownerUserID,
		TenantID:         input.TenantID,
		RelationshipType: relType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func (f OwnerInstitutionFilter) IsEmpty() bool {
	return len(f.IDs) == 0 && len(f.OwnerUserIDs) == 0 && len(f.TenantIDs) == 0
}

// ImpersonateInput for switch account feature
type ImpersonateInput struct {
	TenantID string `json:"tenant_id" validate:"required"`
	ViewMode string `json:"view_mode,omitempty"` // sekolah, pesantren, or empty for auto
}

// ImpersonateResponse contains impersonation token and tenant info
type ImpersonateResponse struct {
	ImpersonateToken string  `json:"impersonate_token"`
	ExpiresAt        int64   `json:"expires_at"`
	Tenant           *Tenant `json:"tenant"`
	ViewMode         string  `json:"view_mode"`
}
