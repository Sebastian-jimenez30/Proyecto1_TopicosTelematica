package handlers

import (
	"context"
	"net/http"
	"time"

	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func CrearTopicoHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Nombre string `json:"nombre"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Nombre == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre de tópico requerido"})
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
		res, err := client.CrearTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tópico"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func EliminarTopicoHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreTopico := c.Param("nombre")
		if nombreTopico == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre del tópico requerido"})
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
			Nombre: nombreTopico,
		}

		res, err := client.EliminarTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar tópico"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func SuscribirseTopicoHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreTopico := c.Param("nombre")
		if nombreTopico == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre del tópico requerido"})
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
			Nombre: nombreTopico,
		}

		res, err := client.SuscribirseTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al suscribirse al tópico"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func PublicarTopicoHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreTopico := c.Param("nombre")
		if nombreTopico == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre del tópico requerido"})
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
			Nombre:    nombreTopico,
			Contenido: body.Contenido,
		}

		res, err := client.PublicarEnTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al publicar mensaje"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func ConsumirTopicoHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombreTopico := c.Param("nombre")
		if nombreTopico == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre del tópico requerido"})
			return
		}

		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.AccionConToken{
			Token:  token,
			Nombre: nombreTopico,
		}

		res, err := client.ConsumirDesdeTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consumir mensajes"})
			return
		}

		var mensajes []gin.H
		for _, m := range res.Mensajes {
			mensajes = append(mensajes, gin.H{
				"remitente": m.Remitente,
				"contenido": m.Contenido,
				"canal":     m.Canal,
				"timestamp": m.Timestamp,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"mensajes": mensajes,
		})
	}
}

func ListarTopicosHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.Token{Token: token}
		res, err := client.ListarTopicos(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar tópicos"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"topicos": res.Nombres})
	}
}
