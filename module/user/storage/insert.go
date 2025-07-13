package storage

import (
	"context"
	"todolist/common"
	"todolist/module/user/model"
)

func (s *sqlStore) CreateUser(ctx context.Context, data *model.UserCreate) error {
	db := s.db.Begin()
	//data.PrepareForInsert()  for transaction
	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		db.Rollback()
		return common.ErrDB(err)
	}
	if err := db.Commit().Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
