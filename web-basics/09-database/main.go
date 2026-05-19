// 09 - Database với database/sql + SQLite
// Dùng pure-Go SQLite driver (modernc.org/sqlite) — không cần cgo.
//
// Lưu ý: bài này sử dụng dependency external. Bạn cần `go mod tidy` để fetch.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite" // register driver "sqlite"
)

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// === Repository ===

type NoteRepo struct {
	db *sql.DB
}

func (r *NoteRepo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS notes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL,
			body       TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	return err
}

func (r *NoteRepo) List(ctx context.Context) ([]Note, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, body, created_at FROM notes ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (r *NoteRepo) Create(ctx context.Context, n Note) (Note, error) {
	// Dùng prepared statement (qua ExecContext) — chống SQL injection
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO notes (title, body) VALUES (?, ?)`,
		n.Title, n.Body)
	if err != nil {
		return Note{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Note{}, err
	}
	return r.Get(ctx, id)
}

func (r *NoteRepo) Get(ctx context.Context, id int64) (Note, error) {
	var n Note
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, body, created_at FROM notes WHERE id = ?`, id).
		Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt)
	return n, err
}

func (r *NoteRepo) Delete(ctx context.Context, id int64) (bool, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// === Helpers ===

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	// In-memory SQLite (file: "notes.db" nếu muốn persist)
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Tuning pool
	db.SetMaxOpenConns(1) // SQLite single-writer
	db.SetMaxIdleConns(1)

	repo := &NoteRepo{db: db}
	ctx := context.Background()
	if err := repo.Init(ctx); err != nil {
		log.Fatal(err)
	}

	// Seed
	repo.Create(ctx, Note{Title: "First", Body: "Hello SQLite from Go"})
	repo.Create(ctx, Note{Title: "Second", Body: "database/sql is great"})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /notes", func(w http.ResponseWriter, r *http.Request) {
		notes, err := repo.List(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, notes)
	})

	mux.HandleFunc("POST /notes", func(w http.ResponseWriter, r *http.Request) {
		var n Note
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		created, err := repo.Create(r.Context(), n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		n, err := repo.Get(r.Context(), id)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, n)
	})

	mux.HandleFunc("DELETE /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		ok, err := repo.Delete(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
