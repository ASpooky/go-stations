package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	response := model.TODO{}

	smtm, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//保存
	result, err := smtm.ExecContext(ctx, subject, description)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//resultからIDを取得
	lastId, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	response.ID = int64(lastId)

	smtm, err = s.db.PrepareContext(ctx, confirm)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//確認
	err = smtm.QueryRowContext(ctx, lastId).Scan(&response.Subject, &response.Description, &response.CreatedAt, &response.UpdatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &response, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos := []*model.TODO{}

	if prevID == 0 {
		rows, err := s.db.QueryContext(ctx, read, size)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			todo := model.TODO{}
			if err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	} else {
		rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			todo := model.TODO{}
			if err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	response := model.TODO{}
	response.ID = id

	smtm, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//保存
	result, err := smtm.ExecContext(ctx, subject, description, id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//変更確認
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Println("in service" + err.Error())
		return nil, err
	}

	if n == 0 {
		errNotFound := model.ErrNotFound{
			What: "update record not found",
			When: time.Now(),
		}
		fmt.Println(errNotFound)
		return nil, &errNotFound
	}

	smtm, err = s.db.PrepareContext(ctx, confirm)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//確認
	err = smtm.QueryRowContext(ctx, id).Scan(&response.Subject, &response.Description, &response.CreatedAt, &response.UpdatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &response, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
