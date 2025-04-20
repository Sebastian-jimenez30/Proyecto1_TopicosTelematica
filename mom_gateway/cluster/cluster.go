package cluster

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"time"

	pb "mom_gateway/pb"

	"google.golang.org/grpc"
)

// Nodo representa un MOM individual
type Nodo struct {
	Nombre   string
	Cliente  pb.MomServiceClient
	Endpoint string
}

// Cluster contiene todos los nodos del sistema distribuido
type Cluster struct {
	Nodos []Nodo
}

// Crea un nuevo clúster conectándose a los tres MOMs
func NuevoCluster(mom1 string, mom2 string, mom3 string) *Cluster {
	endpoints := []string{
		mom1,
		mom2,
		mom3,
	}

	var nodos []Nodo
	for i, ep := range endpoints {
		conn, err := grpc.Dial(ep, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("❌ Error al conectar con MOM en %s: %v", ep, err)
		}
		nombre := fmt.Sprintf("mom%d", i+1)
		nodos = append(nodos, Nodo{
			Nombre:   nombre,
			Cliente:  pb.NewMomServiceClient(conn),
			Endpoint: ep,
		})
	}

	return &Cluster{Nodos: nodos}
}

// Determina el nodo responsable de un nombre (cola o tópico) usando hashing
func (c *Cluster) NodoResponsable(nombre string) Nodo {
	hash := sha1.Sum([]byte(nombre))
	idx := int(hash[0]) % len(c.Nodos)
	return c.Nodos[idx]
}

// Devuelve el siguiente nodo circularmente (para replicación simple)
func (c *Cluster) NodoSiguiente(actual Nodo) Nodo {
	for i, nodo := range c.Nodos {
		if nodo.Endpoint == actual.Endpoint {
			return c.Nodos[(i+1)%len(c.Nodos)]
		}
	}
	return c.Nodos[0]
}

// 🔐 Registra un usuario en todos los nodos (replicación fuerte)
func (c *Cluster) RegistrarUsuarioEnTodos(username, password string) (bool, string) {
	okCount := 0
	var lastMsg string

	for _, nodo := range c.Nodos {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.Credenciales{Username: username, Password: password}
		res, err := nodo.Cliente.RegistrarUsuario(ctx, req)
		if err != nil {
			log.Printf("⚠️ Error al registrar en %s: %v", nodo.Nombre, err)
			continue
		}

		if res.GetExito() {
			okCount++
		}
		lastMsg = res.GetMensaje()
	}

	return okCount > 0, lastMsg
}

// 🔑 Intenta autenticar al usuario en los nodos hasta que uno responda exitosamente
func (c *Cluster) AutenticarUsuario(username, password string) (bool, string) {
	for _, nodo := range c.Nodos {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req := &pb.Credenciales{Username: username, Password: password}
		res, err := nodo.Cliente.AutenticarUsuario(ctx, req)
		if err != nil {
			log.Printf("⚠️ Fallo en %s: %v", nodo.Nombre, err)
			continue
		}

		if res.GetToken() != "" {
			return true, res.GetToken()
		}
	}
	return false, ""
}

// 🧠 Replica el token en todos los nodos para permitir login global
func (c *Cluster) ReplicarTokenEnTodos(username, token, expiracion string) {
	for _, nodo := range c.Nodos {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := nodo.Cliente.GuardarTokenReplica(ctx, &pb.TokenConExpiracion{
			Username:   username,
			Token:      token,
			Expiracion: expiracion,
		})

		if err != nil {
			log.Printf("⚠️ Fallo al replicar token en %s: %v", nodo.Nombre, err)
		} else {
			log.Printf("✅ Token replicado en %s", nodo.Nombre)
		}
	}
}

// Devuelve todos los clientes gRPC disponibles (para usos especiales)
func (c *Cluster) TodosLosClientes() []pb.MomServiceClient {
	var clientes []pb.MomServiceClient
	for _, nodo := range c.Nodos {
		clientes = append(clientes, nodo.Cliente)
	}
	return clientes
}

// ReplicarEnNodosSiguientes aplica una operación gRPC en los 2 nodos siguientes del nodo responsable
func (c *Cluster) ReplicarEnNodosSiguientes(nombre string, replicador func(pb.MomServiceClient) error) {
	principal := c.NodoResponsable(nombre)
	secundario1 := c.NodoSiguiente(principal)
	secundario2 := c.NodoSiguiente(secundario1)

	nodos := []Nodo{secundario1, secundario2}

	for _, nodo := range nodos {
		go func(n Nodo) {
			if err := replicador(n.Cliente); err != nil {
				log.Printf("⚠️ Error replicando en %s: %v", n.Nombre, err)
			} else {
				log.Printf("✅ Replicado correctamente en %s", n.Nombre)
			}
		}(nodo)
	}
}
