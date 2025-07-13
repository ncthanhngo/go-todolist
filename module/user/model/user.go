package model

import (
	"errors"
	"fmt"
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
	}
	*role = r
	return nil
}

//
