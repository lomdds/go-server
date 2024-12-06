package main

import (
	"database/sql"
	"fmt"
	"log"
	// "net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "your_user"
	password = "your_password"
	dbname = "your_name"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Ошибка Подключения к бд: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при проверке подключения: %v", err)
	}

	r := gin.Default()

	r.Run(":8080")
}