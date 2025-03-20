Для запуска приложения необходимо создать конфигурационный файл .env со следующим содержимым:

DATABASE_URL=postgres://user:password@host:port/database?search_path=schema
PORT=8080
API_ENDPOINT=http://localhost:8081/info