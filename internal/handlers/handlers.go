package handlers

import (
	"github.com/gorilla/mux"
)

// SongHandler структура для обработки запросов
type SongHandler struct {
	usecase SongUsecase
}

func NewSongHandler(usecase SongUsecase) *SongHandler {
	return &SongHandler{usecase: usecase}
}

// RegisterRoutes регистрирует маршруты
func (h *SongHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/songs", h.GetSongs).Methods("GET")
	router.HandleFunc("/songs/{id}", h.GetSong).Methods("GET")
	router.HandleFunc("/songs/{id}/verses", h.GetSongVerses).Methods("GET")
	router.HandleFunc("/songs", h.AddSong).Methods("POST")
	router.HandleFunc("/songs/{id}", h.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", h.DeleteSong).Methods("DELETE")
}
