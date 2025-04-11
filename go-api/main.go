package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/db" // Importar el paquete db
	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/middleware"
	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Inicializar la conexión a la base de datos
	db.InitDB()
}

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Configura el enrutador
	r := mux.NewRouter()

	// Ruta para login (sin protección)
	r.HandleFunc("/auth/login", middleware.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/register", middleware.RegisterHandler).Methods("POST")

	// Rutas para colas
	r.HandleFunc("/queues", routes.CreateQueueHandler).Methods("POST")
	r.HandleFunc("/queues", routes.ListQueuesHandler).Methods("GET")
	r.HandleFunc("/queues/{name}", routes.DeleteQueueHandler).Methods("DELETE")

	// Rutas para manejar los mensajes en colas
	r.HandleFunc("/queues/{name}/messages", routes.SendMessageToQueueHandler).Methods("POST")
	r.HandleFunc("/queues/{name}/messages", routes.GetMessagesFromQueueHandler).Methods("GET")
	r.HandleFunc("/messages/{id}", routes.UpdateMessageStatusHandler).Methods("PUT")
	r.HandleFunc("/messages/{id}", routes.DeleteMessageHandler).Methods("DELETE")

	// Rutas protegidas para tópicos
	protectedRoutes := r.PathPrefix("/topics").Subrouter()
	protectedRoutes.Use(middleware.TokenValidationMiddleware) // Aplica el middleware de validación de JWT solo a estas rutas
	protectedRoutes.HandleFunc("", routes.CreateTopicHandler).Methods("POST")
	protectedRoutes.HandleFunc("", routes.ListTopicsHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{name}", routes.DeleteTopicHandler).Methods("DELETE")

	// Rutas para manejar los mensajes en tópicos
	protectedRoutes.HandleFunc("/{name}/messages", routes.SendMessageToTopicHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{name}/messages", routes.GetMessagesFromTopicHandler).Methods("GET")

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Valor por defecto si no se especifica en el .env
	}
	fmt.Printf("Servidor escuchando en el puerto %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
