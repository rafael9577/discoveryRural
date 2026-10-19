package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	adminKey := os.Getenv("ADMIN_KEY")

	registrarAnimais(adminKey)
	registrarUsuarios(adminKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor rodando na porta " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
