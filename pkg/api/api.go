package api

import (
	"TODO_List/pkg/db"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/signin", authHandler)
	mux.HandleFunc("GET /api/nextdate", nextDateHandler)
	mux.HandleFunc("POST /api/task", auth(taskPostHandler))
	mux.HandleFunc("GET /api/tasks", auth(tasksHandler))
	mux.HandleFunc("GET /api/task", auth(taskGetHandler))
	mux.HandleFunc("PUT /api/task", auth(taskPutHandler))
	mux.HandleFunc("POST /api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("DELETE /api/task", auth(taskDeleteHandler))
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if now == "" {
		now = time.Now().Format(formatTime)
	}

	tNow, err := time.Parse(formatTime, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextD, err := NextDate(tNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextD))

}

func taskPostHandler(w http.ResponseWriter, r *http.Request) {
	var info db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &info); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if len(info.Title) == 0 {
		writeJson(w, map[string]string{"error": "Title field is empty"})
		return
	}

	err = checkDate(&info)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&info)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]any{"id": strconv.FormatInt(id, 10)})

}

func taskGetHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		writeJson(w, map[string]string{"error": "id is empty"})
		return
	}

	res, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, res)
}

func taskPutHandler(w http.ResponseWriter, r *http.Request) {
	var info db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	defer r.Body.Close()

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &info); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if len(info.Title) == 0 {
		writeJson(w, map[string]string{"error": "Title field is empty"})
		return
	}

	if len(info.ID) == 0 {
		writeJson(w, map[string]string{"error": "ID field is empty"})
		return
	}

	err = checkDate(&info)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&info)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, struct{}{})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	today, _ := time.Parse(formatTime, now.Format(formatTime))

	if len(task.Date) == 0 {
		task.Date = today.Format(formatTime)
	}

	t, err := time.Parse(formatTime, task.Date)
	if err != nil {
		return err
	}

	if t.Before(today) {
		if len(task.Repeat) == 0 {
			task.Date = today.Format(formatTime)
		} else {
			next, err := NextDate(today, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
