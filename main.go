package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/services"
	"kasir-api/repositories"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port 	string `mapstructure:"PORT"`
	DBConn	string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig(); 
	}

	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal dalam menginisialisasi database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)

	productService := services.NewProductService(productRepo)

	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Jalan",
		})
	})

	fmt.Println("Server running di localhost:8080")
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Gagal running server")
	}
}
