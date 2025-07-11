package biz

import (
	"context"
	"todolist/module/item/model"
)

type UpdateItemStore interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error

	// can use many condictions for this
}
type UpdateItemBiz struct {
	store UpdateItemStore
}

func NewUpdateItemBiz(store UpdateItemStore) *UpdateItemBiz {
	return &UpdateItemBiz{store: store}
}
func (biz *UpdateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}
	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return err
	}
	return nil
}
