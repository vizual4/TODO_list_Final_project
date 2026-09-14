package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

func getDBPath() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		return "scheduler.db"
	}
	return dbFile
}

func AddTask(task *Task) (int64, error) {
	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat);"
	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	res := make([]*Task, 0, limit)

	rows, err := DB.Query("SELECT * FROM scheduler ORDER BY date DESC LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t Task
		var id int64

		err = rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, 10)

		res = append(res, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func GetTask(id string) (*Task, error) {

	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id))

	var result Task
	var idInt int64

	err := row.Scan(&idInt, &result.Date, &result.Title, &result.Comment, &result.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	result.ID = strconv.FormatInt(idInt, 10)

	return &result, nil
}

func UpdateTask(task *Task) error {

	query := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id"

	res, err := DB.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title),
		sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat), sql.Named("id", task.ID))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("incorrect id for updating task")
	}

	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"

	res, err := DB.Exec(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("incorrect id for deletion task")
	}

	return nil
}

func UpdateDate(next, id string) error {
	query := "UPDATE scheduler SET date = :date WHERE id = :id"

	res, err := DB.Exec(query, sql.Named("date", next), sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("incorrect id for updating date")
	}

	return nil
}

func FindWithDate(date string, limit int) ([]*Task, error) {

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date DESC LIMIT :limit"

	rows, err := DB.Query(query, sql.Named("date", date), sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*Task, 0)

	for rows.Next() {
		var t Task
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		t.ID = strconv.FormatInt(id, 10)

		result = append(result, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func FindWord(word string, limit int) ([]*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title 
	LIKE :word OR comment LIKE :word ORDER BY date DESC LIMIT :limit`

	rows, err := DB.Query(query, sql.Named("word", "%"+word+"%"), sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*Task, 0)

	for rows.Next() {
		var t Task
		var id int64
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}

		t.ID = strconv.FormatInt(id, 10)
		result = append(result, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
