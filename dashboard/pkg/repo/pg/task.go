package pg

import (
	"context"
	"fmt"

	"dashboard.dishdash.ru/pkg/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (tr *TaskRepo) CreateTask(ctx context.Context, task *domain.Task) (int64, error) {
	query := `INSERT INTO dashboard.task (title, description, status, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	var id int64
	err := tr.db.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not insert task: %w", err)
	}
	return id, nil
}

func (tr *TaskRepo) GetTaskByID(ctx context.Context, taskId int64) (*domain.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM dashboard.task WHERE id = $1`
	task := &domain.Task{}
	err := tr.db.QueryRow(ctx, query, taskId).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not get task by ID: %w", err)
	}
	return task, nil
}

func (tr *TaskRepo) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM dashboard.task`
	rows, err := tr.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("could not scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (tr *TaskRepo) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `UPDATE dashboard.task SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4 RETURNING id, title, description, status, created_at, updated_at`
	updatedTask := &domain.Task{}
	err := tr.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.ID).
		Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Status, &updatedTask.CreatedAt, &updatedTask.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not update task: %w", err)
	}
	return updatedTask, nil
}

func (tr *TaskRepo) DeleteTask(ctx context.Context, taskId int64) error {
	query := `DELETE FROM dashboard.task WHERE id = $1`
	_, err := tr.db.Exec(ctx, query, taskId)
	if err != nil {
		return fmt.Errorf("could not delete task: %w", err)
	}
	return nil
}
