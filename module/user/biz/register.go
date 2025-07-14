package biz

import (
	"context"
	"todolist/common"
	"todolist/module/user/model"
)

type RegisterStorage interface {
	FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error)
	CreateUser(ctx context.Context, data *model.UserCreate) error
}

// type Hasher interface { //>> tao chuoi salt cho moi user la khac nhau
//
//	Hash(data string) string
type Hasher interface {
	Hash(data string) string
	Compare(hashedData, plainData string) bool
}

type registerBusiness struct {
	registerStorage RegisterStorage
	hasher          Hasher
}

// contructor
func NewRegisterBusiness(registerStorage RegisterStorage, hasher Hasher) *registerBusiness {
	return &registerBusiness{
		registerStorage: registerStorage,
		hasher:          hasher,
	}
}

func (business *registerBusiness) Register(ctx context.Context, data *model.UserCreate) error {
	user, _ := business.registerStorage.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
		//if user.status == 0 {
		//	return error user has been disbale
		//}
		return model.ErrEmailExisted
	}
	salt := common.GenSalt(50)
	data.Password = business.hasher.Hash(data.Password + salt)
	data.Salt = salt
	role := model.RoleUser //Default user role
	data.Role = &role

	if err := business.registerStorage.CreateUser(ctx, data); err != nil {
		return common.ErrCanNotCreateEntity(model.EntityName, err)
	}
	return nil
}
