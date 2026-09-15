package api

import (
	"TODO_List/pkg/db"
	"net/http"
	"time"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	defer r.Body.Close()

	var tasks []*db.Task
	var err error

	if len(search) != 0 {

		if dateSearch, parseErr := time.Parse("02.01.2006", search); parseErr == nil {

			tasks, err = db.FindWithDate(dateSearch.Format("20060102"), 50)
		} else {

			tasks, err = db.FindWord(search, 50)
		}
	} else {

		tasks, err = db.Tasks(50)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
		return
	}

	if len(task.Repeat) == 0 {
		err := db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
			return
		}

		writeJson(w, struct{}{}, http.StatusOK)
		return
	}

	nextD, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = db.UpdateDate(nextD, task.ID)
	if err != nil {
		writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}

func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "internal error"}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
