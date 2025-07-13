package storage

import (
	"context"
	"gorm.io/gorm"
	"todolist/common"
	"todolist/module/user/model"
)

func (s *sqlStore) FindUser(ctx context.Context, cond map[string]interface{}, moreInfo ...string) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())
	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}
	var user model.User
	if err := db.Where(cond).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}
	return &user, nil
}
