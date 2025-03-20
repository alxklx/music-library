package handlers

// AddSongRequest DTO для входных данных добавления песни
// @Description Структура для добавления новой песни
// SongRequest DTO для входных данных добавления песни
// @Description Структура для добавления новой песни (только группа и название)
type AddSongRequest struct {
	Group       string `json:"group" example:"Muse" description:"Название группы" validation:"requared"`
	Song        string `json:"song" example:"Supermassive Black Hole" description:"Название песни" validation:"requared"`
	ReleaseDate string `json:"releaseDate,omitempty" example:"16.07.2006" description:"Дата релиза (опционально)"`
	Text        string `json:"text,omitempty" example:"Ooh baby, don't you know I suffer?" description:"Текст песни (опционально)"`
	Link        string `json:"link,omitempty" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw" description:"Ссылка на песню (опционально)"`
}
