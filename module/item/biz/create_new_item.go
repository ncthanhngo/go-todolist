package biz

import (
	"context"
	"todolist/common"
	"todolist/module/item/model"
)

// Handler -> Business [-> Repository] -> Storage
// the relations between layers through interface
type CreateItemStorage interface {
	CreateItem(ctx context.Context, data *model.TodoItemCreation) error
}
type createItemBiz struct {
	store CreateItemStorage
}

// Contructor for struct
func NewCreateItemBiz(store CreateItemStorage) *createItemBiz {
	return &createItemBiz{store: store}
}

func (biz *createItemBiz) CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := data.Validate(); err != nil {
		return err
	}
	if err := biz.store.CreateItem(ctx, data); err != nil {
		return common.ErrCanNotCreateEntity(model.EntityName, err)
	}
	return nil
}
