package schema

import (
	"time"

	"github.com/supermicah/go-framework-admin/internal/config"
	"github.com/supermicah/go-framework-admin/pkg/util"
)

// UserRole User roles for RBAC
type UserRole struct {
	ID        int64     `json:"id" gorm:"size:64;primarykey;autoIncrement;"` // Unique ID
	UserID    int64     `json:"user_id" gorm:"size:64;index"`                // From User.ID
	RoleID    int64     `json:"role_id" gorm:"size:64;index"`                // From Role.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`                    // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`                    // Update time
	RoleName  string    `json:"role_name" gorm:"<-:false;-:migration;"`      // From Role.Name
}

func (a *UserRole) TableName() string {
	return config.C.FormatTableName("user_role")
}

// UserRoleQueryParam Defining the query parameters for the `UserRole` struct.
type UserRoleQueryParam struct {
	util.PaginationParam
	InUserIDs []int64 `form:"-"` // From User.ID
	UserID    int64   `form:"-"` // From User.ID
	RoleID    int64   `form:"-"` // From Role.ID
}

// UserRoleQueryOptions Defining the query options for the `UserRole` struct.
type UserRoleQueryOptions struct {
	util.QueryOptions
	JoinRole bool // Join role table
}

// UserRoleQueryResult Defining the query result for the `UserRole` struct.
type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *util.PaginationResult
}

// UserRoles Defining the slice of `UserRole` struct.
type UserRoles []*UserRole

func (a UserRoles) ToUserIDMap() map[int64]UserRoles {
	m := make(map[int64]UserRoles)
	for _, userRole := range a {
		m[userRole.UserID] = append(m[userRole.UserID], userRole)
	}
	return m
}

func (a UserRoles) ToRoleIDs() []int64 {
	var ids []int64
	for _, item := range a {
		ids = append(ids, item.RoleID)
	}
	return ids
}

// UserRoleForm Defining the data structure for creating a `UserRole` struct.
type UserRoleForm struct {
}

// Validate A validation function for the `UserRoleForm` struct.
func (a *UserRoleForm) Validate() error {
	return nil
}

func (a *UserRoleForm) FillTo(userRole *UserRole) error {
	return nil
}
