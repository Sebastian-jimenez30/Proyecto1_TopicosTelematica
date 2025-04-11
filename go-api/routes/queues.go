package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/db"
	"github.com/gorilla/mux"
)

// Estructura para la Cola
type Queue struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	CreatedAt   string `json:"created_at"`
}

// Crear una nueva cola
func CreateQueueHandler(w http.ResponseWriter, r *http.Request) {
	var newQueue Queue

	// Parsear el cuerpo JSON a la estructura Queue
	if err := json.NewDecoder(r.Body).Decode(&newQueue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar que el user_id existe en la base de datos
	var userExists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)", newQueue.UserID).Scan(&userExists)
	if err != nil || !userExists {
		http.Error(w, "El user_id proporcionado no existe", http.StatusBadRequest)
		return
	}

	// Validar si la cola ya existe
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM queues WHERE name=$1 AND user_id=$2)", newQueue.Name, newQueue.UserID).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al verificar la cola: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "La cola ya existe", http.StatusConflict)
		return
	}

	// Insertar la nueva cola en la base de datos
	_, err = db.DB.Exec("INSERT INTO queues (name, description, user_id) VALUES ($1, $2, $3)", newQueue.Name, newQueue.Description, newQueue.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al crear la cola: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Cola '%s' creada con éxito", newQueue.Name)
}

// Listar todas las colas de un usuario
func ListQueuesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id") // Obtener el ID del usuario de la query string

	// Consultar las colas asociadas al usuario
	rows, err := db.DB.Query("SELECT name, description, user_id, created_at FROM queues WHERE user_id=$1", userID)
	if err != nil {
		http.Error(w, "Error al obtener las colas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var queues []Queue
	for rows.Next() {
		var queue Queue
		if err := rows.Scan(&queue.Name, &queue.Description, &queue.UserID, &queue.CreatedAt); err != nil {
			http.Error(w, "Error al leer las colas", http.StatusInternalServerError)
			return
		}
		queues = append(queues, queue)
	}

	// Devolver la lista de colas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(queues)
}

// Eliminar una cola
func DeleteQueueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueName := vars["name"] // Obtener el nombre de la cola de la ruta

	// Verificar si la cola existe
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM queues WHERE name=$1)", queueName).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Cola no encontrada", http.StatusNotFound)
		return
	}

	// Eliminar la cola de la base de datos
	_, err = db.DB.Exec("DELETE FROM queues WHERE name=$1", queueName)
	if err != nil {
		http.Error(w, "Error al eliminar la cola", http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cola '%s' eliminada con éxito", queueName)
}
