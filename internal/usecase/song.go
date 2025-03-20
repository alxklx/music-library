package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/alxklx/music-library/internal/domain"
)

type SongUsecase struct {
	repo       domain.SongRepository
	apiBaseURL string // Базовый URL внешнего API
	httpClient *http.Client
}

func NewSongUsecase(repo domain.SongRepository, apiBaseURL string) *SongUsecase {
	return &SongUsecase{
		repo:       repo,
		apiBaseURL: apiBaseURL,
		httpClient: &http.Client{},
	}
}

// SongDetail структура для ответа внешнего API
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (u *SongUsecase) AddSong(group, song string) (*domain.Song, error) {
	log.Printf("INFO: Adding song: group=%s, song=%s", group, song)

	// 1. Запрос к внешнему API для обогащения данных
	songDetail, err := u.fetchSongDetail(group, song)
	if err != nil {
		log.Printf("ERROR: Failed to fetch song detail: %v", err)
		return nil, err
	}
	log.Printf("DEBUG: Song detail fetched: %+v", songDetail)

	// 2. Создание сущности песни
	newSong := &domain.Song{
		Group:       group,
		Song:        song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// 3. Сохранение в БД
	if err := u.repo.Create(newSong); err != nil {
		log.Printf("ERROR: Failed to save song to DB: %v", err)
		return nil, err
	}
	log.Printf("INFO: Song saved with ID: %d", newSong.ID)

	return newSong, nil
}

func (u *SongUsecase) fetchSongDetail(group, song string) (*SongDetail, error) {
	query := url.Values{}
	query.Set("group", group)
	query.Set("song", song)
	fullURL := u.apiBaseURL + "?" + query.Encode()
	log.Printf("DEBUG: Fetching song detail from %s", fullURL)

	resp, err := u.httpClient.Get(fullURL)
	if err != nil {
		log.Printf("ERROR: HTTP request failed: %v", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("external API returned status: %d", resp.StatusCode)
		log.Printf("ERROR: %v", err)
		return nil, err
	}

	var detail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		log.Printf("ERROR: Failed to decode API response: %v", err)
		return nil, err
	}

	return &detail, nil
}

func (u *SongUsecase) UpdateSong(song *domain.Song) error {
	return u.repo.Update(song)
}

func (u *SongUsecase) DeleteSong(id int64) error {
	return u.repo.Delete(id)
}

func (u *SongUsecase) GetSong(id int64) (*domain.Song, error) {
	return u.repo.FindByID(id)
}

func (u *SongUsecase) GetSongs(filter domain.SongFilter) ([]*domain.Song, error) {
	return u.repo.FindAll(filter)
}

func (u *SongUsecase) GetSongVerses(id int64, pagination domain.VersePagination) ([]string, error) {
	log.Printf("INFO: Fetching verses for song ID=%d, limit=%d, offset=%d", id, pagination.Limit, pagination.Offset)

	song, err := u.repo.FindByID(id)
	if err != nil {
		log.Printf("ERROR: Failed to fetch song: %v", err)
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")
	if len(verses) == 0 {
		log.Println("DEBUG: Song text is empty")
		return []string{}, nil
	}

	start := pagination.Offset
	end := start + pagination.Limit
	if start >= len(verses) {
		log.Println("DEBUG: Offset exceeds number of verses")
		return []string{}, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	log.Printf("DEBUG: Returning verses %d to %d out of %d", start, end, len(verses))
	return verses[start:end], nil
}
