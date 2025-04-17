#include "cola.hpp"
#include <iostream>

Cola::Cola(const std::string& nombre, Usuario& creador, Persistencia& persistencia)
    : nombre(nombre), creadorUsername(creador.getUsername()), db(persistencia) {

    db.crearCola(nombre, creadorUsername);
    db.autorizarUsuarioParaConsumir(creadorUsername, nombre);
    cargarMensajesDesdeBD();
}

const std::string& Cola::getNombre() const {
    return nombre;
}

const std::string& Cola::getCreadorUsername() const {
    return creadorUsername;
}

bool Cola::encolar(const Mensaje& mensaje) {
    mensajes.push(mensaje);
    return db.guardarMensaje(mensaje);
}

bool Cola::hayMensajes() const {
    return !mensajes.empty();
}

Mensaje Cola::desencolar() {
    if (mensajes.empty()) {
        throw std::runtime_error("❌ No hay mensajes en la cola.");
    }

    Mensaje m = mensajes.front();
    mensajes.pop();
    return m;
}

bool Cola::puedeEliminar(const Usuario& usuario) const {
    return usuario.getUsername() == creadorUsername;
}

bool Cola::eliminar(const Usuario& usuario) {
    if (puedeEliminar(usuario)) {
        return db.eliminarCola(nombre, usuario.getUsername());
    }
    return false;
}

bool Cola::autorizarUsuario(const std::string& usernameAutorizado) {
    return db.autorizarUsuarioParaConsumir(usernameAutorizado, nombre);
}

bool Cola::usuarioPuedeConsumir(const Usuario& usuario) const {
    return db.verificarAutorizacion(usuario.getUsername(), nombre);
}

void Cola::cargarMensajesDesdeBD() {
    std::vector<Mensaje> historico = db.cargarMensajesPorCanal(nombre, "cola");
    for (const auto& msg : historico) {
        mensajes.push(msg);
    }
}
