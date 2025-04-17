package main

import (
	"fmt"
)

func main() {
	fmt.Println("🔌 Inicializando conexión con el MOM vía gRPC...")
	if err := InitGRPCClient(); err != nil {
		panic(fmt.Sprintf("Error al conectar con gRPC: %v", err))
	}

	fmt.Println("🌐 Servidor REST escuchando en http://localhost:8080")
	router := SetupRoutes(grpcClient)
	router.Run(":8080")
}
