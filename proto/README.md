

This folder contains the shared `.proto` definitions used to define the gRPC services and message formats exchanged between Go and C++ components.

These definitions are used to generate:
- gRPC client stubs in Go (`go-api/grpc`)
- gRPC server interfaces in C++ (`cpp-mom/grpc_server`)

Ensure all services, messages, and data structures required for node-to-node and API-to-node communication are described here.
