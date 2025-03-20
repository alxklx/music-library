package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alxklx/music-library/internal/domain"
	"github.com/gorilla/mux"
)

type SongUsecase interface {
	AddSong(group, song string) (*domain.Song, error)
	UpdateSong(song *domain.Song) error
	DeleteSong(id int64) error
	GetSong(id int64) (*domain.Song, error)
	GetSongs(filter domain.SongFilter) ([]*domain.Song, error)
	GetSongVerses(id int64, pagination domain.VersePagination) ([]string, error)
}

// AddSong godoc
// @Summary Добавление новой песни
// @Description Добавляет новую песню и обогащает её данными из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body AddSongRequest true "Данные песни"
// @Success 201 {object} domain.Song
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var input AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	song, err := h.usecase.AddSong(input.Group, input.Song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(song)
}

// GetSongs godoc
// @Summary Получение списка песен
// @Description Возвращает список песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param releaseDate query string false "Фильтр по дате релиза"
// @Param limit query int false "Лимит записей" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} domain.Song
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	filter := domain.SongFilter{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 10
	}
	filter.Limit = limit
	filter.Offset = offset

	songs, err := h.usecase.GetSongs(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSong godoc
// @Summary Получение песни по ID
// @Description Возвращает данные песни по её идентификатору
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} domain.Song
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [get]
func (h *SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	song, err := h.usecase.GetSong(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// GetSongVerses godoc
// @Summary Получение текста песни с пагинацией по куплетам
// @Description Возвращает куплеты песни с пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param limit query int false "Лимит куплетов" default(1)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} string
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id}/verses [get]
func (h *SongHandler) GetSongVerses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 1
	}

	pagination := domain.VersePagination{
		Limit:  limit,
		Offset: offset,
	}

	verses, err := h.usecase.GetSongVerses(id, pagination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

// UpdateSong godoc
// @Summary Обновление данных песни
// @Description Обновляет информацию о песне по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body domain.Song true "Новые данные песни"
// @Success 200
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var song domain.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	song.ID = id

	if err := h.usecase.UpdateSong(&song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSong godoc
// @Summary Удаление песни
// @Description Удаляет песню по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteSong(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
