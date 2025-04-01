
# Milestone 3: Replicación, Particionamiento y Tolerancia a Fallos

---

## FR8. Replicate Messages Between Nodes

### Description
**As** a MOM node,  
**I want** to replicate messages to other nodes,  
**So that** data is available even if a node fails.

### Acceptance Criteria
- **Given that** a message is stored locally, then it must be sent to all replica nodes.
- **Given that** the replicas confirm receipt, then the message can be considered fully persisted.
- **Given that** a replica is down, the system must retry or queue the replication.

### Tasks
- [x] Track replica nodes per partition
- [x] Send messages via gRPC to replicas
- [x] Confirm acknowledgment and retry on failure

---

## FR9. Support Stateless and Persistent Connections

### Description
**As** a MOM node,  
**I want** to handle both persistent and stateless client interactions,  
**So that** clients can interact according to their context.

### Acceptance Criteria
- **Given that** a client uses HTTP without sessions, the API must function using stateless requests.
- **Given that** a client maintains an open connection (e.g., streaming), the system must handle it gracefully.

### Tasks
- [x] Ensure REST API operates in a stateless manner
- [x] Optionally support persistent connections using websockets or gRPC streaming (optional)

---

## FR10. Monitor Nodes and Recover from Failures

### Description
**As** a MOM cluster,  
**I want** to detect when a node fails,  
**So that** leadership can be reassigned and service continues.

### Acceptance Criteria
- **Given that** a node stops responding, it must be marked as unavailable.
- **Given that** a node holding a leader role fails, then the system must promote a replica.
- **Given that** a failed node recovers, it must sync with the latest data.

### Tasks
- [x] Implement heartbeat mechanism between nodes
- [x] Detect node failure and trigger leader election
- [x] Sync state on node recovery

