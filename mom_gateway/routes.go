package main

import (
	"mom_gateway/handlers"
	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(client pb.MomServiceClient) *gin.Engine {
	r := gin.Default()

	// Rutas de autenticación
	r.POST("/register", handlers.RegisterHandler(client))
	r.POST("/login", handlers.LoginHandler(client))
	r.POST("/colas", handlers.CrearColaHandler(client))
	r.DELETE("/colas/:nombre", handlers.EliminarColaHandler(client))
	r.POST("/colas/:nombre/autorizar", handlers.AutorizarColaHandler(client))
	r.POST("/colas/:nombre/enviar", handlers.EnviarMensajeColaHandler(client))
	r.GET("/colas/:nombre/consumir", handlers.ConsumirColaHandler(client))
	r.POST("/topicos", handlers.CrearTopicoHandler(client))
	r.DELETE("/topicos/:nombre", handlers.EliminarTopicoHandler(client))
	r.POST("/topicos/:nombre/suscribir", handlers.SuscribirseTopicoHandler(client))
	r.POST("/topicos/:nombre/publicar", handlers.PublicarTopicoHandler(client))
	r.GET("/topicos/:nombre/consumir", handlers.ConsumirTopicoHandler(client))
	r.GET("/colas", handlers.ListarColasHandler(client))
	r.GET("/topicos", handlers.ListarTopicosHandler(client))

	return r
}
