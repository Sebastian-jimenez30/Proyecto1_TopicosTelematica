package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Importar el driver PostgreSQL
)

var DB *sql.DB

// Inicializar la conexión a la base de datos
func InitDB() {
	connStr := fmt.Sprintf("postgres://user:password@localhost:5432/messaging_db?sslmode=disable")
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	// Verificar la conexión
	if err = DB.Ping(); err != nil {
		log.Fatal("Error al verificar la conexión a la base de datos:", err)
	}
}
