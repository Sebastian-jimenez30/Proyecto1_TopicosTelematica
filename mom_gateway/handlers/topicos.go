package handlers

import (
	"context"
	"net/http"
	"time"

	"mom_gateway/cluster"
	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func CrearTopicoHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		nodo := cl.NodoResponsable(body.Nombre)
		req := &pb.AccionConToken{Token: token, Nombre: body.Nombre}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.CrearTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear tópico"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(body.Nombre, func(client pb.MomServiceClient) error {
				_, err := client.CrearTopico(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func EliminarTopicoHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombre := c.Param("nombre")
		token := ExtraerToken(c)
		if nombre == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
			return
		}

		nodo := cl.NodoResponsable(nombre)
		req := &pb.AccionConToken{Token: token, Nombre: nombre}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.EliminarTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar tópico"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombre, func(client pb.MomServiceClient) error {
				_, err := client.EliminarTopico(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func SuscribirseTopicoHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombre := c.Param("nombre")
		token := ExtraerToken(c)
		if nombre == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros incompletos"})
			return
		}

		nodo := cl.NodoResponsable(nombre)
		req := &pb.AccionConToken{Token: token, Nombre: nombre}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.SuscribirseTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al suscribirse"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombre, func(client pb.MomServiceClient) error {
				_, err := client.SuscribirseTopico(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func PublicarTopicoHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombre := c.Param("nombre")
		token := ExtraerToken(c)

		var body struct {
			Contenido string `json:"contenido"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || nombre == "" || body.Contenido == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incompletos"})
			return
		}

		principal := cl.NodoResponsable(nombre)
		req := &pb.MensajeConToken{
			Token:     token,
			Nombre:    nombre,
			Contenido: body.Contenido,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := principal.Cliente.PublicarEnTopico(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al publicar mensaje"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombre, func(client pb.MomServiceClient) error {
				_, err := client.PublicarEnTopico(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func ConsumirTopicoHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		nombre := c.Param("nombre")
		token := ExtraerToken(c)
		if nombre == "" || token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incompletos"})
			return
		}

		principal := cl.NodoResponsable(nombre)
		replica := cl.NodoSiguiente(principal)

		req := &pb.AccionConToken{Token: token, Nombre: nombre}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := principal.Cliente.ConsumirDesdeTopico(ctx, req)
		if err != nil || len(res.Mensajes) == 0 {
			ctxFallback, cancelFallback := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFallback()
			res, err = replica.Cliente.ConsumirDesdeTopico(ctxFallback, req)
		}

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

		c.JSON(http.StatusOK, gin.H{"mensajes": mensajes})
	}
}

func ListarTopicosHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		principal := cl.NodoResponsable("topicos")
		replica := cl.NodoSiguiente(principal)

		req := &pb.Token{Token: token}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := principal.Cliente.ListarTopicos(ctx, req)
		if err != nil || len(res.Nombres) == 0 {
			ctxFallback, cancelFallback := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFallback()
			res, err = replica.Cliente.ListarTopicos(ctxFallback, req)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar tópicos"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"topicos": res.Nombres})
	}
}
