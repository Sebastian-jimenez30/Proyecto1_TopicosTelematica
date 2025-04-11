package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/db" // Paquete para manejar la base de datos
	"github.com/gorilla/mux"
)

// Estructura para el Mensaje
type Message struct {
	Content   string `json:"content"`
	Status    string `json:"status"`
	UserID    int    `json:"user_id"`
	QueueID   int    `json:"queue_id"`
	TopicID   int    `json:"topic_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// Enviar un mensaje a una cola
func SendMessageToQueueHandler(w http.ResponseWriter, r *http.Request) {
	var newMessage Message

	// Parsear el cuerpo JSON a la estructura Message
	if err := json.NewDecoder(r.Body).Decode(&newMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el ID de la cola de la ruta
	vars := mux.Vars(r)
	queueName := vars["name"]

	// Verificar si la cola existe
	var queueID int
	err := db.DB.QueryRow("SELECT id FROM queues WHERE name=$1", queueName).Scan(&queueID)
	if err != nil {
		http.Error(w, "Cola no encontrada", http.StatusNotFound)
		return
	}

	// Obtener el ID del usuario desde el JWT (suponiendo que el JWT contiene esta información)
	// Si no es así, asegúrate de obtenerlo de la solicitud o del JWT
	userID := newMessage.UserID

	// Insertar el nuevo mensaje en la base de datos
	_, err = db.DB.Exec("INSERT INTO messages (content, user_id, queue_id, status) VALUES ($1, $2, $3, $4)",
		newMessage.Content, userID, queueID, "pending")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al crear el mensaje: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Mensaje enviado a la cola '%s' con éxito", queueName)
}

// Enviar un mensaje a un tópico
func SendMessageToTopicHandler(w http.ResponseWriter, r *http.Request) {
	var newMessage Message

	// Parsear el cuerpo JSON a la estructura Message
	if err := json.NewDecoder(r.Body).Decode(&newMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener el ID del tópico de la ruta
	vars := mux.Vars(r)
	topicName := vars["name"]

	// Verificar si el tópico existe
	var topicID int
	err := db.DB.QueryRow("SELECT id FROM topics WHERE name=$1", topicName).Scan(&topicID)
	if err != nil {
		http.Error(w, "Tópico no encontrado", http.StatusNotFound)
		return
	}

	// Obtener el ID del usuario desde el JWT
	userID := newMessage.UserID // O el ID extraído del JWT

	// Insertar el nuevo mensaje en la base de datos
	_, err = db.DB.Exec("INSERT INTO messages (content, user_id, topic_id, status) VALUES ($1, $2, $3, $4)",
		newMessage.Content, userID, topicID, "pending")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al crear el mensaje: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Mensaje enviado al tópico '%s' con éxito", topicName)
}

// Obtener los mensajes de una cola
func GetMessagesFromQueueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queueName := vars["name"]

	// Verificar si la cola existe
	var queueID int
	err := db.DB.QueryRow("SELECT id FROM queues WHERE name=$1", queueName).Scan(&queueID)
	if err != nil {
		http.Error(w, "Cola no encontrada", http.StatusNotFound)
		return
	}

	// Obtener los mensajes de la cola
	rows, err := db.DB.Query("SELECT id, content, status, user_id, created_at FROM messages WHERE queue_id=$1", queueID)
	if err != nil {
		http.Error(w, "Error al obtener los mensajes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.QueueID, &message.Content, &message.Status, &message.UserID, &message.CreatedAt); err != nil {
			http.Error(w, "Error al leer los mensajes", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Devolver la lista de mensajes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Obtener los mensajes de un tópico
func GetMessagesFromTopicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicName := vars["name"]

	// Verificar si el tópico existe
	var topicID int
	err := db.DB.QueryRow("SELECT id FROM topics WHERE name=$1", topicName).Scan(&topicID)
	if err != nil {
		http.Error(w, "Tópico no encontrado", http.StatusNotFound)
		return
	}

	// Obtener los mensajes del tópico
	rows, err := db.DB.Query("SELECT topic_id, content, status, user_id, created_at FROM messages WHERE topic_id=$1", topicID)
	if err != nil {
		http.Error(w, "Error al obtener los mensajes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.TopicID, &message.Content, &message.Status, &message.UserID, &message.CreatedAt); err != nil {
			http.Error(w, "Error al leer los mensajes", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	// Devolver la lista de mensajes
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Marcar un mensaje como consumido
func UpdateMessageStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	// Actualizar el estado del mensaje a 'consumed'
	_, err := db.DB.Exec("UPDATE messages SET status='consumed' WHERE id=$1", messageID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al actualizar el estado del mensaje: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mensaje con ID '%s' marcado como consumido", messageID)
}

// Eliminar un mensaje
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	// Verificar si el mensaje existe
	_, err := db.DB.Exec("DELETE FROM messages WHERE id=$1", messageID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al eliminar el mensaje: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Mensaje con ID '%s' eliminado con éxito", messageID)
}
