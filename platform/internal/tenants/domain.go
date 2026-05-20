package tenants

import "time"

type TenantID string
type UserID string
type RoleID string

type TenantType string

const (
	TenantTypePersonal TenantType = "personal"
	TenantTypeTeam     TenantType = "team"
	TenantTypeCompany  TenantType = "company"
	TenantTypeMandate  TenantType = "mandate"
)

type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusInactive TenantStatus = "inactive"
)

type MemberStatus string

const (
	MemberStatusActive   MemberStatus = "active"
	MemberStatusInactive MemberStatus = "inactive"
)

type Tenant struct {
	ID           TenantID
	Type         TenantType
	Name         string
	Slug         string
	OwnerUserID  UserID
	Status       TenantStatus
	BillingEmail string
	CreatedAt    time.Time
}

func (t Tenant) IsActive() bool {
	return t.Status == TenantStatusActive
}

type TenantMember struct {
	TenantID TenantID
	UserID   UserID
	RoleID   RoleID
	Status   MemberStatus
	JoinedAt time.Time
}

func (m TenantMember) IsActive() bool {
	return m.Status == MemberStatusActive
}
