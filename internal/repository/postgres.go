package repository

import (
	"context"
	"strconv"

	"github.com/alxklx/music-library/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(song *domain.Song) error {
	query := `INSERT INTO songs ("group", song, release_date, text, link) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(context.Background(), query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
}

func (r *PostgresRepo) Update(song *domain.Song) error {
	query := `UPDATE songs SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5 WHERE id = $6`
	_, err := r.db.Exec(context.Background(), query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	return err
}

func (r *PostgresRepo) Delete(id int64) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *PostgresRepo) FindByID(id int64) (*domain.Song, error) {
	song := &domain.Song{}
	query := `SELECT id, "group", song, release_date, text, link FROM songs WHERE id = $1`
	err := r.db.QueryRow(context.Background(), query, id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (r *PostgresRepo) FindAll(filter domain.SongFilter) ([]*domain.Song, error) {
	query := `SELECT id, "group", song, release_date, text, link FROM songs WHERE 1=1`
	args := []interface{}{}
	i := 1

	if filter.Group != "" {
		query += " AND \"group\" = $" + strconv.Itoa(i)
		args = append(args, filter.Group)
		i++
	}
	if filter.Song != "" {
		query += " AND song = $" + strconv.Itoa(i)
		args = append(args, filter.Song)
		i++
	}
	if filter.ReleaseDate != "" {
		query += " AND release_date = $" + strconv.Itoa(i)
		args = append(args, filter.ReleaseDate)
		i++
	}

	query += " LIMIT $" + strconv.Itoa(i) + " OFFSET $" + strconv.Itoa(i+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []*domain.Song
	for rows.Next() {
		song := &domain.Song{}
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}
