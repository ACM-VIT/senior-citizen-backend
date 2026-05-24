package models

type PermissionCode string

const (
	ConclaveKickout     PermissionCode = "conclave.kickout"
	ConclaveRemoveAdmin PermissionCode = "conclave.remove_admin"
	ConclaveMakeAdmin   PermissionCode = "conclave.make_admin"
	ConclaveAddMember   PermissionCode = "conclave.add_member"
	ConclaveMuteAnyone  PermissionCode = "conclave.mute_anyone"

	ExamcookerWrite       PermissionCode = "examcooker.write"
	ExamcookerApprove     PermissionCode = "examcooker.approve_exam"
	ExamcookerUpload      PermissionCode = "examcooker.upload_exam"
	ExamcookerAddMember   PermissionCode = "examcooker.add_member"
	ExamcookerRemoveAdmin PermissionCode = "examcooker.remove_admin"

	AddMember             PermissionCode = "add_member"
	UpdateMemberAccess    PermissionCode = "update_member_access"
	CreateAdminController PermissionCode = "create_admin_controller"
	CreateUpdateProfile   PermissionCode = "create_update_profile"
	EventUpdate           PermissionCode = "event_update"
	EventCreate           PermissionCode = "event_create"
	EventDelete           PermissionCode = "event_delete"
)

type Permission struct {
	Base
	Code        PermissionCode `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
}

type RolePermission struct {
	Base
	Role         Role       `gorm:"type:varchar(20);not null;index:idx_role_permission,unique" json:"role"`
	PermissionID string     `gorm:"type:uuid;not null;index:idx_role_permission,unique" json:"permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}

var DefaultPermissions = []Permission{
	{Code: ConclaveKickout, Name: "Conclave kickout"},
	{Code: ConclaveRemoveAdmin, Name: "Conclave remove admin"},
	{Code: ConclaveMakeAdmin, Name: "Conclave make admin"},
	{Code: ConclaveAddMember, Name: "Conclave add member"},
	{Code: ConclaveMuteAnyone, Name: "Conclave mute anyone"},
	{Code: ExamcookerWrite, Name: "Examcooker write"},
	{Code: ExamcookerApprove, Name: "Examcooker approve exam"},
	{Code: ExamcookerUpload, Name: "Examcooker upload exam"},
	{Code: ExamcookerAddMember, Name: "Examcooker add member"},
	{Code: ExamcookerRemoveAdmin, Name: "Examcooker remove admin"},
	{Code: AddMember, Name: "Add member"},
	{Code: UpdateMemberAccess, Name: "Update member access"},
	{Code: CreateAdminController, Name: "Create admin controller"},
	{Code: CreateUpdateProfile, Name: "Create/update profile"},
	{Code: EventUpdate, Name: "Event update"},
	{Code: EventCreate, Name: "Event create"},
	{Code: EventDelete, Name: "Event delete"},
}

func DefaultRolePermissionCodes(role Role) []PermissionCode {
	switch role {
	case Board:
		return []PermissionCode{
			ConclaveKickout,
			ConclaveRemoveAdmin,
			ConclaveMakeAdmin,
			ConclaveAddMember,
			ConclaveMuteAnyone,
			ExamcookerWrite,
			ExamcookerApprove,
			ExamcookerAddMember,
			ExamcookerRemoveAdmin,
			AddMember,
			UpdateMemberAccess,
			CreateAdminController,
			CreateUpdateProfile,
			EventUpdate,
			EventCreate,
			EventDelete,
		}
	case SC:
		return []PermissionCode{
			ConclaveAddMember,
			ExamcookerWrite,
			ExamcookerAddMember,
		}
	case JC:
		return []PermissionCode{
			ConclaveAddMember,
			ExamcookerAddMember,
		}
	default:
		return []PermissionCode{}
	}
}

func HasPermission(role Role, permission PermissionCode) bool {
	for _, allowed := range DefaultRolePermissionCodes(role) {
		if allowed == permission {
			return true
		}
	}
	return false
}
