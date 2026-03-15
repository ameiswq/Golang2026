package internal
import (
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	repo *UserRepository
}

func NewUserHandler(repo *UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	var idPtr *int
	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			idPtr = &id
		}
	}

	filter := UserFilter{
		ID: idPtr,
		Name: q.Get("name"),
		Email: q.Get("email"),
		Gender: q.Get("gender"),
		BirthDate: q.Get("birth_date"),
		OrderBy: q.Get("order_by"),
		OrderDir: q.Get("order_dir"),
		Page: page,
		PageSize: pageSize,
	}

	resp, err := h.repo.GetPaginatedUsers(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	user1ID, err1 := strconv.Atoi(q.Get("user1_id"))
	user2ID, err2 := strconv.Atoi(q.Get("user2_id"))

	if err1 != nil || err2 != nil {
		http.Error(w, "user1_id and user2_id must be valid integers", http.StatusBadRequest)
		return
	}

	if user1ID == user2ID {
		http.Error(w, "users must be different", http.StatusBadRequest)
		return
	}

	friends, err := h.repo.GetCommonFriends(user1ID, user2ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user1_id":       user1ID,
		"user2_id":       user2ID,
		"common_friends": friends,
		"count":          len(friends),
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}