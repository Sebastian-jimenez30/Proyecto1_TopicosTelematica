# Documento de Diseño del Sistema (SDD)

## Proyecto 1: Middleware Orientado a Mensajes (MOM)
Curso: ST0263 - Tópicos Especiales en Telemática  
Universidad EAFIT - 2025-1  
Fecha de entrega: 13 de abril de 2025

## 1. Arquitectura General

El sistema MOM será desarrollado con una arquitectura híbrida distribuida. Esta arquitectura no depende de un único nodo maestro. En su lugar, todos los nodos del sistema pueden asumir roles de liderazgo dependiendo del recurso que estén gestionando (una cola o un tópico). Esto permite alta disponibilidad, balanceo de carga, y una tolerancia a fallos sólida.

Cada nodo del sistema puede:

- Exponer una API REST para recibir solicitudes de clientes.
- Comunicarse con otros nodos mediante gRPC.
- Gestionar recursos como colas y tópicos.
- Almacenar mensajes en almacenamiento local.
- Replicar sus recursos en otros nodos para asegurar redundancia.

El sistema será desplegado en instancias EC2 dentro de Amazon AWS, permitiendo pruebas y demostración de despliegue en la nube.

## 2. Diagramas de Arquitectura

### 2.1 Diagrama de Componentes

Cada nodo MOM incluye los siguientes componentes internos:

- API REST Server: permite que los clientes interactúen con el middleware.
- Gestor de Recursos: maneja las operaciones de creación, eliminación y listado de colas y tópicos.
- Broker de Mensajes: administra el almacenamiento y entrega de mensajes.
- gRPC Peer Manager: facilita la comunicación entre nodos MOM.
- Replicador: se encarga de replicar mensajes y recursos en otros nodos.
- Persistencia Local: maneja la base de datos o archivos donde se almacena la información.

[Insertar aquí imagen del diagrama de componentes]

### 2.2 Diagrama de Despliegue

Cada nodo MOM se ejecuta en una instancia EC2 de AWS. Los clientes acceden a los nodos a través de direcciones IP públicas o mediante un balanceador de carga. Internamente, las instancias se comunican entre sí dentro de una VPC (Virtual Private Cloud) usando direcciones privadas. Los datos se almacenan en volúmenes EBS conectados a cada nodo.

[Insertar aquí imagen del diagrama de despliegue]

### 2.3 Diagrama de Contenedores (C4 - Nivel 2)

Este diagrama muestra cómo se separan los contenedores funcionales de alto nivel entre Go (API) y C++ (MOM), junto con sus tecnologías y roles.

[Insertar aquí imagen del diagrama de contenedores C4 nivel 2]

### 2.4 Diagrama de Arquitectura Técnica General

Representa los módulos funcionales, tecnologías clave, y su comunicación.

[Insertar aquí imagen del diagrama de arquitectura técnica general]

## 3. Diseño de APIs

### 3.1 API REST (Cliente ↔ Nodo MOM)

Los clientes interactúan con el sistema a través de una API RESTful que permite autenticación, gestión de recursos y envío/recepción de mensajes.

Operaciones soportadas:

- Autenticación: `/auth/login`
- Gestión de tópicos: `POST`, `GET`, `DELETE` en `/topics`
- Gestión de colas: `POST`, `GET`, `DELETE` en `/queues`
- Envío de mensajes: `POST` en `/topics/{name}/messages` o `/queues/{name}/messages`
- Recepción de mensajes: `GET` en `/topics/{name}/messages` o `/queues/{name}/messages`

### 3.2 API gRPC (Nodo MOM ↔ Nodo MOM)

Los nodos MOM se comunican entre sí usando gRPC para replicar mensajes, gestionar particiones y mantener el estado sincronizado. El archivo `mom_cluster.proto` define los siguientes servicios:

- `ReplicateTopic`
- `ReplicateQueue`
- `SyncMessage`
- `Heartbeat`
- `DiscoverResources`

Estas operaciones permiten la replicación segura y eficiente entre nodos distribuidos.

## 4. Especificación de Endpoints

