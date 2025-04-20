package main

import (
	"fmt"
	"mom_gateway/cluster"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso: mom_gateway <nodo1> <nodo2> <nodo3>")
		return
	}

	nodo1 := os.Args[1]
	nodo2 := os.Args[2]
	nodo3 := os.Args[3]

	fmt.Println("🔌 Inicializando conexiones con los nodos MOM vía gRPC...")
	cl := cluster.NuevoCluster(nodo1, nodo2, nodo3)

	fmt.Println("🌐 Servidor REST escuchando en http://localhost:8080")
	router := SetupRoutes(cl)
	router.Run(":8080")
}
