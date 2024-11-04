package schema

import (
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/supermicah/go-framework-admin/internal/config"
	"github.com/supermicah/go-framework-admin/pkg/crypto/hash"
	"github.com/supermicah/go-framework-admin/pkg/errors"
	"github.com/supermicah/go-framework-admin/pkg/util"
)

const (
	UserStatusActivated = "activated"
	UserStatusFreezed   = "freezed"
)

// User management for RBAC
type User struct {
	ID        int64     `json:"id" gorm:"size:64;primarykey;autoIncrement;"` // Unique ID
	Username  string    `json:"username" gorm:"size:64;index"`               // Username for login
	Name      string    `json:"name" gorm:"size:64;index"`                   // Name of user
	Password  string    `json:"-" gorm:"size:64;"`                           // Password for login (encrypted)
	Phone     string    `json:"phone" gorm:"size:32;"`                       // Phone number of user
	Email     string    `json:"email" gorm:"size:128;"`                      // Email of user
	Remark    string    `json:"remark" gorm:"size:1024;"`                    // Remark of user
	Status    string    `json:"status" gorm:"size:20;index"`                 // Status of user (activated, freezed)
	CreatedAt time.Time `json:"created_at" gorm:"index;"`                    // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`                    // Update time
	Roles     UserRoles `json:"roles" gorm:"-"`                              // Roles of user
}

func (a *User) TableName() string {
	return config.C.FormatTableName("user")
}

// UserQueryParam Defining the query parameters for the `User` struct.
type UserQueryParam struct {
	util.PaginationParam
	LikeUsername string `form:"username"`                                    // Username for login
	LikeName     string `form:"name"`                                        // Name of user
	Status       string `form:"status" binding:"oneof=activated freezed ''"` // Status of user (activated, freezed)
}

// UserQueryOptions Defining the query options for the `User` struct.
type UserQueryOptions struct {
	util.QueryOptions
}

// UserQueryResult Defining the query result for the `User` struct.
type UserQueryResult struct {
	Data       Users
	PageResult *util.PaginationResult
}

// Users Defining the slice of `User` struct.
type Users []*User

func (a Users) ToIDs() []int64 {
	var ids []int64
	for _, item := range a {
		ids = append(ids, item.ID)
	}
	return ids
}

// UserForm Defining the data structure for creating a `User` struct.
type UserForm struct {
	Username string    `json:"username" binding:"required,max=64"`                // Username for login
	Name     string    `json:"name" binding:"required,max=64"`                    // Name of user
	Password string    `json:"password" binding:"max=64"`                         // Password for login (md5 hash)
	Phone    string    `json:"phone" binding:"max=32"`                            // Phone number of user
	Email    string    `json:"email" binding:"max=128"`                           // Email of user
	Remark   string    `json:"remark" binding:"max=1024"`                         // Remark of user
	Status   string    `json:"status" binding:"required,oneof=activated freezed"` // Status of user (activated, freezed)
	Roles    UserRoles `json:"roles" binding:"required"`                          // Roles of user
}

// Validate A validation function for the `UserForm` struct.
func (a *UserForm) Validate() error {
	if a.Email != "" && validator.New().Var(a.Email, "email") != nil {
		return errors.BadRequest("", "Invalid email address")
	}
	return nil
}

// FillTo Convert `UserForm` to `User` object.
func (a *UserForm) FillTo(user *User) error {
	user.Username = a.Username
	user.Name = a.Name
	user.Phone = a.Phone
	user.Email = a.Email
	user.Remark = a.Remark
	user.Status = a.Status

	if pass := a.Password; pass != "" {
		hashPass, err := hash.GeneratePassword(pass)
		if err != nil {
			return errors.BadRequest("", "Failed to generate hash password: %s", err.Error())
		}
		user.Password = hashPass
	}

	return nil
}
