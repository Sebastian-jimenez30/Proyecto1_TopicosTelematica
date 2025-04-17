package main

import (
	"fmt"
	"mom_gateway/cluster"
)

func main() {
	fmt.Println("🔌 Inicializando conexiones con los nodos MOM vía gRPC...")
	cl := cluster.NuevoCluster()

	fmt.Println("🌐 Servidor REST escuchando en http://localhost:8080")
	router := SetupRoutes(cl)
	router.Run(":8080")
}
