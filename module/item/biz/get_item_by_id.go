package biz

import (
	"context"
	"todolist/module/item/model"
)

type GetItemStorage interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	// can use many condictions for this
}
type GetItemBiz struct {
	store GetItemStorage
}

func NewGetItemBiz(store GetItemStorage) *GetItemBiz {
	return &GetItemBiz{store: store}
}
func (biz *GetItemBiz) GetItem(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	return data, nil
}
