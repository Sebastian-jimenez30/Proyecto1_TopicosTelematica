package handlers

import (
	"context"
	"net/http"
	"time"

	"mom_gateway/cluster"
	pb "mom_gateway/pb"

	"github.com/gin-gonic/gin"
)

func CrearColaHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		nodo := cl.NodoResponsable(body.Nombre)
		req := &pb.AccionConToken{Token: token, Nombre: body.Nombre}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.CrearCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear cola"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(body.Nombre, func(client pb.MomServiceClient) error {
				_, err := client.CrearCola(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func EliminarColaHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		nodo := cl.NodoResponsable(nombreCola)
		req := &pb.AccionConToken{Token: token, Nombre: nombreCola}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.EliminarCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar cola"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombreCola, func(client pb.MomServiceClient) error {
				_, err := client.EliminarCola(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func AutorizarColaHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		nodo := cl.NodoResponsable(nombreCola)
		req := &pb.AutorizacionColaRequest{
			Token:           token,
			Nombre:          nombreCola,
			UsuarioObjetivo: body.UsuarioObjetivo,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := nodo.Cliente.AutorizarUsuarioCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al autorizar usuario"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombreCola, func(client pb.MomServiceClient) error {
				_, err := client.AutorizarUsuarioCola(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": res.Mensaje})
		}
	}
}

func EnviarMensajeColaHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		principal := cl.NodoResponsable(nombreCola)
		req := &pb.MensajeConToken{
			Token:     token,
			Nombre:    nombreCola,
			Contenido: body.Contenido,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := principal.Cliente.EnviarMensajeCola(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al enviar mensaje"})
			return
		}

		if res.Exito {
			go cl.ReplicarEnNodosSiguientes(nombreCola, func(client pb.MomServiceClient) error {
				_, err := client.EnviarMensajeCola(context.Background(), req)
				return err
			})
			c.JSON(http.StatusOK, gin.H{"mensaje": res.Mensaje})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": res.Mensaje})
		}
	}
}

func ConsumirColaHandler(cl *cluster.Cluster) gin.HandlerFunc {
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

		principal := cl.NodoResponsable(nombreCola)
		replica := cl.NodoSiguiente(principal)
		replica2 := cl.NodoSiguiente(replica)

		req := &pb.AccionConToken{Token: token, Nombre: nombreCola}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := principal.Cliente.ConsumirMensajeCola(ctx, req)

		if err != nil || res.GetContenido() == "" {
			ctxFallback, cancelFallback := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFallback()
			res, err = replica.Cliente.ConsumirMensajeCola(ctxFallback, req)
			res, err = replica2.Cliente.ConsumirMensajeCola(ctxFallback, req)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consumir mensaje"})
			return
		}

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

func ListarColasHandler(cl *cluster.Cluster) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ExtraerToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			return
		}

		principal := cl.NodoResponsable("colas")
		replica := cl.NodoSiguiente(principal)

		req := &pb.Token{Token: token}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		res, err := principal.Cliente.ListarColas(ctx, req)
		if err != nil || len(res.GetNombres()) == 0 {
			ctxFallback, cancelFallback := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFallback()
			res, err = replica.Cliente.ListarColas(ctxFallback, req)
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar colas"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"colas": res.Nombres})
	}
}
