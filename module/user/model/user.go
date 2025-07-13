package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"todolist/common"
)

const EntityName = "User"

type UserRole int

const (
	RoleUser UserRole = 1 << iota
	RoleAdmin
	RoleShipper
	RoleMod
)

func (role UserRole) String() string {
	switch role {
	case RoleAdmin:
		return "admin"
	case RoleShipper:
		return "shipper"
	case RoleMod:
		return "mod"
	default:
		return "user"
	}
}

// 2 Ham scan va value dung de thao tac voi DB
func (role *UserRole) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprintf("Faile to unmarshal JSON value: %s", value))
	}
	var r UserRole
	roleValue := string(bytes)
	if roleValue == "user" {
		r = RoleUser
	} else if roleValue == "admin" {
		r = RoleAdmin
	} else if roleValue == "shipper" {
		r = RoleShipper
	} else if roleValue == "mod" {
		r = RoleMod
	}
	*role = r
	return nil
}

func (role *UserRole) Value() (driver.Value, error) {
	if role == nil {
		return nil, nil
	}
	return role.String(), nil
}

type User struct {
	common.SQLModel
	Email     string   `json:"email" gorm:"column:email;"`
	Password  string   `json:"password" gorm:"column:password;"`
	Salt      string   `json:"_" gorm:"column:salt;"`
	LastName  string   `json:"last_name" gorm:"column:last_name;"`
	FirstName string   `json:"first_name" gorm:"column:first_name;"`
	Phone     string   `json:"phone" gorm:"column:phone;"`
	Role      UserRole `json:"role" gorm:"column:role;"`
}

func (role *UserRole) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", role.String())), nil

}
func (u *User) GetUserId() int {
	return u.Id
}
func (u *User) GetUserEmail() string {
	return u.Email
}
func (u *User) GetUserRole() string {
	return u.Role.String()
}
func (User) TableName() string { // User se tro toi bang users trong database
	return "users"
}

type UserCreate struct {
	common.SQLModel `json:",inline"`
	Email           string   `json:"email" gorm:"column:email;"`
	Password        string   `json:"password" gorm:"column:password;"`
	Salt            string   `json:"_" gorm:"column:salt;"`
	LastName        string   `json:"last_name" gorm:"column:last_name;"`
	FirstName       string   `json:"first_name" gorm:"column:first_name;"`
	Role            UserRole `json:"_" gorm:"column:role;"` //Role va salt khong cho truyen tu ngoai vao "_"
}

func (UserCreate) TableName() string {
	return User{}.TableName()
}

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrEmailOrPasswordInvalid = common.NewCustomError(
		errors.New("Email or password invalid"),
		"email or password invalid",
		"ErrUsernameOrPasswordInvalid")
	ErrEmailExisted = common.NewCustomError(
		errors.New("Email has already Existed"),
		"Email has already existed",
		"ErrUserEmailExisted")
)
