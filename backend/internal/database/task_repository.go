package database

import (
	"database/sql"
	"time"

	"github.com/Ghahremanialireza/routinetask-backend/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	var dueDate sql.NullString
	if task.DueDate != nil {
		dueDate = sql.NullString{String: task.DueDate.Format(time.RFC3339), Valid: true}
	}

	query := `
		INSERT INTO tasks (title, description, type, frequency, due_date, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		task.Title,
		task.Description,
		task.Type,
		task.Frequency,
		dueDate,
		task.Completed,
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id

	return nil
}

func scanTask(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Task, error) {
	var t models.Task
	var description, frequency sql.NullString
	var dueDate sql.NullString
	var createdAt, updatedAt string

	err := scanner.Scan(&t.ID, &t.Title, &description, &t.Type, &frequency, &dueDate, &t.Completed, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	t.Description = description.String
	t.Frequency = frequency.String

	if dueDate.Valid {
		parsed, err := time.Parse(time.RFC3339, dueDate.String)
		if err == nil {
			t.DueDate = &parsed
		}
	}

	if parsed, err := time.Parse(time.RFC3339, createdAt); err == nil {
		t.CreatedAt = parsed
	}
	if parsed, err := time.Parse(time.RFC3339, updatedAt); err == nil {
		t.UpdatedAt = parsed
	}

	return &t, nil
}

func (r *TaskRepository) GetAll() ([]models.Task, error) {
	query := `SELECT id, title, description, type, frequency, due_date, completed, created_at, updated_at FROM tasks ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	query := `SELECT id, title, description, type, frequency, due_date, completed, created_at, updated_at FROM tasks WHERE id = ?`
	row := r.db.QueryRow(query, id)
	return scanTask(row)
}

func (r *TaskRepository) UpdateCompleted(id int64, completed bool) error {
	query := `UPDATE tasks SET completed = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, completed, time.Now().Format(time.RFC3339), id)
	return err
}

func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
