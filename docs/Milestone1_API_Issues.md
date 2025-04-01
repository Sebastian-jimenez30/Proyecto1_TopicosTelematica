
# Milestone 1: Módulo API y Autenticación

---

## FR1. User Authentication (Login)

### Description
**As** a user,  
**I want** to authenticate using my credentials,  
**So that** I can securely access the messaging system.

### Acceptance Criteria
- **Given that** a user provides valid credentials, when they log in, then the system must issue a JWT token.
- **Given that** credentials are invalid, when login is attempted, then the system must reject the request with an appropriate error.
- **Given that** a user is authenticated, when making API calls, then they must include a valid JWT token.

### Tasks
- [x] Implement login endpoint `/auth/login`
- [x] Generate JWT tokens upon successful login
- [x] Validate JWT on protected routes

---

## FR2. Create, List, and Delete Topics

### Description
**As** an authenticated user,  
**I want** to create, list and delete messaging topics,  
**So that** I can manage logical channels for pub-sub messaging.

### Acceptance Criteria
- **Given that** a user is authenticated, when they send a POST to `/topics`, then a new topic is created and associated to their account.
- **Given that** a user accesses `/topics`, then the system must return all topics they created.
- **Given that** a user sends a DELETE request to `/topics/{name}`, then the topic must be deleted only if it belongs to them.

### Tasks
- [x] Implement `POST /topics`
- [x] Implement `GET /topics`
- [x] Implement `DELETE /topics/{name}`
- [x] Validate topic ownership before deletion

---

## FR3. Create, List, and Delete Queues

### Description
**As** an authenticated user,  
**I want** to create, list and delete message queues,  
**So that** I can manage point-to-point message delivery.

### Acceptance Criteria
- **Given that** a user is authenticated, when they create a queue, then it must be stored and owned by them.
- **Given that** a user accesses `/queues`, then only their queues must be listed.
- **Given that** a user attempts to delete another user's queue, then the system must deny access.

### Tasks
- [x] Implement `POST /queues`
- [x] Implement `GET /queues`
- [x] Implement `DELETE /queues/{name}`
- [x] Validate queue ownership before deletion

---

## FR4. Send Messages to Topics and Queues

### Description
**As** an authenticated user,  
**I want** to send messages to topics and queues,  
**So that** they are delivered to subscribers or consumers.

### Acceptance Criteria
- **Given that** a user posts to `/topics/{name}/messages`, then the system must send the message to the correct node for that topic.
- **Given that** a user posts to `/queues/{name}/messages`, then the system must enqueue the message correctly.
- **Given that** a topic or queue does not exist or is not owned by the user, then the system must reject the request.

### Tasks
- [x] Implement `POST /topics/{name}/messages`
- [x] Implement `POST /queues/{name}/messages`
- [x] Validate resource ownership and routing

---

## FR5. Receive Messages from Topics and Queues

### Description
**As** an authenticated user,  
**I want** to retrieve messages from my queues or topics,  
**So that** I can process or display received data.

### Acceptance Criteria
- **Given that** a user calls `GET /topics/{name}/messages`, then all available messages should be returned (pub-sub style).
- **Given that** a user calls `GET /queues/{name}/messages`, then a single message should be dequeued.
- **Given that** the resource does not exist or the user is not the owner, the system must return an error.

### Tasks
- [x] Implement `GET /topics/{name}/messages`
- [x] Implement `GET /queues/{name}/messages`
- [x] Handle message retrieval logic

---

## FR6. Associate Every Action with an Authenticated User

### Description
**As** a system administrator,  
**I want** every operation to be associated with the user performing it,  
**So that** we can ensure proper access control and accountability.

### Acceptance Criteria
- **Given that** a request is received, when it includes a valid token, then the system must extract the user identity.
- **Given that** a resource is being created or deleted, then it must be tagged to the authenticated user.
- **Given that** a user accesses resources, then the system must filter by ownership.

### Tasks
- [x] Extract user identity from JWT
- [x] Associate resource creation with user
- [x] Enforce ownership in read/write/delete actions

---

## FR7. RESTful API Exposure

### Description
**As** a developer or tester,  
**I want** the middleware to expose a well-defined RESTful API,  
**So that** I can interact with it using standard HTTP tools.

### Acceptance Criteria
- **Given that** a request is made to a valid endpoint, then the system must respond using standard HTTP semantics.
- **Given that** an invalid method or path is accessed, then an appropriate error (404 or 405) must be returned.
- **Given that** the server starts, it must expose documented endpoints.

### Tasks
- [x] Set up base HTTP server with Go
- [x] Configure all endpoints with proper routing
- [x] Add basic OpenAPI or README-based documentation
