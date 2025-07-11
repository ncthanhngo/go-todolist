package biz

import (
	"context"
	"todolist/common"
	"todolist/module/item/model"
)

type DeleteItemStore interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	DeleteItem(ctx context.Context, cond map[string]interface{}) error

	// can use many condictions for this
}
type DeleteItemBiz struct {
	store DeleteItemStore
}

func NewDeleteItemBiz(store DeleteItemStore) *DeleteItemBiz {
	return &DeleteItemBiz{store: store}
}

func (biz *DeleteItemBiz) DeleteItem(ctx context.Context, id int) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCanNotGetEntity(model.EntityName, err)
	}
	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}
	if err := biz.store.DeleteItem(ctx, map[string]interface{}{"id": id}); err != nil {
		return common.ErrCanNotDeleteEntity(model.EntityName, err)
	}
	return nil
}
