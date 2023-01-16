package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateNoteParams struct {
	Title string `json:"title"`
}

type App struct {
	ctx context.Context
	sql *sql.DB
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := sql.Open("sqlite3", "file:db.sqlite?mode=rwc")
	if err != nil {
		panic(err)
	}

	a.sql = db

	err = a.ensureMigrations()
	if err != nil {
		panic(err)
	}
}

func (a *App) GetAllNotes() ([]Note, error) {
	rows, err := a.sql.QueryContext(a.ctx, `SELECT id, title, note, created_at FROM notes ORDER BY created_at DESC;`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	notes := make([]Note, 0)

	for rows.Next() {
		note := Note{}

		err := rows.Scan(&note.ID, &note.Title, &note.Note, &note.CreatedAt)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (a *App) GetNoteByID(id int) (*Note, error) {
	row := a.sql.QueryRowContext(a.ctx, `SELECT id, title, note, created_at FROM notes WHERE id = ?`, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var note Note

	err := row.Scan(&note.ID, &note.Title, &note.Note, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (a *App) CreateNote(params CreateNoteParams) (*Note, error) {
	row := a.sql.QueryRowContext(a.ctx, `INSERT INTO notes(title, note, created_at) VALUES(?,?,?) RETURNING id, title, note, created_at`, params.Title, "", time.Now())
	if err := row.Err(); err != nil {
		fmt.Printf("error creating note: %v", err)
		return nil, err
	}

	var note Note

	err := row.Scan(&note.ID, &note.Title, &note.Note, &note.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (a *App) UpdateNoteTitle(id int, title string) error {
	_, err := a.sql.ExecContext(a.ctx, "UPDATE notes SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) UpdateNoteText(id int, note string) error {
	_, err := a.sql.ExecContext(a.ctx, "UPDATE notes SET note = ? WHERE id = ?", note, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteArticle(id int) error {
	_, err := a.sql.ExecContext(a.ctx, "DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ensureMigrations() error {
	_, err := a.sql.ExecContext(a.ctx, `CREATE TABLE IF NOT EXISTS notes(id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, note TEXT, created_at DATETIME);`)
	return err
}
