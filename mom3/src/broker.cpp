#include "broker.hpp"
#include "mensaje.hpp"
#include <iostream>
#include <chrono>
#include <sstream>
#include <iomanip>     
#include <ctime> 

Broker::Broker(const std::string& db_path)
    : persistencia(db_path) {
    persistencia.inicializarBaseDeDatos();
}

// ======================= USUARIOS =======================

bool Broker::registrarUsuario(const std::string& username, const std::string& password) {
    return persistencia.crearUsuario(username, password);
}

bool Broker::autenticarUsuario(const std::string& username, const std::string& password, std::string& token_out) {
    if (persistencia.verificarCredenciales(username, password)) {
        Usuario user(username);
        token_out = user.generarToken();

        // 🔁 Conversión de time_point a string
        std::time_t exp_time = std::chrono::system_clock::to_time_t(user.getExpiracion());
        std::stringstream ss;
        ss << std::put_time(std::localtime(&exp_time), "%F %T");
        std::string expiracion = ss.str();

        persistencia.guardarToken(username, token_out, expiracion);
        return true;
    }
    return false;
}


bool Broker::verificarToken(const std::string& token, std::string& username_out) {
    if (!Usuario::verificarToken(token, username_out)) return false;
    return persistencia.validarTokenActivo(token);
}

bool Broker::autenticarYObtenerUsuario(const std::string& token, Usuario& usuario_out) {
    std::string username;
    if (!verificarToken(token, username)) return false;
    usuario_out = Usuario(username);
    return true;
}

bool Broker::guardarTokenReplica(const std::string& username, const std::string& token, const std::string& expiracion) {
    return persistencia.guardarToken(username, token, expiracion);
}

// ======================= COLAS =======================

bool Broker::crearCola(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    if (colas.find(nombre) == colas.end()) {
        auto nueva = std::make_shared<Cola>(nombre, usuario, persistencia);
        colas[nombre] = nueva;
        persistencia.registrarLog(usuario.getUsername(), "Creó cola '" + nombre + "'");
        return true;
    }
    return false;
}

bool Broker::eliminarCola(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = colas.find(nombre);
    if (it != colas.end() && it->second->puedeEliminar(usuario)) {
        colas.erase(it);
        persistencia.eliminarCola(nombre, usuario.getUsername());
        persistencia.registrarLog(usuario.getUsername(), "Eliminó cola '" + nombre + "'");
        return true;
    }
    return false;
}

bool Broker::autorizarCola(const std::string& nombre, const std::string& usernameObjetivo, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = colas.find(nombre);
    if (it != colas.end() && it->second->puedeEliminar(usuario)) {
        return it->second->autorizarUsuario(usernameObjetivo);
    }
    return false;
}

bool Broker::enviarMensajeACola(const std::string& nombre, const std::string& contenido, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = colas.find(nombre);
    if (it != colas.end()) {
        Mensaje msg(usuario.getUsername(), contenido, nombre, "cola", std::chrono::system_clock::now());
        bool ok = it->second->encolar(msg);
        if (ok) persistencia.registrarLog(usuario.getUsername(), "Encoló mensaje en '" + nombre + "'");
        return ok;
    }
    return false;
}

std::optional<Mensaje> Broker::consumirMensajeDeCola(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return std::nullopt;

    auto it = colas.find(nombre);
    if (it != colas.end() && it->second->usuarioPuedeConsumir(usuario) && it->second->hayMensajes()) {
        Mensaje m = it->second->desencolar();
        persistencia.registrarLog(usuario.getUsername(), "Consumió mensaje de '" + nombre + "'");
        return m;
    }
    return std::nullopt;
}

// ======================= TÓPICOS =======================

bool Broker::crearTopico(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    if (topicos.find(nombre) == topicos.end()) {
        auto nuevo = std::make_shared<Topico>(nombre, usuario, persistencia);
        topicos[nombre] = nuevo;
        persistencia.registrarLog(usuario.getUsername(), "Creó tópico '" + nombre + "'");
        return true;
    }
    return false;
}

bool Broker::eliminarTopico(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = topicos.find(nombre);
    if (it != topicos.end() && it->second->puedeEliminar(usuario)) {
        it->second->eliminar(usuario);
        topicos.erase(it);
        persistencia.registrarLog(usuario.getUsername(), "Eliminó tópico '" + nombre + "'");
        return true;
    }
    return false;
}

bool Broker::suscribirATopico(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = topicos.find(nombre);
    if (it != topicos.end()) {
        bool ok = it->second->suscribir(usuario);
        if (ok) persistencia.registrarLog(usuario.getUsername(), "Se suscribió a tópico '" + nombre + "'");
        return ok;
    }
    return false;
}

bool Broker::publicarEnTopico(const std::string& nombre, const std::string& contenido, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return false;

    auto it = topicos.find(nombre);
    if (it != topicos.end()) {
        Mensaje msg(usuario.getUsername(), contenido, nombre, "topico", std::chrono::system_clock::now());
        bool ok = it->second->publicarMensaje(msg);
        if (ok) persistencia.registrarLog(usuario.getUsername(), "Publicó en tópico '" + nombre + "'");
        return ok;
    }
    return false;
}

std::vector<Mensaje> Broker::consumirDesdeTopico(const std::string& nombre, const std::string& token) {
    Usuario usuario("");
    if (!autenticarYObtenerUsuario(token, usuario)) return {};

    auto it = topicos.find(nombre);
    if (it != topicos.end()) {
        auto mensajes = it->second->obtenerMensajesPara(usuario);
        if (!mensajes.empty()) {
            persistencia.registrarLog(usuario.getUsername(), "Consumió mensajes desde tópico '" + nombre + "'");
        }
        return mensajes;
    }
    return {};
}

std::vector<std::string> Broker::listarTopicos() {
    return persistencia.obtenerNombresTopicos();
}

std::vector<std::string> Broker::listarColas() {
    return persistencia.obtenerNombresColas();
}
