# MOM Project Structure

This document explains the folder structure and organization of the Middleware Oriented Messaging (MOM) system. The project is divided into two main codebases: one in Go for the REST API, and one in C++ for the MOM core logic. Additionally, shared protocol definitions, deployment scripts, and documentation are included.

---

## 📁 go-api/

**Language:** Go  
**Purpose:** Implements the RESTful API used by clients to interact with the MOM system.

### Subfolders:
- `routes/`: Handlers for HTTP routes (topics, queues, messages, auth).
- `middleware/`: Logic for JWT authentication and request filtering.
- `grpc/`: gRPC client stubs to interact with MOM nodes.
- `config/`: Configuration and environment loading.
- `proto/`: Shared `.proto` files for gRPC (linked from root `/proto`).

---

## 📁 cpp-mom/

**Language:** C++  
**Purpose:** Contains the MOM node implementation responsible for core message processing, persistence, and inter-node replication.

### Subfolders:
- `src/`: Main entry point and execution flow for the node.
- `grpc_server/`: gRPC server implementation to receive commands from API or other nodes.
- `messaging/`: Management of topics, queues, and messages.
- `replication/`: Handles message replication to peer nodes.
- `persistence/`: Local database or file-based message storage.
- `utils/`: Utility functions (e.g., JSON, hashing, logging).
- `proto/`: Shared `.proto` files for gRPC (linked from root `/proto`).

---

## 📁 proto/

**Language:** Protocol Buffers  
**Purpose:** Defines shared interfaces and data structures for gRPC communication.

### Files:
- `mom_cluster.proto`: Core protocol definition used by both Go and C++.

Used to generate:
- gRPC client in Go (go-api)
- gRPC server in C++ (cpp-mom)

---

## 📁 deployment/

**Language:** Bash / Terraform / Docker / YAML  
**Purpose:** Contains infrastructure scripts for deploying the system to AWS.

### Contents:
- EC2 provisioning scripts
- `.pem` SSH access configuration
- VPC and security group templates
- Optional Dockerfiles or CI/CD configs

---

## 📁 docs/

**Purpose:** Contains documentation related to the analysis and design of the project.

### Files:
- `SRS.md`: Software Requirements Specification
- `SDD.md`: Software Design Document
- Architecture diagrams (UML, C4, PlantUML)
- Wiki-ready content for the GitHub project

---

Each module is independently maintainable, and all communication is handled via gRPC using the shared definitions in `/proto`. This design allows parallel development and a clean separation of responsibilities between components.