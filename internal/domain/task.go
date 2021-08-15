package domain

import (
	"context"
)

const (
	TaskStatusWait = iota + 1
	TaskStatusInProcess
	TaskStatusCompleted
)

type Task struct {
	ID            int64  `db:"id" json:"id"`
	Title         string `db:"title" json:"title"`
	Description   string `db:"description" json:"description"`
	EstimatedTime *int64 `db:"estimated_time" json:"estimatedTime"`
	UserID        *int64 `db:"user_id" json:"userId"`
	Status        int    `db:"status" json:"status"`
	DueDate       *int64 `db:"due_date" json:"dueDate"`
	MailSent      bool   `db:"mail_sent" json:"-"`
	User          *User  `json:"-"`
}

type TaskUsecase interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, ID int64) (Task, error)
	Create(ctx context.Context, task Task) (int64, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, ID int64) error
	AssignUser(ctx context.Context, taskID, userID int64) error
	GetOverdueTasks(ctx context.Context, dueDate int64) ([]Task, error)
}

type TaskRepository interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, ID int64) (Task, error)
	Create(ctx context.Context, task Task) (int64, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, ID int64) error
	AssignUser(ctx context.Context, taskID, userID int64) error
	GetOverdueTasks(ctx context.Context, dueDate int64) ([]Task, error)
}

type TaskValidator interface {
	ValidateUserExist(userID interface{}) error
}
