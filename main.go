package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-server/models"
)

var (
	db           *gorm.DB
	mySigningKey = []byte("secret")
)

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

var ExamPreparationCourseHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var ExamPreparationCourse []models.ExamPreparationCourse
    if err := db.Preload("User").Find(&ExamPreparationCourse).Error; err != nil {
        http.Error(w, "Failed to get course", http.StatusInternalServerError)
        return
    }

    payload, err := json.Marshal(ExamPreparationCourse)
    if err != nil {
        http.Error(w, "Failed to marshal course", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(payload)
})

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var user models.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	newExamPreparationCourse := models.ExamPreparationCourse{
		Subject:   "Математика",
		UserID:  3,
		Relevance: 2026,
		NumberOfClasses: 65,
		ContactTheTeacher: true,
		Individuality: false,
		Price: 48930,
	}
	
	if err := db.Create(&newExamPreparationCourse).Error; err != nil {
		log.Fatal(":(")
	}

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/get-token", http.HandlerFunc(GetTokenHandler)).Methods("GET")
	r.Handle("/courses", jwtMiddleware.Handler(ExamPreparationCourseHandler)).Methods("GET")

	log.Println("Server starting on :3000")
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}