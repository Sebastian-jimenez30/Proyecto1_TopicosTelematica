#include "topic.hpp"
#include <iostream>

Topico::Topico(const std::string& nombre, Usuario& creador, Persistencia& persistencia)
    : nombre(nombre), creadorUsername(creador.getUsername()), db(persistencia) {

    if (!db.crearTopico(nombre, creadorUsername)) {
        std::cerr << "l tópico '" << nombre << "' ya existe o hubo un error.\n";
    }

    suscribir(creador);
}

const std::string& Topico::getNombre() const {
    return nombre;
}

const std::string& Topico::getCreadorUsername() const {
    return creadorUsername;
}

bool Topico::publicarMensaje(const Mensaje& mensaje) {
    return db.guardarMensaje(mensaje);
}

bool Topico::suscribir(Usuario& usuario) {
    return db.suscribirUsuarioATopico(usuario.getUsername(), nombre);
}

std::vector<Mensaje> Topico::obtenerMensajesPara(Usuario& usuario, int maxMensajes) {
    
    int offset = db.obtenerOffset(usuario.getUsername(), nombre);
    std::vector<Mensaje> mensajes = db.cargarMensajesDesdeOffset(nombre, offset, maxMensajes);

    if (!mensajes.empty()) {
        db.actualizarOffset(usuario.getUsername(), nombre, offset + mensajes.size());
    }

    return mensajes;
}

bool Topico::puedeEliminar(const Usuario& usuario) const {
    return usuario.getUsername() == creadorUsername;
}

bool Topico::eliminar(const Usuario& usuario) {
    if (!puedeEliminar(usuario)) return false;
    return db.eliminarTopico(nombre, usuario.getUsername());
}
