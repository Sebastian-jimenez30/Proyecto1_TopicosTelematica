Claro, aquí tienes un **texto para la página principal (Home) de la wiki** del repositorio, que proporciona una descripción general del proyecto y enlaza directamente a los documentos SRS y SDD:

---

# Middleware Orientado a Mensajes (MOM)

## Descripción General

Este proyecto consiste en el diseño e implementación de un **Middleware Orientado a Mensajes (MOM)**, desarrollado como parte del curso **ST0263 - Tópicos Especiales en Telemática** en la Universidad EAFIT (2025-1). El sistema tiene como objetivo permitir la **comunicación asincrónica entre aplicaciones distribuidas**, soportando tanto **colas** como **tópicos** para el intercambio de mensajes.

El MOM se despliega en un entorno **clúster** distribuido sobre **Amazon AWS**, integrando aspectos de **replicación, particionamiento, tolerancia a fallos, autenticación, persistencia** y escalabilidad. Está compuesto por dos tecnologías principales:

- **Go**: utilizado para implementar la API REST que permite la interacción del cliente con el middleware.
- **C++**: responsable de toda la lógica interna del MOM, incluyendo gestión de colas/tópicos, replicación entre nodos, almacenamiento y sincronización.

El sistema fue diseñado teniendo en cuenta principios de sistemas distribuidos modernos, inspirado en arquitecturas como Apache Kafka, con soporte para múltiples nodos MOM, detección de fallos, y consistencia eventual a través de gRPC.

---

## Documentación

Puedes consultar la documentación completa del proyecto en los siguientes enlaces:

- 📄 [Documento de Requerimientos de Software (SRS)](https://github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/wiki/SRS)
- 📐 [Documento de Diseño del Sistema (SDD)](https://github.com/Sebastian-jimenez30/Proyecto1_TopicosTelematica/wiki/SDD)

