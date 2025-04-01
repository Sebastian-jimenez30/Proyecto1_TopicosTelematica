#  Documento de Requerimientos de Software (SRS)

## Proyecto 1: Middleware Orientado a Mensajes (MOM)  
**Curso:** ST0263 - Tópicos Especiales en Telemática  
**Universidad EAFIT - 2025-1**  
**Fecha de entrega:** 13 de abril de 2025

---

## 1. Introducción

El presente documento detalla los requerimientos funcionales y no funcionales del proyecto de desarrollo de un Middleware Orientado a Mensajes (MOM). Este middleware facilitará la comunicación asincrónica entre aplicaciones distribuidas, soportando tópicos y colas como mecanismos de entrega de mensajes. El sistema estará desplegado en un entorno clúster en Amazon AWS y permitirá tanto conexiones REST (cliente–MOM) como gRPC (MOM–MOM).

---

## 2. Descripción general

El sistema será capaz de:
- Autenticar usuarios.
- Gestionar colas y tópicos (crear, listar, eliminar).
- Enviar y recibir mensajes.
- Replicar y particionar recursos en un clúster de nodos MOM.
- Garantizar disponibilidad y tolerancia a fallos.
- Operar de forma transparente y segura.

---

## 3. Requerimientos del sistema

### 3.1 Requerimientos funcionales

| Código | Descripción |
|--------|-------------|
| RF1 | Autenticación de usuarios para acceder al sistema. |
| RF2 | Crear, listar y eliminar tópicos. |
| RF3 | Crear, listar y eliminar colas. |
| RF4 | Enviar mensajes a tópicos y colas. |
| RF5 | Recibir mensajes desde tópicos y colas. |
| RF6 | Asociar toda acción con un usuario autenticado. |
| RF7 | Exponer API REST para comunicación cliente–MOM. |
| RF8 | Usar gRPC para comunicación MOM–MOM. |
| RF9 | Soportar conexiones persistentes y sin estado. |
| RF10 | Implementar persistencia de mensajes. |

---

### 3.2 Requerimientos no funcionales

| Código | Descripción |
|--------|-------------|
| RNF1 | Cifrado en la comunicación (opcional pero recomendado). |
| RNF2 | Implementar tolerancia a fallos (servidores y mensajes). |
| RNF3 | Permitir escalabilidad horizontal con múltiples nodos. |
| RNF4 | Soportar replicación y particionamiento de colas/tópicos. |
| RNF5 | Soporte para múltiples usuarios simultáneos. |
| RNF6 | Transparencia frente a ubicación y cantidad de nodos. |
| RNF7 | Diseño extensible para incorporar nuevas funcionalidades. |
| RNF8 | Procesamiento eficiente y baja latencia. |
| RNF9 | Despliegue en entorno virtualizado AWS. |

---

## 4. Casos de uso

Ver diagrama UML en la documentación gráfica. A continuación, se describen los principales casos de uso del sistema agrupados por actor:

### Actor: Usuario
- Autenticarse.
- Crear, listar y eliminar colas.
- Crear, listar y eliminar tópicos.
- Enviar y recibir mensajes.
- Establecer conexión persistente o sin estado.

### Actor: Nodo MOM
- Replicar colas y tópicos.
- Particionar mensajes.
- Recuperarse de fallos.

[Insertar aquí imagen del diagrama de casos de uso UML]

---

## 5. Reglas de negocio

1. Solo usuarios autenticados pueden operar sobre el sistema.
2. Un usuario solo puede eliminar recursos que haya creado.
3. Toda acción debe asociarse a un usuario.
4. Los nombres de colas/tópicos deben ser únicos.
5. Mensajes deben tener control de acceso.
6. Los mensajes en colas son consumidos una sola vez.
7. Los mensajes en tópicos son entregados a múltiples suscriptores.
8. Cada partición debe tener al menos una réplica.
9. Fallos de nodos no deben implicar pérdida de mensajes.

---

## 6. Restricciones

- La comunicación cliente–MOM será por REST.
- La comunicación MOM–MOM será por gRPC.
- Se debe desplegar en AWS Academy con VMs.
- El código estará en GitHub, con roles bien definidos.
- No se requiere interfaz gráfica.
- La app de prueba será basada en el ejemplo de RabbitMQ.
- Se deben documentar todas las fases del desarrollo.
- La entrega final debe incluir IP pública y archivo `.pem`.

---

## 7. Identificación de usuarios y roles

### Usuario Cliente
- Usa API REST para enviar/recibir mensajes.
- Crea y gestiona recursos propios (colas, tópicos).
- Accede solo a su información.

### Nodo MOM
- Usa gRPC para replicar, particionar y sincronizar nodos.
- Administra colas/tópicos internamente.
- Gestiona la tolerancia a fallos y disponibilidad.

### (Opcional) Administrador
- Supervisión general del sistema.
- Revisión de logs y métricas.
- Acciones administrativas o de emergencia.

---

## 8. Diagramas complementarios

### Diagrama de Clases (UML)
Representa las entidades principales del modelo de datos y sus relaciones.
[Insertar aquí imagen del diagrama de clases UML]

### Diagrama de Actividades
Describe el flujo detallado de procesos como el envío de mensajes.
[Insertar aquí imagen del diagrama de actividades]

### Diagrama de Contexto (C4 - Nivel 1)
Visualiza la relación del sistema MOM con los actores y sistemas externos.
[Insertar aquí imagen del diagrama de contexto C4 nivel 1]

