package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ameiswq/Golang2026/assignment-01/internal/models"
	"github.com/ameiswq/Golang2026/assignment-01/internal/storage"
)

type TaskHandler struct {
	store *storage.Store
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	doneParam := r.URL.Query().Get("done")

	if idParam == "" {
		tasks := h.store.List()

		if doneParam != "" {
			done, errDone := strconv.ParseBool(doneParam)
			if errDone != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid param"})
				return
			}

			filtered := make([]models.Task, 0)
			for _, task := range tasks {
				if task.Done == done {
					filtered = append(filtered, task)
				}	
			}
			tasks = filtered
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	task, ok := h.store.Get(id)
	if ok != true {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title string `json:"title"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}
	if reqBody.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "empty title"})
		return
	}
	if len(reqBody.Title) > 100{
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid title. max length - 100 ch."})
		return
	}
	task := h.store.Create(reqBody.Title)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	var reqBody struct {
		Done *bool `json:"done"`
	}

	err2 := json.NewDecoder(r.Body).Decode(&reqBody)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid json"})
		return
	}

	if reqBody.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "done field required"})
		return
	}

	ok := h.store.UpdateDone(id, *reqBody.Done)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}

func NewTaskHandler(store *storage.Store) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	ok := h.store.Delete(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"deleted": true,
	})
}

