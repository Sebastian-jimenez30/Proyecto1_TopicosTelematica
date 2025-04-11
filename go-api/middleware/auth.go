package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/go-api/db"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Función para registrar un usuario
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser Credentials

	// Parsear el cuerpo JSON a las credenciales
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cifrar la contraseña antes de guardarla
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al cifrar la contraseña", http.StatusInternalServerError)
		return
	}

	// Guardar el usuario en la base de datos
	_, err = db.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", newUser.Username, hash)
	if err != nil {
		http.Error(w, "Error al registrar el usuario", http.StatusInternalServerError)
		return
	}

	// Devolver una respuesta exitosa
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuario '%s' registrado con éxito", newUser.Username)
}

// Función para hacer login de un usuario
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// Parsear el cuerpo JSON a las credenciales
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener la contraseña cifrada desde la base de datos
	var storedHash string
	err := db.DB.QueryRow("SELECT password_hash FROM users WHERE username=$1", creds.Username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	// Comparar la contraseña proporcionada con la almacenada
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	// Crear el token JWT
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = creds.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Expiración en 72 horas

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString([]byte("mysecretkey"))
	if err != nil {
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	}

	// Enviar el token como respuesta
	w.Write([]byte(tokenString))
}
