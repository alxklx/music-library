package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/alxklx/music-library/docs" // Импорт сгенерированных Swagger-файлов
	"github.com/alxklx/music-library/internal/config"
	"github.com/alxklx/music-library/internal/handlers"
	"github.com/alxklx/music-library/internal/repository"
	"github.com/alxklx/music-library/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

func applyMigrations(db *pgxpool.Pool) error {
	log.Println("INFO: Applying migrations...")
	migrationDir := filepath.Join(".", "migrations")
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Printf("ERROR: Failed to read migration directory: %v", err)
		return err
	}
	log.Printf("INFO: Found %d files in migration directory", len(files))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			log.Printf("INFO: Applying migration %s", file.Name())
			content, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
			if err != nil {
				log.Printf("ERROR: Failed to read migration file %s: %v", file.Name(), err)
				return err
			}
			_, err = db.Exec(context.Background(), string(content))
			if err != nil {
				log.Printf("ERROR: Failed to apply migration %s: %v", file.Name(), err)
				return err
			}
			log.Printf("INFO: Migration %s applied successfully", file.Name())
		}
	}
	return nil
}

// @title Music Library API
// @version 1.0
// @description API для управления онлайн-библиотекой песен
// @host localhost:8080
// @BasePath /
func main() {
	log.Println("INFO: Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ERROR: Failed to load config: %v", err)
	}
	log.Printf("DEBUG: Config loaded: %+v", cfg)

	log.Println("INFO: Connecting to database...")
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("ERROR: Failed to connect to database: %v", err)
	}
	defer db.Close()
	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("ERROR: Failed to ping the database: %v", err)
	}
	log.Println("INFO: Database connected successfully")

	if err := applyMigrations(db); err != nil {
		log.Fatalf("ERROR: Failed to apply migrations: %v", err)
	}

	repo := repository.NewPostgresRepo(db)
	usecase := usecase.NewSongUsecase(repo, cfg.APIEndpoint)
	handler := handlers.NewSongHandler(usecase)

	log.Println("INFO: Setting up router...")
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Добавляем Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("INFO: Server starting on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("ERROR: Server failed: %v", err)
	}
}
