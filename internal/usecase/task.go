package usecase

import (
	"cat-test/internal/domain"
	"cat-test/internal/errors"
	"context"
)

type taskUsecase struct {
	repository domain.TaskRepository
}

func (t taskUsecase) List(ctx context.Context) ([]domain.Task, error) {
	return t.repository.List(ctx)
}

func (t taskUsecase) Get(ctx context.Context, ID int64) (domain.Task, error) {
	return t.repository.Get(ctx, ID)
}

func (t taskUsecase) Create(ctx context.Context, task domain.Task) (int64, error) {
	return t.repository.Create(ctx, task)
}

func (t taskUsecase) Update(ctx context.Context, task domain.Task) error {
	return t.repository.Update(ctx, task)
}

func (t taskUsecase) Delete(ctx context.Context, ID int64) error {
	return t.repository.Delete(ctx, ID)
}

func (t taskUsecase) AssignUser(ctx context.Context, taskID, userID int64) error {
	identity := ctx.Value("identity").(domain.User)
	if identity.ID != userID && identity.Type != domain.UserTypeAdmin {
		return errors.TaskCannotBeAssigned
	}

	return t.repository.AssignUser(ctx, taskID, userID)
}

func (t taskUsecase) GetOverdueTasks(ctx context.Context, dueDate int64) ([]domain.Task, error) {
	return t.repository.GetOverdueTasks(ctx, dueDate)
}

func NewTaskUsecase(repository domain.TaskRepository) domain.TaskUsecase {
	return taskUsecase{repository: repository}
}
