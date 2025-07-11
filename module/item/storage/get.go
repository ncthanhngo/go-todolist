package storage

import (
	"context"
	"todolist/module/item/model"
)

func (s *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem
	if err := s.db.Where(cond).First(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
