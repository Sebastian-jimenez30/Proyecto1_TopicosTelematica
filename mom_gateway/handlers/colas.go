package handlers

import (
	"context"
	"net/http"
	"time"

	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func CrearColaHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Nombre string `json:"nombre"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Nombre == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de cola requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.AccionConToken{
			Token:  token,
			Nombre: body.Nombre,
		}
		res, err := client.CrearCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cola"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func EliminarColaHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreCola := c.Param("nombre")
		if nombreCola == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de cola requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.AccionConToken{
			Token:  token,
			Nombre: nombreCola,
		}
		res, err := client.EliminarCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cola"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func AutorizarColaHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreCola := c.Param("nombre")
		if nombreCola == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de cola requerido"})
			return
		}

		var body struct {
			UsuarioObjetivo string `json:"usuario"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.UsuarioObjetivo == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Usuario objetivo requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.AutorizacionColaRequest{
			Token:           token,
			Nombre:          nombreCola,
			UsuarioObjetivo: body.UsuarioObjetivo,
		}

		res, err := client.AutorizarUsuarioCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al autorizar usuario"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func EnviarMensajeColaHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreCola := c.Param("nombre")
		if nombreCola == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de cola requerido"})
			return
		}

		var body struct {
			Contenido string `json:"contenido"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Contenido == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Contenido requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.MensajeConToken{
			Token:     token,
			Nombre:    nombreCola,
			Contenido: body.Contenido,
		}

		res, err := client.EnviarMensajeCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar mensaje"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func ConsumirColaHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreCola := c.Param("nombre")
		if nombreCola == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de cola requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.AccionConToken{
			Token:  token,
			Nombre: nombreCola,
		}

		res, err := client.ConsumirMensajeCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consumir mensaje"})
			return
		}

		// Si no hay contenido, devolvemos 204 No Content
		if res.Contenido == "" {
			c.JSON(http.StatusNoContent, gin.H{})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"remitente": res.Remitente,
			"contenido": res.Contenido,
			"canal":     res.Canal,
			"timestamp": res.Timestamp,
		})
	}
}

func ListarColasHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.Token{Token: token}
		res, err := client.ListarColas(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar colas"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"colas": res.Nombres})
	}
}
