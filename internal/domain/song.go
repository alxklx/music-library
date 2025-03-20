package domain

type Song struct {
	ID          int64  `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongFilter struct {
	Group       string
	Song        string
	ReleaseDate string
	Limit       int
	Offset      int
}

type VersePagination struct {
	Limit  int
	Offset int
}

type SongRepository interface {
	Create(song *Song) error
	Update(song *Song) error
	Delete(id int64) error
	FindByID(id int64) (*Song, error)
	FindAll(filter SongFilter) ([]*Song, error)
}
