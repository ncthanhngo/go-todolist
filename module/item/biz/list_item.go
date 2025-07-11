package biz

import (
	"context"
	"todolist/common"
	"todolist/module/item/model"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}
type ListItemBiz struct {
	store ListItemStorage
}

func NewListItemBiz(store ListItemStorage) *ListItemBiz {
	return &ListItemBiz{store: store}

}
func (s *ListItemBiz) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	data, err := s.store.ListItem(ctx, filter, paging)
	if err != nil {
		return nil, common.ErrCanNotListEntity(model.EntityName, err)
	}
	return data, nil

}