Todos los endpoints REST utilizan tokens JWT para autenticación. Los mensajes se transmiten en formato JSON, y los códigos de estado HTTP siguen los estándares usuales (`200 OK`, `201 Created`, `403 Forbidden`, `404 Not Found`, etc.).

Ejemplos:

- Crear tópico:
  - `POST /topics`
  - Cuerpo: `{ "name": "notificaciones" }`
- Enviar mensaje:
  - `POST /topics/notificaciones/messages`
  - Cuerpo: `{ "payload": "Hola", "headers": {"tipo": "alerta"} }`

## 5. Modelo de Datos

El modelo de datos del sistema incluye las siguientes entidades:

- Usuario: nombre de usuario, contraseña hasheada, fecha de creación.
- Tópico: nombre único, creador, nodos donde está replicado.
- Cola: nombre único, creador, partición asignada.
- Mensaje: identificador, payload, headers, timestamp, estado, partición, réplicas.
- ReplicaStatus: estado de sincronización por nodo.
- Log de eventos: trazabilidad de acciones como envío de mensajes.

[Insertar aquí imagen del diagrama de clases UML]

Este modelo garantiza trazabilidad, control de propiedad y recuperación en caso de fallos.

## 6. Diagramas de Secuencia

### Envío de mensaje a tópico:

1. El cliente envía el mensaje a un nodo MOM.
2. El nodo valida si gestiona el tópico.
3. Si no lo gestiona, reenvía la solicitud al nodo correspondiente.
4. El nodo receptor almacena el mensaje y lo replica a otros nodos.
5. Se devuelve confirmación al cliente.

### Replicación entre nodos MOM:

1. Nodo líder recibe un mensaje.
2. Lo replica a las réplicas definidas vía gRPC.
3. Las réplicas almacenan el mensaje.
4. Responden al nodo líder con confirmación.
5. El nodo líder registra el mensaje como replicado exitosamente.

## 7. Estrategias de Particionamiento y Replicación

El sistema utiliza una estrategia basada en hash para asignar particiones:

`partition_id = hash(nombre_recurso) % total_particiones`

Cada recurso es gestionado por un nodo líder y replicado en al menos un nodo adicional. La replicación ocurre mediante gRPC y requiere confirmación antes de considerar un mensaje como entregado.

Esta estrategia permite escalabilidad, balanceo de carga y tolerancia a fallos.

## 8. Mecanismos de Tolerancia a Fallos

- Detección de fallos mediante heartbeats entre nodos.
- Failover automático a réplicas si un nodo falla.
- Replicación sincronizada antes de confirmar escritura.
- Persistencia local para recuperación tras reinicio.
- Reintentos y buffers en caso de caída de red.
- Logs de operación para auditoría y reconstrucción de estado.

## 9. Mecanismo de Autenticación y Autorización

La autenticación se realiza mediante JWT. Los usuarios deben iniciar sesión usando `/auth/login`, y luego adjuntar el token a cada solicitud.

La autorización se basa en propiedad: un usuario solo puede eliminar colas y tópicos que haya creado. El sistema registra cada acción sensible en un log de auditoría que incluye usuario, recurso, operación y estado.

Entre nodos MOM, la seguridad puede implementarse mediante autenticación mútua TLS (mTLS) o tokens compartidos.

## 10. Esquema de Persistencia

Cada nodo MOM mantiene su propia base de datos local, la cual puede implementarse con SQLite o estructuras JSON.

Datos persistidos incluyen:

- Usuarios
- Tópicos
- Colas
- Mensajes
- Logs de eventos

El sistema puede organizar los datos en carpetas por recurso, o emplear una base de datos embebida. En ambos casos, se asegura que los datos sobrevivan reinicios del sistema.

Adicionalmente, los logs de auditoría pueden ser centralizados en una base de datos para monitoreo y análisis, aunque esto es opcional.

Ejemplo de log de evento:

```json
{
  "event_type": "send_message",
  "user": "usuario1",
  "resource": "notificaciones",
  "timestamp": "2025-03-25T13:16:00Z",
  "details": {
    "status": "ok",
    "message_id": "uuid"
  }
}
```

Estos logs permiten trazabilidad completa y detección de errores o mal uso del sistema.

