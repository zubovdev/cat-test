package repository

import (
	"cat-test/internal/domain"
	"context"
	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db *sqlx.DB
}

func (t taskRepository) List(ctx context.Context) (tasks []domain.Task, err error) {
	sql := `SELECT * FROM task ORDER BY id`
	return tasks, t.db.SelectContext(ctx, &tasks, sql)
}

func (t taskRepository) Get(ctx context.Context, ID int64) (task domain.Task, err error) {
	sql := `SELECT * FROM task WHERE task.id = $1 LIMIT 1`
	return task, t.db.GetContext(ctx, &task, sql, ID)
}

func (t taskRepository) Create(ctx context.Context, task domain.Task) (taskID int64, err error) {
	sql := `INSERT INTO task(title, description, estimated_time, user_id, due_date) 
				VALUES(:title, :description, :estimated_time, :user_id, :due_date) RETURNING id`
	rows, err := t.db.NamedQueryContext(ctx, sql, task)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()

	return taskID, rows.Scan(&taskID)
}

func (t taskRepository) Update(ctx context.Context, task domain.Task) error {
	sql := `UPDATE task SET title = :title, description = :description, status = :status, due_date = :due_date,
                estimated_time = :estimated_time, mail_sent = :mail_sent WHERE id = :id`
	_, err := t.db.NamedExecContext(ctx, sql, task)
	return err
}

func (t taskRepository) Delete(ctx context.Context, ID int64) error {
	sql := `DELETE FROM task WHERE task.id = $1`
	_, err := t.db.ExecContext(ctx, sql, ID)
	return err
}

func (t taskRepository) AssignUser(ctx context.Context, taskID, userID int64) error {
	sql := `UPDATE task SET user_id = $2 WHERE id = $1`
	_, err := t.db.ExecContext(ctx, sql, taskID, userID)
	return err
}

func (t taskRepository) GetOverdueTasks(ctx context.Context, dueDate int64) (tasks []domain.Task, err error) {
	sql := `SELECT task.*, u.id, u.email, u.first_name, u.last_name FROM task 
    			INNER JOIN "user" u ON task.user_id = u.id WHERE due_date < $1 AND mail_sent = false`

	rows, err := t.db.QueryxContext(ctx, sql, dueDate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		task := domain.Task{}
		user := domain.User{}

		if err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.EstimatedTime, &task.UserID,
			&task.Status, &task.DueDate, &task.MailSent, &user.ID, &user.Email, &user.FirstName, &user.LastName); err != nil {
			return nil, err
		}

		task.User = &user
		tasks = append(tasks, task)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewTaskRepository(db *sqlx.DB) domain.TaskRepository {
	return taskRepository{db: db}
}
