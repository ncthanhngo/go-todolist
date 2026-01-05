package biz

import (
	"context"
	"errors"
	"todolist/common"
	"todolist/module/item/model"
)

type UpdateItemStore interface {
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error

	// can use many condictions for this
}
type UpdateItemBiz struct {
	store     UpdateItemStore
	requester common.Requester
}

func NewUpdateItemBiz(store UpdateItemStore, requester common.Requester) *UpdateItemBiz {
	return &UpdateItemBiz{store: store, requester: requester}
}
func (biz *UpdateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return common.ErrCanNotGetEntity(model.EntityName, err)
	}
	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}
	//check owner: Neu dung moi cho update
	isOwner := biz.requester.GetUserId() == data.UserId

	//check Admin hoa Mod : Neu dung thi cho update
	if !isOwner || !common.IsAdminOrMode(biz.requester) {
		return common.ErrNoPermission(errors.New("No permission to update item"))
	}
	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return common.ErrCanNotUpdateEntity(model.EntityName, err)
	}
	return nil
}
