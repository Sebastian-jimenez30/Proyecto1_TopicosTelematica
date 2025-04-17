package cluster

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"sync"

	pb "mom_gateway/pb"

	"google.golang.org/grpc"
)

type Nodo struct {
	ID     string `json:"id"`
	Host   string `json:"host"`
	Puerto int    `json:"puerto"`
}

type Cluster interface {
	SeleccionarPrimario(nombre string) Nodo
	SeleccionarSecundario(primario Nodo) Nodo
	GetCliente(id string) pb.MomServiceClient
}

type ClusterManager struct {
	Clientes map[string]pb.MomServiceClient
	Nodos    []Nodo
	mu       sync.Mutex
}

func CargarCluster(path string) (*ClusterManager, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config struct {
		Nodos []Nodo `json:"nodos"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	cm := &ClusterManager{
		Clientes: make(map[string]pb.MomServiceClient),
		Nodos:    config.Nodos,
	}

	for _, nodo := range config.Nodos {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", nodo.Host, nodo.Puerto), grpc.WithInsecure())
		if err != nil {
			log.Printf("⚠️ Error conectando a %s: %v", nodo.ID, err)
			continue
		}
		cm.Clientes[nodo.ID] = pb.NewMomServiceClient(conn)
	}

	return cm, nil
}

func (cm *ClusterManager) SeleccionarPrimario(nombre string) Nodo {
	hash := crc32.ChecksumIEEE([]byte(nombre))
	index := int(hash) % len(cm.Nodos)
	return cm.Nodos[index]
}

func (cm *ClusterManager) SeleccionarSecundario(primario Nodo) Nodo {
	for i, n := range cm.Nodos {
		if n.ID == primario.ID {
			return cm.Nodos[(i+1)%len(cm.Nodos)]
		}
	}
	return primario
}

func (cm *ClusterManager) GetCliente(id string) pb.MomServiceClient {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.Clientes[id]
}
