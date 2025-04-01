
# Milestone 2: MOM Core en C++

---

## FR4. Process Incoming Messages Internally

### Description
**As** a MOM node,  
**I want** to process incoming messages,  
**So that** they are stored locally and prepared for replication.

### Acceptance Criteria
- **Given that** a message is received, then it must be assigned to the correct queue or topic partition.
- **Given that** the message is valid, then it must be stored locally in the node.
- **Given that** the node is not the leader for the resource, then it must forward the message to the correct leader node.

### Tasks
- [x] Implement routing of incoming messages by resource name
- [x] Store messages in local storage (SQLite or JSON)
- [x] Forward to leader node if necessary

---

## FR5. Retrieve Messages from Internal Storage

### Description
**As** a MOM node,  
**I want** to retrieve and serve messages from local storage,  
**So that** they can be consumed by clients through the API.

### Acceptance Criteria
- **Given that** a client requests messages from a queue or topic, the system must return the relevant data.
- **Given that** a message is retrieved from a queue, then it must be marked as consumed or removed.
- **Given that** a message is from a topic, it must be delivered to all subscribers without being deleted immediately.

### Tasks
- [x] Implement queue and topic access handlers
- [x] Support topic: deliver without delete
- [x] Support queue: deliver and delete (or mark consumed)

---

## FR8. gRPC Server for MOM–MOM Communication

### Description
**As** a MOM node,  
**I want** to expose a gRPC server,  
**So that** I can receive messages and synchronization requests from other nodes.

### Acceptance Criteria
- **Given that** a gRPC message is received for replication or discovery, then it must be processed appropriately.
- **Given that** a peer node sends a sync request, then it must receive the latest state or data needed.

### Tasks
- [x] Define gRPC service in `.proto` file
- [x] Implement gRPC server in C++
- [x] Handle incoming messages and replication commands

---

## FR10. Local Persistence of Messages and Resources

### Description
**As** a MOM node,  
**I want** to persist all relevant data locally,  
**So that** I can recover from crashes or reboots.

### Acceptance Criteria
- **Given that** a message or resource is created or modified, it must be written to disk.
- **Given that** the system restarts, then it must reload all persisted data to restore state.

### Tasks
- [x] Define schema or structure for local storage
- [x] Persist messages, queues, topics, and metadata
- [x] Load persisted data on startup

