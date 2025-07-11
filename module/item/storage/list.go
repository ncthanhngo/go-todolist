package storage

import (
	"context"
	"todolist/common"
	"todolist/module/item/model"
)

func (s *sqlStore) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	var result []model.TodoItem
	db := s.db.Table(model.TodoItem{}.TableName()).Where("status <> ?", "deleted")

	if filter != nil && filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	// Đếm tổng
	if err := db.Select("id").Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	// Lấy dữ liệu
	if err := db.
		Select("*").
		Order("id DESC").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
