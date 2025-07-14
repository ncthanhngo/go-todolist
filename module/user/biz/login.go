package biz

import (
	"context"
	"todolist/common"
	"todolist/component/tokenprovider"
	"todolist/module/user/model"
)

type LoginStorage interface {
	FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error)
}
type LoginBisiness struct {
	storeUser     LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBusiness(storeUser LoginStorage, tokenprovider tokenprovider.Provider, hasher Hasher, expiry int) *LoginBisiness {
	return &LoginBisiness{
		storeUser:     storeUser,
		tokenProvider: tokenprovider,
		hasher:        hasher,
		expiry:        expiry,
	}
}

// 1. Find user, email
// 2. Hash pass from input and compare with pass in db
// 3. Provider: issue JWT token from client
// 3.1 Access token and refresh token
// 4. Return token(s)
func (business *LoginBisiness) Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error) {
	user, err := business.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
		if err == common.RecordNotFound {
			return nil, model.ErrEmailOrPasswordInvalid
		}
		return nil, common.ErrDB(err)
	}
	if !business.hasher.Compare(user.Password, data.Password+user.Salt) {
		return nil, model.ErrEmailOrPasswordInvalid
	}

	payload := &common.TokenPayload{
		UId:   user.Id,
		URole: user.Role.String(),
	}
	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}
	//refreshToken, err := business.tokenProvider.Generate(payload, business.tkCfg.GetRtExp())
	//if err != nil {
	//	return nil, common.ErrInternal(err)
	//}
	//account := model.User.NewAccount(accessToken, refreshToken)
	return accessToken, nil
}
