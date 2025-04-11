package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/db" // Importar el paquete db
	"github.com/gorilla/mux"
)

// Estructura para el Tópico
type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Crear un nuevo tópico en la base de datos
func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	var newTopic Topic

	// Parsear el cuerpo JSON a la estructura Topic
	if err := json.NewDecoder(r.Body).Decode(&newTopic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Verificar si el tópico ya existe en la base de datos
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM topics WHERE name=$1)", newTopic.Name).Scan(&exists)
	if err != nil {
		http.Error(w, "Error al verificar el tópico", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Tópico ya existe", http.StatusConflict)
		return
	}

	// Insertar el nuevo tópico en la base de datos
	_, err = db.DB.Exec("INSERT INTO topics (name, description) VALUES ($1, $2)", newTopic.Name, newTopic.Description)
	if err != nil {
		http.Error(w, "Error al crear el tópico", http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Tópico '%s' creado con éxito", newTopic.Name)
}

// Listar todos los tópicos de la base de datos
func ListTopicsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT name, description FROM topics")
	if err != nil {
		http.Error(w, "Error al obtener los tópicos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var topic Topic
		if err := rows.Scan(&topic.Name, &topic.Description); err != nil {
			http.Error(w, "Error al leer los tópicos", http.StatusInternalServerError)
			return
		}
		topics = append(topics, topic)
	}

	// Devolver la lista de tópicos
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

// Eliminar un tópico de la base de datos
func DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["name"] // Obtener el parámetro de la ruta {name}

	// Verificar si el tópico existe
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM topics WHERE name=$1)", topicName).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Tópico no encontrado", http.StatusNotFound)
		return
	}

	// Eliminar el tópico de la base de datos
	_, err = db.DB.Exec("DELETE FROM topics WHERE name=$1", topicName)
	if err != nil {
		http.Error(w, "Error al eliminar el tópico", http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Tópico '%s' eliminado con éxito", topicName)
}
