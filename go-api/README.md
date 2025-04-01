

This directory contains the RESTful API implemented in Go. It handles client interaction, authentication, and forwarding of operations to the MOM nodes via gRPC.

## Contents
- `main.go`: Entry point of the REST API server.
- `routes/`: Defines HTTP routes for authentication, topics, queues, and messaging.
- `middleware/`: JWT validation and request filtering.
- `grpc/`: gRPC clients used to communicate with the C++ MOM nodes.
- `config/`: Environment and configuration loading.
- `proto/`: Shared `.proto` definitions used to generate gRPC clients.

The API exposes endpoints for login, managing topics and queues, sending and receiving messages, and ensures all operations are tied to authenticated users.
