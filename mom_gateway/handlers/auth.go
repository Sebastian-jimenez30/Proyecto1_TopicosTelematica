package handlers

import (
	"context"
	"net/http"
	"time"

	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := client.RegistrarUsuario(ctx, &pb.Credenciales{
			Username: body.Username,
			Password: body.Password,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al registrar usuario"})
			return
		}

		if res.Exito {
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func LoginHandler(client pb.MomServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := client.AutenticarUsuario(ctx, &pb.Credenciales{
			Username: body.Username,
			Password: body.Password,
		})

		if err != nil || res.Token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": res.Token})
	}
}
