package main

import (
	"encoding/json"
	"errors"
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
	"go-server/schemas"
)

var (
	db           *gorm.DB
	mySigningKey = []byte("secret")
)

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})

var ProductCardHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var productCards []models.ProductCard
    if err := db.Preload("User").Find(&productCards).Error; err != nil {
        http.Error(w, "Ошибка при получении товаров: " + err.Error(), http.StatusInternalServerError)
        return
    }

    payload, err := json.Marshal(productCards)
    if err != nil {
        http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
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

var GetCartHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Не удалось получить данные пользователя", http.StatusBadRequest)
        return
    }
    userID := uint(claims["user_id"].(float64))
    

    var cart models.Cart
    err := db.Preload("CartItems.ProductCard").Where("user_id = ?", userID).First(&cart).Error
    

    if errors.Is(err, gorm.ErrRecordNotFound) {
        cart = models.Cart{UserID: userID}
        if err := db.Create(&cart).Error; err != nil {
            http.Error(w, "Ошибка при создании корзины", http.StatusInternalServerError)
            return
        }
    } else if err != nil {
        http.Error(w, "Ошибка при получении корзины", http.StatusInternalServerError)
        return
    }
    

    var total int
    var items []schemas.CartItemResponse
    
    for _, item := range cart.CartItems {
        subtotal := item.ProductCard.Price * item.Quantity
        total += subtotal
        
        items = append(items, schemas.CartItemResponse{
            ID:          item.ID,
            ProductCard: item.ProductCard,
            Quantity:    item.Quantity,
            Subtotal:    subtotal,
        })
    }
    

    response := schemas.CartResponse{
        ID:     cart.ID,
        UserID: cart.UserID,
        Items:  items,
        Total:  total,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
})


var AddToCartHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Не удалось получить данные пользователя", http.StatusBadRequest)
        return
    }
    userID := uint(claims["user_id"].(float64))
    

    var req schemas.AddToCartRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }
    

    var cart models.Cart
    if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            cart = models.Cart{UserID: userID}
            if err := db.Create(&cart).Error; err != nil {
                http.Error(w, "Ошибка при создании корзины", http.StatusInternalServerError)
                return
            }
        } else {
            http.Error(w, "Ошибка при получении корзины", http.StatusInternalServerError)
            return
        }
    }
    

    var product models.ProductCard
    if err := db.First(&product, req.ProductCardID).Error; err != nil {
        http.Error(w, "Товар не найден", http.StatusNotFound)
        return
    }
    

    var cartItem models.CartItem
    result := db.Where("cart_id = ? AND product_card_id = ?", cart.ID, req.ProductCardID).First(&cartItem)
    
    if result.Error == nil {
        cartItem.Quantity += req.Quantity
        if err := db.Save(&cartItem).Error; err != nil {
            http.Error(w, "Ошибка при обновлении корзины", http.StatusInternalServerError)
            return
        }
    } else {
        cartItem = models.CartItem{
            CartID:        cart.ID,
            ProductCardID: req.ProductCardID,
            Quantity:      req.Quantity,
        }
        if err := db.Create(&cartItem).Error; err != nil {
            http.Error(w, "Ошибка при добавлении в корзину", http.StatusInternalServerError)
            return
        }
    }
    

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Товар добавлен в корзину"})
})


var RemoveFromCartHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Не удалось получить данные пользователя", http.StatusBadRequest)
        return
    }
    userID := uint(claims["user_id"].(float64))


    vars := mux.Vars(r)
    itemID := vars["id"]
    

    var cart models.Cart
    if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        http.Error(w, "Корзина не найдена", http.StatusNotFound)
        return
    }
    

    result := db.Where("cart_id = ? AND id = ?", cart.ID, itemID).Delete(&models.CartItem{})
    if result.Error != nil {
        http.Error(w, "Ошибка при удалении товара", http.StatusInternalServerError)
        return
    }
    
    if result.RowsAffected == 0 {
        http.Error(w, "Товар не найден в корзине", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Товар удалён из корзины"})
})


var UpdateCartItemHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Не удалось получить данные пользователя", http.StatusBadRequest)
        return
    }
    userID := uint(claims["user_id"].(float64))


    vars := mux.Vars(r)
    itemID := vars["id"]
    

    var req schemas.UpdateCartItemRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }
    

    var cart models.Cart
    if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        http.Error(w, "Корзина не найдена", http.StatusNotFound)
        return
    }
    

    var item models.CartItem
    result := db.Model(&item).
        Where("id = ? AND cart_id = ?", itemID, cart.ID).
        Update("quantity", req.Quantity)
    
    if result.Error != nil {
        http.Error(w, "Ошибка при обновлении количества", http.StatusInternalServerError)
        return
    }
    
    if result.RowsAffected == 0 {
        http.Error(w, "Товар не найден в корзине", http.StatusNotFound)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Количество обновлено"})
})

var ClearCartHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
    if !ok {
        http.Error(w, "Не удалось получить данные пользователя", http.StatusBadRequest)
        return
    }
    userID := uint(claims["user_id"].(float64))
    

    var cart models.Cart
    if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
        http.Error(w, "Корзина не найдена", http.StatusNotFound)
        return
    }
    

    result := db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
    if result.Error != nil {
        http.Error(w, "Ошибка при очистке корзины", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Корзина очищена"})
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

	// newUser := models.User{
	// 	Username: "Jack",
	// 	Email: "Alala@gmail.com",
	// 	Password: "ololo123",
	// }

	// if err := db.Create(&newUser).Error; err != nil {
	// 	log.Fatal(":(")
	// }

	newProductCard := models.ProductCard{
		UserID: 1,
		Brand:   "BMW",
		BikeModel: "R1200 GS",
		EngineCapacity: 1200,
		Power: 186,
		Color: "White",
		Price: 1390000,
	}
	
	if err := db.Create(&newProductCard).Error; err != nil {
		log.Fatal(":(")
	}

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/status", StatusHandler).Methods("GET")
	r.Handle("/get-token", http.HandlerFunc(GetTokenHandler)).Methods("GET")
	r.Handle("/products", jwtMiddleware.Handler(ProductCardHandler)).Methods("GET")

	secured := r.PathPrefix("/").Subrouter()
    secured.Use(jwtMiddleware.Handler)

	secured.HandleFunc("/cart", GetCartHandler).Methods("GET")
    secured.HandleFunc("/cart/add", AddToCartHandler).Methods("POST")
	secured.HandleFunc("/cart/items/{id}", RemoveFromCartHandler).Methods("DELETE")
	secured.HandleFunc("/cart/items/{id}", UpdateCartItemHandler).Methods("PUT")
	secured.HandleFunc("/cart/clear", ClearCartHandler).Methods("DELETE")

	log.Println("Server starting on :3000")
	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}